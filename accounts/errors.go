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
	return fmt.Sprintf("Account not found: %s", e.ID)
}

type InvalidVersionError struct {
	Ver int64
}

func (e *InvalidVersionError) Error() string {
	return fmt.Sprintf("Invalid version: %d", e.Ver)
}
