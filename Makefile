build:
	docker compose build simple-blog

run:
	docker compose up simple-blog

migrate_up:
	migrate -database "postgres://postgres:qwerty@localhost:5432/db_postgres?sslmode=disable" -path ./ up

migrate_down:
	migrate -database "postgres://postgres:qwerty@localhost:5432/db_postgres?sslmode=disable" -path ./ down

swag:
	swag init -g cmd/main.go