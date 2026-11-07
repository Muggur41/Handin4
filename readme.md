# Handin4

## How to Run

1. Generate the protobuf files (if not already done):
```bash
protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative
--go-grpc_opt=paths=source_relative proto/mutex.proto

2. Install dependencies:
go mod tidy

3. Start 3 nodes in separate terminals:

Terminal 1:
go run main.go --id=1 --port=5001 --peers=localhost:5002,localhost:5003

Terminal 2:
go run main.go --id=2 --port=5002 --peers=localhost:5001,localhost:5003

Terminal 3:
go run main.go --id=3 --port=5003 --peers=localhost:5001,localhost:5002

Start all nodes within a few seconds of each other.