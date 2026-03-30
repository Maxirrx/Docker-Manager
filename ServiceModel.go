package main

import (
	"github.com/uptrace/bun"
)

type Service struct {
	bun.BaseModel `bun:"table:services,alias:s"`
	Uuid          string `bun:"uuid,pk" json:"uuid"`
	Image         string `bun:"image" json:"image"`
	StartedSince  string `bun:"started_since" json:"started_since"`
	Name          string `bun:"name" json:"name"`
	ProjectId     int    `bun:"project_id" json:"projectid"`
	StatusId      int    `bun:"status_id" json:"status_id"`
	Ports         []Port `bun:"m2m:services_ports,join:Service=Port" json:"ports"`
}
