package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	chatlib "github.com/walterjwhite/go-code/lib/chat"
)


func TestOriginCheckIgnoresXForwardedHost(t *testing.T) {
	req := &http.Request{
		Host: "127.0.0.1:8080",
		Header: http.Header{
			"Origin":           []string{"http://evil.example"},
			"X-Forwarded-Host": []string{"evil.example"},
		},
	}
	if isAllowedWebSocketOrigin(req) {
		t.Fatal("origin check must reject evil.example even when X-Forwarded-Host matches")
	}
}

func TestOriginCheckUsesHostHeader(t *testing.T) {
	req := &http.Request{
		Host:   "myapp.example.com",
		Header: http.Header{"Origin": []string{"https://myapp.example.com"}},
	}
	if !isAllowedWebSocketOrigin(req) {
		t.Fatal("origin check must accept an origin that matches r.Host")
	}
}

func TestOriginCheckIPv6LinkLocal(t *testing.T) {
	req := &http.Request{
		Host:   "[fe80::3c00:69ff:fe78:cba1]:8180",
		Header: http.Header{"Origin": []string{"http://[fe80::3c00:69ff:fe78:cba1]:8180"}},
	}
	if !isAllowedWebSocketOrigin(req) {
		t.Fatal("IPv6 link-local with explicit port must be accepted")
	}
}

func TestOriginCheckIPv6PortMismatch(t *testing.T) {
	req := &http.Request{
		Host:   "[fe80::3c00:69ff:fe78:cba1]", // $host — port stripped
		Header: http.Header{"Origin": []string{"http://[fe80::3c00:69ff:fe78:cba1]:8180"}},
	}
	if isAllowedWebSocketOrigin(req) {
		t.Fatal("host-without-port must not match origin-with-port — confirms $host is wrong variable")
	}
}


func TestSecurityHeadersPresent(t *testing.T) {
	router := gin.New()
	router.Use(securityHeadersMiddleware)
	router.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	router.ServeHTTP(w, req)

	required := map[string]string{
		"X-Content-Type-Options":    "nosniff",
		"X-Frame-Options":           "DENY",
		"Referrer-Policy":           "strict-origin-when-cross-origin",
		"Permissions-Policy":        "camera=(), microphone=(), geolocation=()",
		"Strict-Transport-Security": "max-age=63072000; includeSubDomains; preload",
	}
	for header, want := range required {
		got := w.Header().Get(header)
		if got != want {
			t.Errorf("header %q: got %q want %q", header, got, want)
		}
	}

	if v := w.Header().Get("X-XSS-Protection"); v != "" {
		t.Errorf("deprecated header X-XSS-Protection must be absent, got %q", v)
	}
}


func TestIPRateLimiterMaxEntries(t *testing.T) {
	const maxIPs = 3
	limiter := chatlib.NewIPRateLimiter(chatlib.RateLimitConfig{
		Messages: 10,
		Window:   time.Minute,
		MaxIPs:   maxIPs,
	})

	for i := range maxIPs {
		ip := "10.0.0." + string(rune('1'+i))
		if !limiter.Allow(ip) {
			t.Fatalf("expected Allow for IP %s within cap", ip)
		}
	}

	newIP := "10.0.1.1"
	if limiter.Allow(newIP) {
		t.Fatal("expected Allow=false for IP beyond maxIPs cap")
	}
}

func TestIPRateLimiterExistingIPsStillAllowedAtCap(t *testing.T) {
	limiter := chatlib.NewIPRateLimiter(chatlib.RateLimitConfig{
		Messages: 5,
		Window:   time.Minute,
		MaxIPs:   1,
	})

	if !limiter.Allow("1.1.1.1") {
		t.Fatal("first request from 1.1.1.1 should be allowed")
	}
	if !limiter.Allow("1.1.1.1") {
		t.Fatal("subsequent requests from 1.1.1.1 should be allowed (cap applies to new IPs)")
	}
}


func TestHubPrunesExpiredRoutes(t *testing.T) {
	h := newHub()

	route := pendingRoute{
		requestID:  "req-expired",
		replyToken: "tok-expired",
	}
	h.track(route) // sets expiresAt = now + pendingRouteTTL

	h.mu.Lock()
	tracked := len(h.requestsByID)
	h.mu.Unlock()
	if tracked != 1 {
		t.Fatalf("expected 1 tracked route, got %d", tracked)
	}

	h.pruneExpiredRoutes(time.Now().Add(pendingRouteTTL + time.Second))

	h.mu.Lock()
	remaining := len(h.requestsByID) + len(h.requestsByReplyKey)
	h.mu.Unlock()
	if remaining != 0 {
		t.Fatalf("expected 0 routes after TTL expiry, got %d", remaining)
	}
}

func TestHubDoesNotPruneActiveRoutes(t *testing.T) {
	h := newHub()

	route := pendingRoute{
		requestID:  "req-active",
		replyToken: "tok-active",
	}
	h.track(route)

	h.pruneExpiredRoutes(time.Now().Add(-time.Second))

	h.mu.Lock()
	remaining := len(h.requestsByID)
	h.mu.Unlock()
	if remaining != 1 {
		t.Fatalf("expected 1 active route to survive pruning, got %d", remaining)
	}
}


func TestValidateContactRequest_MaxLengths(t *testing.T) {
	tests := []struct {
		name      string
		req       *ContactRequest
		wantError string
	}{
		{
			name: "name at max length accepted",
			req: &ContactRequest{
				Name:    strings.Repeat("a", maxNameLength),
				Email:   "alice@example.com",
				Subject: "Hi",
				Message: "Hello",
			},
		},
		{
			name: "name exceeds max length",
			req: &ContactRequest{
				Name:    strings.Repeat("a", maxNameLength+1),
				Email:   "alice@example.com",
				Subject: "Hi",
				Message: "Hello",
			},
			wantError: "name exceeds maximum length",
		},
		{
			name: "subject at max length accepted",
			req: &ContactRequest{
				Name:    "Alice",
				Email:   "alice@example.com",
				Subject: strings.Repeat("s", maxSubjectLength),
				Message: "Hello",
			},
		},
		{
			name: "subject exceeds max length",
			req: &ContactRequest{
				Name:    "Alice",
				Email:   "alice@example.com",
				Subject: strings.Repeat("s", maxSubjectLength+1),
				Message: "Hello",
			},
			wantError: "subject exceeds maximum length",
		},
		{
			name: "email at max length accepted",
			req: &ContactRequest{
				Name:    "Alice",
				Email:   strings.Repeat("a", maxEmailLength-len("@x.co")) + "@x.co",
				Subject: "Hi",
				Message: "Hello",
			},
		},
		{
			name: "CRLF in name rejected",
			req: &ContactRequest{
				Name:    "Alice\nEvil",
				Email:   "alice@example.com",
				Subject: "Hi",
				Message: "Hello",
			},
			wantError: "invalid characters in request",
		},
		{
			name: "CRLF in email rejected",
			req: &ContactRequest{
				Name:    "Alice",
				Email:   "alice@example.com\rBCC: victim@evil.com",
				Subject: "Hi",
				Message: "Hello",
			},
			wantError: "invalid email address",
		},
		{
			name: "multibyte runes counted correctly at message limit",
			req: &ContactRequest{
				Name:    "Alice",
				Email:   "alice@example.com",
				Subject: "Hi",
				Message: strings.Repeat("é", maxMessageLength),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateContactRequest(tt.req)
			if tt.wantError == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected error %q, got nil", tt.wantError)
			}
			if !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("error %q does not contain %q", err.Error(), tt.wantError)
			}
		})
	}
}
