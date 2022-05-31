GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)
INTERNAL_PROTO_FILES=$(shell find internal/conf -name *.proto)
API_PROTO_FILES=$(shell find api -name *.proto)
GRPC_CLIENT_PROTO_FILES=$(shell find internal/pkg/grpc_client -name *.proto)
HTTP_CLIENT_PROTO_FILES=$(shell find internal/pkg/http_client -name *.proto)

.PHONY: config
# generate config proto
config:
	protoc --proto_path=./internal \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)

.PHONY: api
# generate api proto
api:
	protoc --proto_path=./api \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./api \
 	       --go-grpc_out=paths=source_relative:./api \
	       $(API_PROTO_FILES)

.PHONY: grpc-client
# generate grpc-client proto
grpc-client:
	protoc --proto_path=./api \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./api \
 	       --go-grpc_out=paths=source_relative:./api \
	       $(GRPC_CLIENT_PROTO_FILES)

.PHONY: http-client
# generate http-client proto
http-client:
	protoc --proto_path=./api \
	       --proto_path=./third_party \
 	       --go_out=paths=source_relative:./api \
 	       --go-http_out=paths=source_relative:./api \
	       $(HTTP_CLIENT_PROTO_FILES)

.PHONY: generate
# generate wire
generate:
	cd cmd/at-kratos && wire

.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...

# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z0-9_-]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help