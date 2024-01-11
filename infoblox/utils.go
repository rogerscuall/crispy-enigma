package infoblox

import "strings"

func isDataConflictError(err error) bool {
	return strings.Contains(err.Error(), "IB.Data.Conflict")
}
