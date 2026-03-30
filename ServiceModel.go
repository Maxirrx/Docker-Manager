package main

import (
	"github.com/uptrace/bun"
)

type Service struct {
	bun.BaseModel `bun:"table:services,alias:s"`
	Uuid          string        `bun:"uuid,pk"`
	Image         string    `bun:"image"`
	StartedSince  string 	`bun:"started_since"`
	Name          string    `bun:"name"`
	ProjectId     int       `bun:"project_id"`
	StatusId      int       `bun:"status_id"`
	Ports         []Port 	`bun:"m2m:services_ports,join:Service=Port"`
}