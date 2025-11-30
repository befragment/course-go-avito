package core

import (
	"context"
	"fmt"
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
	PhoneRegex 					= `^\+[0-9]{11}$`
	CheckFreeCouriersInterval 	= 10 * time.Second
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
