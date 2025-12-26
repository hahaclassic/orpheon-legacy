ifneq (,$(wildcard ./.env))
    include .env
    export
endif

test:
	go test ./...

unit:
	go test -v ./internal/domain/services/...

unit-random:
	go test -shuffle=on -v ./internal/domain/services/...

report-unit-gen:
	go clean -testcache
	go test -json ./internal/domain/... ./internal/repository/... | golurectl -l -e -s -a -o report/allure-results
	allure generate ./report/allure-results -o ./report/allure-report --clean          
	allure open ./report/allure-report

report-all-gen:
	./scripts/allure-report.sh

# Old approach
# go test -covermode=atomic -coverprofile=./coverage/coverage.out ./internal/... >> /dev/null
# go tool cover -func=./coverage/coverage.out | grep "total:"
cov:
	mkdir -p ./coverage
	./scripts/cov.sh

cov-v: 
	mkdir -p ./coverage
	./scripts/cov.sh -v

integration:
	go test -v ./tests/integration

e2e:
	go test -v ./tests/e2e

clear:
	rm -rf ./coverage

mock:
	mockery --all --recursive --output=./mocks --with-expecter

lint: 
	golangci-lint run
	
migrate-up:
	goose -dir db/migrations postgres "host=${POSTGRES_HOST} port=${POSTGRES_PORT} user=${POSTGRES_USER} \
		password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=${POSTGRES_SSL_MODE}" up

migrate-down:
	goose -dir db/migrations postgres "host=${POSTGRES_HOST} port=${POSTGRES_PORT} user=${POSTGRES_USER} \
		password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} sslmode=${POSTGRES_SSL_MODE}" down

containers-up:
	docker compose -f docker-compose.dev.backend.yml up
	
build-image:
	docker build -t orpheon-e2e:latest -f ./Dockerfile .