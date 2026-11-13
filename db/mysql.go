// Copyright 2024 eve.  All rights reserved.

package db

import (
	"time"

	gosqlmysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQLOptions defines optsions for mysql database.
type MySQLOptions struct {
	Host                  string
	Username              string
	Password              string
	Database              string
	MaxIdleConnections    int
	MaxOpenConnections    int
	MaxConnectionLifeTime time.Duration
	MaxConnectionIdleTime time.Duration
	// +optional
	Logger logger.Interface
}

// DSN return DSN from MySQLOptions.
// 通过 go-sql-driver 的 Config 构造 DSN，保证密码中的特殊字符（如 @、/）被正确转义。
func (o *MySQLOptions) DSN() string {
	cfg := gosqlmysql.NewConfig()
	cfg.User = o.Username
	cfg.Passwd = o.Password
	cfg.Net = "tcp"
	cfg.Addr = o.Host
	cfg.DBName = o.Database
	cfg.ParseTime = true
	cfg.Loc = time.Local
	cfg.Params = map[string]string{"charset": "utf8mb4"}
	return cfg.FormatDSN()
}

// NewMySQL create a new gorm db instance with the given options.
func NewMySQL(opts *MySQLOptions) (*gorm.DB, error) {
	// Set default values to ensure all fields in opts are available.
	setDefaults(opts)

	db, err := gorm.Open(mysql.Open(opts.DSN()), &gorm.Config{
		// PrepareStmt executes the given query in cached statement.
		// This can improve performance.
		PrepareStmt: true,
		Logger:      opts.Logger,
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(opts.MaxOpenConnections)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(opts.MaxConnectionLifeTime)

	sqlDB.SetConnMaxIdleTime(opts.MaxConnectionIdleTime)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(opts.MaxIdleConnections)

	return db, nil
}

// setDefaults set available default values for some fields.
func setDefaults(opts *MySQLOptions) {
	if opts.Host == "" {
		opts.Host = "127.0.0.1:3306"
	}
	if opts.MaxIdleConnections == 0 {
		opts.MaxIdleConnections = 100
	}
	if opts.MaxOpenConnections == 0 {
		opts.MaxOpenConnections = 100
	}
	// 与 options.MySQLOptions 的默认值保持一致，避免两层默认值不一致导致行为差异
	if opts.MaxConnectionLifeTime == 0 {
		opts.MaxConnectionLifeTime = time.Hour
	}
	if opts.MaxConnectionIdleTime == 0 {
		opts.MaxConnectionIdleTime = 30 * time.Minute
	}
	if opts.Logger == nil {
		opts.Logger = logger.Default
	}
}
