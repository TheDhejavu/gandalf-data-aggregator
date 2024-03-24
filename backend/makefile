include $(if $(wildcard .env), .env)
export

MODULE = $(shell go list -m)
VERSION ?= $(shell echo "1.0.0")
UNIX_TIMESTAMP := $(shell date +%s)

.PHONY: db-start
db-start: ## start the database server
	docker run --name data_aggregator_postgres15 -p 5440:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=data_aggregator -d postgres:15-alpine

.PHONY: start
start:
	go run main.go

.PHONY: worker
worker:
	go run ./worker/worker.go
