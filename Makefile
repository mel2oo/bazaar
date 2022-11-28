build:
	GOOS=linux go build -o bin/bazaar cmd/bazaar/main.go
	GOOS=linux go build -o bin/bazaar-cli cmd/bazaar-cli/main.go
