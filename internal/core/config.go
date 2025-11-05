package core

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type Config struct {
	Port  string
	DBCfg DBConfig
}

var (
	PhoneRegex = `^\+[0-9]{11}$`
)

func LoadConfig() (*Config, error) {
	_ = godotenv.Load(".env")

	cfg := &Config{}

	if err := getCmd(cfg).Run(context.Background(), os.Args); err != nil {
		return nil, err
	}
	cfg.Port = ":" + cfg.Port

	cfg.DBCfg = DBConfig{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Name:     os.Getenv("POSTGRES_DB"),
		SSLMode:  os.Getenv("POSTGRES_SSLMODE"),
	}

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
	ssl := c.DBCfg.SSLMode
	if ssl == "" {
		ssl = "disable"
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBCfg.User,
		c.DBCfg.Password,
		c.DBCfg.Host,
		c.DBCfg.Port,
		c.DBCfg.Name,
		ssl,
	)
}
