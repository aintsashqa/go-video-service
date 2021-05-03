.SILENT:
.PHONY:

deps:
	go mod tidy && go mod vendor

docker-down:
	docker-compose -f "deploy\docker-compose.dev.yaml" down

docker-up: deps docker-down
	docker-compose -f "deploy\docker-compose.dev.yaml" up -d --build
