build:
	@go build cmd/main.go

proto-v1: ### generate source files from proto
	protoc --go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		docs/proto/v1/*.proto
.PHONY: proto-v1

update:
	@git pull && make build

install prereq mac brew:
	@/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
	@brew install go
	@go get -u github.com/pressly/goose/v3/cmd/goose@latest
	@brew install memcached
	@brew install sqlite3
	@go get google.golang.org/grpc
	@brew install grpc
	@brew install protoc-gen-go
	@brew install protoc-gen-go-grpc

create migration:
	@go run github.com/pressly/goose/v3/cmd/goose@latest create create_links sql -dir migration