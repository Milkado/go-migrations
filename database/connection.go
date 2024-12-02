package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Milkado/go-migrations/config"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func buildDSN(config *config.DBConfig) (string, error) {
	var dsn string
	switch config.Driver {
	case "mysql":
		dsn = config.User + `:` +
			config.Password + `@tcp(` +
			config.Host + `:` +
			config.Port + `)/` +
			config.DBName + `?charset=utf8mb4&parseTime=True&loc=Local`
	case "postgres":
		dsn = `host=` + config.Host +
			` port=` + config.Port +
			` user=` + config.User +
			` password=` + config.Password +
			` dbname=` + config.DBName +
			` sslmode=` + config.SSLMode
	case "sqlite3":
	default:
		return "", fmt.Errorf("driver not supported")

	}

	return dsn, nil
}

func setConnectionPool(db *sql.DB, config *config.DBConfig) {
	pool := config.Pool

	//Set pool setts based on driver defaults if not set
	if pool.MaxOpenConns == 0 {
		switch config.Driver {
		case "mysql":
			pool.MaxOpenConns = 100
		case "postgres":
			pool.MaxOpenConns = 90
		case "sqlite3":
			pool.MaxOpenConns = 1
		}
	}

	if pool.MaxIdleConns == 0 {
		pool.MaxIdleConns = pool.MaxOpenConns / 4
	}

	if pool.ConnMaxLifetime == 0 {
		pool.ConnMaxLifetime = time.Hour
	}

	if pool.ConnMaxIdleTime == 0 {
		pool.ConnMaxIdleTime = time.Minute * 30
	}

	db.SetMaxOpenConns(pool.MaxOpenConns)
	db.SetMaxIdleConns(pool.MaxIdleConns)
	db.SetConnMaxLifetime(pool.ConnMaxLifetime)
	db.SetConnMaxIdleTime(pool.ConnMaxIdleTime)
}

func NewConnection(config *config.DBConfig) (*sql.DB, error) {
	dsn, err := buildDSN(config)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open(config.Driver, dsn)
	if err != nil {
		return nil, err
	}

	setConnectionPool(db, config)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return sql.Open(config.Driver, dsn)
}

func NewConnectionWithMonitoring(config *config.DBConfig) (*sql.DB, error) {
	db, err := NewConnection(config)
	if err != nil {
		return nil, err
	}

	//Start monitoring with 30s interval
	MonitorConnectionPool(db, 30*time.Second)

	//Initial health check
	if err := HealthCheck(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("health check failed, error: %v", err)
	}

	//Create migrations table
	if err := CreateMigrationsTable(db, config.Driver); err != nil {
		db.Close()
		return nil, fmt.Errorf("error creating migrations table, error: %v", err)
	}

	//Start periodic health checks
	StartHealthChecks(db, time.Minute)

	return db, nil
}
