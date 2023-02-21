package server

import (
	"errors"
	"sync"
	"time"

	"github.com/sina-ghaderi/rabbitio"
)

// configuration struct for vpn cookie
//

const _expireDefualt = 5 * time.Minute

type cookieConfig struct {
	rabbitKey   []byte        // rabbit cipher key value
	expireation time.Duration // cookie expiration time

}

// errors associated with vpn cookies
var (
	ErrCookieBadKey = errors.New("cookieConfig: rabbit key must be exactly 16 byte len")
)

// NewCookieConfig returns CookieConfig struct type filled with defualt values
// default value for expiration is 5 min and for rabbit key is 0x00*16
// DO NOT use defualt key value in production, set your own key with SetRabbitKey method
func NewCookieConfig() *cookieConfig {
	return &cookieConfig{
		expireation: _expireDefualt,
		rabbitKey:   make([]byte, rabbitio.KeyLen),
	}
}

// SetExpiration set expiration duration for vpn cookies
func (p *cookieConfig) SetExpiration(v time.Duration) { p.expireation = v }

// GetExpiration returns current vpn cookie expiration duration
func (p *cookieConfig) GetExpiration() time.Duration { return p.expireation }

// SetExpiration set rabbit cipher key for vpn cookies
func (p *cookieConfig) SetRabbitKey(v []byte) (err error) {
	if len(v) != rabbitio.KeyLen {
		return ErrCookieBadKey
	}

	copy(p.rabbitKey, v)
	return
}

// GetRabbitKey returns a copy of current rabbit key value
func (p *cookieConfig) GetRabbitKey() []byte {
	v := make([]byte, rabbitio.KeyLen)
	copy(v, p.rabbitKey)
	return v
}

type cookieJar struct {
	// user last update map -- int: userid, time.Time: last update/delete...
	cookm map[int]time.Time
	mutex *sync.Mutex
}
