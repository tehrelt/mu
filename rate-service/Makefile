ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build: gen
	go build -o ./bin/app ./cmd/app

run: build
	./bin/app -env .env

PB=pkg/pb
PROTO=proto/rate.proto
protobuf:
	cp -r ../proto .

	protoc  \
		--go_out=. \
		--go_opt=M$(PROTO)=$(PB)/ratepb \
		--go-grpc_out=. \
		--go-grpc_opt=M$(PROTO)=$(PB)/ratepb \
		$(PROTO)


gen: protobuf wire

wire:
	wire ./internal/app

migrate.up:
	migrate -path ./migrations -database 'postgres://$(PG_USER):$(PG_PASS)@$(PG_HOST):$(PG_PORT)/$(PG_NAME)?sslmode=disable' up

migrate.down:
	migrate -path ./migrations -database 'postgres://$(PG_USER):$(PG_PASS)@$(PG_HOST):$(PG_PORT)/$(PG_NAME)?sslmode=disable' down
