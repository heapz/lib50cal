package server

import "lib50cal/server/config"

type serverConfig struct {
	cookieConfig    *config.CookieConfig
	ratelimitConfig *config.RateLimitConfig
	loggerConfig    *config.LoggerConfig

	// http server config
	// tls server config
	// web socket config
	// tun/tap networking config
	// some configs that i can't think of right now
}

func NewServerConfig() *serverConfig {

	// default configuration for all modules
	return &serverConfig{
		cookieConfig:    config.NewCookieConfig(),
		ratelimitConfig: config.NewRateLimitConfig(),
		loggerConfig:    config.NewLoggerConfig(),
	}
}

func (p *serverConfig) SetCookieConfig(val *config.CookieConfig) {
	// config validation -- > somthing like val.Validate() method
	// to check whether config is valid or not
	p.cookieConfig = val
}

func (p *serverConfig) GetCookieConfig() *config.CookieConfig {
	return p.cookieConfig
}

func (p *serverConfig) SetRateLimitConfig(val *config.RateLimitConfig) {
	// config validation -- > somthing like val.Validate() method
	// to check whether config is valid or not

	p.ratelimitConfig = val
}

func (p *serverConfig) GetRateLimitConfig() *config.RateLimitConfig {
	return p.ratelimitConfig
}

type vpnServer struct{}

// New VPN Server
func NewVPNServer(config *serverConfig) *vpnServer {
	if config == nil {
		config = NewServerConfig()
	}

	// do stuff
	return &vpnServer{}
}

func (p *vpnServer) ListenAndServe() error {
	// start vpn server
	return nil
}

func (p *vpnServer) Shutdown() error {
	// shutdown vpn server
	return nil
}

func (p *vpnServer) ReloadConfig() error {
	// reload config on the fly
	// or restart server if necessary
	return nil
}
