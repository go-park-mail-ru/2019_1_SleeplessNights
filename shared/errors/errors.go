package errors

import "errors"

var (
	DataBaseUniqueViolation = errors.New("ERROR: unique_violation (SQLSTATE 23505)")
	DataBaseNoDataFound     = errors.New("ERROR: no_data_found (SQLSTATE P0002)")
	AuthWrongPassword       = errors.New("authentication failed because of wong password")
)
