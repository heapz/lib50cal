package config

import (
	"errors"
	"time"

	"github.com/sina-ghaderi/rabbitio"
)

// configuration struct for vpn cookie
//

const _exdefualt = 5 * time.Minute

type CookieConfig struct {
	rabbitKey   []byte        // rabbit cipher key value
	expireation time.Duration // cookie expiration time

}

// errors associated with vpn cookies
var (
	ErrCookieBadKey = errors.New("rabbit key must be exactly 16 byte len")
)

// NewCookieConfig returns CookieConfig struct type filled with defualt values
// default value for expiration is 5 min and for rabbit key is 0x00*16
// DO NOT use defualt key value in production, set your own key with SetRabbitKey method
func NewCookieConfig() *CookieConfig {
	return &CookieConfig{
		expireation: _exdefualt,
		rabbitKey:   make([]byte, rabbitio.KeyLen),
	}
}

// SetExpiration set expiration duration for vpn cookies
func (p *CookieConfig) SetExpiration(v time.Duration) { p.expireation = v }

// GetExpiration returns current vpn cookie expiration duration
func (p *CookieConfig) GetExpiration() time.Duration { return p.expireation }

// SetExpiration set rabbit cipher key for vpn cookies
func (p *CookieConfig) SetRabbitKey(v []byte) (err error) {
	if len(v) != rabbitio.KeyLen {
		return ErrCookieBadKey
	}

	copy(p.rabbitKey, v)
	return
}

// GetRabbitKey returns a copy of current rabbit key value
func (p *CookieConfig) GetRabbitKey() []byte {
	v := make([]byte, rabbitio.KeyLen)
	copy(v, p.rabbitKey)
	return v
}
