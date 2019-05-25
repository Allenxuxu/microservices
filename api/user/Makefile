
.PHONY: build
build:
	go build -o ../../build/userAPI main.go

.PHONY: run
run:build
	../../build/userAPI  --registry=consul 

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	CGO_ENABLED=0 GOOS=linux go build -o user-api main.go
	docker build . -t xuxu123/user-api:latest
	rm user-api