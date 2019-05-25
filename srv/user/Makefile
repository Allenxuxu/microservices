
.PHONY: build
build:
	go build -o ../../build/userSRV main.go

.PHONY: run
run:build
	../../build/userSRV  --registry=consul 

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	CGO_ENABLED=0 GOOS=linux go build -o user-srv main.go
	docker build . -t xuxu123/user-srv:latest
	rm user-srv