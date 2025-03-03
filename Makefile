# ==============================================================================
# Install external tools

install: install-migrate install-sqlc

.PHONY: install-migrate
install-migrate:
	go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: install-sqlc
install-sqlc:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# ==============================================================================
# Database support

migrate-create:
	goose -dir storage/repo/migrations -s create $(NAME) sql

generate:
	sqlc generate

dev-clean:
	rm cmd/app/dimdim.db*
	
# ==============================================================================
# Git management

git-clean:
	git checkout main
	git remote update origin --prune
	git branch | grep -v "\smain\b" | xargs git branch -D

# ==============================================================================
# Building app

.PHONY: dependencies
dependencies:
	go mod tidy
	go mod vendor

build: dependencies
	go build -ldflags="-s -w" -o ./bin/app ./cmd/app/main.go

# ==============================================================================
# Metrics and Tracing

metrics:
	expvarmon -ports="localhost:3010" -vars="build,requests,goroutines,errors,panics,mem:memstats.HeapAlloc,mem:memstats.HeapSys,mem:memstats.Sys"

statsviz:
	open "Google Chrome" http://localhost:3010/debug/statsviz

# ==============================================================================
# SMTP support

smtp:
	MailHog