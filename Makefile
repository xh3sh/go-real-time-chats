default: up

up:
	docker network create projects-network || true
	docker compose up -d --build

test:
	go test ./... -v

fmt:
	go fmt ./...