# build server
# build client

## Server
run-server:
	go run ./server/server.go

test-server:
	go test ./server -v

## client
run-client:
	go run ./client/server.go

test-client:
	go test ./client -v

## Routes
test-routers:
	go test ./routers -v

## General
coverage:
	go test --cover ./...

compile: clean
	go build -o bin/server/server ./server/server.go
	go build -o bin/client/client ./client/server.go
	cd client && npm run build:bin

clean:
	cd client && npm run clean

## build docker
## publish docker
