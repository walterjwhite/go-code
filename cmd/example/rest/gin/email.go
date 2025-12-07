package main

import (
	"fmt"

	"strings"

	gomail "gopkg.in/gomail.v2"
)

func sendContactEmail(cfg *EmailConfig, req ContactRequest) error {
	m := gomail.NewMessage()
	if cfg.From != "" {
		m.SetHeader("From", cfg.From)
	}
	m.SetHeader("To", cfg.To)
	subj := strings.TrimSpace(req.Subject)
	if subj == "" {
		subj = "Contact form message"
	}
	m.SetHeader("Subject", subj)

	plainBody := fmt.Sprintf("Name: %s\nEmail: %s\n\nMessage:\n%s", req.Name, req.Email, req.Message)
	htmlBody := fmt.Sprintf(`<p><strong>Name:</strong> %s<br><strong>Email:</strong> %s</p><hr><p>%s</p>`,
		htmlEscape(req.Name), htmlEscape(req.Email), htmlEscape(req.Message))

	m.SetBody("text/plain", plainBody)
	m.AddAlternative("text/html", htmlBody)

	d := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.Username, cfg.Password)
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
