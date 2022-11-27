build:
	go build -o bin/bazaar cmd/bazaar/main.go
	docker build -t bazaar:1.0 .