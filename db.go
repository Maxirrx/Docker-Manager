package main

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)
var DB *bun.DB  

func ConnectDB() (*bun.DB, error) {

	sqldb := sql.OpenDB(pgdriver.NewConnector(
		pgdriver.WithDSN("postgres://postgres:example@api-culteur-db:5432/db_api_culteur?sslmode=disable"),
	))
	
	db := bun.NewDB(sqldb, pgdialect.New())

	db.RegisterModel((*ServicePort)(nil))

	return db, nil
}
