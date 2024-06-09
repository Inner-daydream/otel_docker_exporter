# otel_docker_exporter

This application collects metrics from Docker containers and exports them to an OpenTelemetry collector using gRPC.

## Metrics

The following table describes the metrics that are exported:
| Metric Name | Description | Possible Values |
|-------------|-------------|-----------------|
| `container.memory.usage.percentage` | The percentage of host memory used by the container. | Any value between 0 and 100 |
| `container.memory.usage.bytes` | The amount of memory used by the container in bytes. | Any non-negative integer |
| `container.memory.total` | The total amount of memory available to the container in bytes. | Any non-negative integer |
| `container.cpu.usage` | The percentage of host CPU used by the container. | Any value between 0 and 100 |
| `container.restart.count` | The number of times the container has been restarted. | Any non-negative integer |
| `container.state` | The current state of the container. | -1: unknown, 0: created, 1: restarting, 2: running, 3: removing, 4: paused, 5: exited, 6: dead |
| `container.health` | The health status of the container. | -1: unknown, 0: starting, 1: healthy, 2: unhealthy |
| `container.uptime` | The uptime of the container in seconds. | Any non-negative integer |


Each metric is associated with the following labels:

| Label Name | Description |
|------------|-------------|
| `container_id` | The ID of the container. |
| `name` | The name of the container. |
| `image` | The image used by the container. |


## Additional Labels

The application can also export additional labels from the Docker containers. Any Docker label that has a key starting with `otlp.label.` will be exported as an additional label with the metric data. The `otlp.label.` prefix will be removed from the key when it is exported.

For example, if you have a Docker container with the following labels:

```shell
docker run -d --label otlp.label.description=WebServer --label otlp.label.department=IT my-web-app
```

These labels will be exported as description=WebServer and department=IT with the metric data.

This allows you to add arbitrary labels to your Docker containers and have those labels exported with your metric data.


## Configuration

The application is configured using the standard OpenTelemetry exporter configuration:

[OpenTelemetry SDK Configuration - OTLP Exporter](https://opentelemetry.io/docs/languages/sdk-configuration/otlp-exporter/)

Note that it only exports metrics using gRPC.

You can configure the application using the following environment variables:

| Environment Variable | Description | Default Value |
| -------------------- | ----------- | ------------- |
| `SERVICE_NAME`       | The name of the service. | "otel-docker-exporter" |
| `SERVICE_NAMESPACE`  | The namespace of the service. | "default" |
| `INTERVAL`           | The interval for metrics export, in seconds. | 15 |
| `PREFIX`             | The prefix for the metric names. For example, if the prefix is set to "docker", the metric names will be like "docker.container.memory.usage.percentage". | "" |


## docker compose

here is a sample docker compose file to run the exporter

```yaml
name: monitoring
services:
    otel_docker_exporter:
      restart: unless-stopped
      image: ghcr.io/inner-daydream/otel_docker_exporter:2.0.0
      environment:
        - OTEL_EXPORTER_OTLP_ENDPOINT=http://my-grpc-otlp-endpoint:12345
        - INTERVAL=120
        - PREFIX=docker
      volumes:
        - /var/run/docker.sock:/var/run/docker.sock:ro
```

## Building

To build the application you need go installed, run the following command:

```shell
go build -o otel-docker-exporter cmd/otel-docker-exporter/main.go
```

you can also run make build to build the application for every platform.

## Running
You can run the application using the following command:

```shell
./otel-docker-exporter
```