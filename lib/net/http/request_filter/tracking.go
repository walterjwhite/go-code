package request_filter

import (
	"time"
)

func (i *Conf) markVisitorSeen(ip string) {
	if ip == "" {
		return
	}

	if i.TTL <= 0 {
		i.TTL = 5 * time.Minute
	}

	i.visitorsMu.Lock()
	expire := time.Now().Add(i.TTL)
	i.visitors[ip] = expire

	if t, ok := i.timers[ip]; ok {
		t.Stop()
	}

	i.timers[ip] = time.AfterFunc(i.TTL, func() {
		i.visitorsMu.Lock()
		defer i.visitorsMu.Unlock()
		if exp, ok := i.visitors[ip]; ok && time.Now().After(exp) {
			delete(i.visitors, ip)
			if tm, ok := i.timers[ip]; ok {
				if tm != nil {
					tm.Stop()
				}
				delete(i.timers, ip)
			}
		}
	})
	i.visitorsMu.Unlock()
}

func (i *Conf) visitorSeen(ip string) bool {
	if ip == "" {
		return false
	}
	i.visitorsMu.RLock()
	t, ok := i.visitors[ip]
	i.visitorsMu.RUnlock()
	if !ok {
		return false
	}
	if time.Now().After(t) {
		i.visitorsMu.Lock()
		delete(i.visitors, ip)
		if tm, ok := i.timers[ip]; ok {
			tm.Stop()
			delete(i.timers, ip)
		}
		i.visitorsMu.Unlock()
		return false
	}
	return true
}
