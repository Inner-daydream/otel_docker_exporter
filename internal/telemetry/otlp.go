package telemetry

import (
	"context"
	"fmt"

	"github.com/Inner-daydream/otel_docker_exporter/internal/metrics"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

type ContainerStatusMetrics struct {
	memoryUsagePerMetric   metric.Float64Gauge
	MemoryUsageBytesMetric metric.Int64Gauge
	totalMemoryMetric      metric.Int64Gauge
	cpuUsageMetric         metric.Float64Gauge
	restartCountMetric     metric.Int64Gauge
	stateMetric            metric.Int64Gauge
	healthMetric           metric.Int64Gauge
	uptimeMetric           metric.Int64Gauge
}

type MeterConfiig struct {
	ServiceName      string
	ServiceNamespace string
	Prefix           string
}

func InitTelemetry(ctx context.Context, config MeterConfiig) (*ContainerStatusMetrics, func(context.Context) error, error) {
	metricExporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metrics exporter: %w", err)
	}

	meterProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(metricExporter)),
		sdkmetric.WithResource(resource.NewWithAttributes(
			"",
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.ServiceNamespaceKey.String(config.ServiceNamespace),
		)),
	)
	otel.SetMeterProvider(meterProvider)
	meter := otel.Meter("container_statuses")

	// Initialize metrics
	prefix := config.Prefix
	if len(prefix) > 0 {
		prefix += "."
	}
	memoryUsagePerMetric, err := meter.Float64Gauge(prefix + "container.memory.usage.percentage")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create memory usage metric: %w", err)
	}
	memoryUsageBytesMetric, err := meter.Int64Gauge(prefix + "container.memory.usage.bytes")

	if err != nil {
		return nil, nil, fmt.Errorf("failed to create memory usage metric: %w", err)
	}
	totalMemoryMetric, err := meter.Int64Gauge(prefix + "container.memory.total")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create total memory metric: %w", err)
	}

	cpuUsageMetric, err := meter.Float64Gauge(prefix + "container.cpu.usage")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create CPU usage metric: %w", err)
	}

	restartCountMetric, err := meter.Int64Gauge(prefix + "container.restart.count")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create restart count metric: %w", err)
	}
	healthMetric, err := meter.Int64Gauge(prefix + "container.health")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create health metric: %w", err)
	}
	stateMetric, err := meter.Int64Gauge(prefix + "container.state")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create state metric: %w", err)
	}
	uptimeMetric, err := meter.Int64Gauge(prefix + "container.uptime")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create uptime metric: %w", err)
	}
	return &ContainerStatusMetrics{
		memoryUsagePerMetric:   memoryUsagePerMetric,
		MemoryUsageBytesMetric: memoryUsageBytesMetric,
		totalMemoryMetric:      totalMemoryMetric,
		cpuUsageMetric:         cpuUsageMetric,
		restartCountMetric:     restartCountMetric,
		stateMetric:            stateMetric,
		healthMetric:           healthMetric,
		uptimeMetric:           uptimeMetric,
	}, meterProvider.Shutdown, nil
}

func SendContainerStatuses(ctx context.Context, containerStatuses []metrics.ContainerStatus, metrics *ContainerStatusMetrics) {
	// Record the status values
	for _, status := range containerStatuses {

		// Create a set of attributes to associate with the measurements
		attrs := []attribute.KeyValue{
			attribute.String("container_id", status.ContainerID),
			attribute.String("name", status.Name),
			attribute.String("image", status.Image),
		}
		for k, v := range status.AdditionalLabels {
			attrs = append(attrs, attribute.String(k, v))
		}
		commonAttributes := metric.WithAttributes(attrs...)
		// Record the measurements
		metrics.memoryUsagePerMetric.Record(ctx, status.MemoryUsagePer, commonAttributes)
		metrics.MemoryUsageBytesMetric.Record(ctx, status.MemoryUsageBytes, commonAttributes)
		metrics.totalMemoryMetric.Record(ctx, status.TotalMemory, commonAttributes)
		metrics.cpuUsageMetric.Record(ctx, status.CpuUsage, commonAttributes)
		metrics.restartCountMetric.Record(ctx, status.RestartCount, commonAttributes)
		metrics.healthMetric.Record(ctx, status.Health, commonAttributes)
		metrics.stateMetric.Record(ctx, status.State, commonAttributes)
		metrics.uptimeMetric.Record(ctx, status.Uptime, commonAttributes)
	}
}
