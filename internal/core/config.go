package core

import (
	"fmt"
	// "log"	
	"os"
	"context"	

	"github.com/urfave/cli/v3"
	"github.com/joho/godotenv"

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
	Port 	string
	DBCfg	DBConfig
}

func initEnv() {
	_ = godotenv.Load("configs/.env")
}

func LoadConfig() (*Config, error) {
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

func getDBConnectionString() string {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")

	sslmode := os.Getenv("POSTGRES_SSLMODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbname, sslmode,
	)
}