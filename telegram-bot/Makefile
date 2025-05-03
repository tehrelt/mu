ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build: gen
	go build -o ./bin/app.exe ./cmd/app

run: build
	./bin/app -local

PB=pkg/pb
NOTIFICATION_PROTO=proto/notification.proto
protobuf:
	cp -r ../proto .

	protoc  \
		--go_out=. \
		--go_opt=M$(PROTO)=$(PB)/registerpb \
		--go-grpc_out=. \
		--go-grpc_opt=M$(PROTO)=$(PB)/registerpb \
		$(PROTO)





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
