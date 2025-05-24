package oauth

import "errors"

var (
	ErrInvalidProvider = errors.New("invalid provider")
	ErrInvalidCode     = errors.New("invalid authorization code")
)
