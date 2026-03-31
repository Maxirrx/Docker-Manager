package main

import (
	"github.com/uptrace/bun"
)

type Port struct {
	bun.BaseModel `bun:"table:ports,alias:p"`
	ID            int    `bun:"id,pk,autoincrement" json:"Id"`
	Libelle       string `bun:"libelle" json:"Libelle"`
}
