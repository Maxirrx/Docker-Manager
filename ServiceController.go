package main

import "context"

func StartService(uuid string) Result {
	StartDocker(uuid)
	var serviceRepo ServiceRepository
	service, err := serviceRepo.FindServiceByUUID(context.Background(), uuid)
	if err != nil {
		panic(err)
	}
	service.StatusId = 1
	err = serviceRepo.UpdateService(context.Background(), service)
	if err != nil {
		panic(err)
	}

	return Result{
		Success: true,
		Data:    service,
		Message: "Le service a bien été start",
	}
}

func RestartService(uuid string) Result {
	RestartDocker(uuid)
	var serviceRepo ServiceRepository
	service, err := serviceRepo.FindServiceByUUID(context.Background(), uuid)
	if err != nil {
		panic(err)
	}
	service.StatusId = 2
	err = serviceRepo.UpdateService(context.Background(), service)
	if err != nil {
		panic(err)
	}

	return Result{
		Success: true,
		Data:    service,
		Message: "Le service a bien été restart",
	}
}

func StopService(uuid string) Result {
	StopDocker(uuid)
	var serviceRepo ServiceRepository
	service, err := serviceRepo.FindServiceByUUID(context.Background(), uuid)
	if err != nil {
		panic(err)
	}
	service.StatusId = 3
	err = serviceRepo.UpdateService(context.Background(), service)
	if err != nil {
		panic(err)
	}

	return Result{
		Success: true,
		Data:    service,
		Message: "Le service a bien été stoppé",
	}
}

func CreateService(uuid string) {

}

func DeleteService(uuid string) Result {
	DeleteDocker(uuid)
	var serviceRepo ServiceRepository
	err := serviceRepo.DeleteService(context.Background(), uuid)
	if err != nil {
		panic(err)
	}

	return Result{
		Success: true,
		Message: "Le service a bien été remove",
	}
}
