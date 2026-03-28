package main

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func ConnectDB() (*bun.DB, error) {

	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN("postgres://postgres:exemple@localhost:5433/db_api_culteur?sslmode=disable"),
	))

	db := bun.NewDB(sqldb, pgdialect.New())

	return db, nil
}
