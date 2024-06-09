package metrics

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type DockerMetricsProvider struct {
	Client *client.Client
}

func healthStatusToInt(status string) int64 {
	switch status {
	case "starting":
		return 0
	case "healthy":
		return 1
	case "unhealthy":
		return 2
	default:
		return -1 // unknown state
	}
}

func stateStatusToInt(state string) int64 {
	switch state {
	case "created":
		return 0
	case "restarting":
		return 1
	case "running":
		return 2
	case "removing":
		return 3
	case "paused":
		return 4
	case "exited":
		return 5
	case "dead":
		return 6
	default:
		return -1 // unknown state
	}
}

func NewDockerMetricsProvider() (*DockerMetricsProvider, error) {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}
	return &DockerMetricsProvider{Client: client}, nil
}

func (d *DockerMetricsProvider) getContainerDetails(ctx context.Context, containerID string) (types.ContainerJSON, error) {
	containerDetails, err := d.Client.ContainerInspect(ctx, containerID)
	if err != nil {
		return types.ContainerJSON{}, fmt.Errorf("failed to inspect container %s: %w", containerID, err)
	}
	return containerDetails, nil
}

func (d *DockerMetricsProvider) getContainerStats(ctx context.Context, containerID string) (types.StatsJSON, error) {
	containerStats, err := d.Client.ContainerStats(ctx, containerID, false)
	if err != nil {
		return types.StatsJSON{}, fmt.Errorf("failed to get container stats for %s: %w", containerID, err)
	}
	defer containerStats.Body.Close()

	var stats types.StatsJSON
	if err := json.NewDecoder(containerStats.Body).Decode(&stats); err != nil {
		return types.StatsJSON{}, fmt.Errorf("failed to decode container stats for %s: %w", containerID, err)
	}
	return stats, nil
}

func (d *DockerMetricsProvider) createContainerStatus(c types.Container, containerDetails types.ContainerJSON, stats types.StatsJSON, totalMemory int64, totalCpu int, labels map[string]string) ContainerStatus {
	startedAt, _ := time.Parse(time.RFC3339, containerDetails.State.StartedAt)
	var (
		health int64 = -1
		state  int64 = -1
	)
	if containerDetails.State.Health != nil {
		health = healthStatusToInt(containerDetails.State.Health.Status)
	}
	if containerDetails.State != nil {
		state = stateStatusToInt(containerDetails.State.Status)
	}
	return ContainerStatus{
		ContainerID:      c.ID,
		Name:             strings.TrimPrefix(c.Names[0], "/"),
		Health:           health,
		State:            state,
		RestartCount:     int64(containerDetails.RestartCount),
		MemoryUsagePer:   float64(stats.MemoryStats.Usage) / float64(totalMemory) * 100, // Memory usage as a percentage of total memory
		MemoryUsageBytes: int64(stats.MemoryStats.Usage),                                // Memory usage in bytes
		TotalMemory:      totalMemory,
		CpuUsage:         float64(stats.CPUStats.CPUUsage.TotalUsage) / float64(totalCpu*1000000000) * 100, // CPU usage as a percentage of total CPU
		Image:            c.Image,
		Uptime:           int64(time.Since(startedAt).Seconds()),
		AdditionalLabels: labels,
	}
}

func (d *DockerMetricsProvider) GetContainersStatus() ([]ContainerStatus, error) {
	ctx := context.Background()

	// Get system info
	info, err := d.Client.Info(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Docker info: %w", err)
	}
	totalMemory := info.MemTotal
	totalCpu := info.NCPU

	containers, err := d.Client.ContainerList(ctx, container.ListOptions{
		All: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	var statuses []ContainerStatus
	var wg sync.WaitGroup
	statusesMutex := &sync.Mutex{}
	errChan := make(chan error, len(containers))

	for _, container := range containers {
		wg.Add(1)
		go func(c types.Container) {
			defer wg.Done()

			containerDetails, err := d.getContainerDetails(ctx, c.ID)
			if err != nil {
				errChan <- err
				return
			}

			stats, err := d.getContainerStats(ctx, c.ID)
			if err != nil {
				errChan <- err
				return
			}

			additionalLabels := make(map[string]string)

			for k, v := range containerDetails.Config.Labels {
				if strings.HasPrefix(k, "otlp.label.") {
					newKey := strings.TrimPrefix(k, "otlp.label.")
					additionalLabels[newKey] = v
				}
			}

			status := d.createContainerStatus(c, containerDetails, stats, totalMemory, totalCpu, additionalLabels)
			statusesMutex.Lock()
			statuses = append(statuses, status)
			statusesMutex.Unlock()
		}(container)
	}

	wg.Wait()
	close(errChan)

	// Check if any errors occurred in the goroutines
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return statuses, nil
}

var _ ContainerMetricsProvider = &DockerMetricsProvider{}
