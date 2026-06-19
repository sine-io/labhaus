package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	MinIO    MinIOConfig
	Log      LogConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port            string
	Environment     string
	ShutdownTimeout int // seconds
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

type LogConfig struct {
	Level  string
	Format string // "json" or "console"
}

type JWTConfig struct {
	SecretKey       string
	TokenDuration   int // hours
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.environment", "development")
	v.SetDefault("server.shutdown_timeout", 30)

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.user", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.dbname", "labhaus")
	v.SetDefault("database.sslmode", "disable")

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)
	v.SetDefault("redis.password", "")
	v.SetDefault("redis.db", 0)

	v.SetDefault("minio.endpoint", "localhost:9000")
	v.SetDefault("minio.access_key", "minioadmin")
	v.SetDefault("minio.secret_key", "minioadmin")
	v.SetDefault("minio.use_ssl", false)

	v.SetDefault("log.level", "info")
	v.SetDefault("log.format", "json")

	v.SetDefault("jwt.secret_key", "your-secret-key-change-in-production")
	v.SetDefault("jwt.token_duration", 24) // 24 hours

	// Read from environment variables
	v.SetEnvPrefix("LABHAUS")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Bind specific environment variables for nested structures
	v.BindEnv("jwt.secret_key", "LABHAUS_JWT_SECRET_KEY")
	v.BindEnv("jwt.token_duration", "LABHAUS_JWT_TOKEN_DURATION")
	v.BindEnv("minio.endpoint", "LABHAUS_MINIO_ENDPOINT")
	v.BindEnv("minio.access_key", "LABHAUS_MINIO_ACCESS_KEY")
	v.BindEnv("minio.secret_key", "LABHAUS_MINIO_SECRET_KEY")
	v.BindEnv("minio.use_ssl", "LABHAUS_MINIO_USE_SSL")

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Manual override for JWT config (Viper nested struct issue)
	if v.IsSet("jwt.secret_key") {
		config.JWT.SecretKey = v.GetString("jwt.secret_key")
	}
	if v.IsSet("jwt.token_duration") {
		config.JWT.TokenDuration = v.GetInt("jwt.token_duration")
	}

	// Manual override for MinIO config (Viper nested struct issue)
	if v.IsSet("minio.endpoint") {
		config.MinIO.Endpoint = v.GetString("minio.endpoint")
	}
	if v.IsSet("minio.access_key") {
		config.MinIO.AccessKey = v.GetString("minio.access_key")
	}
	if v.IsSet("minio.secret_key") {
		config.MinIO.SecretKey = v.GetString("minio.secret_key")
	}
	if v.IsSet("minio.use_ssl") {
		config.MinIO.UseSSL = v.GetBool("minio.use_ssl")
	}

	return &config, nil
}
