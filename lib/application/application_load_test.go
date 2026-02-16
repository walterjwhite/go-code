package application

import (
	"context"
	"testing"
)

type testCfg struct {
	Calls []string
}

func (t *testCfg) PreLoad() {
	t.Calls = append(t.Calls, "pre")
}

func (t *testCfg) PostLoad(ctx context.Context) error {
	t.Calls = append(t.Calls, "post")
	return nil
}

func TestLoadCallsPreAndPost(t *testing.T) {
	cfg := &testCfg{}

	Load(cfg)

	if len(cfg.Calls) < 2 {
		t.Fatalf("expected PreLoad and PostLoad to be called, got calls=%v", cfg.Calls)
	}

	if cfg.Calls[0] != "pre" {
		t.Fatalf("expected first call to be pre, got %v", cfg.Calls[0])
	}
	if cfg.Calls[1] != "post" {
		t.Fatalf("expected second call to be post, got %v", cfg.Calls[1])
	}
}
