# build:
# 	docker compose build simple-blog

# run:
# 	docker compose up simple-blog


.PHONY: databaseinit
databaseinit:
	docker run --name blog -e POSTGRES_PASSWORD='example' -p 5436:5432 -d --rm postgres

.PHONY: migrate_up
migrate_up:
	migrate -database "postgres://postgres:example@localhost:5432/blog?sslmode=disable" -path migrations/ up

.PHONY: migrate_down
migrate_down:
	migrate -database "postgres://postgres:example@localhost:5432/blog?sslmode=disable" -path migrations/ down

.PHONY: swag
swag:
	swag init -g cmd/main.go

.PHONY: run
run:
	go run cmd/main.go

.PHONY: test
test:
	go test -v -cover -race -timeout 15s ./...
	
.PHONY: database
database:
	docker exec -it postgres_db /bin/bash && psql -U postgres