
.PHONY: build
build:
	go build -o ../build/microGW main.go

.PHONY: run
run:build
	../build/microGW  --registry=consul api -handler=http

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	CGO_ENABLED=0 GOOS=linux go build -o micro main.go
	docker build . -t xuxu123/micro:latest
	rm micro