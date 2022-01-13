.PHONY: clean test lint

all: clean bin/budget bin/budgetview test lint check-arch

clean:
	rm -rf bin/*

bin/%:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o ./bin/$(notdir $@) ./cmd/$(notdir $@)

test:
	go test ./...

lint:
	golangci-lint run

check-arch:
	go-cleanarch -ignore-package=github.com/klwxsrx/budget-tracker/pkg/common