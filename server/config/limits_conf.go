package config

import (
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

type RateLimitConfig struct {
	// TODO: slice of forwardHeader? maybe

	forwardHeader string        // custom forwarded-for (in revers-proxy) header to read real client ip
	bantime       time.Duration // ban time for clients that exceed http rate limit
	reqLimit      rate.Limit    // limit
	reqBurst      int           // burst
	behindProxy   bool          // whether the server situated behind a revers proxy
}

// NewRateLimitConfig returns default RateLimitConfig which is http request limiter to pervent dos attack
// default values: ratelimit=2, burst=8, behindProxy=false, forwardHeader="X-Forwarded-For"
// and bantime is set to one minute
func NewRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		forwardHeader: _realHeaderIP,
		bantime:       _blockDefault,
		reqLimit:      _limitDefualt,
		reqBurst:      _burstDefualt,
		behindProxy:   _xbehindProxy,
	}
}

// SetForwardHeader set Forwarded-For header to extract real client ip
// if server is behind a revers proxy
func (p *RateLimitConfig) SetForwardHeader(h string) { p.forwardHeader = h }

// SetBanTime set new bantime duration
func (p *RateLimitConfig) SetBanTime(v time.Duration) { p.bantime = v }

// SetLimit set new rate limit value for http rate limiter
func (p *RateLimitConfig) SetLimit(x rate.Limit) { p.reqLimit = x }

// setBurst set new burst value for http rate limiter
func (p *RateLimitConfig) SetBurst(x int) { p.reqBurst = x }

// SetBehindProxy determine whether server is situated behind a revers proxy
func (p *RateLimitConfig) SetBehindProxy(b bool) { p.behindProxy = b }

// return current forwardHeader value
func (p *RateLimitConfig) GetForwardHeader() (h string) { return p.forwardHeader }

// return current bantime value
func (p *RateLimitConfig) GetBanTime() (v time.Duration) { return p.bantime }

// return current rate limit value
func (p *RateLimitConfig) GetLimit() (x rate.Limit) { return p.reqLimit }

// return current burst value
func (p *RateLimitConfig) GetBurst() (x int) { return p.reqBurst }

// return current behindProxy value
func (p *RateLimitConfig) GetBehindProxy() (b bool) { return p.behindProxy }
