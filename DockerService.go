package main

import (
	"context"


	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"fmt"
	"time"
	"encoding/json"
)


func NewDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return cli, nil
}



func StartDocker(uuid string) error{
	cli, err := NewDockerClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	containerInfo, err := cli.ContainerInspect(ctx, uuid)
	if err != nil {
    	return err
	}
	if containerInfo.Config.Labels["CssSexy"] != "true"{
		return fmt.Errorf("ce conteneur n'appartient pas a CSSSexy")
	}
	err = cli.ContainerStart(ctx, uuid, container.StartOptions{})
	if err != nil {
		return err
	}

	var repo ServiceRepository
	service, err := repo.FindServiceByUUID(ctx, uuid)
	if err != nil {
	    return err
	}

	service.StatusId = 2

	err = repo.UpdateService(ctx, service)
	if err != nil {
	    return err
	}

	return nil
}

func RestartDocker(uuid string) error{
	cli, err := NewDockerClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	containerInfo, err := cli.ContainerInspect(ctx, uuid)
	if err != nil {
    	return err
	}
	if containerInfo.Config.Labels["CssSexy"] != "true"{
		return fmt.Errorf("ce conteneur n'appartient pas a CSSSexy")
	}
	err = cli.ContainerRestart(ctx, uuid, container.StopOptions{})
	if err != nil {
		return err
	}

	var repo ServiceRepository
	service, err := repo.FindServiceByUUID(ctx, uuid)
	if err != nil {
	    return err
	}
	
	service.StatusId = 2
	
	err = repo.UpdateService(ctx, service)
	if err != nil {
	    return err
	}
	return nil
}

func StopDocker(uuid string) error{
	cli, err := NewDockerClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	containerInfo, err := cli.ContainerInspect(ctx, uuid)
	if err != nil {
    	return err
	}
	if containerInfo.Config.Labels["CssSexy"] != "true"{
		return fmt.Errorf("ce conteneur n'appartient pas a CSSSexy")
	}
	err = cli.ContainerStop(ctx, uuid, container.StopOptions{})
	if err != nil {
		return err
	}
	
	var repo ServiceRepository
	service, err := repo.FindServiceByUUID(ctx, uuid)
	if err != nil {
	    return err
	}
	
	service.StatusId = 3
	
	err = repo.UpdateService(ctx, service)
	if err != nil {
	    return err
	}
	return nil
}

func DeleteDocker(uuid string) error{
	
	cli, err := NewDockerClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	containerInfo, err := cli.ContainerInspect(ctx, uuid)
	if err != nil {
    	return err
	}
	if containerInfo.Config.Labels["CssSexy"] != "true"{
		return fmt.Errorf("ce conteneur n'appartient pas a CSSSexy")
	}
	err = cli.ContainerRemove(ctx, uuid, container.RemoveOptions{})
	if err != nil {
		return err
	}
	StopDocker(uuid)	
	var repo ServiceRepository
	err = repo.DeleteService(ctx, uuid)
	if err != nil {
	    return err
	}
	return nil
}


func GetMonitoring() error {
	cli, err := NewDockerClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return err
	}
	for _, c := range containers{
    if _, ok := c.Labels["CssSexy"]; ok{
	var repo ServiceRepository
	err, ram, cpu := repo.GetMonitoringID(ctx, c.ID)
	if err != nil {
		return err
	}

	statsResponse, err := cli.ContainerStats(ctx, c.ID, false)
    if err != nil {
        return err
    }
    defer statsResponse.Body.Close()

    var stats container.StatsResponse
    if err := json.NewDecoder(statsResponse.Body).Decode(&stats); err != nil {
        return err
    }

    measureRam := Measure  {
			ID: 1,
			MonitoringServiceID: ram, 
			Value:               int(stats.MemoryStats.Usage) / 1024 / 1024,
            MeasuredAt:          time.Now().String(),
	}
	measureCpu := Measure  {
			ID: 1,
			MonitoringServiceID: cpu, 
			Value:               int(stats.CPUStats.CPUUsage.TotalUsage),
            MeasuredAt:          time.Now().String(),
	}
	fmt.Println(measureRam, measureCpu)
    	}
	}
	return nil

}



//func GetMonitoringg() error {
//	cli, err := NewDockerClient()
//	if err != nil {
//		return err
//	}
//
//	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
//	if err != nil {
//		return err
//	}
//
//	re := regexp.MustCompile(`^coucou`)
//
//	serviceRepo := ServiceRepository{DB: DB}
//
//	for _, container := range containers {
//
//		match := false
//		for key := range container.Labels {
//			if re.MatchString(key) {
//				match = true
//				break
//			}
//		}
//		if !match {
//			continue
//		}
//
//		service, err := serviceRepo.FindServiceByUUID(context.Background(), container.ID)
//		if err != nil {
//			fmt.Println("service not found:", err)
//			continue
//		}
//
//		stats, err := cli.ContainerStats(context.Background(), container.ID, false)
//		if err != nil {
//			fmt.Println("stats error:", err)
//			continue
//		}
//
//		var data types.StatsJSON
//		if err := json.NewDecoder(stats.Body).Decode(&data); err != nil {
//			fmt.Println("decode error:", err)
//			stats.Body.Close()
//			continue
//		}
//		stats.Body.Close()
//
//		now := time.Now()
//
//		cpu := int(data.CPUStats.CPUUsage.TotalUsage)
//
//		mem := int(data.MemoryStats.Usage / 1024 / 1024)
//
//		err = MonitoringSave(Measure{
//			MonitoringServiceID: service.ID,
//			Value:               cpu,
//			MeasuredAt:          now,
//		})
//		if err != nil {
//			fmt.Println("save cpu:", err)
//		}
//
//		err = MonitoringSave(Measure{
//			MonitoringServiceID: service.ID,
//			Value:               mem,
//			MeasuredAt:          now,
//		})
//		if err != nil {
//			fmt.Println("save mem:", err)
//		}
//	}
//
//	return nil
//}
//

func WatchContainers(repo *ServiceRepository) error{
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
f := filters.NewArgs()
f.Add("type", "container")
f.Add("event", "die")
f.Add("event", "stop")
f.Add("event", "start")

msgs, errs := cli.Events(context.Background(), events.ListOptions{
	Filters: f,
})
	ctx := context.Background()

	for {
		select {
		case msg := <-msgs:

			if _, ok := msg.Actor.Attributes["CSSexy"]; ok {

				service, err := repo.FindServiceByUUID(ctx, msg.ID)
				if err != nil {
					return err
					continue
				}

				service.StatusId = 4

				err = repo.UpdateService(ctx, service)
				if err != nil {
					return err
				}
			}

		case err := <-errs:
			return err
		}
	}
}


func GetAllDocker() []types.Container {
	cli, err := NewDockerClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	var repo ServiceRepository
	containerDB := repo.GetAllService(ctx)

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
	if err != nil {
		return err
	}
	for _, c := range containers{
    	if _, ok := c.Labels["CssSexy"]; ok{
			for _, cdb := range containerDB{
				if c.ID == cdb.Uuid{
					
				}
			}
		}
	}
	return containers
}