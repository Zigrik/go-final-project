FROM golang:1.24 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server_app

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/server_app .
COPY web ./web

ENV TODO_PORT=7540
ENV TODO_DBFILE=scheduler.db
ENV TODO_PASSWORD=yandex

EXPOSE $TODO_PORT

CMD ["./server_app"] 