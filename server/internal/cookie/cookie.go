package cookie

import (
	"lib50cal/server/config"
	"sync"
	"time"
)

type cookieJar struct {
	// user last update map -- int: userid, time.Time: last update/delete...
	cookm map[int]time.Time
	mutex *sync.Mutex

	expire time.Duration
	enckey []byte
}

func NewCookieJar(conf *config.CookieConfig) *cookieJar {
	return &cookieJar{
		expire: conf.GetExpiration(), enckey: conf.GetRabbitKey(),
	}
}

// methods goes here
