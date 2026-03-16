package main

import (
	"time"

	"github.com/uptrace/bun"
)

type Service struct {
	bun.BaseModel `bun:"table:users,alias:u"`
	Uuid          string    `bun:"uuid"`
	Image         string    `bun:"image"`
	StartedSince  time.Time `bun:"started_since"`
	Name          string    `bun:"name"`
	ProjectId     int       `bun:"project_id"`
	StatusId      int       `bun:"status_id"`
}
