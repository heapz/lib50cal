package ratelimit

import (
	"lib50cal/server/config"
	"lib50cal/server/internal/logger"
	"lib50cal/server/internal/misc"

	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type hashLimit map[string]*singleRequest
type headerExtractor func(*http.Request) string

type singleRequest struct {
	ticker  *time.Ticker
	limiter *rate.Limiter
	banned  bool
}

type rateLimiter struct {
	requests hashLimit
	headar   string
	mu       *sync.Mutex
	limit    rate.Limit
	extract  headerExtractor
	burst    int
	bantime  time.Duration
}

func NewRateLimiter(conf *config.RateLimitConfig) *rateLimiter {
	rt := &rateLimiter{
		limit: conf.GetLimit(), burst: conf.GetBurst(),
		bantime: conf.GetBanTime(),
	}

	rt.mu = &sync.Mutex{}
	rt.requests = make(hashLimit)

	behindProxyHeader := func(r *http.Request) string {
		var str string
		vsz := strings.Split(r.Header.Get(conf.GetForwardHeader()), ",")
		if len(vsz) <= 0 {
			return str
		}
		str = string(net.ParseIP(strings.Replace(vsz[0], " ", "", -1)))
		return str
	}

	if conf.GetBehindProxy() {
		rt.extract = behindProxyHeader
		return rt
	}

	rt.extract = func(r *http.Request) string { return r.RemoteAddr }
	return rt
}

func (p *rateLimiter) inspectIP(ip string) *singleRequest {
	p.mu.Lock()
	defer p.mu.Unlock()

	v, exists := p.requests[ip]
	if !exists {
		rq := new(singleRequest)
		rq.limiter = rate.NewLimiter(p.limit, p.burst)
		return rq
	}

	return v
}

// vacuum cleanup malicious user address from jail map
func (p *rateLimiter) vacuum(ticker *time.Ticker, ip string) {

	<-ticker.C // malicious user learned his/her lesson, time to let go
	//defer p.mu.Unlock() // funcion overhead for no reason

	ticker.Stop() // dont fuck with memory
	p.mu.Lock()
	delete(p.requests, ip)
	p.mu.Unlock()

}

// http middleware to limit number of request per second, usually we shouldn't have
// too many connection from single address, this is not a regular http server
// connections may stay active as long as vpn tunnel stay connected, bottom line is for
// each user we should have one active http connection, nothing more.
func (p *rateLimiter) HttpRateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ipReal := p.extract(r)
		if len(ipReal) == 0 {
			logger.Printf("error reading ip from %s proxy header", p.headar)
			misc.HttpErr(w, http.StatusInternalServerError)
			return
		}

		rateLimit := p.inspectIP(ipReal)

		// if we already caught this bad boy
		if rateLimit.banned {
			rateLimit.ticker.Reset(p.bantime) // first thing first, extent sentence
			logger.Debugf("too many request from %s reset bantime period", ipReal)
			misc.HttpErr(w, http.StatusTooManyRequests) // bad news for clinet
			return
		}

		// slowdown cowboy, we have limited resources
		if !rateLimit.limiter.Allow() {
			rateLimit.banned = true                      // we have a suspect
			rateLimit.ticker = time.NewTicker(p.bantime) // start counter to free up space and client
			go p.vacuum(rateLimit.ticker, ipReal)        // vacuum clients from jail map
			logger.Debugf("too many request from %s detected", ipReal)
			misc.HttpErr(w, http.StatusTooManyRequests)
			return
		}

		// and the journey begins
		next.ServeHTTP(w, r)
	}
}
