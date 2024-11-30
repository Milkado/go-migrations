## go-migrations

Building a database migration package with Go
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
- Create new migration files
- Monitor connection pool
- Health checks
- Name validation to ensure pattern is followed

### TODO
- Add migration logic with type safety
- Add migrate command
- Add alter migration
- Add rollback migration
- Add seeding with type safety
- Add tests


### Final
- Make a full blown db standalone package for go modules
- Use in ```"go-api-blueprint/boilerplate"``` future project