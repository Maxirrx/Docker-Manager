package main

import (
	"github.com/uptrace/bun"
)

type ServicePort struct {
	bun.BaseModel `bun:"table:services_ports,alias:sp"`
	PortId        int      `bun:"port_id"`
	ServiceUuid   string   `bun:"service_uuid"`
	Port          *Port    `bun:"rel:belongs-to,join:port_id=id"`
	Service       *Service `bun:"rel:belongs-to,join:service_uuid=uuid"`
}
