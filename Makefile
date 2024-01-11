build:
	go build .

test:
	go test ./...

format:
	gofmt -s -w .

run:
	go run . $(CONFIG)
