package main

type Measure struct {
	ID                   int       `bun:",pk,autoincrement"`
	MonitoringServiceID  int       `bun:"monitoring_service_id"`
	Value                int       `bun:"value"`
	MeasuredAt           string 	`bun:"measured_at"`
}