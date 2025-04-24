ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build: gen 
	go build -o ./bin/app.exe ./cmd/app

run: build
	./bin/app -local

PB=pkg/pb
PROTO=proto/account.proto
USER_PROTO=proto/user.proto
HOUSE_PROTO=proto/house.proto
BILLING_PROTO=proto/billing.proto
protobuf:
	cp -r ../proto .

	protoc  \
		--go_out=. \
		--go_opt=M$(PROTO)=$(PB)/accountpb \
		--go-grpc_out=. \
		--go-grpc_opt=M$(PROTO)=$(PB)/accountpb \
		$(PROTO)

	protoc  \
		--go_out=. \
		--go_opt=M$(USER_PROTO)=$(PB)/userpb \
		--go-grpc_out=. \
		--go-grpc_opt=M$(USER_PROTO)=$(PB)/userpb \
		$(USER_PROTO)

	protoc  \
		--go_out=. \
		--go_opt=M$(HOUSE_PROTO)=$(PB)/housepb \
		--go-grpc_out=. \
		--go-grpc_opt=M$(HOUSE_PROTO)=$(PB)/housepb \
		$(HOUSE_PROTO)

	protoc  \
		--go_out=. \
		--go_opt=M$(BILLING_PROTO)=$(PB)/billingpb \
		--go-grpc_out=. \
		--go-grpc_opt=M$(BILLING_PROTO)=$(PB)/billingpb \
		$(BILLING_PROTO)

gen: protobuf wire

wire: 
	wire ./internal/app

migrate.up:
	migrate -path ./migrations -database 'postgres://$(PG_USER):$(PG_PASS)@$(PG_HOST):$(PG_PORT)/$(PG_NAME)?sslmode=disable' up

migrate.down:
	migrate -path ./migrations -database 'postgres://$(PG_USER):$(PG_PASS)@$(PG_HOST):$(PG_PORT)/$(PG_NAME)?sslmode=disable' down

