
build image:
	docker build -t localhost:5000/structx/ddns:latest .

deps:
	go mod tidy
	go mod vendor

clean:
	go clean -modcache

rpc:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	proto/ddns/ddns_service.proto

lint:
	golangci-lint run
