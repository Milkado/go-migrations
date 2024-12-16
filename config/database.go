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

var Config DBConfig

func Init() {
	Config.Driver = Env("DB_DRIVER")
	Config.Host = Env("DB_HOST")
	Config.Port = Env("DB_PORT")
	Config.User = Env("DB_USER")
	Config.Password = Env("DB_PASS")
	Config.DBName = Env("DB_NAME")
	Config.SSLMode = Env("DB_SSL_MODE")
	Config.FilePath = Env("DB_FILE_PATH")

	Config.Pool.MaxOpenConns = 1
	Config.Pool.MaxIdleConns = 1
	Config.Pool.ConnMaxLifetime = time.Second * 30
	Config.Pool.ConnMaxIdleTime = time.Second * 30
}