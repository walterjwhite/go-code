package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net"
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

	body, err := io.ReadAll(resp.Body)
	if err != nil && err != io.EOF {
		return "", err
	}

	responseText := string(body)

	var jsonResp map[string]any
	if err := json.Unmarshal([]byte(responseText), &jsonResp); err == nil {
		if origin, ok := jsonResp["origin"]; ok {
			if ipStr, ok := origin.(string); ok {
				if net.ParseIP(ipStr) != nil {
					return ipStr, nil
				}
			}
		}
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
	parts := strings.FieldsSeq(text)
	for part := range parts {
		cleaned := ""
		for _, ch := range part {
			if (ch >= '0' && ch <= '9') || ch == '.' {
				cleaned += string(ch)
			}
		}
		if net.ParseIP(cleaned) != nil {
			return cleaned
		}
	}
	return ""
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
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		},
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(timeoutSeconds) * time.Second,
	}

	return client, nil
}

func closeBody(body io.ReadCloser) {
	logging.Warn(body.Close(), false, "Failed to close body")
}
