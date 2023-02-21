package server

type serverConfig struct {
	cookieConfig    *cookieConfig
	ratelimitConfig *rateLimitConfig
	loggerConfig    *loggerConfig

	// http server config
	// tls server config
	// web socket config
	// tun/tap networking config
	// some configs that i can't think of right now
}

func NewServerConfig() *serverConfig {

	// default configuration for all modules
	return &serverConfig{
		cookieConfig:    NewCookieConfig(),
		ratelimitConfig: NewrateLimitConfig(),
		loggerConfig:    NewLoggerConfig(),
	}
}

func (p *serverConfig) SetCookieConfig(val *cookieConfig) {
	if val == nil {
		return
	}
	p.cookieConfig = val
}

func (p *serverConfig) GetCookieConfig() *cookieConfig {
	return p.cookieConfig
}

func (p *serverConfig) SetRateLimitConfig(val *rateLimitConfig) {
	if val == nil {
		return
	}
	p.ratelimitConfig = val
}

func (p *serverConfig) GetRateLimitConfig() *rateLimitConfig {
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
