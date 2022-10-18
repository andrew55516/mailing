package pg

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/url"
)

// Config struct represents the configuration for connection to database
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DbName   string
	Timeout  int
}

// NewPoolConfig creates a new pool configuration
func NewPoolConfig(cfg *Config) (*pgxpool.Config, error) {
	connStr :=
		fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%d",
			"postgres",
			url.QueryEscape(cfg.Username),
			url.QueryEscape(cfg.Password),
			cfg.Host,
			cfg.Port,
			cfg.DbName,
			cfg.Timeout)

	poolConfig, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	return poolConfig, nil
}

// NewConnection creates a connection with poolConfig
func NewConnection(poolConfig *pgxpool.Config) (*pgxpool.Pool, error) {
	conn, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
