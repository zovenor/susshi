build:
	go build -o ./bin/susshi ./cmd/susshi.go

run:
	./bin/susshi

run-test:
	go run ./cmd/susshi.go

.PHONY: run, build, run-test
