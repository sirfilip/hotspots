package errors

import "errors"

var (
	RecordNotFound          = errors.New("Record Not Found")
	ForbidenMarshalingError = errors.New("Marshaling is forbiden")
	TokenExpiredError       = errors.New("Token is expired")
	InvalidGeometryError    = errors.New("Invalid geometry")
)

func Is(err, other error) bool {
	return errors.Is(err, other)
}
