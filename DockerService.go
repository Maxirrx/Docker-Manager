package main

import (
	"context"

	"strings"

    "io"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/go-connections/nat"
	//"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"fmt"
	"time"
	"encoding/json"
	"strconv"
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
	for _, c := range containers{
    if _, ok := c.Labels["CssSexy"]; ok{
	repo := &ServiceRepository{DB: DB}
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

func CreateDocker(service *Service) error {
	repo := &ServiceRepository{DB: DB}
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

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

	if err := repo.Create(ctx, service); err != nil {
		return err
	}

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
    f.Add("event", "pause")
    f.Add("event", "unpause")
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
            case "start", "unpause":
                statusId = 2 
            case "stop", "pause":
                statusId = 3 
            case "die", "destroy":
                statusId = 4 
            default:
                continue
            }

            service, err := repo.FindServiceByUUID(ctx, msg.Actor.ID)
			if err != nil {

				projectId := 0
    			if val, ok := msg.Actor.Attributes["project"]; ok {
    			    projectId, _ = strconv.Atoi(val)
    			}

				newService := Service{
			        Uuid:      msg.Actor.ID,
			        Name:      msg.Actor.Attributes["name"],
			        Image:     msg.Actor.Attributes["image"],
					StartedSince: time.Now().Format("2006-01-02 15:04:05"),
			        ProjectId:    projectId,
					StatusId:  statusId,
			    }
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

    containers, err := cli.ContainerList(context.Background(), container.ListOptions{})
    if err != nil {
        panic(err)
    }

    containersDocker := make(map[string]Service)
    for _, c := range containers {
        if _, ok := c.Labels["CssSexy"]; !ok {
            continue
        }

        projectId := 0
        if val, ok := c.Labels["project"]; ok {
            projectId, _ = strconv.Atoi(val)
        }

        containersDocker[c.ID] = Service{
            Uuid:         c.ID,
            Image:        c.Image,
            StartedSince: time.Unix(c.Created, 0).Format("2006-01-02 15:04:05"),
            Name:         c.Names[0],
            ProjectId:    projectId,
			StatusId: 1,
        }
    }

    for _, cdb := range containersDB {
        docker, exists := containersDocker[cdb.Uuid]
        if !exists {
            continue
        }

        if docker.Image == cdb.Image && docker.StartedSince == cdb.StartedSince {
            continue
        }

        updated := docker 
        if err := repo.UpdateService(ctx, &updated); err != nil {
            panic(err)
        }
    }
	for uuid, docker := range containersDocker {
   		found := false
   		for _, cdb := range containersDB {
   		    if cdb.Uuid == uuid {
   		        found = true
   		        break
   		    }
   		}

   		if !found {
   		    newService := docker
   		    if err := repo.Create(ctx, &newService); err != nil {
   		        panic(err)
   		    }
   		}
	}
    fmt.Println(containersDB)
}