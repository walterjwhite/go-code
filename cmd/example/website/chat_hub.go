package main

import (
	"context"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
	chatlib "github.com/walterjwhite/go-code/lib/chat"
)

var chatHub = newHub()

type hub struct {
	mu                 sync.Mutex
	requestsByID       map[string]pendingRoute
	requestsByReplyKey map[string]pendingRoute
	allSessions        map[*wsSession]struct{}
}

func newHub() *hub {
	return &hub{
		requestsByID:       make(map[string]pendingRoute),
		requestsByReplyKey: make(map[string]pendingRoute),
		allSessions:        make(map[*wsSession]struct{}),
	}
}

func (h *hub) register(s *wsSession) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.allSessions[s] = struct{}{}
	log.Debug().Str("ip", s.ip).Msg("hub: session registered")
}

func (h *hub) unregister(s *wsSession) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.allSessions, s)

	for id, route := range h.requestsByID {
		if route.session == s {
			delete(h.requestsByID, id)
		}
	}

	for key, route := range h.requestsByReplyKey {
		if route.session == s {
			delete(h.requestsByReplyKey, key)
		}
	}

	log.Debug().Str("ip", s.ip).Msg("hub: session unregistered")
}

const pendingRouteTTL = 5 * time.Minute

func (h *hub) track(route pendingRoute) {
	route.expiresAt = time.Now().Add(pendingRouteTTL)
	h.mu.Lock()
	defer h.mu.Unlock()
	h.requestsByID[route.requestID] = route
	h.requestsByReplyKey[route.replyToken] = route
}

func (h *hub) route(msg chatlib.ServerMessage) {
	route, ok := h.popByRequestID(msg.CorrelationID)
	if !ok {
		route, ok = h.popByReplyToken(msg.Metadata[metadataKeyReplyToken])
	}
	if !ok {
		route, ok = h.popUniqueByIP(msg.Metadata[metadataKeyIP])
	}
	if !ok {
		log.Warn().
			Str("correlation_id", msg.CorrelationID).
			Msg("hub: cannot route response without a recognized request id or reply token")
		return
	}

	msg.CorrelationID = route.clientCorrelationID
	delete(msg.Metadata, metadataKeyReplyToken)
	route.session.sendServerMessage(msg)
}

func (h *hub) start(ctx context.Context, svc chatService) {
	svc.Start(ctx)
	go h.dispatch(ctx, svc)
	go h.pruneLoop(ctx)
}

func (h *hub) pruneLoop(ctx context.Context) {
	ticker := time.NewTicker(pendingRouteTTL / 2)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			h.pruneExpiredRoutes(time.Now())
		}
	}
}

func (h *hub) pruneExpiredRoutes(now time.Time) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for id, route := range h.requestsByID {
		if !route.expiresAt.IsZero() && now.After(route.expiresAt) {
			delete(h.requestsByID, id)
			delete(h.requestsByReplyKey, route.replyToken)
			log.Debug().Str("request_id", id).Msg("hub: pruned expired pending route")
		}
	}
}

func (h *hub) dispatch(ctx context.Context, svc chatService) {
	for {
		select {
		case <-ctx.Done():
			return

		case msg, ok := <-svc.Responses():
			if !ok {
				log.Warn().Msg("hub: responses channel closed")
				return
			}
			h.route(msg)

		case err, ok := <-svc.Errors():
			if !ok {
				return
			}
			log.Error().Err(err).Msg("hub: pubsub error")
		}
	}
}

func (h *hub) popByRequestID(requestID string) (pendingRoute, bool) {
	if requestID == "" {
		return pendingRoute{}, false
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	route, ok := h.requestsByID[requestID]
	if !ok {
		return pendingRoute{}, false
	}

	delete(h.requestsByID, requestID)
	delete(h.requestsByReplyKey, route.replyToken)
	return route, true
}

func (h *hub) popByReplyToken(replyToken string) (pendingRoute, bool) {
	if replyToken == "" {
		return pendingRoute{}, false
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	route, ok := h.requestsByReplyKey[replyToken]
	if !ok {
		return pendingRoute{}, false
	}

	delete(h.requestsByReplyKey, replyToken)
	delete(h.requestsByID, route.requestID)
	return route, true
}

func (h *hub) popUniqueByIP(ip string) (pendingRoute, bool) {
	if ip == "" {
		return pendingRoute{}, false
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	var matched pendingRoute
	found := false

	for _, route := range h.requestsByID {
		if route.session == nil || route.session.ip != ip {
			continue
		}
		if found {
			log.Warn().
				Str("ip", ip).
				Msg("hub: cannot safely route response by ip with multiple pending requests")
			return pendingRoute{}, false
		}
		matched = route
		found = true
	}

	if !found {
		return pendingRoute{}, false
	}

	delete(h.requestsByID, matched.requestID)
	delete(h.requestsByReplyKey, matched.replyToken)
	return matched, true
}
