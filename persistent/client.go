package persistent

import (
	"database/sql"

	"github.com/jczornik/glacier_backup/config"
	_ "github.com/mattn/go-sqlite3"
)

type DBClient struct {
	config config.DB
}

func NewDBClient(config config.DB) DBClient {
	return DBClient{config}
}

func (c DBClient) OpenDB() (*sql.DB, error) {
	return sql.Open("sqlite3", c.config.Path)
}
