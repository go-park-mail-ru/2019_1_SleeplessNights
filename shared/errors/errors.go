package errors

import "errors"

var (
	DataBaseUniqueViolation     = errors.New("ERROR: unique violation exception in database")
	DataBaseNoDataFound         = errors.New("ERROR: no data found exception in database")
	DataBaseForeignKeyViolation = errors.New("ERROR: foreign key violation exception in database")
	AuthWrongPassword           = errors.New("ERROR: authentication failed, because of wong password")
)
