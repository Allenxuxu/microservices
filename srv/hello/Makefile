
GOPATH:=$(shell go env GOPATH)


.PHONY: proto
proto:
	protoc --proto_path=${GOPATH}/src:. --micro_out=. --go_out=. proto/example/example.proto

.PHONY: build
build: proto
	go build -o ../../build/helloSRV main.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: run
run:build
	../../build/helloSRV  --registry=consul 

.PHONY: docker
docker:
	CGO_ENABLED=0 GOOS=linux go build -o hello-srv main.go
	docker build . -t xuxu123/hello-srv:latest
	rm hello-srv