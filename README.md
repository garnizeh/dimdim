# dimdim

Personal project to manage domestic finances.

## Tech Stack

- main language: [go](https://go.dev)
- ui: [pico css](https://picocss.com)
- database: [sqlite](https://www.sqlite.org)

### Libraries

- migration: [goose](https://github.com/pressly/goose)
- sql: [sqlc](https://sqlc.dev)
- http: [echo](https://echo.labstack.com/)

### Make Rules

To install goose and sqlc:
```
make install
```

To create a new migration file in sql: `make migrate-create NAME=<migration-related-name>`
```
make migrate-create NAME=add-table-users
```
