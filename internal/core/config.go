package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
}

var (
	PhoneRegex                = `^\+[0-9]{11}$`
	CheckFreeCouriersInterval = 10 * time.Second
)

func LoadConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	cfg := &Config{}

	if err := getCmd(cfg).Run(context.Background(), os.Args); err != nil {
		return nil, err
	}
	cfg.Port = ":" + cfg.Port

	cfg.DBHost = os.Getenv("POSTGRES_HOST")
	cfg.DBPort = os.Getenv("POSTGRES_PORT")
	cfg.DBUser = os.Getenv("POSTGRES_USER")
	cfg.DBPassword = os.Getenv("POSTGRES_PASSWORD")
	cfg.DBName = os.Getenv("POSTGRES_DB")
	cfg.DBSSLMode = os.Getenv("POSTGRES_SSLMODE")

	return cfg, nil
}

func getCmd(cfg *Config) *cli.Command {
	return &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Value:       cfg.Port,
				Usage:       "server port",
				Sources:     cli.EnvVars("PORT"),
				Destination: &cfg.Port,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error { return nil },
	}
}

func (c *Config) DBConnString() string {
	ssl := c.DBSSLMode
	if ssl == "" {
		ssl = "disable"
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		ssl,
	)
}

func DBConnStringFromEnv() string {
	cfg, _ := LoadConfig()
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode,
	)
}

func TestDBConnString() string {
	// Игнорируем ошибку если .env не найден - переменные могут быть в окружении
	_ = godotenv.Load(".env")

	host := getEnvOrDefault("POSTGRES_HOST_TEST", "localhost")
	port := getEnvOrDefault("POSTGRES_PORT_TEST", "5432")
	user := getEnvOrDefault("POSTGRES_USER_TEST", "postgres")
	password := getEnvOrDefault("POSTGRES_PASSWORD_TEST", "postgres")
	dbname := getEnvOrDefault("POSTGRES_DB_TEST", "courier_service_test")
	sslmode := "disable"

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user,
		password,
		host,
		port,
		dbname,
		sslmode,
	)

	log.Printf("TestDBConnString: %s", connStr)
	return connStr
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
