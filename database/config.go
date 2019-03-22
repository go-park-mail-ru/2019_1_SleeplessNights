package database

import (
	"fmt"
	"github.com/jackc/pgx"
	"time"
)

var connConfig = pgx.ConnConfig {
	Host: "localhost",
	Port: 5001,
	Database: "forum",
	User: "maxim",
	Password: "starwars",
	TLSConfig: nil,
	UseFallbackTLS: false,
	FallbackTLSConfig: nil,
	//Logger: ,
	LogLevel: pgx.LogLevelInfo,
	//Dial: ,
	RuntimeParams: nil,
	OnNotice: nil,
	CustomConnInfo: nil,
	CustomCancel: nil,
}

var poolConfig = pgx.ConnPoolConfig {
	ConnConfig:     connConfig,
	MaxConnections: 10,
	AfterConnect:   afterConnect,
	AcquireTimeout: 5 * time.Second,
}

func afterConnect(connection *pgx.Conn) error {
	//TODO WRITE FUNC
	fmt.Println("Connected")
	return nil
}