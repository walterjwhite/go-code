package main

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	chatlib "github.com/walterjwhite/go-code/lib/chat"
)

func init() {
	gin.SetMode(gin.TestMode)
}


type stubChatService struct {
	sendFn    func(ctx context.Context, msg chatlib.ClientMessage) error
	responses chan chatlib.ServerMessage
	errs      chan error
}

func (s *stubChatService) Send(ctx context.Context, msg chatlib.ClientMessage) error {
	if s.sendFn != nil {
		return s.sendFn(ctx, msg)
	}
	return nil
}
func (s *stubChatService) Responses() <-chan chatlib.ServerMessage { return s.responses }
func (s *stubChatService) Errors() <-chan error                    { return s.errs }
func (s *stubChatService) Start(_ context.Context)                 {}

func newStub(fn func(context.Context, chatlib.ClientMessage) error) *stubChatService {
	return &stubChatService{
		sendFn:    fn,
		responses: make(chan chatlib.ServerMessage, 4),
		errs:      make(chan error, 4),
	}
}


func wsTestServer(t *testing.T, svc chatService) (*websocket.Conn, func()) {
	t.Helper()

	prev := chatClient
	chatClient = svc
	t.Cleanup(func() { chatClient = prev })

	router := gin.New()
	router.GET("/ws/chat", onChatWSRequest)
	srv := httptest.NewServer(router)

	ctx, cancel := context.WithCancel(context.Background())
	chatHub.start(ctx, svc)
	t.Cleanup(cancel)

	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws/chat"
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		srv.Close()
		t.Fatalf("dial websocket: %v", err)
	}

	return conn, func() {
		_ = conn.Close()
		srv.Close()
	}
}

func recv(t *testing.T, conn *websocket.Conn) wsOutbound {
	t.Helper()
	_ = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	_, raw, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read message: %v", err)
	}
	var out wsOutbound
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatalf("unmarshal %q: %v", raw, err)
	}
	return out
}

func send(t *testing.T, conn *websocket.Conn, v any) {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := conn.WriteMessage(websocket.TextMessage, b); err != nil {
		t.Fatalf("write message: %v", err)
	}
}

func recvWithin(conn *websocket.Conn, timeout time.Duration) ([]byte, error) {
	_ = conn.SetReadDeadline(time.Now().Add(timeout))
	_, raw, err := conn.ReadMessage()
	return raw, err
}


func TestWSServiceUnavailable(t *testing.T) {
	prev := chatClient
	chatClient = nil
	defer func() { chatClient = prev }()

	router := gin.New()
	router.GET("/ws/chat", onChatWSRequest)
	srv := httptest.NewServer(router)
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/ws/chat")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("want 503, got %d", resp.StatusCode)
	}
}

func TestWSAck(t *testing.T) {
	conn, cleanup := wsTestServer(t, newStub(nil))
	defer cleanup()

	send(t, conn, wsInbound{Message: "hello", CorrelationID: "abc"})

	out := recv(t, conn)
	if out.Type != "ack" {
		t.Fatalf("type=%q want ack", out.Type)
	}
	if out.CorrelationID != "abc" {
		t.Fatalf("correlation_id=%q want abc", out.CorrelationID)
	}
}

func TestWSRejectsCrossOrigin(t *testing.T) {
	prev := chatClient
	chatClient = newStub(nil)
	defer func() { chatClient = prev }()

	router := gin.New()
	router.GET("/ws/chat", onChatWSRequest)
	srv := httptest.NewServer(router)
	defer srv.Close()

	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws/chat"
	header := http.Header{"Origin": []string{"https://evil.example"}}
	_, resp, err := websocket.DefaultDialer.Dial(url, header)
	if err == nil {
		t.Fatal("expected cross-origin websocket handshake to fail")
	}
	if resp == nil || resp.StatusCode != http.StatusForbidden {
		t.Fatalf("status=%v want 403", resp)
	}
}

func TestWSAcceptsMatchingHostOrigin(t *testing.T) {
	prev := chatClient
	chatClient = newStub(nil)
	defer func() { chatClient = prev }()

	router := gin.New()
	router.GET("/ws/chat", func(c *gin.Context) {
		c.Request.Host = "myapp.example.com"
		onChatWSRequest(c)
	})
	srv := httptest.NewServer(router)
	defer srv.Close()

	url := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws/chat"
	header := http.Header{"Origin": []string{"http://myapp.example.com"}}
	conn, resp, err := websocket.DefaultDialer.Dial(url, header)
	if err != nil {
		t.Fatalf("dial websocket: %v (resp=%v)", err, resp)
	}
	_ = conn.Close()
}

