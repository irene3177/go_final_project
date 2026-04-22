FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o todo-scheduler .

FROM ubuntu:latest

WORKDIR /app
COPY --from=builder /app/todo-scheduler .
COPY --from=builder /app/web ./web 

RUN mkdir -p /data

EXPOSE 7540
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db  

CMD ["./todo-scheduler"]