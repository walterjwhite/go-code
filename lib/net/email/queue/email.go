package queue

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log/slog"
	"net/smtp"
	"strings"
	"sync"
	"time"

	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/mail"
)


type IMAPConfig struct {
	Host     string // e.g. "imap.example.com"
	Port     int    // typically 993 (TLS) or 143 (STARTTLS)
	Username string
	Password string
	Mailbox  string // folder to watch, e.g. "INBOX"
	TLS      bool   // true → implicit TLS on connect; false → STARTTLS
}

type SMTPConfig struct {
	Host     string // e.g. "smtp.example.com"
	Port     int    // typically 587 (STARTTLS) or 465 (TLS)
	Username string
	Password string
	From     string // envelope / header From address
}


type Client struct {
	smtp SMTPConfig
	icfg IMAPConfig

	log *slog.Logger

	mu      sync.Mutex
	imapC   *imapclient.Client
	running bool
}

func New(imapCfg IMAPConfig, smtpCfg SMTPConfig, logger *slog.Logger) *Client {
	if logger == nil {
		logger = slog.Default()
	}
	return &Client{
		icfg: imapCfg,
		smtp: smtpCfg,
		log:  logger,
	}
}


type MessageHandler func(payload []byte) error

func (c *Client) Listen(ctx context.Context, handler MessageHandler) error {
	c.mu.Lock()
	if c.running {
		c.mu.Unlock()
		return fmt.Errorf("emailclient: Listen is already running")
	}
	c.running = true
	c.mu.Unlock()

	defer func() {
		c.mu.Lock()
		c.running = false
		c.mu.Unlock()
	}()

	backoff := 2 * time.Second
	const maxBackoff = 2 * time.Minute

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		c.log.Info("emailclient: connecting to IMAP", "host", c.icfg.Host)
		err := c.runIDLELoop(ctx, handler)
		if err == nil || ctx.Err() != nil {
			return ctx.Err()
		}

		c.log.Warn("emailclient: IMAP error, will reconnect",
			"err", err, "backoff", backoff)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
		}
		backoff = min(backoff*2, maxBackoff)
	}
}

func (c *Client) runIDLELoop(ctx context.Context, handler MessageHandler) error {
	newMail := make(chan struct{}, 1)
	unilateralHandler := &imapclient.UnilateralDataHandler{
		Mailbox: func(data *imapclient.UnilateralDataMailbox) {
			select {
			case newMail <- struct{}{}:
			default:
			}
		},
	}

	imapClient, err := c.dialWithHandler(unilateralHandler)
	if err != nil {
		return err
	}
	defer func() {
		if err := imapClient.Close(); err != nil {
			c.log.Warn("emailclient: IMAP close error", "err", err)
		}
	}()

	c.mu.Lock()
	c.imapC = imapClient
	c.mu.Unlock()

	if _, err = imapClient.Select(c.icfg.Mailbox, nil).Wait(); err != nil {
		return fmt.Errorf("SELECT %s: %w", c.icfg.Mailbox, err)
	}

	if err = c.fetchAndHandle(ctx, imapClient, handler); err != nil {
		return err
	}

	for {
		if err = ctx.Err(); err != nil {
			return nil // context cancelled, clean exit
		}

		idleCmd, err := imapClient.Idle()
		if err != nil {
			return fmt.Errorf("IDLE: %w", err)
		}
		c.log.Debug("emailclient: IDLE – waiting for new mail")

		select {
		case <-ctx.Done():
			_ = idleCmd.Close()
			return nil
		case <-newMail:
		}

		if err = idleCmd.Close(); err != nil {
			return fmt.Errorf("IDLE close: %w", err)
		}

		if err = c.fetchAndHandle(ctx, imapClient, handler); err != nil {
			return err
		}
	}
}

func (c *Client) fetchAndHandle(
	ctx context.Context,
	cl *imapclient.Client,
	handler MessageHandler,
) error {
	searchData, err := cl.Search(&imap.SearchCriteria{
		NotFlag: []imap.Flag{imap.FlagSeen},
	}, nil).Wait()
	if err != nil {
		return fmt.Errorf("SEARCH UNSEEN: %w", err)
	}
	if len(searchData.AllSeqNums()) == 0 {
		return nil
	}

	seqSet := new(imap.SeqSet)
	for _, num := range searchData.AllSeqNums() {
		seqSet.AddNum(num)
	}

	fetchCmd := cl.Fetch(seqSet, &imap.FetchOptions{
		BodySection: []*imap.FetchItemBodySection{{}}, // entire RFC 5322 message
	})
	defer func() {
		if err := fetchCmd.Close(); err != nil {
			c.log.Warn("emailclient: fetch close error", "err", err)
		}
	}()

	for {
		msg := fetchCmd.Next()
		if msg == nil {
			break
		}
		if err = c.processMessage(ctx, cl, msg, handler); err != nil {
			c.log.Error("emailclient: handler error", "err", err)
		}
	}
	return fetchCmd.Close()
}

