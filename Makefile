build:
	GOOS=linux go build -o bin/bazaar cmd/bazaar/main.go
	GOOS=linux go build -o bin/bazaar-cli cmd/bazaar-cli/main.go
	docker build -t bazaar:1.0 .
