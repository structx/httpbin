
lint:
	golangci-lint run ./...

build:
	docker build -t trevatk/httpbin:latest .