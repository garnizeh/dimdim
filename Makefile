install: install-migrate install-sqlc

.PHONY: install-migrate
install-migrate:
	go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: install-sqlc
install-sqlc:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

migrate-create:
	goose -dir storage/repo/migrations -s create $(NAME) sql

generate:
	sqlc generate

git-clean:
	git checkout main
	git remote update origin --prune
	git branch | grep -v "\smain\b" | xargs git branch -D

go:
	go mod tidy
	go mod vendor
