package database

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func NewPostgreSQLConnection(uri string, connections int32) error {
	config, err := pgxpool.ParseConfig(uri)
	if err != nil {
		return err
	}
	config.MaxConns = connections
	//create a new connection pool
	PostgreSQLDB, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return err
	}
	return nil
}
