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
	service := &Service{}
	_, err := r.DB.NewDelete().
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

func(r *ServiceRepository) MonitoringSave(ctx context.Context, measure Measure) error {
	_, err := r.DB.NewInsert().
	Model(measure).
	Exec(ctx)

	return err
}

func(r *ServiceRepository) GetMonitoringID(ctx context.Context,uuidService string)(error, int, int){
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
	return nil, ram["id"].(int), cpu["id"].(int)
}