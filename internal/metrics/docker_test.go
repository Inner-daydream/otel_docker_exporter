package metrics

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestDockerMetricsProvider_GetContainersStatus(t *testing.T) {
	ctx := context.Background()

	// Start a container using testcontainers-go
	test_prefix := "gotest-"
	container_name := test_prefix + "redis"
	req := testcontainers.ContainerRequest{
		Image:      "redis:latest",
		WaitingFor: wait.ForLog("Ready to accept connections tcp"),
		Name:       container_name,
		Labels: map[string]string{
			"otlp.label.test-service": "my-test",
		},
	}
	redis, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer redis.Terminate(ctx)

	// Create a Docker client
	provider, err := NewDockerMetricsProvider()
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name    string
		want    []ContainerStatus
		wantErr bool
	}{
		{
			name: "Get status of redis container",
			want: []ContainerStatus{
				{
					ContainerID:  redis.GetContainerID(),
					Name:         container_name,
					Health:       -1,
					State:        2,
					RestartCount: 0,
					Image:        "redis:latest",
					AdditionalLabels: map[string]string{
						"test-service": "my-test",
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		got, err := provider.GetContainersStatus()
		if (err != nil) != tt.wantErr {
			t.Errorf("DockerMetricsProvider.GetContainersStatus() error = %v, wantErr %v", err, tt.wantErr)
			return
		}

		// Create a new slice that excludes testcontainers containers
		var filteredGot []ContainerStatus
		for _, g := range got {
			if !strings.HasPrefix(g.Name, "testcontainers/") && strings.HasPrefix(g.Name, test_prefix) {
				filteredGot = append(filteredGot, g)
			}

		}

		for _, g := range filteredGot {
			for _, w := range tt.want {
				if g.ContainerID != w.ContainerID ||
					g.Name != w.Name ||
					g.Health != w.Health ||
					g.State != w.State ||
					g.RestartCount != w.RestartCount ||
					g.Image != w.Image {
					t.Errorf("DockerMetricsProvider.GetContainersStatus() = %v, want %v", g, w)
				}
				if !reflect.DeepEqual(g.AdditionalLabels, w.AdditionalLabels) {
					t.Errorf("DockerMetricsProvider.GetContainersStatus() = %v, want %v", g.AdditionalLabels, w.AdditionalLabels)
				}
				// Check if the values are within the expected range
				if g.MemoryUsagePer < 0.0 || g.MemoryUsagePer > 100.0 {
					t.Errorf("MemoryUsagePer is out of range: %v", g.MemoryUsagePer)
				}
				if g.MemoryUsageBytes < 1000 {
					t.Errorf("MemoryUsageBytes is out of range: %v", g.MemoryUsageBytes)
				}
				if g.TotalMemory < 1000 {
					t.Errorf("TotalMemory is out of range: %v", g.TotalMemory)
				}
				if g.CpuUsage < 0.1 || g.CpuUsage > 80 {
					t.Errorf("CpuUsage is out of range: %v", g.CpuUsage)
				}
			}
		}
	}
}
