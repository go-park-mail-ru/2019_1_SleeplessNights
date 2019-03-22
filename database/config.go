package database

import (
	"github.com/jackc/pgx"
)

var ConnConfig = pgx.ConnConfig {
	Host: "localhost",
	Port: 5432,
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

