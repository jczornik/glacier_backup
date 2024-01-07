build:
	go build .

mod-tidy:
	go mod tidy

test: mod-tidy
	go test ./...

run: mod-tidy
	go run . $(CONFIG)
