package main

import (
	"github.com/uptrace/bun"
)

type Measure struct {
	bun.BaseModel       `bun:"table:measures,alias:m"`
	ID                  int    `bun:",pk,autoincrement"`
	MonitoringServiceID int    `bun:"monitoring_service_id"`
	Value               int    `bun:"value"`
	MeasuredAt          string `bun:"measured_at"`
}
