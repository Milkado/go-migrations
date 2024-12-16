## go-migrations

Building a database migration package with Go.

For study purposes only.

### Usage
```clone the repo
git clone https://github.com/Milkado/go-migrations.git
cd go-migrations
```

```build the binary
go build -o go-migrations
```

```run the binary
./go-migrations --c migration:create --name create-users-table
```

### Features
- Migration files as you like (raw sql or builder)
- Monitor connection pool
- Health checks
- Name validation to ensure pattern is followed
- Full fledged sql builder (to add more database support)


### Drawbacks
- Every new migration created needs a new build (from using go files)

### TODO
- [X] Environment variables support
- [X] Migration option using query builder
- [X] Change sqlgen to support multiple databases (postgres, mysql, sqlite)
- [X] Add Drop and Alter sqlgen
- [X] Migrate command
- [ ] Alter migration
- [X] Rollback migrations
- [ ] Seeding with type safety
- [X] Refactor alterations to fucntion like references
- [ ] Separate migrate command to a standalone build
- [ ] Add option to use SQL files
- [ ] Add tests
    - [X] Add tests for SQL generation
        - [X] MySQL
        - [X] Postgres
        - [ ] SQLite
    - [X] Add tests for generating files
    - [ ] Add tests for scanning files
    - [ ] Add tests for getting pending migrations
    - [ ] Add tests for migrate and rollback