package core

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v3"
)

type Config struct {
	Port string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	CheckFreeCouriersInterval time.Duration
	OrderCheckCursorDelta     time.Duration

	KafkaPort    string
	KafkaBrokers []string
	KafkaGroupID string
	KafkaTopic   string

	GRPCServiceOrderServer string

	TokenBucketCapacity   int
	TokenBucketRefillRate int

	RetryMaxAttempts int

	PprofAddress string
}

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

	cfg.OrderCheckCursorDelta = secondsStringToDuration(
		os.Getenv("ORDER_CHECK_CURSOR_DELTA_SECONDS"))
	cfg.CheckFreeCouriersInterval = secondsStringToDuration(
		os.Getenv("CHECK_FREE_COURIERS_INTERVAL_SECONDS"))

	cfg.KafkaBrokers = strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	cfg.KafkaGroupID = os.Getenv("KAFKA_GROUP_ID")
	cfg.KafkaTopic = os.Getenv("KAFKA_TOPIC")
	cfg.KafkaPort = os.Getenv("KAFKA_PORT")

	cfg.GRPCServiceOrderServer = os.Getenv("GRPC_SERVICE_ORDER_SERVER")

	cfg.TokenBucketCapacity = toInt(os.Getenv("TOKEN_BUCKET_CAPACITY"))
	cfg.TokenBucketRefillRate = toInt(os.Getenv("TOKEN_BUCKET_REFILL_RATE"))
	cfg.RetryMaxAttempts = toInt(os.Getenv("RETRY_MAX_ATTEMPTS"))

	cfg.PprofAddress = os.Getenv("PPROF_ADDR")
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

func toInt(value string) int {
	integer, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("invalid integer config value %q: %v", value, err))
	}
	return integer
}

func secondsStringToDuration(value string) time.Duration {
	duration := toInt(value)
	return time.Duration(duration) * time.Second
}

func (c *Config) PostgresDSN() string {
	ssl := c.DBSSLMode
	if ssl == "" {
		ssl = "disable"
	}
	host := c.DBHost
	if host == "" {
		host = "localhost"
	}
	port := c.DBPort
	if port == "" {
		port = "5432"
	}

	user := url.QueryEscape(c.DBUser)
	pass := url.QueryEscape(c.DBPassword)
	db := url.PathEscape(c.DBName)

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, pass, host, port, db, url.QueryEscape(ssl),
	)
}