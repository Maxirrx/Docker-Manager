package main

import (
	"github.com/uptrace/bun"
)

type MonitoringService struct {
	bun.BaseModel `bun:"table:monitorings_services"`
	MonitoringID  int    `bun:"monitoring_id"`
	ServiceUUID   string `bun:"service_uuid"`
	MinValue      int    `bun:"min_value"`
	MaxValue      int    `bun:"max_value"`
}
