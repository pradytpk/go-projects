package main

import "net/http"

var ErrUserInvalid = apiError{Err: "user not valid", Status: http.StatusForbidden}

func (e apiError) Error() string {
	return e.Err
}
