GO_VERSION_SHORT:=$(shell echo `go version` | sed -E 's/.* go(.*) .*/\1/g')
ifneq ("1.16","$(shell printf "$(GO_VERSION_SHORT)\n1.16" | sort -V | head -1)")
$(error NEED GO VERSION >= 1.16. Found: $(GO_VERSION_SHORT))
endif

export GO111MODULE=on
export GOPROXY=https://proxy.golang.org|direct

PGV_VERSION:="v0.6.1"
GOOGLEAPIS_VERSION="master"
DB_STRING:= "host=localhost port=5432 user=ova_exam_api_user password=ova_exam_api_password dbname=ova_exam_api sslmode=disable"

all: generate build

.PHONY: vendor-proto
vendor-proto:
	$(eval THIRD_PARTY:=$(CURDIR)/third_party)
	@[ -f $(THIRD_PARTY)/validate/validate.proto ] || (mkdir -p $(THIRD_PARTY)/validate/ && curl -sSL0 https://raw.githubusercontent.com/envoyproxy/protoc-gen-validate/$(PGV_VERSION)/validate/validate.proto -o $(THIRD_PARTY)/validate/validate.proto)
	@[ -f $(THIRD_PARTY)/google/api/http.proto ] || (mkdir -p $(THIRD_PARTY)/google/api/ && curl -sSL0 https://raw.githubusercontent.com/googleapis/googleapis/$(GOOGLEAPIS_VERSION)/google/api/http.proto -o $(THIRD_PARTY)/google/api/http.proto)
	@[ -f $(THIRD_PARTY)/google/api/annotations.proto ] || (mkdir -p $(THIRD_PARTY)/google/api/ && curl -sSL0 https://raw.githubusercontent.com/googleapis/googleapis/$(GOOGLEAPIS_VERSION)/google/api/annotations.proto -o $(THIRD_PARTY)/google/api/annotations.proto)

GOBIN?=$(GOPATH)/bin

.PHONY: build
build: deps
	go build -o $(CURDIR)/bin/main $(CURDIR)/cmd/main.go

bin-deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.5.0
	go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.5.0
	go install github.com/envoyproxy/protoc-gen-validate@$(PGV_VERSION)

run:
	go run cmd/main.go

test:
	go test ova-exam-api/internal/app
	go test ova-exam-api/internal/flusher

generate-mocks:
	go generate ./...

generate:
	@protoc -I vendor.protogen --go_out=pkg \
 		--go_opt=paths=import --go-grpc_out=pkg \
		--go-grpc_opt=paths=import --proto_path=api/ova-exam-api \
		api/ova-exam-api/users.proto

LOCAL_BIN:=$(CURDIR)/bin

.PHONY: deps
deps: .install-go-deps

.PHONY: .install-go-deps
.install-go-deps:
	ls go.mod || go mod init github.com/ozonva/ova-exam-api
	GOBIN=$(LOCAL_BIN) go get -u github.com/DATA-DOG/go-sqlmock
	GOBIN=$(LOCAL_BIN) go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
	GOBIN=$(LOCAL_BIN) go get -u github.com/golang/protobuf/proto
	GOBIN=$(LOCAL_BIN) go get -u github.com/golang/protobuf/protoc-gen-go
	GOBIN=$(LOCAL_BIN) go get -u github.com/opentracing/opentracing-go
	GOBIN=$(LOCAL_BIN) go get -u github.com/jmoiron/sqlx
	GOBIN=$(LOCAL_BIN) go get -u github.com/jackc/pgx/stdlib
	GOBIN=$(LOCAL_BIN) go get -u github.com/Masterminds/squirrel
	GOBIN=$(LOCAL_BIN) go get -u github.com/pressly/goose/v3/cmd/goose
	GOBIN=$(LOCAL_BIN) go get -u github.com/prometheus/client_golang/prometheus/promhttp
	GOBIN=$(LOCAL_BIN) go get -u github.com/slok/go-http-metrics/metrics/prometheus
 	GOBIN=$(LOCAL_BIN) go get -u github.com/slok/go-http-metrics/middleware
	GOBIN=$(LOCAL_BIN) go get -u github.com/rs/zerolog/log
	GOBIN=$(LOCAL_BIN) go get -u github.com/segmentio/kafka-go
	GOBIN=$(LOCAL_BIN) go get -u google.golang.org/grpc
	GOBIN=$(LOCAL_BIN) go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
	GOBIN=$(LOCAL_BIN) go get -u google.golang.org/protobuf/types/known/emptypb
	GOBIN=$(LOCAL_BIN) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger

.PHONY: goose-up
goose-up: .goose-up

.PHONY: .goose-up
.goose-up:
	GOBIN=$(LOCAL_BIN) bin/goose -dir $(CURDIR)/migrations/ postgres $(DB_STRING) up

.PHONY: goose-status
goose-status: .goose-status

.PHONY: .goose-status
.goose-status:
	$(LOCAL_BIN)/goose -dir $(CURDIR)/migrations/ postgres $(DB_STRING) status