func TestIsAllowedWebSocketOrigin(t *testing.T) {
	tests := []struct {
		name string
		req  *http.Request
		want bool
	}{
		{
			name: "same host accepted",
			req: &http.Request{
				Host:   "localhost:18180",
				Header: http.Header{"Origin": []string{"http://localhost:18180"}},
			},
			want: true,
		},
		{
			name: "IPv6 link-local host accepted",
			req: &http.Request{
				Host:   "[fe80::3c00:69ff:fe78:cba1]:8180",
				Header: http.Header{"Origin": []string{"http://[fe80::3c00:69ff:fe78:cba1]:8180"}},
			},
			want: true,
		},
		{
			name: "X-Forwarded-Host injection rejected (header absent in Go)",
			req: &http.Request{
				Host: "127.0.0.1:8080",
				Header: http.Header{
					"Origin": []string{"http://evil.example"},
				},
			},
			want: false,
		},
		{
			name: "cross origin rejected",
			req: &http.Request{
				Host:   "localhost:18180",
				Header: http.Header{"Origin": []string{"https://evil.example"}},
			},
			want: false,
		},
		{
			name: "invalid origin rejected",
			req: &http.Request{
				Host:   "localhost:18180",
				Header: http.Header{"Origin": []string{"://bad"}},
			},
			want: false,
		},
		{
			name: "missing origin allowed",
			req: &http.Request{
				Host:   "localhost:18180",
				Header: http.Header{},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isAllowedWebSocketOrigin(tt.req); got != tt.want {
				t.Fatalf("isAllowedWebSocketOrigin()=%v want %v", got, tt.want)
			}
		})
	}
}

func TestWSMessageIncludesClientIP(t *testing.T) {
	var sent chatlib.ClientMessage
	stub := newStub(func(_ context.Context, msg chatlib.ClientMessage) error {
		sent = msg
		return nil
	})
	conn, cleanup := wsTestServer(t, stub)
	defer cleanup()

	send(t, conn, wsInbound{Message: "hello", CorrelationID: "abc"})

	out := recv(t, conn)
	if out.Type != "ack" {
		t.Fatalf("type=%q want ack", out.Type)
	}
	if sent.IP == "" {
		t.Fatal("expected websocket handler to pass client IP to chat service")
	}
	if sent.CorrelationID == "abc" {
		t.Fatal("expected websocket handler to replace client correlation_id before publishing")
	}
	if sent.Metadata[metadataKeyReplyToken] == "" {
		t.Fatal("expected websocket handler to stamp reply token into metadata")
	}
}

func TestWSErrorOnInvalidJSON(t *testing.T) {
	conn, cleanup := wsTestServer(t, newStub(nil))
	defer cleanup()

	_ = conn.WriteMessage(websocket.TextMessage, []byte("{bad json"))
	out := recv(t, conn)
	if out.Type != "error" {
		t.Fatalf("type=%q want error", out.Type)
	}
	if !strings.Contains(out.Message, "invalid JSON") {
		t.Fatalf("message=%q missing 'invalid JSON'", out.Message)
	}

	send(t, conn, wsInbound{Message: "still alive"})
	out = recv(t, conn)
	if out.Type != "ack" {
		t.Fatalf("connection closed after bad frame; got type=%q", out.Type)
	}
}

func TestWSErrorOnBlankMessage(t *testing.T) {
	conn, cleanup := wsTestServer(t, newStub(nil))
	defer cleanup()

	send(t, conn, wsInbound{Message: "   "})
	out := recv(t, conn)
	if out.Type != "error" {
		t.Fatalf("type=%q want error", out.Type)
	}
	if !strings.Contains(out.Message, "message is required") {
		t.Fatalf("message=%q missing 'message is required'", out.Message)
	}
}

func TestWSRateLimited(t *testing.T) {
	stub := newStub(func(_ context.Context, _ chatlib.ClientMessage) error {
		return chatlib.ErrRateLimited
	})
	conn, cleanup := wsTestServer(t, stub)
	defer cleanup()

	send(t, conn, wsInbound{Message: "hi", CorrelationID: "r1"})
	out := recv(t, conn)
	if out.Type != "error" {
		t.Fatalf("type=%q want error", out.Type)
	}
	if out.CorrelationID != "r1" {
		t.Fatalf("correlation_id=%q want r1", out.CorrelationID)
	}
	if !strings.Contains(out.Message, "rate limit") {
		t.Fatalf("message=%q missing 'rate limit'", out.Message)
	}
}

