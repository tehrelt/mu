ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build: gen
	go build -o ./bin/app ./cmd/app

run: build
	./bin/app -env .env

PB=pkg/pb
PROTO=proto/consumption.proto
BILLING_PROTO=proto/billing.proto
RATE_PROTO=proto/rate.proto
ACCOUNT_PROTO=proto/account.proto
protobuf:
	cp -r ../proto .

	protoc  \
		--go_out=. \
		--go_opt=M$(PROTO)=$(PB)/consumptionpb \
		--go-grpc_out=. \
		--go-grpc_opt=M$(PROTO)=$(PB)/consumptionpb \
		$(PROTO)

	protoc  \
		--go_out=. \
		--go_opt=M$(BILLING_PROTO)=$(PB)/billingpb \
		--go-grpc_out=. \
		--go-grpc_opt=M$(BILLING_PROTO)=$(PB)/billingpb \
		$(BILLING_PROTO)

	protoc  \
		--go_out=. \
		--go_opt=M$(RATE_PROTO)=$(PB)/ratepb \
		--go-grpc_out=. \
		--go-grpc_opt=M$(RATE_PROTO)=$(PB)/ratepb \
		$(RATE_PROTO)

	protoc  \
		--go_out=. \
		--go_opt=M$(ACCOUNT_PROTO)=$(PB)/accountpb \
		--go-grpc_out=. \
		--go-grpc_opt=M$(ACCOUNT_PROTO)=$(PB)/accountpb \
		$(ACCOUNT_PROTO)

gen: protobuf wire

wire:
	wire ./internal/app

migrate.up:
	migrate -path ./migrations -database 'postgres://$(PG_USER):$(PG_PASS)@$(PG_HOST):$(PG_PORT)/$(PG_NAME)?sslmode=disable' up

migrate.down:
	migrate -path ./migrations -database 'postgres://$(PG_USER):$(PG_PASS)@$(PG_HOST):$(PG_PORT)/$(PG_NAME)?sslmode=disable' down
