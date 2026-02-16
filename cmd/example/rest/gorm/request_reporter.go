package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

func StartDailyRequestReportWorker(db *sqlx.DB) {
	go func() {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 5, 0, 0, now.Location())
		time.Sleep(time.Until(next))

		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for {
			runDailyReport(db)

			<-ticker.C
		}
	}()
}

func runDailyReport(db *sqlx.DB) {
	rows := []RequestLog{}
	err := db.Select(&rows, `SELECT ts, ip, method, request_uri, user_agent, status FROM http_requests WHERE ts >= now() - interval '24 hours' ORDER BY ts DESC`)
	if err != nil {
		log.Printf("daily report query error: %v", err)
		return
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)
	_ = w.Write([]string{"ts", "ip", "method", "request_uri", "user_agent", "status"})
	for _, r := range rows {
		_ = w.Write([]string{r.TS.Format(time.RFC3339), r.IP, r.Method, r.RequestURI, r.UserAgent, fmt.Sprintf("%d", r.Status)})
	}
	w.Flush()
	if err := w.Error(); err != nil {
		log.Printf("csv write error: %v", err)
		return
	}

	from := os.Getenv("REPORT_EMAIL_FROM")
	toList := os.Getenv("REPORT_EMAIL_TO")
	smtpHost := os.Getenv("REPORT_SMTP_HOST")
	smtpPort := os.Getenv("REPORT_SMTP_PORT")
	smtpUser := os.Getenv("REPORT_SMTP_USER")
	smtpPass := os.Getenv("REPORT_SMTP_PASS")

	if from == "" || toList == "" || smtpHost == "" || smtpPort == "" {
		log.Printf("daily report: missing REPORT_EMAIL_FROM/REPORT_EMAIL_TO/REPORT_SMTP_HOST/REPORT_SMTP_PORT; skipping email")
		return
	}

	tos := strings.Split(toList, ",")
	for i := range tos {
		tos[i] = strings.TrimSpace(tos[i])
	}

	subject := fmt.Sprintf("Daily HTTP requests report - %s", time.Now().Format("2006-01-02"))
	body := buf.String()

	msg := bytes.Buffer{}
	msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(tos, ",")))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
	msg.WriteString("\r\n")
	msg.WriteString(body)

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)

	var auth smtp.Auth
	if smtpUser != "" || smtpPass != "" {
		auth = smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)
	}

	if err := smtp.SendMail(addr, auth, from, tos, msg.Bytes()); err != nil {
		log.Printf("daily report send error: %v", err)
		return
	}

	log.Printf("daily report sent: %d rows", len(rows))
}
