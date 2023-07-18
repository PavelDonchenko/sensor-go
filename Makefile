.PHONY: lint run swag createTestDB dropDB createDBtest int_test test

lint:
	golangci-lint run ./...

run:
	go run cmd/sensor/main.go

compose_up:
	docker-compose -f docker-compose.yml up --build

swag:
	swag init -g cmd/sensor/main.go

createTestDB:
	docker exec -it db createdb --username=root --owner=root test_sensor

dropDB:
	docker exec -it db dropdb --username=root sensor_db

createDB:
	docker exec -it db createdb --username=root --owner=root sensor_db

int_test:
	docker-compose -f docker-compose.test.yml up --build -d
	go test -v ./...
	docker-compose -f docker-compose.test.yml down --volumes
