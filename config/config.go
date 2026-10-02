package config

type Config struct {
	DatabasePath string
	AWSRegion    string
	AWSProfile   string
	Port         int
	LogLevel     string
}

func DefaultConfig() *Config {
	return &Config{
		DatabasePath: "./data/autoscaler.db",
		AWSRegion:    "us-east-1",
		AWSProfile:   "default",
		Port:         8080,
		LogLevel:     "info",
	}
}
