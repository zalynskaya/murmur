package middleware

import "encoding/json"

type CustomError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
}

var (
	ErrEntityNotFound = NewCustomError(nil, "entity not found", "")
	ErrUserDuplicate  = NewCustomError(nil, "user already exists", "")
)

func NewCustomError(err error, message, developerMessage string) *CustomError {
	return &CustomError{
		Err:              err,
		Message:          message,
		DeveloperMessage: developerMessage,
	}
}

func (c CustomError) Unwrap() error {
	return c.Err
}

func (c CustomError) Error() string {
	return c.Message
}

func (c CustomError) Marshal() []byte {
	marshal, err := json.Marshal(c)
	if err != nil {
		return nil
	}

	return marshal
}

func systemError(err error) *CustomError {
	return NewCustomError(err, "internal system error", err.Error())
}
