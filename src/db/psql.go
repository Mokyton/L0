package db

import (
	"bytes"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
)

var initScript = `
	CREATE TABLE IF NOT EXISTS orders (
		id BIGSERIAL PRIMARY KEY,
		order_info json
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

func (m *model) Insert(src []byte) error {
	m.db.Exec(`INSERT INTO orders (order_info) VALUES ($1)`, src)
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
func (m *model) Get() []byte {
	var buf bytes.Buffer
	_ = m.db.Select(&buf, `SELECT order_info FROM orders`)
	return buf.Bytes()
}
