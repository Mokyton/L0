package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

var initScript = `
	CREATE TABLE IF NOT EXISTS orders (
		id text PRIMARY KEY,
		order_info jsonb
	);
`

type model struct {
	db *sqlx.DB
}

func Connect(conn *string) (*model, error) {
	connConfig, err := pgx.ParseConfig(*conn)
	if err != nil {
		return nil, err
	}

	connStr := stdlib.RegisterConnConfig(connConfig)
	dbh, err := sqlx.Connect("pgx", connStr)
	if err != nil {
		return nil, err
	}

	if err := dbh.Ping(); err != nil {
		return nil, err
	}

	return &model{db: dbh}, nil
}

func (m *model) Insert(uuid string, src []byte) error {
	_, err := m.db.Exec(`INSERT INTO orders (id, order_info) VALUES ($1, $2)`, uuid, src)
	if err != nil {
		return err
	}
	return nil
}

func (m *model) Close() error {
	err := m.db.Close()
	if err != nil {
		return err
	}
	return nil
}

func (m *model) InitTable() error {
	_, err := m.db.ExecContext(context.Background(), initScript)
	if err != nil {
		return err
	}
	return nil
}

func (m *model) Get() map[string][]byte {
	var buf []struct {
		Uid  string `db:"id"`
		Json []byte `db:"order_info"`
	}
	_ = m.db.Select(&buf, `SELECT id, order_info FROM orders`)
	cache := make(map[string][]byte)
	for _, v := range buf {
		cache[v.Uid] = v.Json
	}
	return cache
}
