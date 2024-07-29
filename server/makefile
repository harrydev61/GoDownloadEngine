VERSION := $(shell cat VERSION)
PROJECT_NAME := downloader

.PHONY: start
start:
	go run main.go standalone-server

.PHONY: gen-doc
gen-doc:
	 swag init --output docs

.PHONY: generate
generate:
	cd proto && protoc --go_out=pb --go_opt=paths=source_relative \
                       --go-grpc_out=pb --go-grpc_opt=require_unimplemented_servers=false \
                       --go-grpc_opt=paths=source_relative *.proto

.PHONY: build-web
build-web:
	cd web && npm run build && cd ..

.PHONY: build
build:
	go build \
		-ldflags "-X main.version=$(VERSION)" \
		-o build/$(PROJECT_NAME) main.go

.PHONY: clean
clean:
	rm -rf build/

.PHONY: lint
lint:
	golangci-lint run ./... 