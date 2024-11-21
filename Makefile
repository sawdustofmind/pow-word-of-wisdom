dep:
	go mod tidy
	go mod vendor

lint:
	golangci-lint run

up:
	docker-compose up -d --build
