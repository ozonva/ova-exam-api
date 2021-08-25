build:
	go build -o bin/main cmd/main.go

run:
	go run cmd/main.go

test:
	go test ova-exam-api/internal/flusher

generate-mocks:
	go generate ./...
