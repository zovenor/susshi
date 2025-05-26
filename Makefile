install:
	go install ./cmd/susshi

vet:
	go vet ./...

.PHONY: install
