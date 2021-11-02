package accounts

import "fmt"

type InvalidDataError struct {
	Code string `json:"error_code"`
	Msg  string `json:"error_message"`
}

func (e *InvalidDataError) Error() string {
	return e.Msg
}

type AccountNotFoundError struct {
	ID string
}

func (e *AccountNotFoundError) Error() string {
	return fmt.Sprintf("account not found: %s", e.ID)
}

type AccountAlreadyExistsError struct {
	ID string
}

func (e *AccountAlreadyExistsError) Error() string {
	return fmt.Sprintf("account already exists: %s", e.ID)
}

type InvalidVersionError struct {
	Ver int64
}

func (e *InvalidVersionError) Error() string {
	return fmt.Sprintf("invalid version: %d", e.Ver)
}

type HttpStatusError struct {
	StatusCode int
}

func (e *HttpStatusError) Error() string {
	return fmt.Sprintf("error code returned: %d", e.StatusCode)
}
