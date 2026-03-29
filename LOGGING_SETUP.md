# Logging Setup Guide

This guide explains how to configure your Go application to automatically write logs to `/logs/otel-logs.json` and integrate with your observability stack.

## Current Setup

- **Go Application**: Runs directly on host (not containerized)
- **Observability Stack**: Promtail, OTEL Collector, Loki running in containers
- **Log Path**: `./logs/otel-logs.json` (host) → `/logs/otel-logs.json` (containers)

## Quick Start

### 1. Start the Observability Stack

```bash
make observability-up
```

This starts:
- **Loki** on `http://localhost:3100`
- **Promtail** (scrapes logs from `./logs/`)
- **OTEL Collector** on `http://localhost:4318`
- **Tempo** on `http://localhost:3200`
- **Grafana** on `http://localhost:3000` (admin/admin)

### 2. Run Your Go Application

```bash
make run-with-logs
```

Or manually:
```bash
# Start observability stack
make observability-up

# Run your application
go run main.go
```

### 3. View Logs

- **File**: `./logs/otel-logs.json`
- **Grafana**: http://localhost:3000 → Explore → Loki
- **Direct Loki**: http://localhost:3100

## Configuration

### Environment Variables

The application reads configuration from `app.env`:

```env
SERVICE_NAME=sale-service
LOG_FILE_PATH=./logs/otel-logs.json
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318
# ... other config
```

### Log Format

Logs are written in structured JSON format:

```json
{
  "timestamp": "2025-11-05T10:13:01Z",
  "level": "INFO",
  "source": {
    "function": "main.main",
    "file": "/path/to/file.go",
    "line": 23
  },
  "msg": "Application started",
  "service": "sale-service",
  "version": "1.0.0"
}
```

## Log Collection Flow

```
Go App → ./logs/otel-logs.json → Promtail → Loki → Grafana
```

1. **Go Application** writes structured JSON logs to `./logs/otel-logs.json`
2. **Promtail** scrapes the log file and sends to Loki
3. **Loki** stores and indexes the logs
4. **Grafana** provides visualization and querying

## Useful Commands

```bash
# Start observability stack
make observability-up

# Stop observability stack
make observability-down

# View observability logs
make observability-logs

# Run app with logging
make run-with-logs

# Test logging setup
make testlogging
```

## Grafana Setup

1. Open http://localhost:3000 (admin/admin)
2. Go to **Connections** → **Data Sources**
3. Add **Loki** data source: `http://loki:3100`
4. Go to **Explore** and query logs:
   ```
   {job="otel-app-logs"}
   ```

## Log Queries

### Basic Queries
```
# All application logs
{job="otel-app-logs"}

# Error logs only
{job="otel-app-logs"} |= "ERROR"

# Specific service
{job="otel-app-logs", service="sale-service"}

# Database operations
{job="otel-app-logs"} |= "database"
```

### Advanced Queries
```
# Count errors per minute
count_over_time({job="otel-app-logs"} |= "ERROR" [1m])

# Parse JSON and filter
{job="otel-app-logs"} | json | level="ERROR"
```

## Troubleshooting

### Logs Not Appearing

1. Check if logs directory exists: `ls -la ./logs/`
2. Check if log file is being written: `tail -f ./logs/otel-logs.json`
3. Check Promtail logs: `docker logs promtail`
4. Check Loki logs: `docker logs loki`

### Permission Issues

```bash
# Fix permissions
chmod 755 ./logs
chmod 644 ./logs/otel-logs.json
```

### Container Access

The `./logs` directory is mounted to containers as `/logs`, so:
- Host path: `./logs/otel-logs.json`
- Container path: `/logs/otel-logs.json`

## Production Considerations

1. **Log Rotation**: Implement log rotation for large deployments
2. **Retention**: Configure Loki retention policies
3. **Security**: Secure log files and endpoints
4. **Performance**: Monitor log volume and processing
5. **Backup**: Backup important logs

## Integration with OpenTelemetry

The application also sends traces and metrics to the OTEL Collector:
- **Traces**: Go to Grafana → Tempo
- **Metrics**: Available via OTEL Collector metrics endpoint
- **Logs**: File-based collection via Promtail