func (c *Client) processMessage(
	ctx context.Context,
	cl *imapclient.Client,
	msg *imapclient.FetchMessageData,
	handler MessageHandler,
) error {
	var rawBody []byte

	for {
		item := msg.Next()
		if item == nil {
			break
		}
		if bs, ok := item.(imapclient.FetchItemDataBodySection); ok {
			data, err := io.ReadAll(bs.Literal)
			if err != nil {
				return fmt.Errorf("read body: %w", err)
			}
			rawBody = data
			break // we only requested one section
		}
	}

	payload := extractTextPayload(rawBody)

	var wg sync.WaitGroup
	wg.Add(1)
	var handlerErr error
	go func() {
		defer wg.Done()
		handlerErr = handler(payload)
	}()

	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
	}

	if handlerErr != nil {
		return fmt.Errorf("handler returned error: %w", handlerErr)
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(msg.SeqNum)
	storeCmd := cl.Store(seqSet, &imap.StoreFlags{
		Op:    imap.StoreFlagsAdd,
		Flags: []imap.Flag{imap.FlagDeleted},
	}, nil)
	if err := storeCmd.Close(); err != nil {
		return fmt.Errorf("STORE \\Deleted: %w", err)
	}
	expungeCmd := cl.Expunge()
	if err := expungeCmd.Close(); err != nil {
		return fmt.Errorf("EXPUNGE: %w", err)
	}
	return nil
}

func (c *Client) dialWithHandler(unilateralHandler *imapclient.UnilateralDataHandler) (*imapclient.Client, error) {
	addr := fmt.Sprintf("%s:%d", c.icfg.Host, c.icfg.Port)
	tlsCfg := &tls.Config{ServerName: c.icfg.Host}

	var (
		cl  *imapclient.Client
		err error
	)
	opts := &imapclient.Options{
		TLSConfig:             tlsCfg,
		UnilateralDataHandler: unilateralHandler,
	}

	if c.icfg.TLS {
		cl, err = imapclient.DialTLS(addr, opts)
	} else {
		cl, err = imapclient.DialStartTLS(addr, opts)
	}
	if err != nil {
		return nil, fmt.Errorf("IMAP dial %s: %w", addr, err)
	}

	if err = cl.Login(c.icfg.Username, c.icfg.Password).Wait(); err != nil {
		if closeErr := cl.Close(); closeErr != nil {
			c.log.Warn("emailclient: IMAP close error", "err", closeErr)
		}
		return nil, fmt.Errorf("IMAP LOGIN: %w", err)
	}
	return cl, nil
}

func extractTextPayload(raw []byte) []byte {
	if len(raw) == 0 {
		return raw
	}
	mr, err := mail.CreateReader(strings.NewReader(string(raw)))
	if err != nil {
		return raw
	}
	for {
		part, err := mr.NextPart()
		if err != nil {
			break
		}
		if _, ok := part.Header.(*mail.InlineHeader); ok {
			ct := part.Header.Get("Content-Type")
			if strings.HasPrefix(ct, "text/plain") {
				data, err := io.ReadAll(part.Body)
				if err == nil {
					return data
				}
			}
		}
	}
	return raw // fallback: return the whole raw message
}


func (c *Client) PublishMessage(to []string, subject string, payload []byte) error {
	addr := fmt.Sprintf("%s:%d", c.smtp.Host, c.smtp.Port)
	auth := smtp.PlainAuth("", c.smtp.Username, c.smtp.Password, c.smtp.Host)

	header := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n",
		c.smtp.From,
		strings.Join(to, ", "),
		subject,
	)
	body := append([]byte(header), payload...)

	if err := smtp.SendMail(addr, auth, c.smtp.From, to, body); err != nil {
		return fmt.Errorf("emailclient: SMTP send: %w", err)
	}
	return nil
}

