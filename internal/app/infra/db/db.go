package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type DbConfig struct {
	Host         string
	Port         int
	Username     string
	Password     string
	DatabaseName string
}

type Db struct {
	Config DbConfig
}

func (db *Db) Connect() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%d/%s", db.Config.Username, db.Config.Password, db.Config.Host, db.Config.Port, db.Config.DatabaseName))
	if err != nil {
		return nil, err
	}
	return conn, nil
}

type Pair[T, V any] struct {
	First  T
	Second V
}
