run:
	go run ./cmd/server

test:
	go test ./...

fmt:
	go fmt ./...

lint:
	go vet ./...

build:
	go build -o smarthome ./cmd/server

clean:
	rm -f smarthome