# Business Metrics for Sale Service

## 🎯 Business-Specific Metrics Added

### Sale Unit Operations
- **`sale_unit_operations_total`** - Counter tracking all sale unit operations
  - Labels: `operation` (get, list, create), `sale_unit_id`, `sale_unit_name`
  
- **`sale_units_created_total`** - Counter for new sale units created
  - Labels: `sale_unit_name`, `sale_unit_id`
  
- **`sale_unit_retrievals_total`** - Counter for individual sale unit retrievals
  - Labels: `sale_unit_name`, `status`
  
- **`sale_unit_list_requests_total`** - Counter for list operations
  - Labels: `limit`, `offset`

### Database Operations
- **`database_operation_duration_seconds`** - Histogram of DB operation times
  - Labels: `operation`, `table`
  
- **`database_operation_errors_total`** - Counter of DB errors
  - Labels: `operation`, `error_type`

### Authentication & Security
- **`authentication_attempts_total`** - Counter of auth attempts
  - Labels: `operation`, `status`

### Business Logic
- **`active_sale_units_count`** - Gauge of active sale units
  - Labels: `operation`

## 📊 Example Grafana Queries

### Business KPIs
```promql
# Sale unit creation rate
rate(sale_units_created_total[5m])

# Most popular sale unit operations
sum by (operation) (rate(sale_unit_operations_total[5m]))

# Database operation performance
histogram_quantile(0.95, database_operation_duration_seconds)

# Error rates by operation
rate(database_operation_errors_total[5m]) / rate(sale_unit_operations_total[5m])

# Authentication success rate
(rate(authentication_attempts_total[5m]) - rate(authentication_attempts_total{status="failed"}[5m])) / rate(authentication_attempts_total[5m])
```

### Operational Metrics
```promql
# Active sale units trend
active_sale_units_count

# Database operation errors by type
sum by (error_type) (rate(database_operation_errors_total[5m]))

# Slowest database operations
topk(5, histogram_quantile(0.95, database_operation_duration_seconds) by (operation))
```

## 🚀 Test Your Metrics

1. **Start your service**
2. **Make some API calls:**
   ```bash
   # Health check
   curl http://localhost:15350/healthz
   
   # List sale units (requires auth)
   curl -H "Authorization: Bearer YOUR_TOKEN" http://localhost:15350/v1/sale/list
   
   # Create a sale unit (requires auth)
   curl -X POST -H "Authorization: Bearer YOUR_TOKEN" \
        -d "Name=TestUnit" http://localhost:15350/v1/sale/create
   ```

3. **Check Grafana Cloud** - Your metrics should appear within 30 seconds

## 🔧 Next Steps

You can extend these metrics by:
- Adding more business operations (update, delete)
- Tracking user-specific metrics
- Adding inventory-related metrics
- Monitoring performance SLAs