.PHONY: build test lint integration-test e2e-test

build:
	go build ./...

test:
	./scripts/test.sh

integration-test:
	docker-compose -f tests/integration/docker-compose.yml up -d
	go test ./tests/integration/...

e2e-test:
	go test ./tests/e2e/...

lint:
	./scripts/lint.sh 