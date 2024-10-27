package models

type ErrInvalidPasscode struct {
	Err       string
	CsrfToken string
}

func (e ErrInvalidPasscode) Error() string {
	return e.Err
}
