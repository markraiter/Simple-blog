build:
	docker-compose build

run:
	docker-compose up

migrate_up:
	migrate -database "postgres://postgres:qwerty@localhost:5432/db_postgres?sslmode=disable" -path ./ up

migrate_down:
	migrate -database "postgres://postgres:qwerty@localhost:5432/db_postgres?sslmode=disable" -path ./ down

swag:
	swag init -g cmd/main.go