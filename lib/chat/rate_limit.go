package chat

import (
	"context"
	"time"
)

func NewIPRateLimiter(config RateLimitConfig) *IPRateLimiter {
	maxIPs := config.MaxIPs
	if maxIPs <= 0 {
		maxIPs = DefaultMaxTrackedIPs
	}
	return &IPRateLimiter{
		config:  config,
		maxIPs:  maxIPs,
		windows: make(map[string]ipWindow),
	}
}

func (l *IPRateLimiter) Start(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				l.prune()
			}
		}
	}()
}

func (l *IPRateLimiter) prune() {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	for ip, window := range l.windows {
		if l.windowExpired(now, window) {
			delete(l.windows, ip)
		}
	}
}

func (l *IPRateLimiter) Allow(ip string) bool {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	window, exists := l.windows[ip]
	if l.windowExpired(now, window) {
		if !exists && len(l.windows) >= l.maxIPs {
			return false
		}
		l.windows[ip] = ipWindow{start: now, count: 1}
		return true
	}
	if window.count >= l.config.Messages {
		return false
	}
	window.count++
	l.windows[ip] = window
	return true
}

func (l *IPRateLimiter) windowExpired(now time.Time, window ipWindow) bool {
	return window.start.IsZero() || now.Sub(window.start) >= l.config.Window
}
