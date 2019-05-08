package errors

import "errors"

var (
	DataBaseUniqueViolation = errors.New("rpc error: code = Unknown desc = ERROR: unique_violation (SQLSTATE 23505)")
	DataBaseNoDataFound     = errors.New("rpc error: code = Unknown desc = ERROR: no_data_found (SQLSTATE P0002)")
	AuthWrongPassword       = errors.New("rpc error: code = Unknown desc = authentication failed because of wong password")
)
