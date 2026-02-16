package main

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/walterjwhite/go-code/lib/application/logging"
	"golang.org/x/net/proxy"
)

type Session struct {
	ID        int       `json:"id"`
	ProxyAddr string    `json:"proxy_addr"`
	PublicIP  string    `json:"public_ip"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type DeviceInfo struct {
	PublicIP string
}

func createSessionTable(db *sql.DB) error {
	_, err := db.Exec(SessionTableSchema)
	return err
}

func getPublicIP(client *http.Client, serviceURL string) (string, error) {
	resp, err := client.Get(serviceURL)
	if err != nil {
		return "", err
	}
	defer closeBody(resp.Body)

	if resp.StatusCode != HTTPStatusOK {
		return "", fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	body := make([]byte, ResponseBufferSize) // Read only first ResponseBufferSize bytes
	n, err := resp.Body.Read(body)
	if err != nil && err.Error() != "EOF" {
		return "", err
	}

	responseText := string(body[:n])

	if strings.Contains(responseText, OriginField) {
		startIdx := strings.Index(responseText, `"origin"`) + len(`"origin"`)
		startIdx = strings.Index(responseText[startIdx:], ":") + startIdx + 1
		endIdx := strings.Index(responseText[startIdx:], ",")
		if endIdx == -1 {
			endIdx = strings.Index(responseText[startIdx:], "}")
		}

		ipPart := responseText[startIdx : endIdx+startIdx]
		ipPart = strings.TrimSpace(ipPart)
		ipPart = strings.Trim(ipPart, `" `)

		return ipPart, nil
	}

	return extractIP(responseText), nil
}

func getDevicePublicIP(serviceURL string, timeoutSeconds int) (string, error) {
	client := &http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}

	return getPublicIP(client, serviceURL)
}

func extractIP(text string) string {
	parts := strings.Split(text, ".")
	if len(parts) != 4 {
		return ""
	}

	for _, part := range parts {
		num := 0
		for _, c := range part {
			if c >= '0' && c <= '9' {
				num = num*10 + int(c-'0')
			} else {
				newPart := ""
				for _, ch := range part {
					if ch >= '0' && ch <= '9' {
						newPart += string(ch)
					}
				}
				num = 0
				for _, ch := range newPart {
					num = num*10 + int(ch-'0')
				}
				break
			}
		}
		if num > 255 {
			return ""
		}
	}

	for i, part := range parts {
		num := 0
		for _, c := range part {
			if c >= '0' && c <= '9' {
				num = num*10 + int(c-'0')
			} else {
				newPart := ""
				for _, ch := range part {
					if ch >= '0' && ch <= '9' {
						newPart += string(ch)
					}
				}
				parts[i] = newPart
				num = 0
				for _, digit := range newPart {
					num = num*10 + int(digit-'0')
				}
				break
			}
		}
		if num > 255 {
			return ""
		}
	}

	return strings.Join(parts, ".")
}

func saveSession(db *sql.DB, session Session) error {
	_, err := db.Exec(InsertSessionQuery, session.ProxyAddr, session.PublicIP, session.StartTime, session.EndTime)
	return err
}

func setupProxyClient(socksProxy string, timeoutSeconds int) (*http.Client, error) {
	dialer, err := proxy.SOCKS5(SOCKSProtocol, socksProxy, nil, proxy.Direct)
	if err != nil {
		return nil, fmt.Errorf("failed to create SOCKS5 dialer: %w", err)
	}

	transport := &http.Transport{
		Dial: dialer.Dial,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeoutSeconds) * time.Second,
	}

	return client, nil
}

func closeBody(body io.ReadCloser) {
	logging.Warn(body.Close(), "Failed to close body")
}
