# REST API for Simple Blog Project

[![Go Report Card](https://goreportcard.com/badge/github.com/markraiter/simple-blog)](https://goreportcard.com/report/github.com/markraiter/simple-blog)

### The App:

This is a simple-blog application with jwt authentication and CRUD operations for posts and comments. Only registered users can create posts and comments.

### To run the blog please proceed next:

1. Clone the repository.
2. Install the dependencies with `go mod download`
3. Create `.env` file and copy values from `.env_example`
4. Follow the instructions to install [Taskfile](https://taskfile.dev/ru-ru/installation/) utility
5. Run the app with `task run`

### Running the tests

1. Run the tests with `task test`
2. Also you can proceed with the [OpenAPI](https://swagger.io/) docs by link `localhost:8080/swagger`

### Built With

- [Go](https://golang.org/) - The programming language used.
- [net/http](https://pkg.go.dev/net/http) - Package used for HTTP client and server implementations.
- [REST](https://en.wikipedia.org/wiki/Representational_state_transfer) - Architectural style for the API.
- [Clean Architecture](https://8thlight.com/blog/uncle-bob/2012/08/13/the-clean-architecture.html) - Architectural pattern used.
- [Postgres](https://www.postgresql.org/) - Database used.
- [Golang-Migrate](https://github.com/golang-migrate/migrate) - Database migrations tool.
- [JWT](https://jwt.io/) - Used for authentication.