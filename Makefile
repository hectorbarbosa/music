# Create local database
.PHONY: createdb 
createdb:
	psql -h localhost -U postgres \
        -c "CREATE DATABASE music ENCODING='UTF8'";

# Drop local database
.PHONY: dropdb 
dropdb:
	psql -h localhost -U postgres \
        -c "DROP DATABASE IF EXISTS music";

.PHONY: build
build:
	go build -o bin/music -v ./cmd/music

.PHONY: run 
run:
	bin/music

.PHONY: buildapi
buildapi:
	go build -o bin/api -v ./cmd/api

.PHONY: runapi 
runapi:
	bin/api

.PHONY: fix
fix:
	migrate -path db/migrations/ -database "postgresql://postgres:password@localhost:5432/music?sslmode=disable" force 1

.PHONY: swag 
swag:
	swag init -d ./cmd/music,./internal/rest,./internal/app/models,internal/rest/models

migratedown:
	migrate -path db/migrations/ -database "postgresql://postgres:password@localhost:5432/music?sslmode=disable" down 


.DEFAULT_GOAL := build