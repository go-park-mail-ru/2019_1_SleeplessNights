package database

type SQL interface {
	GetQuery() string
	GetArgs() []interface{}
}
