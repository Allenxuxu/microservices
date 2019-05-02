
.PHONY: build
build:
	go build -o ../../build/userAPI main.go

.PHONY: run
run:build
	../../build/userAPI  --registry=consul 

.PHONY: test
test:
	go test -v ./... -cover