func TestWSUpstreamError(t *testing.T) {
	stub := newStub(func(_ context.Context, _ chatlib.ClientMessage) error {
		return errors.New("pubsub down")
	})
	conn, cleanup := wsTestServer(t, stub)
	defer cleanup()

	send(t, conn, wsInbound{Message: "hi"})
	out := recv(t, conn)
	if out.Type != "error" {
		t.Fatalf("type=%q want error", out.Type)
	}
	if strings.Contains(out.Message, "pubsub") {
		t.Fatalf("internal error detail leaked: %q", out.Message)
	}
}

func TestWSMultipleMessages(t *testing.T) {
	conn, cleanup := wsTestServer(t, newStub(nil))
	defer cleanup()

	for i, id := range []string{"m1", "m2", "m3"} {
		send(t, conn, wsInbound{Message: "msg", CorrelationID: id})
		out := recv(t, conn)
		if out.Type != "ack" {
			t.Fatalf("msg %d: type=%q want ack", i, out.Type)
		}
		if out.CorrelationID != id {
			t.Fatalf("msg %d: correlation_id=%q want %q", i, out.CorrelationID, id)
		}
	}
}


func TestWSMessageTooLong(t *testing.T) {
	conn, cleanup := wsTestServer(t, newStub(nil))
	defer cleanup()

	send(t, conn, wsInbound{Message: strings.Repeat("a", maxChatMessageLength+1)})
	out := recv(t, conn)
	if out.Type != "error" || !strings.Contains(out.Message, "message exceeds maximum length") {
		t.Fatalf("unexpected: %+v", out)
	}
}

func TestWSMessageAtMaxLength(t *testing.T) {
	conn, cleanup := wsTestServer(t, newStub(nil))
	defer cleanup()

	send(t, conn, wsInbound{Message: strings.Repeat("a", maxChatMessageLength)})
	out := recv(t, conn)
	if out.Type != "ack" {
		t.Fatalf("type=%q want ack", out.Type)
	}
}

func TestWSMultibyteMessageAtLimit(t *testing.T) {
	conn, cleanup := wsTestServer(t, newStub(nil))
	defer cleanup()

	send(t, conn, wsInbound{Message: strings.Repeat("é", maxChatMessageLength)})
	out := recv(t, conn)
	if out.Type != "ack" {
		t.Fatalf("type=%q want ack (multibyte at limit)", out.Type)
	}
}

func TestWSChannelTooLong(t *testing.T) {
	conn, cleanup := wsTestServer(t, newStub(nil))
	defer cleanup()

	send(t, conn, wsInbound{
		Message: "hi",
		Channel: strings.Repeat("c", maxChatChannelLength+1),
	})
	out := recv(t, conn)
	if out.Type != "error" || !strings.Contains(out.Message, "channel exceeds") {
		t.Fatalf("unexpected: %+v", out)
	}
}

func TestWSCorrelationIDTooLong(t *testing.T) {
	conn, cleanup := wsTestServer(t, newStub(nil))
	defer cleanup()

	send(t, conn, wsInbound{
		Message:       "hi",
		CorrelationID: strings.Repeat("x", maxChatCorrelationIDLen+1),
	})
	out := recv(t, conn)
	if out.Type != "error" || !strings.Contains(out.Message, "correlation_id exceeds") {
		t.Fatalf("unexpected: %+v", out)
	}
}

func TestWSMetadataTooManyKeys(t *testing.T) {
	conn, cleanup := wsTestServer(t, newStub(nil))
	defer cleanup()

	meta := make(map[string]string, maxChatMetadataKeys+1)
	for i := 0; i <= maxChatMetadataKeys; i++ {
		meta[string(rune('a'+i))] = "v"
	}
	send(t, conn, wsInbound{Message: "hi", Metadata: meta})
	out := recv(t, conn)
	if out.Type != "error" || !strings.Contains(out.Message, "metadata exceeds maximum") {
		t.Fatalf("unexpected: %+v", out)
	}
}

func TestWSMetadataKeyTooLong(t *testing.T) {
	conn, cleanup := wsTestServer(t, newStub(nil))
	defer cleanup()

	send(t, conn, wsInbound{
		Message:  "hi",
		Metadata: map[string]string{strings.Repeat("k", maxChatMetadataKeyLen+1): "v"},
	})
	out := recv(t, conn)
	if out.Type != "error" || !strings.Contains(out.Message, "metadata key exceeds") {
		t.Fatalf("unexpected: %+v", out)
	}
}

