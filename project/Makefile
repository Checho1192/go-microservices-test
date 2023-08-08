SHELL=cmd.exe
STATS_BINARY=statsMsApp

## up: starts all containers in the background without forcing build
up:
	@echo Starting Docker images...
	docker-compose up -d
	@echo Docker images started!

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_stats
	@echo Stopping docker images (if running...)
	docker-compose down
	@echo Building (when required) and starting docker images...
	docker-compose up --build -d
	@echo Docker images built and started!

## down: stop docker compose
down:
	@echo Stopping docker compose...
	docker-compose down
	@echo Done!

## build_stats: builds the stats binary as a linux executable
build_stats:
	@echo Building stats binary...
	chdir ..\stats-ms && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ${STATS_BINARY} ./cmd/api
	@echo Done!