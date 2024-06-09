package telemetry

import (
	"context"
	"testing"

	container_metrics "github.com/Inner-daydream/otel_docker_exporter/internal/metrics"
)

func TestInitTelemetry(t *testing.T) {
	ctx := context.Background()
	config := MeterConfiig{
		ServiceName:      "test-service",
		ServiceNamespace: "test-namespace",
		Prefix:           "test",
	}

	_, _, err := InitTelemetry(ctx, config)
	if err != nil {
		t.Fatalf("failed to initialize telemetry: %v", err)
	}
}

func TestSendContainerStatuses(t *testing.T) {
	ctx := context.Background()
	config := MeterConfiig{
		ServiceName:      "test-service",
		ServiceNamespace: "test-namespace",
		Prefix:           "test",
	}

	metrics, _, err := InitTelemetry(ctx, config)
	if err != nil {
		t.Fatalf("failed to initialize telemetry: %v", err)
	}

	containerStatuses := []container_metrics.ContainerStatus{
		{
			ContainerID:      "test-id",
			Name:             "test-name",
			Image:            "test-image",
			MemoryUsagePer:   50.0,
			MemoryUsageBytes: 500,
			TotalMemory:      1000,
			CpuUsage:         0.5,
			RestartCount:     1,
			Health:           1,
			State:            1,
			Uptime:           100,
			AdditionalLabels: map[string]string{"label": "value"},
		},
	}
	// not really testing anything here as I'm not sure how to do it, just making sure the function doesn't panic
	SendContainerStatuses(ctx, containerStatuses, metrics)

}
