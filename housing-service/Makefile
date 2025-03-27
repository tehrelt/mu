ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build: gen 
	go build -o ./bin/app.exe ./cmd/app

run: build
	./bin/app -local

PB=pkg/pb
PROTO=proto/house.proto
RATE_PROTO=proto/rate.proto
protobuf:
	cp -r ../proto .

	protoc  \
		--go_out=. \
		--go_opt=M$(PROTO)=$(PB)/housepb \
		--go-grpc_out=. \
		--go-grpc_opt=M$(PROTO)=$(PB)/housepb \
		$(PROTO)

	protoc  \
		--go_out=. \
		--go_opt=M$(RATE_PROTO)=$(PB)/ratespb \
		--go-grpc_out=. \
		--go-grpc_opt=M$(RATE_PROTO)=$(PB)/ratespb \
		$(RATE_PROTO)


gen: protobuf wire

wire: 
	wire ./internal/app

migrate.up:
	migrate -path ./migrations -database 'postgres://$(PG_USER):$(PG_PASS)@$(PG_HOST):$(PG_PORT)/$(PG_NAME)?sslmode=disable' up

migrate.down:
	migrate -path ./migrations -database 'postgres://$(PG_USER):$(PG_PASS)@$(PG_HOST):$(PG_PORT)/$(PG_NAME)?sslmode=disable' down

