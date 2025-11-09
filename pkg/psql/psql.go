package psql

import (
	"context"
	"fmt"
	"net/url"

	"github.com/askaroe/dockify-backend/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client struct {
	*pgxpool.Pool
}

func New(cfg config.Config) (*Client, error) {

	dbQuery := url.Values{}
	dbQuery.Set("sslmode", cfg.DbSslmode)

	dbURL := url.URL{
		Scheme:   "postgresql",
		User:     url.UserPassword(cfg.PostgresConfig.DbUsername, cfg.PostgresConfig.DbPassword),
		Host:     fmt.Sprintf("%s:%s", cfg.PostgresConfig.DbHost, cfg.PostgresConfig.DbPort),
		Path:     "/" + cfg.PostgresConfig.DbName,
		RawQuery: dbQuery.Encode(),
	}

	fmt.Println("Database URL:", dbURL.String())

	dbConfig, err := pgxpool.ParseConfig(dbURL.String())
	if err != nil {
		return nil, err
	}

	db, err := pgxpool.NewWithConfig(context.Background(), dbConfig)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			db.Close()
		}
	}()

	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}
	return &Client{db}, nil
}
