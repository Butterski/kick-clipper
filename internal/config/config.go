package config

// Config holds the application configuration
type Config struct {
	ProxyFile     string
	NoproxyMode   bool
	BotCount      int
	ViewsPerBot   int
	ClipURL       string
	NoInteractive bool
	Timeout       int
	MinDelay      int
	MaxDelay      int
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() *Config {
	return &Config{
		ProxyFile:     "proxies.txt",
		NoproxyMode:   false,
		BotCount:      0, // Will be auto-calculated
		ViewsPerBot:   0, // Will be prompted
		ClipURL:       "",
		NoInteractive: false,
		Timeout:       15, // seconds
		MinDelay:      2,  // seconds (with proxy)
		MaxDelay:      8,  // seconds (with proxy)
	}
}

// GetDelay returns appropriate delay range based on proxy mode
func (c *Config) GetDelayRange() (int, int) {
	if c.NoproxyMode {
		return 5, 15 // Longer delays without proxy
	}
	return c.MinDelay, c.MaxDelay
}