func TestWSMetadataValueTooLong(t *testing.T) {
	conn, cleanup := wsTestServer(t, newStub(nil))
	defer cleanup()

	send(t, conn, wsInbound{
		Message:  "hi",
		Metadata: map[string]string{"key": strings.Repeat("v", maxChatMetadataValueLen+1)},
	})
	out := recv(t, conn)
	if out.Type != "error" || !strings.Contains(out.Message, "metadata value") {
		t.Fatalf("unexpected: %+v", out)
	}
}

func TestWSMetadataRejectsReservedKeys(t *testing.T) {
	conn, cleanup := wsTestServer(t, newStub(nil))
	defer cleanup()

	send(t, conn, wsInbound{
		Message:  "hi",
		Metadata: map[string]string{metadataKeyReplyToken: "forged"},
	})
	out := recv(t, conn)
	if out.Type != "error" || !strings.Contains(out.Message, "reserved") {
		t.Fatalf("unexpected: %+v", out)
	}
}


func TestHubRoute(t *testing.T) {
	var published chatlib.ClientMessage
	stub := newStub(func(_ context.Context, msg chatlib.ClientMessage) error {
		published = msg
		return nil
	})
	conn, cleanup := wsTestServer(t, stub)
	defer cleanup()

	send(t, conn, wsInbound{Message: "ping", CorrelationID: "cid-1"})
	ack := recv(t, conn)
	if ack.Type != "ack" {
		t.Fatalf("expected ack, got %+v", ack)
	}

	sentAt := time.Now().UTC().Truncate(time.Second)
	stub.responses <- chatlib.ServerMessage{
		CorrelationID: published.CorrelationID,
		Message:       "pong",
		SentAt:        sentAt,
	}

	out := recv(t, conn)
	if out.Type != "response" {
		t.Fatalf("type=%q want response", out.Type)
	}
	if out.Message != "pong" {
		t.Fatalf("message=%q want pong", out.Message)
	}
	if out.CorrelationID != "cid-1" {
		t.Fatalf("correlation_id=%q want cid-1", out.CorrelationID)
	}
	if out.SentAt == nil || !out.SentAt.Equal(sentAt) {
		t.Fatalf("sent_at=%v want %v", out.SentAt, sentAt)
	}
}

func TestHubRouteByReplyToken(t *testing.T) {
	var published chatlib.ClientMessage
	stub := newStub(func(_ context.Context, msg chatlib.ClientMessage) error {
		published = msg
		return nil
	})
	conn, cleanup := wsTestServer(t, stub)
	defer cleanup()

	send(t, conn, wsInbound{Message: "ping", CorrelationID: "client-1"})
	ack := recv(t, conn)
	if ack.Type != "ack" {
		t.Fatalf("expected ack, got %+v", ack)
	}

	stub.responses <- chatlib.ServerMessage{
		Message: "pong",
		Metadata: map[string]string{
			metadataKeyReplyToken: published.Metadata[metadataKeyReplyToken],
		},
		SentAt: time.Now().UTC().Truncate(time.Second),
	}

	out := recv(t, conn)
	if out.Type != "response" {
		t.Fatalf("type=%q want response", out.Type)
	}
	if out.Message != "pong" {
		t.Fatalf("message=%q want pong", out.Message)
	}
	if out.CorrelationID != "client-1" {
		t.Fatalf("correlation_id=%q want client-1", out.CorrelationID)
	}
}

func TestHubRouteByIPWhenSinglePendingRequest(t *testing.T) {
	conn, cleanup := wsTestServer(t, newStub(nil))
	defer cleanup()

	send(t, conn, wsInbound{Message: "ping", CorrelationID: "client-1"})
	ack := recv(t, conn)
	if ack.Type != "ack" {
		t.Fatalf("expected ack, got %+v", ack)
	}

	stub := chatClient.(*stubChatService)
	stub.responses <- chatlib.ServerMessage{
		Message:  "pong",
		Metadata: map[string]string{metadataKeyIP: "127.0.0.1"},
		SentAt:   time.Now().UTC().Truncate(time.Second),
	}

	out := recv(t, conn)
	if out.Type != "response" {
		t.Fatalf("type=%q want response", out.Type)
	}
	if out.CorrelationID != "client-1" {
		t.Fatalf("correlation_id=%q want client-1", out.CorrelationID)
	}
}

