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
		Exec(ctx)

	return err
}

func (r *ServiceRepository) DeleteService(ctx context.Context, uuid string) error {
	_, err := r.DB.NewDelete().
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
