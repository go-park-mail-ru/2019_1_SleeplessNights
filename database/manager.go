package database

var db *dbManager

type dbManager struct {
	//TODO put connection pool here
}

func init() {
	db = &dbManager{
		//TODO init dbManager fields
	}
}

func GetInstance() *dbManager {
	return db
}
