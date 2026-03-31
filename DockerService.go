package main

import (
	"context"

	"strings"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/go-connections/nat"
	"io"
	//"github.com/docker/docker/api/types"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"strconv"
	"time"
	"log"
)

func NewDockerClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

func StartDocker(uuid string) error {
	cli, err := NewDockerClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	containerInfo, err := cli.ContainerInspect(ctx, uuid)
	if err != nil {
		return err
	}
	if containerInfo.Config.Labels["CssSexy"] != "true" {
		return fmt.Errorf("ce conteneur n'appartient pas a CSSSexy")
	}
	err = cli.ContainerStart(ctx, uuid, container.StartOptions{})
	if err != nil {
		return err
	}

	repo := &ServiceRepository{DB: DB}
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

func RestartDocker(uuid string) error {
	cli, err := NewDockerClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	containerInfo, err := cli.ContainerInspect(ctx, uuid)
	if err != nil {
		return err
	}
	if containerInfo.Config.Labels["CssSexy"] != "true" {
		return fmt.Errorf("ce conteneur n'appartient pas a CSSSexy")
	}
	err = cli.ContainerRestart(ctx, uuid, container.StopOptions{})
	if err != nil {
		return err
	}

	repo := &ServiceRepository{DB: DB}
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

func StopDocker(uuid string) error {
	cli, err := NewDockerClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	containerInfo, err := cli.ContainerInspect(ctx, uuid)
	if err != nil {
		return err
	}
	if containerInfo.Config.Labels["CssSexy"] != "true" {
		return fmt.Errorf("ce conteneur n'appartient pas a CSSSexy")
	}
	err = cli.ContainerStop(ctx, uuid, container.StopOptions{})
	if err != nil {
		return err
	}

	repo := &ServiceRepository{DB: DB}
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

func DeleteDocker(uuid string) error {

	cli, err := NewDockerClient()
	if err != nil {
		return err
	}
	ctx := context.Background()
	containerInfo, err := cli.ContainerInspect(ctx, uuid)
	if err != nil {
		return err
	}
	if containerInfo.Config.Labels["CssSexy"] != "true" {
		return fmt.Errorf("ce conteneur n'appartient pas a CSSSexy")
	}
	StopDocker(uuid)
	err = cli.ContainerRemove(ctx, uuid, container.RemoveOptions{})
	if err != nil {
		return err
	}
	StopDocker(uuid)
	repo := &ServiceRepository{DB: DB}
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
	for _, c := range containers {
		if _, ok := c.Labels["CssSexy"]; ok {
			repo := &ServiceRepository{DB: DB}
			err, ram, cpu := repo.GetMonitoringID(ctx, c.ID)
			if err != nil {
				continue
			}

			statsResponse, err := cli.ContainerStats(ctx, c.ID, false)
			if err != nil {
				log.Printf("stats error: %v", err)
				continue
			}
			statsResponse.Body.Close()

			var stats container.StatsResponse
			if err := json.NewDecoder(statsResponse.Body).Decode(&stats); err != nil {
				log.Printf("stats error: %v", err)
				continue
			}
			cpuDelta := float64(stats.CPUStats.CPUUsage.TotalUsage - stats.PreCPUStats.CPUUsage.TotalUsage)
			systemDelta := float64(stats.CPUStats.SystemUsage - stats.PreCPUStats.SystemUsage)
			cpuPercent := 0.0
			if systemDelta > 0 {
				cpuPercent = (cpuDelta / systemDelta) * float64(len(stats.CPUStats.CPUUsage.PercpuUsage)) * 100.0
			}

			measureRam := Measure{
				MonitoringServiceID: ram,
				Value:               int(stats.MemoryStats.Usage) / 1024 / 1024,
				MeasuredAt:          time.Now().Format("2006-01-02 15:04:05"),
			}
			measureCpu := Measure{
				MonitoringServiceID: cpu,
				Value:               int(cpuPercent),
				MeasuredAt:          time.Now().Format("2006-01-02 15:04:05"),
			}
			err = repo.MonitoringSave(ctx, measureRam)
			if err != nil {
				log.Printf("stats error: %v", err)
				continue
			}
			err = repo.MonitoringSave(ctx, measureCpu)
			if err != nil {
				log.Printf("stats error: %v", err)
				continue
			}
		}
	}
	return nil

}

func CreateDocker(service *Service) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	fmt.Println(service)
	ctx := context.Background()


	reader, err := cli.ImagePull(ctx, service.Image, image.PullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()
	io.Copy(io.Discard, reader)

	exposedPorts := nat.PortSet{}
	portBindings := nat.PortMap{}

	for _, sp := range service.Ports {
		parts := strings.Split(sp.Libelle, ":")
		if len(parts) != 2 {
			return fmt.Errorf("format de port invalide: %s", sp.Libelle)
		}
		hostPort := parts[0]
		containerPort := parts[1]

		port, err := nat.NewPort("tcp", containerPort)
		if err != nil {
			return err
		}
		exposedPorts[port] = struct{}{}
		portBindings[port] = []nat.PortBinding{
			{HostIP: "0.0.0.0", HostPort: hostPort},
		}
	}

	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        service.Image,
			ExposedPorts: exposedPorts,
			Labels: map[string]string{
				"CssSexy": "true",
				"project": strconv.Itoa(service.ProjectId),
			},
		},
		&container.HostConfig{
			PortBindings: portBindings,
		},
		nil,
		nil,
		service.Name,
	)
	if err != nil {
		return err
	}

	err = cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
	if err != nil {
		return err
	}

	service.Uuid = resp.ID
	service.StartedSince = time.Now().Format("2006-01-02 15:04:05")
	service.StatusId = 2

	return nil
}

func WatchContainers() {
	repo := &ServiceRepository{DB: DB}
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		panic(err)
	}

	f := filters.NewArgs()
	f.Add("type", "container")
	f.Add("event", "start")
	f.Add("event", "stop")
	f.Add("event", "die")
	f.Add("event", "starting")
	f.Add("event", "destroy")

	msgs, errs := cli.Events(context.Background(), events.ListOptions{
		Filters: f,
	})

	ctx := context.Background()

	for {
		select {
		case msg := <-msgs:
			if _, ok := msg.Actor.Attributes["CssSexy"]; !ok {
				continue
			}

			statusId := 0
			switch msg.Action {
			case "start":
				statusId = 1
			case "stop", "pause":
				statusId = 3
			case "die":
				statusId = 4
			case "destroy":
				err := repo.DeleteService(ctx, msg.Actor.ID)
				if err != nil {
					panic(err)
				}
				continue
			default:
				continue
			}

			service, err := repo.FindServiceByUUID(ctx, msg.Actor.ID)
			if err != nil {

				inspect, err := cli.ContainerInspect(ctx, msg.Actor.ID)
    		if err != nil {
    		    fmt.Println("erreur inspect:", err)
    		    continue
    		}
		
    		var ports []Port
    		for containerPort, bindings := range inspect.HostConfig.PortBindings {
    		    for _, binding := range bindings {
    		        ports = append(ports, Port{
    		            Libelle: fmt.Sprintf("%s:%s", binding.HostPort, containerPort.Port()),
    		        })
    		    }
    		}

				projectId := 0
				if val, ok := msg.Actor.Attributes["project"]; ok {
					projectId, _ = strconv.Atoi(val)
				}

				newService := Service{
					Uuid:         msg.Actor.ID,
					Name:         msg.Actor.Attributes["name"],
					Image:        msg.Actor.Attributes["image"],
					StartedSince: time.Now().Format("2006-01-02 15:04:05"),
					ProjectId:    projectId,
					StatusId:     statusId,
					Ports: ports,
				}
				fmt.Println(newService)
				if err := repo.Create(ctx, &newService); err != nil {
					fmt.Println("erreur création service:", err)
				}
				continue
			}

			service.StatusId = statusId

			if err := repo.UpdateService(ctx, service); err != nil {
				fmt.Println("erreur update:", err)
				continue
			}

		case err := <-errs:
			panic(err)
		}
	}
}

