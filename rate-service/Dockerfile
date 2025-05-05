FROM golang:1.24.0-alpine3.20 AS builder

RUN apk update --no-cache
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.12.1/wait /app/wait
RUN chmod +x /app/wait
WORKDIR /app
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download && go mod verify

COPY . /app
RUN go build -mod=readonly -o migrator cmd/migrator/main.go &&\
    go build -mod=readonly -o app cmd/app/main.go

FROM alpine
RUN apk update --no-cache
WORKDIR /app
COPY --from=builder /app /app
COPY --from=builder /app/wait /app/wait

CMD ./wait && ./migrator && ./app
