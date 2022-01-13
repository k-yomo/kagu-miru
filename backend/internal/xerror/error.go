package xerror

import "errors"

type Type int

const (
	TypeUnknown Type = iota + 1
	TypeNotFound
)

type internalError struct {
	t   Type
	err error
}

func (i *internalError) Error() string {
	return i.err.Error()
}

func New(t Type, err error) *internalError {
	return &internalError{
		t:   t,
		err: err,
	}
}

func NewNotFound(err error) *internalError {
	return New(TypeNotFound, err)
}

func ErrorType(err error) Type {
	var e *internalError
	if errors.As(err, &e) {
		return e.t
	}
	return TypeUnknown
}
