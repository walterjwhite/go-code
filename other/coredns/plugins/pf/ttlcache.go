package pf

import (
	"github.com/ReneKroon/ttlcache/v2"
	"github.com/walterjwhite/go/lib/application/logging"
	"time"
)

var (
	cache = ttlcache.NewCache()
	ttl   time.Duration

// time.Duration(5 * time.Minute)
)

func init() {
	logging.Panic(cache.SetTTL(ttl))
	cache.SetExpirationReasonCallback(_expire)
}

func add(ip string) {
	// if the key already exists, bypass adding it to pf, that will just slow the operation down
	logging.Panic(cache.Set(ip, ""))
	pfAdd(ip)
}

func remove(ip string) {
	logging.Panic(cache.Remove(ip))
	pfRemove(ip)
}

func _expire(key string, reason ttlcache.EvictionReason, value interface{}) {
	remove(key)
}
