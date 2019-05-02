
.PHONY: build
build:
	go build -o ../build/microGW main.go

.PHONY: run
run:build
	../build/microGW  --registry=consul api -handler=http

.PHONY: test
test:
	go test -v ./... -cover
