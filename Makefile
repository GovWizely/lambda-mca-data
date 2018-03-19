build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/mca mca/main.go
