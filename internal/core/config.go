package core

import (
	// "flag"
	"os"
	"context"	
	"github.com/urfave/cli/v3"
	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load("configs/.env")

	cfg := &Config{Port: ":8080"}

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "port",
				Aliases:     []string{"p"},
				Value:       cfg.Port,
				Usage:       "server port",
				Sources:     cli.EnvVars("PORT"), // вместо EnvVars в v3
				Destination: &cfg.Port,           // сразу кладём значение в cfg.Port
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		return nil, err
	}

	if cfg.Port != "" && cfg.Port[0] != ':' {
		cfg.Port = ":" + cfg.Port
	}
	return cfg, nil
}
