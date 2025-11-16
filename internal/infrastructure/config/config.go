package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"
)

type Config struct {
	DB      DBConfig
	Server  ServerConfig
	Command Command
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type Command struct {
	Name string
	Args []string
}

func (c *Command) IsMigrationCommand() bool {
	return c.Name == "migrate" || c.Name == "rollback" || c.Name == "migration-status"
}

func Load(logger *slog.Logger) *Config {
	command := parseCommand()

	dbConfig := DBConfig{
		Host:            getEnvWithDefault("DB_HOST", "localhost"),
		Port:            getEnvWithDefault("DB_PORT", "5432"),
		User:            getEnvWithDefault("DB_USER", "postgres"),
		Password:        getEnvWithDefault("DB_PASSWORD", "postgres"),
		DBName:          getEnvWithDefault("DB_NAME", "pr_reviewer"),
		SSLMode:         getEnvWithDefault("DB_SSL_MODE", "disable"),
		MigrationsPath:  getEnvWithDefault("MIGRATIONS_PATH", "./migrations"),
		MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
		ConnMaxLifetime: time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME_SECONDS", 300)) * time.Second,
	}

	serverConfig := ServerConfig{
		Port:         getEnvWithDefault("SERVER_PORT", "8080"),
		ReadTimeout:  time.Duration(getEnvInt("SERVER_READ_TIMEOUT", 15)) * time.Second,
		WriteTimeout: time.Duration(getEnvInt("SERVER_WRITE_TIMEOUT", 15)) * time.Second,
		IdleTimeout:  time.Duration(getEnvInt("SERVER_IDLE_TIMEOUT", 60)) * time.Second,
	}

	return &Config{
		DB:      dbConfig,
		Server:  serverConfig,
		Command: command,
	}
}

func parseCommand() Command {
	if len(os.Args) > 1 {
		return Command{
			Name: os.Args[1],
			Args: os.Args[2:],
		}
	}
	return Command{}
}

func getEnvInt(key string, defaultValue int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultValue
}

func getEnvWithDefault(key string, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
