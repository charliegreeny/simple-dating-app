package app

type ErrBadRequest struct {
	Message string
}

func (e ErrBadRequest) Error() string {
	return e.Message
}

type ErrUnauthorized struct {
	Message string
}

func (e ErrUnauthorized) Error() string {
	return e.Message
}

type ErrForbidden struct {
	Message string
}

func (e ErrForbidden) Error() string {
	return e.Message
}

type ErrNotFound struct {
	Message string
}

func (e ErrNotFound) Error() string {
	return e.Message
}
