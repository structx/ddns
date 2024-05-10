
build image:
	docker build -t trevatk/daisy:v0.0.1 .

deps:
	go mod tidy
	go mod vendor

clean:
	go clean -modcache

rpc:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/ddns/v1/ddns_service.proto