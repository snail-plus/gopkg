// Copyright 2024 eve.  All rights reserved.

package config

import (
	"time"
)

// MySQLConfiguration defines the configuration of mysql
// clients for components that can run with mysql database.
type MySQLConfiguration struct {
	// MySQL service host address. If left blank, the following related mysql options will be ignored.
	Host string
	// Username for access to mysql service.
	Username string
	// Password for access to mysql, should be used pair with password.
	Password string
	// Database name for the server to use.
	Database string
	// Maximum idle connections allowed to connect to mysql.
	MaxIdleConnections int32
	// Maximum open connections allowed to connect to mysql.
	MaxOpenConnections int32
	// Maximum connection life time allowed to connect to mysql.
	MaxConnectionLifeTime time.Duration
}

// RedisConfiguration defines the configuration of redis
// clients for components that can run with redis key-value database.
type RedisConfiguration struct {
	// Address of your Redis server(ip:port).
	Addr string
	// Username for access to redis service.
	Username string
	// Optional auth password for Redis db.
	Password string
	// Database to be selected after connecting to the server.
	Database int
	// Maximum number of retries before giving up.
	MaxRetries int
	// Timeout when connecting to redis service.
	Timeout time.Duration
}
