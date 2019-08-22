.PHONY: build clean ecs_deploy install integration lint migrations run_local run_local_api run_local_consumer run_local_dependencies stop_local_dependencies stop_local test

clean:
	rm -rf ./.bin 2>/dev/null || true
	rm ./goldmine 2>/dev/null || true
	go fix ./...
	go clean -i

build: clean
	go fmt ./...
	go build -v -o ./.bin/goldmine_api ./cmd/api
	go build -v -o ./.bin/goldmine_consumer ./cmd/consumer
	go build -v -o ./.bin/goldmine_migrate ./cmd/migrate

ecs_deploy:
	./scripts/ecs_deploy.sh

install: clean
	go install ./...

lint:
	./scripts/lint.sh

migrations: build run_local_dependencies
	./scripts/run_local_migrations.sh

run_local: build run_local_dependencies
	./scripts/run_local.sh

run_local_api: build run_local_dependencies
	./scripts/run_local_api.sh

run_local_consumer: build run_local_dependencies
	./scripts/run_local_consumer.sh

run_local_dependencies:
	./scripts/run_local_dependencies.sh

stop_local_dependencies:
	./scripts/stop_local_dependencies.sh

stop_local:
	./scripts/stop_local.sh

test: build
	NATS_SERVER_PORT=4223 NATS_STREAMING_SERVER_PORT=4224 ./scripts/run_local_dependencies.sh
	NATS_SERVER_PORT=4223 NATS_STREAMING_SERVER_PORT=4224 ./scripts/run_unit_tests.sh

integration: build
	NATS_SERVER_PORT=4223 NATS_STREAMING_SERVER_PORT=4224 ./scripts/run_local_dependencies.sh
	NATS_SERVER_PORT=4223 NATS_STREAMING_SERVER_PORT=4224 ./scripts/run_integration_tests.sh
