package providers

type Config struct {
	Username string
}

type Provider struct {
	Config Config
}

var authP Provider

func AuthProvider() *Provider {
	return &authP
}

func SetAuthProvider(config ...Config) {
	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}
	if cfg.Username == "" {
		cfg.Username = "username"
	}
	authP = Provider{Config: cfg}
}
