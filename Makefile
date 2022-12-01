build:
	GOOS=linux go build -o bin/bazaar cmd/bazaar/main.go

	GOOS=linux go build -o bin/bazaar-cli cmd/bazaar-cli/main.go
	GOOS=darwin go build -o bin/bazaar-cli-macos cmd/bazaar-cli/main.go
	GOOS=windows go build -o bin/bazaar-cli-windows cmd/bazaar-cli/main.go
