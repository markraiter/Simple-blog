# build:
# 	docker compose build simple-blog

# run:
# 	docker compose up simple-blog


databaseinit:
	docker run --name blog -e POSTGRES_PASSWORD='example' -p 5436:5432 -d --rm postgres

migrate_up:
	migrate -database "postgres://postgres:example@localhost:5432/blog?sslmode=disable" -path migrations/ up

migrate_down:
	migrate -database "postgres://postgres:example@localhost:5432/blog?sslmode=disable" -path migrations/ down

run:
	go run cmd/main.go
	
database:
	docker exec -it postgres_db /bin/bash && psql -U postgres

swag:
	swag init -g cmd/main.go