package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServiceName              string `mapstructure:"SERVICE_NAME"`
	OTELExporterOTLPEndpoint string `mapstructure:"OTEL_EXPORTER_OTLP_ENDPOINT"`
	OTELExporterOTLPHeaders  string `mapstructure:"OTEL_EXPORTER_OTLP_HEADERS"`
	OTELResourceAttreibutes  string `mapstructure:"OTEL_RESOURCE_ATTRIBUTES"`
	DBSource                 string `mapstructure:"DB_SOURCE"`
	ClerkKey                 string `mapstructure:"CLERK_KEY"`
	LogFilePath              string `mapstructure:"LOG_FILE_PATH"`
	LokiURL                  string `mapstructure:"LOKI_URL"`
	SyslogAddress            string `mapstructure:"SYSLOG_ADDRESS"`
	SyslogNetwork            string `mapstructure:"SYSLOG_NETWORK"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AutomaticEnv()

	_ = viper.BindEnv("SERVICE_NAME")
	_ = viper.BindEnv("OTEL_EXPORTER_OTLP_ENDPOINT")
	_ = viper.BindEnv("OTEL_EXPORTER_OTLP_HEADERS")
	_ = viper.BindEnv("OTEL_RESOURCE_ATTRIBUTES")
	_ = viper.BindEnv("DB_SOURCE")
	_ = viper.BindEnv("CLERK_KEY")
	_ = viper.BindEnv("LOG_FILE_PATH")
	_ = viper.BindEnv("LOKI_URL")
	_ = viper.BindEnv("SYSLOG_ADDRESS")
	_ = viper.BindEnv("SYSLOG_NETWORK")

	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("Failed to unmarshal config: %w", err)
	}
	return config, nil
}
