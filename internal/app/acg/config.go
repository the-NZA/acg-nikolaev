package acg

// Config for ACG app
type Config struct {
	BindAddr    string `json:"bind_addr"`
	DatabaseURL string `json:"db_url"`
	LogDebug    bool   `json:"log_debug"`
	SecretKey   string `json:"secret_key"`
}

// NewConfig returns config with mocked values
func NewConfig() *Config {
	return &Config{
		BindAddr:    ":9999",
		DatabaseURL: "mongodb://test:27017",
		LogDebug:    false,
		SecretKey:   "Sample_Secret",
	}
}
