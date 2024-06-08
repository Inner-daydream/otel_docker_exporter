# otel_docker_exporter

This application collects metrics from Docker containers and exports them to an OpenTelemetry collector using gRPC.

## Metrics

The following table describes the metrics that are exported:

| Metric Name | Description | Possible Values |
|-------------|-------------|-----------------|
| `memory_usage` | The percentage of host memory used by the container. | Any value between 0 and 100 |
| `cpu_usage` | The percentage of host CPU used by the container. | Any value between 0 and 100 |
| `restart_count` | The number of times the container has been restarted. | Any non-negative integer |
| `state` | The current state of the container. | -1: unknown, 0: created, 1: restarting, 2: running, 3: removing, 4: paused, 5: exited, 6: dead |
| `health` | The health status of the container. | -1: unknown, 0: starting, 1: healthy, 2: unhealthy |
| `uptime` | The uptime of the container in seconds. | Any non-negative integer |

Each metric is associated with the following labels:

| Label Name | Description |
|------------|-------------|
| `container_id` | The ID of the container. |
| `name` | The name of the container. |
| `image` | The image used by the container. |

## Collection Interval

Metrics are collected every 15 seconds.

## Configuration

The application is configured using the standard OpenTelemetry exporrter configuration:

https://opentelemetry.io/docs/languages/sdk-configuration/otlp-exporter/