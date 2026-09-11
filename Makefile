.PHONY: run migrate 

run:
	go run ./cmd/server

migrate:
	docker exec -i live-attendance-backend-postgres-1 psql -U postgres -d postgres < migrations/01_auth.sql
