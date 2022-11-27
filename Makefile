build:
	GOOS=linux go build -o bin/bazaar cmd/bazaar/main.go
	go build -o bin/bazaar-cli cmd/cli/main.go
	docker build -t bazaar:1.0 .

build-dev:
	docker build -f Dockerfile.dev -t bazaar-dev:1.0 .