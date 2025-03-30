ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build: gen 
	go build -o ./bin/app.exe ./cmd/app

run: build
	./bin/app -local

PB=pkg/pb
PROTO=proto/register.proto
USER_PROTO=proto/user.proto
AUTH_PROTO=proto/auth.proto
protobuf:
	cp -r ../proto .

	protoc  \
		--go_out=. \
		--go_opt=M$(PROTO)=$(PB)/registerpb \
		--go-grpc_out=. \
		--go-grpc_opt=M$(PROTO)=$(PB)/registerpb \
		$(PROTO)


	protoc  \
		--go_out=. \
		--go_opt=M$(USER_PROTO)=$(PB)/userpb \
		--go-grpc_out=. \
		--go-grpc_opt=M$(USER_PROTO)=$(PB)/userpb \
		$(USER_PROTO)

	protoc  \
		--go_out=. \
		--go_opt=M$(AUTH_PROTO)=$(PB)/authpb \
		--go-grpc_out=. \
		--go-grpc_opt=M$(AUTH_PROTO)=$(PB)/authpb \
		$(AUTH_PROTO)



gen: protobuf wire

clean:
	rm -rf ./pkg/pb
	rm -rf ./internal/app/wire_gen.go
	

wire: 
	wire ./internal/app

migrate.up:
	migrate -path ./migrations -database 'postgres://$(PG_USER):$(PG_PASS)@$(PG_HOST):$(PG_PORT)/$(PG_NAME)?sslmode=disable' up

migrate.down:
	migrate -path ./migrations -database 'postgres://$(PG_USER):$(PG_PASS)@$(PG_HOST):$(PG_PORT)/$(PG_NAME)?sslmode=disable' down

