package errors

import "errors"

var (
	DataBaseUniqueViolation = errors.New("rpc error: code = Unknown desc = ERROR: unique violation exception in database")
	DataBaseNoDataFound     = errors.New("rpc error: code = Unknown desc = ERROR: no data found exception in database")
	AuthWrongPassword       = errors.New("rpc error: code = Unknown desc = ERROR: authentication failed, because of wong password")
)
