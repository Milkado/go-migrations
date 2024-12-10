## go-migrations

Building a database migration package with Go.

Initial fase of the project, little to none implemented

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
    - Future mitigation: use a separate build to generate migration files

### TODO
- [ ] Environment variables support
- [X] Migration option using query builder
- [ ] Change sqlgen to support multiple databases
- [X] Add Drop and Alter sqlgen
- [X] Migrate command
- [ ] Alter migration
- [ ] Rollback migration (just need testing)
- [ ] Seeding with type safety
- [X] Refactor alterations to fucntion like references
- [ ] Separate migrate command to a standalone build
- [ ] Add option to use SQL files
- [ ] Add tests
    - [X] Add tests for SQL generation
    - [X] Add tests for generating files
    - [ ] Add tests for scanning files
    - [ ] Add tests for getting pending migrations
    - [ ] Add tests for migrate and rollback



### Final
- Make a full blown db standalone package for go modules or change into ```"go-api-blueprint/boilerplate"``` future project