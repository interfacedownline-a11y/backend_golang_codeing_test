APP_NAME = backend_golang_codeing_test
MAIN_FILE = cmd/main.go
PROTO_DIR = proto
PROTO_OUT = backend_golang_codeing_test/proto/userpb
PORT = 8080

.PHONY: all tidy generate test build run migrate docker clean

# Run all key tasks
all: tidy generate test build run

# Clean Go modules
tidy:
	go mod tidy

# Generate gRPC code from .proto
generate:
	protoc \
		--go_out=$(PROTO_OUT) \
		--go-grpc_out=$(PROTO_OUT) \
		--proto_path=$(PROTO_DIR) \
		$(PROTO_DIR)/*.proto

# Run all tests with verbose output
test:
	go test ./... -v

# Build the binary to ./bin/
build:
	go build -o bin/$(APP_NAME) $(MAIN_FILE)

# Run the application
run:
	go run $(MAIN_FILE)

# Run migration and seeding
migrate:
	go run migrations/migration.go && go run migrations/seeder.go

# Build Docker image
docker:
	docker build -t $(APP_NAME):latest .

# Clean build artifacts
clean:
	rm -rf bin/
