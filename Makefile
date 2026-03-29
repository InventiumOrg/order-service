postgres:
	podman run --network inventium --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -p 5432:5432 -d postgres:16-alpine
createdb:
	podman exec -it postgres createdb --username=root --owner=root order-service
dropdb:
	podman exec -it postgres dropdb --username=root order-service
migrateup:
	migrate -path ./models/migration -database "postgresql://root:secret@localhost:5432/order-service?sslmode=disable" -verbose up
migratedown:
	migrate -path ./models/migration -database "postgresql://root:secret@localhost:5432/order-service?sslmode=disable" -verbose down
sqlc:
	sqlc generate --no-remote
loaddata:
	PGPASSWORD=secret psql -h localhost -U root -d order-service -f data/sql/inventium.sql
runcontainer:
	podman run --network inventium --name order-service -p 9820:9820 -d -e DB_SOURCE="postgresql://root:secret@postgres:5432/order-service?sslmode=disable" -e CLERK_KEY="sk_test_XhHg2KNAIqm9I65JwOgQbLajZj6UqeeLTnpjx1p4oa" order-service:1.0.0
.PHONY: postgres createdb dropdb migrateup migratedown sqlc loaddata runcontainer testlogging testlogging-docker observability-up observability-down observability-logs run-with-logs test-otlp-logs test-loki-logs test-syslog-logs