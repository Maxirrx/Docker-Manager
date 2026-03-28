package main

func StartService(uuid string) Result {
	err := StartDocker(uuid)
	if err != nil {
		return Result{
			Success: false,
			Message: err.Error(),
		}
	}
	return Result{
		Success: true,
		Message: "Le service a bien été start",
	}
}
func RestartService(uuid string) Result {
	err := RestartDocker(uuid)
		if err != nil {
		return Result{
			Success: false,
			Message: err.Error(),
		}
	}
	return Result{
		Success: true,
		Message: "Le service a bien été restart",
	}
}

func StopService(uuid string) Result {
	err := StopDocker(uuid)
		if err != nil {
		return Result{
			Success: false,
			Message: err.Error(),
		}
	}
	return Result{
		Success: true,
		Message: "Le service a bien été stoppé",
	}
}


func DeleteService(uuid string) Result {
	err := DeleteDocker(uuid)
		if err != nil {
		return Result{
			Success: false,
			Message: err.Error(),
		}
	}
	return Result{
		Success: true,
		Message: "Le service a bien été remove",
	}
}