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
- Migration files as you like
- Monitor connection pool
- Health checks
- Name validation to ensure pattern is followed

### TODO
- [ ] Environment variables support
- [ ] Migration option using query builder
- [X] Migrate command
- [ ] Alter migration
- [ ] Rollback migration (just need testing)
- [ ] Seeding with type safety
- [ ] Add tests


### Final
- Make a full blown db standalone package for go modules
- Use in ```"go-api-blueprint/boilerplate"``` future project