# go-todo-cleanarch

[![GoDoc](https://godoc.org/github.com/h4ckm03d/go-todo-cleanarch?status.svg)](https://godoc.org/github.com/h4ckm03d/go-todo-cleanarch)
[![Build Status](https://travis-ci.org/h4ckm03d/go-todo-cleanarch.svg?branch=main)](https://travis-ci.org/h4ckm03d/go-todo-cleanarch)
[![Go Report Card](https://goreportcard.com/badge/github.com/h4ckm03d/go-todo-cleanarch)](https://goreportcard.com/report/github.com/h4ckm03d/go-todo-cleanarch)
[![Maintainability](https://api.codeclimate.com/v1/badges/ce812254a495e287b45d/maintainability)](https://codeclimate.com/github/h4ckm03d/go-todo-cleanarch/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/ce812254a495e287b45d/test_coverage)](https://codeclimate.com/github/h4ckm03d/go-todo-cleanarch/test_coverage)

Go Todo Clean Arch Example Using Modular Project Layout for Product Microservice. It's suitable as starting point for a medium to larger project. It's originally fork from https://github.com/Fs02/go-todo-backend

This example uses [Chi](https://github.com/go-chi/chi) for http router and [REL](https://github.com/go-rel/rel) for database access.

Feature:

- Modular Project Structure.
- Full example including tests.
- Docker deployment.
- Compatible with [todobackend](https://www.todobackend.com/specs/index.html).

## Installation

### Prerequisite

1. Install [mockery](https://github.com/vektra/mockery#installation) for interface mock generation.

### Running

1. Prepare `.env`.
    ```
    cp .env.sample .env
    ```
2. Create a database and update `.env`.
2. Prepare database schema.
    ```
    rel migrate -adapter=github.com/go-rel/rel/adapter/sqlite3 -driver=github.com/mattn/go-sqlite3 -dsn="./rel_test.db?_foreign_keys=1&_loc=Local"
    ```
3. Build and Running
    ```
    make
    ```

## Project Structure

```
.
├── api
│   ├── handler
│   │   ├── todos.go
│   │   └── [other handler].go
│   └── middleware
│       └── [other middleware].go
├── bin
│   ├── api
│   └── [other executable]
├── cmd
│   ├── api
│   │   └── main.go
│   └── [other cmd]
│       └── main.go
├── db
│   ├── schema.sql
│   └── migrations
│       └── [migration file]
├── todos
│   ├── todo.go
│   ├── create.go
│   ├── update.go
│   ├── delete.go
│   ├── service.go
│   └── todostest
│       ├── todo.go
│       └── service.go
├── [other domain]
│   ├── [entity a].go
│   ├── [business logic].go
│   ├── [other domain]test
│   │   └── service.go
│   └── service.go
└── [other client]
    ├── [entity b].go
    ├── client.go
    └── [other client]test
        └── client.go
```

This project structure is based on a modular project structure, with loosely coupled dependencies between domain, Think of making libraries under a single repo that only exports certain functionality that used by other service and HTTP handler. One of domain that present in this example is todos.

Loosely coupled dependency between domain is enforced by avoiding the use of shared entity package, therefore any entity struct should be included inside it's own respective domain. This will prevent cyclic dependency between entities. This shouldn't be a problem in most cases, because if you encounter cyclic dependency, there's a huge chance that the entity should belongs to the same domain.

For example, consider three structs: user, transaction and transaction items. transaction and its transaction items might need cyclic dependency and items doesn't works standalone (items without transaction should not exists), thus it should be on the same domain.
In the other hand, user and transaction shouldn't require cyclic dependency, the transaction might have a user field in the struct, but a user shouldn't have a slice of the transaction field, therefore it should be on a separate domain.

### Domain vs Client

Domain and Client folder is very similar, the difference is client folder doesn't actually implement any business logic (service), but instead a client that calls any internal/external API to works with the domain entity.
