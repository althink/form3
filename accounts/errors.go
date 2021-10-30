package accounts

type AccountNotFoundError struct {
	ID  string
	msg string
}

func (e *AccountNotFoundError) Error() string {
	return e.msg
}

type InvalidVersionError struct {
	Ver string
	msg string
}

func (e *InvalidVersionError) Error() string {
	return e.msg
}
