package errors

import "errors"

var (
	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrTokenExpired       = errors.New("token expired")
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailExists        = errors.New("email already exists")
	ErrPhoneExists        = errors.New("phone already exists")

	// Driver errors
	ErrDriverNotFound      = errors.New("driver not found")
	ErrDriverExists        = errors.New("driver already exists")
	ErrLicenseExists       = errors.New("license number already exists")
	ErrInvalidVehicleType  = errors.New("invalid vehicle type")
	ErrInvalidDocumentType = errors.New("invalid document type")
	ErrMissingDocuments    = errors.New("missing required documents")
	ErrDriverNotVerified   = errors.New("driver not verified")
	ErrInvalidLocation     = errors.New("invalid location coordinates")
	ErrDocumentNotFound    = errors.New("document not found")
	ErrUnauthorizedAccess  = errors.New("unauthorized access")
)

type ErrorResponse struct {
	Success bool       `json:"success"`
	Error   *ErrorData `json:"error"`
}

type ErrorData struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func NewErrorResponse(code, message string, details interface{}) *ErrorResponse {
	return &ErrorResponse{
		Success: false,
		Error: &ErrorData{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
}

// Error code mapping
var ErrorCodes = map[error]string{
	ErrInvalidCredentials:  "AUTH001",
	ErrTokenExpired:        "AUTH002",
	ErrInvalidToken:        "AUTH003",
	ErrUserNotFound:        "AUTH004",
	ErrEmailExists:         "AUTH005",
	ErrPhoneExists:         "AUTH006",
	ErrDriverNotFound:      "DRV001",
	ErrInvalidVehicleType:  "DRV002",
	ErrInvalidDocumentType: "DRV003",
	ErrMissingDocuments:    "DRV004",
	ErrDriverNotVerified:   "DRV005",
	ErrInvalidLocation:     "DRV006",
}
