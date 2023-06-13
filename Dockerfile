FROM golang

RUN go version
WORKDIR /app
COPY ./ ./

# build go app
RUN go mod download
RUN go build -o main cmd/main.go
RUN chmod +x main
EXPOSE 8080

CMD ["/app/main"]