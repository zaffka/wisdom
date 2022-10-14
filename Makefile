.PHONY: all mod test gen lint build server client

all: mod gen test build

mod:
	go mod download

test:
	go test -v -count=1 --race ./...

gen:
	go generate -v ./...

lint:
	golangci-lint run ./...

build:
	docker build -t wisdom:latest .

server:
	docker run --network host wisdom server

client:
	docker run --network host wisdom client