package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	container_metrics "github.com/Inner-daydream/otel_docker_exporter/internal/metrics"
	"github.com/Inner-daydream/otel_docker_exporter/internal/telemetry"
)

func main() {
	// Create an instance of a type that implements ContainerMetricsProvider
	// For example, if DockerMetricsProvider implements ContainerMetricsProvider
	provider, err := container_metrics.NewDockerMetricsProvider()
	if err != nil {
		log.Fatalf("Failed to create Docker metrics provider: %v", err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	containerStatusMetrics, shutdown, err := telemetry.InitTelemetry(ctx)
	if err != nil {
		log.Fatalf("Failed to initialize telemetry: %v", err)
	}
	defer func() {
		if err := shutdown(ctx); err != nil {
			log.Printf("Failed to shut down meter provider: %v", err)
		}
	}()

	// Create a ticker that fires every 15 seconds
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			// The context has been cancelled, stop the program
			return
		case <-ticker.C:
			// The ticker has fired, get and send the container statuses
			containerStatuses, err := provider.GetContainersStatus()
			if err != nil {
				log.Printf("Failed to get container statuses: %v", err)
				continue
			}
			log.Printf("sending container statuses: %v", containerStatuses)
			telemetry.SendContainerStatuses(ctx, containerStatuses, containerStatusMetrics)
		}
	}

}
