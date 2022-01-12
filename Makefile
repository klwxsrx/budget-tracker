.PHONY: clean test lint

all: clean bin/budget bin/budgetview test lint

clean:
	rm -rf bin/*

bin/%:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o ./bin/$(notdir $@) ./cmd/$(notdir $@)

test:
	go test ./...

lint:
	golangci-lint run