func GetAllDocker() {
	cli, err := NewDockerClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	repo := &ServiceRepository{DB: DB}

	containersDB, err := repo.GetAllServices(ctx)
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		panic(err)
	}

	containersDocker := make(map[string]Service)
	for _, c := range containers {
		if _, ok := c.Labels["CssSexy"]; !ok {
			continue
		}
		var ports []Port

		for _, port := range c.Ports {
			var p = Port{
				ID:      0,
				Libelle: fmt.Sprintf("%d:%d", port.PublicPort, port.PrivatePort),
			}
			ports = append(ports, p)
		}

		inspect, err := cli.ContainerInspect(context.Background(), c.ID)
		if err != nil {
			continue
		}

		debut := time.Now().Format("2006-01-02 15:04:05")
		if inspect.State.StartedAt != "" {
			t, err := time.Parse(time.RFC3339Nano, inspect.State.StartedAt)
			if err == nil {
				debut = t.Format("2006-01-02 15:04:05")
			}
		}
		projectId := 0
		if val, ok := c.Labels["project"]; ok {
			projectId, _ = strconv.Atoi(val)
		}

		service := Service{
			Uuid:         c.ID,
			Image:        c.Image,
			StartedSince: debut,
			Name:         c.Names[0],
			ProjectId:    projectId,
			StatusId:     1,
			Ports:        ports,
		}

		containersDocker[c.ID] = service
	}

	containersDB, err = repo.GetAllServices(ctx)
	if err != nil {
		panic(err)
	}

	for _, cdb := range containersDB {
		docker, exists := containersDocker[cdb.Uuid]

		if !exists {
			err := repo.DeleteService(ctx, cdb.Uuid)
			if err != nil {
				panic(err)
			}
			continue
		}
		delete(containersDocker, cdb.Uuid)
		if docker.Image == cdb.Image && docker.StartedSince == cdb.StartedSince {
			continue
		}

		updated := docker
		if docker.StartedSince == "" {
			updated.StartedSince = cdb.StartedSince
		}
		if err := repo.UpdateService(ctx, &updated); err != nil {
			panic(err)
		}
	}

	for _, cverif := range containersDocker {
		err = repo.Create(ctx, &cverif)
		if err != nil {
			panic(err)
		}
	}
}


