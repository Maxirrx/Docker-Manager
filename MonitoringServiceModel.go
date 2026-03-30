package main

import (
	"github.com/uptrace/bun"
)

type MonitoringService struct {
	bun.BaseModel `bun:"table:monitorings_services"`
    Uuid         string `json:"uuid"`
    Image        string `json:"image"`
    StartedSince string `json:"started_since"`
    Name         string `json:"name"`
    ProjectId    int    `json:"projectid"`
    StatusId     int    `json:"status_id"`
    Ports        []Port `json:"ports"`
}
