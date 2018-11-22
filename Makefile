BASE_PATH			:= $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
BASE_DIR			:= $(notdir $(BASE_PATH))

ORG					?= chiswicked
SERVICE				?= $(BASE_DIR)
DOCKER_IMAGE		:= $(ORG)/$(SERVICE)
VERSION				?= $(shell cat $(BASE_PATH)/VERSION 2> /dev/null || echo 0.0.1)

BUILD_PATH			:= $(BASE_PATH)/build
PROTOBUF_OUT_DIR	:= $(BASE_PATH)/protobuf
SWAGGER_OUT_DIR		:= $(BASE_PATH)/swagger

GRPC_GW_REPO		:= github.com/grpc-ecosystem/grpc-gateway
GRPC_GW_PROTO_PATH	:= $(GOPATH)/src/$(GRPC_GW_REPO)/third_party/googleapis
SERVICE_PROTO_PATH	:= $(GOPATH)/src/github.com/$(ORG)/go-grpc-crud-protobuf-model
PROTOC_IMPORTS		:= -I/usr/include -I$(GRPC_GW_PROTO_PATH) -I$(SERVICE_PROTO_PATH)

# Local commands

.PHONY: all clean install build test cover cover-clean run docker docker-build

all: clean install test build

clean: cover-clean
	@echo [clean] removing binary and other object files
	@go clean
	@rm -f $(BUILD_PATH)/$(SERVICE)

install:
	@echo [install] installing dependencies
	@go get -v -t -d ./...

build: clean
	@echo [build] building binary
	@go build -o $(BUILD_PATH)/$(SERVICE) -a .

test:
	@echo [test] running unit tests
	@go test -v -cover ./...

cover: cover-clean
	@echo [cover] generating test coverage report
	@go test -coverprofile cover.out ./...
	@go tool cover -html=cover.out -o cover.html
	@open cover.html

cover-clean:
	@echo [cover-clean] removing cover.out cover.html
	@rm -f cover.out cover.html

run:
	@echo [run] executing binary
	@$(BUILD_PATH)/$(SERVICE)

docker: docker-build

docker-build:
	@echo [docker-build] building docker image
	@docker build -t $(DOCKER_IMAGE):local . \
		--build-arg ORG=$(ORG) \
		--build-arg SERVICE=$(SERVICE) \
		--build-arg GITHUB_TOKEN=$(GITHUB_TOKEN)

# Protobuf commands

.PHONY: protoc protoc-clean protoc-install protoc-build

protoc: protoc-clean protoc-install protoc-build

protoc-clean:
	@echo "[protoc-clean] removing $(PROTOBUF_OUT_DIR)/*.pb.go"
	@rm -rf $(PROTOBUF_OUT_DIR)/*.pb.go
	@echo "[protoc-clean] removing $(PROTOBUF_OUT_DIR)/*.pb.gw.go"
	@rm -rf $(PROTOBUF_OUT_DIR)/*.pb.gw.go
	@echo "[protoc-clean] removing $(PROTOBUF_OUT_DIR)/*.swagger.json"
	@rm -rf $(SWAGGER_OUT_DIR)/*.swagger.json

protoc-install:
	@echo [protoc-install] installing dependencies
	@[ -x $(GOPATH)/bin/protoc-gen-grpc-gateway ] || go get -u $(GRPC_GW_REPO)/protoc-gen-grpc-gateway
	@[ -x $(GOPATH)/bin/protoc-gen-swagger ] || go get -u $(GRPC_GW_REPO)/protoc-gen-swagger
	@[ -x $(GOPATH)/bin/protoc-gen-go ] || go get -u github.com/golang/protobuf/protoc-gen-go

protoc-build:
	@mkdir -p $(PROTOBUF_OUT_DIR)
	@mkdir -p $(SWAGGER_OUT_DIR)

	@echo [protoc-build] building protobufs
	@protoc $(PROTOC_IMPORTS) \
		--go_out=plugins=grpc:$(PROTOBUF_OUT_DIR) \
		$(SERVICE_PROTO_PATH)/*.proto

	@echo [protoc-build] building grpc gateway
	@protoc $(PROTOC_IMPORTS) \
		--grpc-gateway_out=logtostderr=true:$(PROTOBUF_OUT_DIR) \
		$(SERVICE_PROTO_PATH)/*.proto

	@echo [protoc-build] generating swagger definitions
	@protoc $(PROTOC_IMPORTS) \
		--swagger_out=logtostderr=true:$(SWAGGER_OUT_DIR) \
		$(SERVICE_PROTO_PATH)/*.proto 

# CI commands

.PHONY: ci-docker-login ci-docker-build ci-docker-push ci-docker-logout

ci-docker-login:
	@echo "[ci-docker-login] logging in to docker hub"
	@echo $(DOCKER_PASSWORD) | docker login -u $(DOCKER_USERNAME) --password-stdin

ci-docker-build:
	@echo "[ci-docker-build] building docker image $(DOCKER_IMAGE):$(VERSION)"
	@docker build -t $(DOCKER_IMAGE):$(VERSION) . \
		--build-arg ORG=$(ORG) \
		--build-arg SERVICE=$(SERVICE) \
		--build-arg GITHUB_TOKEN=$(GITHUB_TOKEN)

ci-docker-push:
	@echo "[ci-docker-push] pushing docker image $(DOCKER_IMAGE):$(VERSION) to repository"
	@docker tag $(DOCKER_IMAGE):$(VERSION) $(DOCKER_IMAGE):latest
	@docker push $(DOCKER_IMAGE)

ci-docker-logout:
	@echo "[ci-docker-logout] logging out of docker hub"
	@docker logout
