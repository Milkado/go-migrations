package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"log"
)

type PoolStats struct {
	OpenConnections int
	InUse int
	Idle int
	WaitCount int
	WaitDuration time.Duration
	MaxIdleClosed int
	MaxLifetimeClosed int
}

func MonitorConnectionPool(db *sql.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			stats := db.Stats()
			log.Printf(
				"\nConnection Pool Stats:\n"+
                "Open Connections: %d\n"+
                "In Use: %d\n"+
                "Idle: %d\n"+
                "Wait Count: %d\n"+
                "Wait Duration: %v\n"+
                "Max Idle Closed: %d\n"+
                "Max Lifetime Closed: %d\n",
				stats.OpenConnections,
				stats.InUse,
				stats.Idle,
				stats.WaitCount,
				stats.WaitDuration,
				stats.MaxIdleClosed,
				stats.MaxLifetimeClosed,
			)
		}
	}()
}

func HealthCheck(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//Check if db is reachable
	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("database is not reachable, error: %v", err)
	}

	//Check if query is successful
	_, err := db.ExecContext(ctx, "SELECT 1")
	if err != nil {
		return fmt.Errorf("query failed, error: %v", err)
	}

	return nil
}

//TODO: Check if this should stay
func StartHealthChecks(db *sql.DB, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			if err := HealthCheck(db); err != nil {
				log.Printf("Health check failed, error: %v", err)
			}
	}
	}()
}
