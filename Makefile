format:
	goimports -w .

test: test-unit test-integration test-e2e

test-unit:
	go test -count=1 ./...

test-integration:
	(TEST_DB_FILE=integration_test.db \
	go test -count=1 ./... && rm -f pkg/integration_test.db) || (rm -f pkg/integration_test.db)


test-e2e:
	(TEST_E2E_DB_FILE=e2e_test.db \
	go test -count=1 e2e/e2e_test.go && rm -f e2e/e2e_test.db) || (rm -f e2e/e2e_test.db)

mocks:
	go generate ./...
	goimports -w .

build:
	go build -o ./bin/telegram-bot/telegram-bot ./cmd/telegram-bot/main.go

build-rpi:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o ./bin/telegram-bot/telegram-bot-rpi ./cmd/telegram-bot/main.go
