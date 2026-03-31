package main

import (
	"context"
	"github.com/uptrace/bun"
)

type ServiceRepository struct {
	DB *bun.DB
}

func (r *ServiceRepository) Create(ctx context.Context, service *Service) error {
	_, err := r.DB.NewInsert().
		Model(service).
		Exec(ctx)
	if err != nil {
		return err
	}
	
	for _, port := range service.Ports{
		var portId int
		err = r.DB.QueryRowContext(ctx,
    	"INSERT INTO ports (libelle) VALUES (?) ON CONFLICT (libelle) DO UPDATE SET libelle = EXCLUDED.libelle RETURNING id",
    	port,
		).Scan(&portId)
		if err != nil {
		    return err
		}

		_, err = r.DB.ExecContext(ctx,
		    "INSERT INTO services_ports (port_id, service_uuid) VALUES (?, ?)",
		    portId, service.Uuid,
		)
	}


	monitorings := []MonitoringService{
		{MonitoringID: 1, ServiceUUID: service.Uuid, MinValue: 0, MaxValue: 100},
		{MonitoringID: 2, ServiceUUID: service.Uuid, MinValue: 0, MaxValue: 1000},
	}

	_, err = r.DB.NewInsert().
		Model(&monitorings).
		Exec(ctx)

	return err
}

func (r *ServiceRepository) UpdateService(ctx context.Context, service *Service) error {
	_, err := r.DB.NewUpdate().
		Model(service).
		Where("uuid = ?", service.Uuid).
		Exec(ctx)

	return err
}

func (r *ServiceRepository) DeleteService(ctx context.Context, uuid string) error {

	_, err := r.DB.NewDelete().
		Model((*Measure)(nil)).
		Where("monitoring_service_id IN (SELECT id FROM monitorings_services WHERE service_uuid = ?)", uuid).
		Exec(ctx)
	if err != nil {
		return err
	}

	_, err = r.DB.NewDelete().
		Model((*MonitoringService)(nil)).
		Where("service_uuid = ?", uuid).
		Exec(ctx)

	if err != nil {
		return err
	}

	service := &Service{}
	_, err = r.DB.NewDelete().
		Model(service).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}

func (r *ServiceRepository) FindServiceByUUID(ctx context.Context, uuid string) (*Service, error) {
	service := new(Service)

	err := r.DB.NewSelect().
		Model(service).
		Where("uuid = ?", uuid).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (r *ServiceRepository) GetAllServices(ctx context.Context) ([]Service, error) {
	services := []Service{}

	err := r.DB.NewSelect().
		Model(&services).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return services, nil
}

func (r *ServiceRepository) MonitoringSave(ctx context.Context, measure Measure) error {
	_, err := r.DB.NewInsert().
		Model(&measure).
		Exec(ctx)

	return err
}

func (r *ServiceRepository) GetMonitoringID(ctx context.Context, uuidService string) (error, int, int) {
	var ram map[string]interface{}
	err := r.DB.NewSelect().
		TableExpr("monitorings_services").
		Column("id").
		Where("monitoring_id = ?", 1).
		Where("service_uuid = ?", uuidService).
		Scan(ctx, &ram)
	if err != nil {
		return err, 0, 0
	}
	var cpu map[string]interface{}
	err = r.DB.NewSelect().
		TableExpr("monitorings_services").
		Column("id").
		Where("monitoring_id = ?", 2).
		Where("service_uuid = ?", uuidService).
		Scan(ctx, &cpu)
	if err != nil {
		return err, 0, 0
	}
	return nil, int(ram["id"].(int64)), int(cpu["id"].(int64))
}
