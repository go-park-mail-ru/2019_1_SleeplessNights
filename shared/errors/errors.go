package errors

import "errors"

var (
	DataBaseUniqueViolation = errors.New("ERROR: the non-unique data")
	DataBaseNoDataFound     = errors.New("ERROR: no data found")
	AuthWrongPassword       = errors.New("rpc error: code = Unknown desc = authentication failed because of wong password")
)
