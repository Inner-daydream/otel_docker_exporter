FROM golang:1.22.4 as builder

WORKDIR /app

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
WORKDIR /app/cmd
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


FROM alpine:latest
WORKDIR /root/

# Set the source repository
LABEL org.opencontainers.image.source https://github.com/Inner-daydream/otel_docker_exporter

# Copy the binary from the builder stage
COPY --from=builder /app/cmd/main .

# Run the application
CMD ["./main"]