package server

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

const (
	_limitDefualt = 2
	_burstDefualt = 8
	_realHeaderIP = "X-Forwarded-For"
	_xbehindProxy = false
	_blockDefault = time.Minute
)

type rateLimitConfig struct {
	// TODO: slice of forwardHeader? maybe

	forwardHeader string        // custom forwarded-for (in revers-proxy) header to read real client ip
	bantime       time.Duration // ban time for clients that exceed http rate limit
	reqLimit      rate.Limit    // limit
	reqBurst      int           // burst
	behindProxy   bool          // whether the server situated behind a revers proxy
}

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

func newRateLimiter(conf *rateLimitConfig) *rateLimiter {
	rt := &rateLimiter{
		limit:   conf.reqLimit,
		burst:   conf.reqBurst,
		bantime: conf.bantime,
	}

	rt.mu = &sync.Mutex{}
	rt.requests = make(hashLimit)

	behindProxyHeader := func(r *http.Request) string {
		var str string
		vsz := strings.Split(r.Header.Get(conf.forwardHeader), ",")
		if len(vsz) <= 0 {
			return str
		}
		str = string(net.ParseIP(strings.Replace(vsz[0], " ", "", -1)))
		return str
	}

	if conf.behindProxy {
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
func (p *rateLimiter) httpRateLimit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ipReal := p.extract(r)
		if len(ipReal) == 0 {
			atprint.Load().Printf("error reading ip from %s proxy header", p.headar)
			httpError(w, http.StatusInternalServerError)
			return
		}

		rateLimit := p.inspectIP(ipReal)

		// if we already caught this bad boy
		if rateLimit.banned {
			rateLimit.ticker.Reset(p.bantime) // first thing first, extent sentence
			atdebug.Load().Printf("too many request from %s reset bantime period", ipReal)
			httpError(w, http.StatusTooManyRequests) // bad news for clinet
			return
		}

		// slowdown cowboy, we have limited resources
		if !rateLimit.limiter.Allow() {
			rateLimit.banned = true                      // we have a suspect
			rateLimit.ticker = time.NewTicker(p.bantime) // start counter to free up space and client
			go p.vacuum(rateLimit.ticker, ipReal)        // vacuum clients from jail map
			atdebug.Load().Printf("too many request from %s detected", ipReal)
			httpError(w, http.StatusTooManyRequests)
			return
		}

		// and the journey begins
		next.ServeHTTP(w, r)
	}
}

// NewrateLimitConfig returns default rateLimitConfig which is http request limiter to pervent dos attack
// default values: ratelimit=2, burst=8, behindProxy=false, forwardHeader="X-Forwarded-For"
// and bantime is set to one minute
func NewrateLimitConfig() *rateLimitConfig {
	return &rateLimitConfig{
		forwardHeader: _realHeaderIP,
		bantime:       _blockDefault,
		reqLimit:      _limitDefualt,
		reqBurst:      _burstDefualt,
		behindProxy:   _xbehindProxy,
	}
}

// SetForwardHeader set Forwarded-For header to extract real client ip
// if server is behind a revers proxy
func (p *rateLimitConfig) SetForwardHeader(h string) { p.forwardHeader = h }

// SetBanTime set new bantime duration
func (p *rateLimitConfig) SetBanTime(v time.Duration) { p.bantime = v }

// SetLimit set new rate limit value for http rate limiter
func (p *rateLimitConfig) SetLimit(x rate.Limit) { p.reqLimit = x }

// setBurst set new burst value for http rate limiter
func (p *rateLimitConfig) SetBurst(x int) { p.reqBurst = x }

// SetBehindProxy determine whether server is situated behind a revers proxy
func (p *rateLimitConfig) SetBehindProxy(b bool) { p.behindProxy = b }

// return current forwardHeader value
func (p *rateLimitConfig) GetForwardHeader() (h string) { return p.forwardHeader }

// return current bantime value
func (p *rateLimitConfig) GetBanTime() (v time.Duration) { return p.bantime }

// return current rate limit value
func (p *rateLimitConfig) GetLimit() (x rate.Limit) { return p.reqLimit }

// return current burst value
func (p *rateLimitConfig) GetBurst() (x int) { return p.reqBurst }

// return current behindProxy value
func (p *rateLimitConfig) GetBehindProxy() (b bool) { return p.behindProxy }
