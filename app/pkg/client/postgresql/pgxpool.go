package postgresql

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"prod/pkg/logging"
	"time"
)

type pgConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func NewPgConfig(host string, port string, username string, password string, database string) *pgConfig {
	return &pgConfig{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Database: database,
	}
}

func NewClient(ctx context.Context, maxAttempts int, maxDelay time.Duration, cfg *pgConfig) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s",
		cfg.Username, cfg.Password,
		cfg.Host, cfg.Port, cfg.Database,
	)

	err = DoWithAttempts(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		pgxCfg, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			log.Fatalf("Unable to parse config: $v\n", err)
		}

		pgxCfg.ConnConfig.Logger = logrusadapter.NewLogger(logging.GetLogger(ctx))
		pool, err = pgxpool.ConnectConfig(ctx, pgxCfg)
		if err != nil {
			log.Println("Failed to connect to postgres... Going to do the next attempt")
			return err
		}

		return nil
	}, maxAttempts, maxDelay)

	if err != nil {
		log.Fatal("All attempts are exceeded. Unable to connect to postgres")
	}

	return pool, nil
}

func DoWithAttempts(fn func() error, maxAttempts int, delay time.Duration) error {
	var err error

	for maxAttempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			maxAttempts--
			continue
		}
		return nil
	}

	return err
}
