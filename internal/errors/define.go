package errors

// Error Define
var (
	ErrBadRequest            = &_error{StatusCode: 400, ErrorCode: "400000", Message: "Invailed Input"}
	ErrPageNotFound          = &_error{StatusCode: 404, ErrorCode: "404000", Message: "Page Not Found"}
	ErrResourceNotFound      = &_error{StatusCode: 404, ErrorCode: "404001", Message: "Resource Not Found"}
	ErrDataConflict          = &_error{StatusCode: 409, ErrorCode: "409000", Message: "Data Conflict"}
	ErrResourceAlreadyExists = &_error{StatusCode: 409, ErrorCode: "409001", Message: "Resource Already Exists"}

	ErrInternalError = &_error{StatusCode: 500, ErrorCode: "500000", Message: "Server Internal Error"}
)
