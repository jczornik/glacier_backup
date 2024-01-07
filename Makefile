build:
	go build .

test:
	go test ./...

run:
	go run . $(CONFIG)
