# Direct Logging Setup (No File Writing)

This guide shows you how to implement logging without writing to files, using direct network-based approaches.

## Available Logging Options (Priority Order)

### 1. 🚀 OTLP Direct Logs (Recommended)
**Best for OpenTelemetry setups**

```bash
# Set in app.env
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318
```

**Flow**: Go App → OTLP HTTP → OTEL Collector → Loki → Grafana

**Pros**: 
- Native OpenTelemetry integration
- Structured data with traces/metrics correlation
- Automatic retries and batching
- Industry standard

### 2. 🎯 Direct Loki HTTP
**Best for Loki-focused setups**

```bash
# Set in app.env
LOKI_URL=http://localhost:3100
```

**Flow**: Go App → Loki HTTP API → Loki → Grafana

**Pros**:
- Direct to Loki, no intermediaries
- Fast and efficient
- Custom labels and structured data
- Real-time log streaming

### 3. 📡 Syslog
**Best for traditional infrastructure**

```bash
# Set in app.env
SYSLOG_ADDRESS=localhost:514
SYSLOG_NETWORK=udp
```

**Flow**: Go App → Syslog → Rsyslog/Syslog-ng → Log aggregator

**Pros**:
- Standard protocol
- Works with existing infrastructure
- Network resilient
- Widely supported

### 4. 📁 File Logging (Fallback)
**Only if network options fail**

```bash
# Set in app.env
LOG_FILE_PATH=./logs/otel-logs.json
```

## Quick Setup

### Option 1: OTLP Logs (Recommended)

1. **Start observability stack**:
```bash
make observability-up
```

2. **Configure environment** (app.env):
```env
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318
```

3. **Run application**:
```bash
go run main.go
```

4. **View logs in Grafana**:
   - Go to http://localhost:3000
   - Explore → Loki
   - Query: `{job="go-direct"}`

### Option 2: Direct Loki

1. **Start Loki**:
```bash
docker run -d --name loki -p 3100:3100 grafana/loki:latest
```

2. **Configure environment**:
```env
LOKI_URL=http://localhost:3100
```

3. **Run and view logs**:
```bash
go run main.go
# Logs go directly to Loki at http://localhost:3100
```

### Option 3: Syslog

1. **Configure syslog server** (rsyslog example):
```bash
# /etc/rsyslog.conf
$ModLoad imudp
$UDPServerRun 514
$UDPServerAddress 127.0.0.1
```

2. **Configure environment**:
```env
SYSLOG_ADDRESS=localhost:514
SYSLOG_NETWORK=udp
```

3. **Run application**:
```bash
go run main.go
# Check /var/log/messages or configured syslog destination
```

## Implementation Details

### OTLP Handler Features
- ✅ Automatic JSON structured logs
- ✅ OpenTelemetry correlation IDs
- ✅ Async sending (non-blocking)
- ✅ Fallback to stdout on failure
- ✅ Configurable retry logic

### Loki Handler Features
- ✅ Direct HTTP API calls
- ✅ Custom labels and streams
- ✅ JSON log parsing
- ✅ Timestamp handling
- ✅ Async sending

### Syslog Handler Features
- ✅ RFC3164/RFC5424 compatible
- ✅ TCP/UDP support
- ✅ Local and remote syslog
- ✅ Priority mapping
- ✅ JSON structured messages

## Configuration Examples

### Complete OTLP Setup
```env
SERVICE_NAME=sale-service
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318
OTEL_EXPORTER_OTLP_HEADERS=
```

### Complete Loki Setup
```env
SERVICE_NAME=sale-service
LOKI_URL=http://localhost:3100
```

### Complete Syslog Setup
```env
SERVICE_NAME=sale-service
SYSLOG_ADDRESS=localhost:514
SYSLOG_NETWORK=udp
```

### Hybrid Setup (Multiple Options)
```env
# Primary: OTLP
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318

# Fallback: Direct Loki
LOKI_URL=http://localhost:3100

# Emergency: Syslog
SYSLOG_ADDRESS=localhost:514
```

## Testing

### Test OTLP Logs
```bash
# Start OTEL Collector + Loki
make observability-up

# Set OTLP endpoint
export OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318

# Run app
go run main.go

# Check Grafana: http://localhost:3000
```

### Test Direct Loki
```bash
# Start only Loki
docker run -d --name loki -p 3100:3100 grafana/loki:latest

# Set Loki URL
export LOKI_URL=http://localhost:3100

# Run app
go run main.go

# Query Loki directly
curl -G -s "http://localhost:3100/loki/api/v1/query" \
  --data-urlencode 'query={job="go-direct"}' | jq
```

### Test Syslog
```bash
# Start local syslog (if not running)
sudo systemctl start rsyslog

# Set syslog config
export SYSLOG_ADDRESS=localhost:514
export SYSLOG_NETWORK=udp

# Run app
go run main.go

# Check syslog
sudo tail -f /var/log/messages
```

## Monitoring and Troubleshooting

### Check Log Delivery
```bash
# OTLP: Check OTEL Collector logs
docker logs otel-collector

# Loki: Check Loki logs
docker logs loki

# Syslog: Check system logs
journalctl -u rsyslog -f
```

### Performance Considerations

1. **OTLP**: Batched, efficient, includes retries
2. **Loki**: Direct HTTP, fast, minimal overhead
3. **Syslog**: UDP is fast but unreliable, TCP is reliable but slower
4. **File**: Disk I/O dependent, requires log rotation

### Failover Behavior

The application tries logging methods in this order:
1. OTLP (if endpoint configured)
2. Direct Loki (if URL configured)
3. Syslog (if address configured)
4. File (if path configured)
5. Stdout (always available)

Each method also outputs to stdout as a fallback, so you'll always see logs in the console.

## Production Recommendations

### For Kubernetes/Cloud Native
```env
OTEL_EXPORTER_OTLP_ENDPOINT=otel-collector:4318
```

### For Traditional Infrastructure
```env
SYSLOG_ADDRESS=log-server:514
SYSLOG_NETWORK=tcp
```

### For Development
```env
LOKI_URL=http://localhost:3100
```

This setup eliminates file I/O completely while providing multiple robust logging options!