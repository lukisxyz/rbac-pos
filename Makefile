include .env
APP_NAME = main
BUILD_DIR = $(PWD)/.tmp
PSQL=postgres://${DB_USER}:${DB_SECRET}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSL}

clean:
	rm -rf ./build

critic:
	gocritic check -enableAll ./...

security:
	.tmp/gosec ./...

lint:
	.tmp/golangci-lint run ./...

test: clean critic security lint
	go test -v -timeout 30s -coverprofile=cover.out -cover ./...
	go tool cover -func=cover.out

build: clean critic security lint
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(APP_NAME) *.go

run: build
	.tmp/air

setup.air:
	curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b .tmp/

setup.gosec:
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b .tmp/ v2.16.0

setup.migrate:
	curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz -C .tmp/

setup.gocritic:
	go install -v github.com/go-critic/go-critic/cmd/gocritic@latest

setup.lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b .tmp/ v1.54.1

setup: setup.air setup.gosec setup.gocritic setup.lint setup.migrate

swagger:
	swag init -g *.go --output docs/

test.coverage: 
	go test ./... -coverprofile docs/coverage.out

test.run:
	go test ./...

test: test.run test.coverage

cmgr:
	.tmp/migrate create -ext sql -dir db/migrations -seq ${name}

migup:
	.tmp/migrate -path db/migrations -database "${PSQL}" -verbose up

migdown:
	.tmp/migrate -path db/migrations -database "${PSQL}" -verbose down

seed: 
	./seed.sh ${token}