package config

import "time"

type PoolConfig struct {
	MaxOpenConns int
	MaxIdleConns int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

type DBConfig struct {
	Driver string
	Host   string
	Port   string
	User   string
	Password  string
	DBName string
	SSLMode string
	FilePath string
	Pool PoolConfig
}