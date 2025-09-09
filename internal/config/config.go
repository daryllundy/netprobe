package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	LogLevel   string          `mapstructure:"log_level"`
	LogFormat  string          `mapstructure:"log_format"`
	LogFile    string          `mapstructure:"log_file"`
	FIPS       bool            `mapstructure:"fips_mode"`
	ConfigFile string          `mapstructure:"-"`
	Logging    LoggingConfig   `mapstructure:"logging"`
	Security   SecurityConfig  `mapstructure:"security"`
	Scan       ScanConfig      `mapstructure:"scan"`
	Bandwidth  BandwidthConfig `mapstructure:"bandwidth"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	File   string `mapstructure:"file"`
}

type SecurityConfig struct {
	Encryption    string   `mapstructure:"encryption"`
	FIPSMode      bool     `mapstructure:"fips_mode"`
	MinTLSVersion string   `mapstructure:"min_tls_version"`
	CipherSuites  []string `mapstructure:"cipher_suites"`
}

type ScanConfig struct {
	Timeout       string `mapstructure:"timeout"`
	MaxConcurrent int    `mapstructure:"max_concurrent"`
}

type BandwidthConfig struct {
	Timeout      string `mapstructure:"timeout"`
	TestDuration string `mapstructure:"test_duration"`
	BufferSize   int    `mapstructure:"buffer_size"`
	Protocol     string `mapstructure:"protocol"`
}

func Load() (*Config, error) {
	cfg := &Config{
		LogLevel:  "info",
		LogFormat: "json",
		FIPS:      false,
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
		Security: SecurityConfig{
			Encryption:    "aes-256-gcm",
			FIPSMode:      false,
			MinTLSVersion: "1.2",
			CipherSuites: []string{
				"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
				"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
			},
		},
		Scan: ScanConfig{
			Timeout:       "5s",
			MaxConcurrent: 100,
		},
		Bandwidth: BandwidthConfig{
			Timeout:      "30s",
			TestDuration: "10s",
			BufferSize:   65536,
			Protocol:     "tcp",
		},
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Add config paths
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.netprobe")
	viper.AddConfigPath("/etc/netprobe")

	// Environment variables
	viper.SetEnvPrefix("NETPROBE")
	viper.AutomaticEnv()

	// Read config file if it exists
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config: %w", err)
		}
	}

	// Unmarshal config
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate config
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	// Add validation logic here
	return nil
}

func (c *Config) SaveDefault(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	viper.SetConfigFile(path)
	return viper.WriteConfig()
}
