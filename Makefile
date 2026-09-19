.PHONY: run migrate 

run:
	go run ./cmd/server

migrate:
	docker exec -i my_db psql -U postgres -d postgres < migrations/tables.sql