func TestHubDoesNotRouteByIPWhenMultiplePendingRequestsExist(t *testing.T) {
	conn, cleanup := wsTestServer(t, newStub(nil))
	defer cleanup()

	send(t, conn, wsInbound{Message: "first", CorrelationID: "client-1"})
	if ack := recv(t, conn); ack.Type != "ack" {
		t.Fatalf("expected first ack, got %+v", ack)
	}

	send(t, conn, wsInbound{Message: "second", CorrelationID: "client-2"})
	if ack := recv(t, conn); ack.Type != "ack" {
		t.Fatalf("expected second ack, got %+v", ack)
	}

	stub := chatClient.(*stubChatService)
	stub.responses <- chatlib.ServerMessage{
		Message:  "pong",
		Metadata: map[string]string{metadataKeyIP: "127.0.0.1"},
		SentAt:   time.Now().UTC().Truncate(time.Second),
	}

	if raw, err := recvWithin(conn, 200*time.Millisecond); err == nil {
		t.Fatalf("unexpected routed response: %s", raw)
	} else if netErr, ok := err.(net.Error); !ok || !netErr.Timeout() {
		t.Fatalf("unexpected read error: %v", err)
	}
}

func TestHubUnregistersOnDisconnect(t *testing.T) {
	stub := newStub(nil)
	conn, cleanup := wsTestServer(t, stub)

	send(t, conn, wsInbound{Message: "hi", CorrelationID: "track-me"})
	recv(t, conn) // consume ack

	chatHub.mu.Lock()
	tracked := len(chatHub.requestsByID) > 0
	chatHub.mu.Unlock()
	if !tracked {
		t.Fatal("expected hub to track pending request after Send")
	}

	cleanup() // close conn and server; triggers readPump teardown

	time.Sleep(100 * time.Millisecond)

	chatHub.mu.Lock()
	still := len(chatHub.requestsByID) > 0 || len(chatHub.requestsByReplyKey) > 0
	chatHub.mu.Unlock()
	if still {
		t.Fatal("expected hub to remove pending routes after disconnect")
	}
}

func TestMarshalClientMessageIncludesIPAddressFallback(t *testing.T) {
	payload, err := marshalClientMessage(chatlib.ClientMessage{
		IP:      "203.0.113.10",
		Message: "hello",
	})
	if err != nil {
		t.Fatalf("marshalClientMessage: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(payload, &decoded); err != nil {
		t.Fatalf("unmarshal payload: %v", err)
	}

	if decoded["ip"] != "203.0.113.10" {
		t.Fatalf("ip=%v want 203.0.113.10", decoded["ip"])
	}
	if decoded["ip_address"] != "203.0.113.10" {
		t.Fatalf("ip_address=%v want 203.0.113.10", decoded["ip_address"])
	}

	meta, ok := decoded["metadata"].(map[string]any)
	if !ok || meta["ip"] != "203.0.113.10" {
		t.Fatalf("metadata.ip=%v want 203.0.113.10", meta["ip"])
	}
}

func TestUnmarshalServerMessageNormalizesIPAddress(t *testing.T) {
	msg, err := unmarshalServerMessage([]byte(`{"ip_address":"203.0.113.11","message":"reply"}`))
	if err != nil {
		t.Fatalf("unmarshalServerMessage: %v", err)
	}
	if msg.Metadata["ip"] != "203.0.113.11" {
		t.Fatalf("metadata.ip=%q want 203.0.113.11", msg.Metadata["ip"])
	}
}

func TestPubsubChatServiceRateLimitsByIP(t *testing.T) {
	published := 0
	limiter := chatlib.NewIPRateLimiter(chatlib.RateLimitConfig{
		Messages: 1,
		Window:   time.Minute,
	})
	limiter.Start(context.Background(), time.Minute)

	svc := &pubsubChatService{
		publishTopic: "chat-requests",
		limiter:      limiter,
		publishFn: func(_ string, _ []byte) error {
			published++
			return nil
		},
		responses: make(chan chatlib.ServerMessage, 1),
		errors:    make(chan error, 1),
	}

	err := svc.Send(context.Background(), chatlib.ClientMessage{
		IP:      "203.0.113.20",
		Message: "first",
	})
	if err != nil {
		t.Fatalf("first send: %v", err)
	}

	err = svc.Send(context.Background(), chatlib.ClientMessage{
		IP:      "203.0.113.20",
		Message: "second",
	})
	if !errors.Is(err, chatlib.ErrRateLimited) {
		t.Fatalf("second send err=%v want ErrRateLimited", err)
	}
	if published != 1 {
		t.Fatalf("published=%d want 1", published)
	}
}
