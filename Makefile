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

install goose:
	@go get -u github.com/pressly/goose/v3/cmd/goose@latest

create migration:
	@go run github.com/pressly/goose/v3/cmd/goose@latest create create_links sql -dir migration