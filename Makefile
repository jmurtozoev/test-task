
proto-gen:
	protoc -I ./proto \
      --go_out ./proto --go_opt paths=source_relative \
      --go-grpc_out ./proto --go-grpc_opt paths=source_relative \
      --grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative \
      ./proto/*.proto

run-server:
	go run cmd/main.go

run-gateway:
	go run gateway/gateway.go

migrate:  ## make sure you installed migrate
	migrate -database ${POSTGRESQL_URL} -path db/migrations up
