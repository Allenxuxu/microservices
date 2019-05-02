
.PHONY: build
build:
	go build -o ../../build/userSRV main.go

.PHONY: run
run:build
	../../build/userSRV  --registry=consul 

.PHONY: test
test:
	go test -v ./... -cover
