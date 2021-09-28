.PHONY: clean test lint

all: clean bin/budget bin/budget-view test lint

clean:
	rm -f $(foreach NAME, $(APP_CMD_NAMES), "bin/$(NAME)")

bin/%:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o ./bin/$(notdir $@) ./cmd/$(notdir $@)

test:
	go test ./...

lint:
	golangci-lint run