package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
	chatlib "github.com/walterjwhite/go-code/lib/chat"
)


const (
	maxChatMessageLength    = 4096
	maxChatChannelLength    = 128
	maxChatCorrelationIDLen = 128
	maxChatMetadataKeys     = 10
	maxChatMetadataKeyLen   = 64
	maxChatMetadataValueLen = 256

	wsMaxMessageSize = 32 * 1024 // 32 KB
	wsWriteWait      = 10 * time.Second
	wsPongWait       = 60 * time.Second
	wsPingPeriod     = (wsPongWait * 9) / 10
	wsSendBuffer     = 32

	metadataKeyIP         = "ip"
	metadataKeyReplyToken = "reply_token"
)


var upgrader = websocket.Upgrader{
	HandshakeTimeout: 10 * time.Second,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	CheckOrigin: func(r *http.Request) bool {
		return isAllowedWebSocketOrigin(r)
	},
}


type wsInbound struct {
	Message       string            `json:"message"`
	Channel       string            `json:"channel,omitempty"`
	CorrelationID string            `json:"correlation_id,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
}

type wsOutbound struct {
	Type          string            `json:"type"`
	CorrelationID string            `json:"correlation_id,omitempty"`
	Message       string            `json:"message,omitempty"`
	Metadata      map[string]string `json:"metadata,omitempty"`
	SentAt        *time.Time        `json:"sent_at,omitempty"`
}


type chatService interface {
	Send(ctx context.Context, message chatlib.ClientMessage) error
	Responses() <-chan chatlib.ServerMessage
	Errors() <-chan error
	Start(ctx context.Context)
}

var chatClient chatService


func onChatWSRequest(c *gin.Context) {
	if chatClient == nil {
		c.AbortWithStatus(http.StatusServiceUnavailable)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error().Err(err).Str("ip", c.ClientIP()).Msg("websocket upgrade failed")
		return
	}


	ctx, cancel := context.WithCancel(context.Background())
	s := &wsSession{
		conn:   conn,
		ip:     c.ClientIP(),
		send:   make(chan []byte, wsSendBuffer),
		ctx:    ctx,
		cancel: cancel,
	}

	log.Info().Str("ip", s.ip).Msg("websocket session opened")
	chatHub.register(s)

	go s.writePump()
	s.readPump() // blocks until the connection closes
}


type wsSession struct {
	conn   *websocket.Conn
	ip     string
	send   chan []byte
	ctx    context.Context
	cancel context.CancelFunc
}


func (s *wsSession) readPump() {
	defer func() {
		chatHub.unregister(s)
		s.cancel()
		close(s.send)
		_ = s.conn.Close()
		log.Info().Str("ip", s.ip).Msg("websocket session closed")
	}()

	s.conn.SetReadLimit(wsMaxMessageSize)
	s.resetReadDeadline()
	s.conn.SetPongHandler(func(string) error {
		s.resetReadDeadline()
		return nil
	})

	for {
		_, raw, err := s.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err,
				websocket.CloseGoingAway,
				websocket.CloseNormalClosure,
			) {
				log.Warn().Err(err).Str("ip", s.ip).Msg("websocket read error")
			}
			return
		}
		s.handleMessage(raw)
	}
}

func (s *wsSession) resetReadDeadline() {
	s.conn.SetReadDeadline(time.Now().Add(wsPongWait)) //nolint:errcheck
}


func (s *wsSession) writePump() {
	ticker := time.NewTicker(wsPingPeriod)
	defer ticker.Stop()

	for {
		select {
		case msg, ok := <-s.send:
			s.conn.SetWriteDeadline(time.Now().Add(wsWriteWait)) //nolint:errcheck
			if !ok {
				s.conn.WriteMessage(websocket.CloseMessage, []byte{}) //nolint:errcheck
				return
			}
			if err := s.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Warn().Err(err).Str("ip", s.ip).Msg("websocket write error")
				return
			}

		case <-ticker.C:
			s.conn.SetWriteDeadline(time.Now().Add(wsWriteWait)) //nolint:errcheck
			if err := s.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case <-s.ctx.Done():
			return
		}
	}
}


func (s *wsSession) handleMessage(raw []byte) {
	var req wsInbound
	if err := json.Unmarshal(raw, &req); err != nil {
		s.sendError("", "invalid JSON payload")
		return
	}

	req.Message = strings.TrimSpace(req.Message)
	req.Channel = strings.TrimSpace(req.Channel)
	req.CorrelationID = strings.TrimSpace(req.CorrelationID)

	if req.Message == "" {
		s.sendError(req.CorrelationID, "message is required")
		return
	}
	if err := validateChatFields(&req); err != nil {
		s.sendError(req.CorrelationID, err.Error())
		return
	}

	route, err := s.newPendingRoute(req)
	if err != nil {
		log.Error().Err(err).Str("ip", s.ip).Msg("failed to create chat route")
		s.sendError(req.CorrelationID, "failed to queue message")
		return
	}

	err = chatClient.Send(s.ctx, chatlib.ClientMessage{
		IP:            s.ip,
		Message:       req.Message,
		Channel:       req.Channel,
		CorrelationID: route.requestID,
		Metadata:      route.metadata,
		SentAt:        time.Now().UTC(),
	})
	if err != nil {
		s.handleSendError(req.CorrelationID, err)
		return
	}

	chatHub.track(route)
	s.sendAck(req.CorrelationID)
}

func (s *wsSession) handleSendError(correlationID string, err error) {
	switch {
	case errors.Is(err, chatlib.ErrRateLimited),
		errors.Is(err, chatlib.ErrEmptyIP),
		errors.Is(err, chatlib.ErrEmptyMessage),
		errors.Is(err, chatlib.ErrNilContext):
		s.sendError(correlationID, err.Error())
	default:
		log.Error().Err(err).Str("ip", s.ip).Msg("chat send failed")
		s.sendError(correlationID, "failed to queue message")
	}
}


func (s *wsSession) sendAck(correlationID string) {
	s.enqueue(wsOutbound{Type: "ack", CorrelationID: correlationID})
}

func (s *wsSession) sendError(correlationID, message string) {
	s.enqueue(wsOutbound{Type: "error", CorrelationID: correlationID, Message: message})
}

func (s *wsSession) sendServerMessage(msg chatlib.ServerMessage) {
	s.enqueue(wsOutbound{
		Type:          "response",
		CorrelationID: msg.CorrelationID,
		Message:       msg.Message,
		Metadata:      msg.Metadata,
		SentAt:        &msg.SentAt,
	})
}

func (s *wsSession) enqueue(v wsOutbound) {
	b, err := json.Marshal(v)
	if err != nil {
		log.Error().Err(err).Msg("marshal outbound websocket message")
		return
	}
	select {
	case s.send <- b:
	default:
		log.Warn().Str("ip", s.ip).Msg("websocket send buffer full — closing session")
		s.cancel()
	}
}


func validateChatFields(r *wsInbound) error {
	if utf8.RuneCountInString(r.Message) > maxChatMessageLength {
		return fmt.Errorf("message exceeds maximum length of %d characters", maxChatMessageLength)
	}
	if utf8.RuneCountInString(r.Channel) > maxChatChannelLength {
		return fmt.Errorf("channel exceeds maximum length of %d characters", maxChatChannelLength)
	}
	if utf8.RuneCountInString(r.CorrelationID) > maxChatCorrelationIDLen {
		return fmt.Errorf("correlation_id exceeds maximum length of %d characters", maxChatCorrelationIDLen)
	}
	if len(r.Metadata) > maxChatMetadataKeys {
		return fmt.Errorf("metadata exceeds maximum of %d keys", maxChatMetadataKeys)
	}
	for k, v := range r.Metadata {
		if isReservedMetadataKey(k) {
			return fmt.Errorf("metadata key %q is reserved", k)
		}
		if utf8.RuneCountInString(k) > maxChatMetadataKeyLen {
			return fmt.Errorf("metadata key exceeds maximum length of %d characters", maxChatMetadataKeyLen)
		}
		if utf8.RuneCountInString(v) > maxChatMetadataValueLen {
			return fmt.Errorf("metadata value for key %q exceeds maximum length of %d characters", k, maxChatMetadataValueLen)
		}
	}
	return nil
}

func isAllowedWebSocketOrigin(r *http.Request) bool {
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		return true
	}

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	allowedHosts := allowedWebSocketHosts(r)
	for _, host := range allowedHosts {
		if strings.EqualFold(u.Host, host) {
			return true
		}
	}

	return false
}

func allowedWebSocketHosts(r *http.Request) []string {
	host := strings.TrimSpace(r.Host)
	if host == "" {
		return nil
	}
	return []string{host}
}

func isReservedMetadataKey(key string) bool {
	switch key {
	case metadataKeyIP, metadataKeyReplyToken:
		return true
	default:
		return false
	}
}

type pendingRoute struct {
	requestID           string
	replyToken          string
	clientCorrelationID string
	session             *wsSession
	metadata            map[string]string
	expiresAt           time.Time // set by hub.track(); zero means no expiry
}

func (s *wsSession) newPendingRoute(req wsInbound) (pendingRoute, error) {
	requestID, err := newOpaqueToken()
	if err != nil {
		return pendingRoute{}, err
	}

	replyToken, err := newOpaqueToken()
	if err != nil {
		return pendingRoute{}, err
	}

	metadata := cloneMetadata(req.Metadata)
	metadata[metadataKeyReplyToken] = replyToken

	return pendingRoute{
		requestID:           requestID,
		replyToken:          replyToken,
		clientCorrelationID: req.CorrelationID,
		session:             s,
		metadata:            metadata,
	}, nil
}

func cloneMetadata(in map[string]string) map[string]string {
	if len(in) == 0 {
		return map[string]string{}
	}

	out := make(map[string]string, len(in))
	maps.Copy(out, in)

	return out
}

func newOpaqueToken() (string, error) {
	var buf [16]byte
	if _, err := rand.Read(buf[:]); err != nil {
		return "", fmt.Errorf("generate opaque token: %w", err)
	}

	return hex.EncodeToString(buf[:]), nil
}
