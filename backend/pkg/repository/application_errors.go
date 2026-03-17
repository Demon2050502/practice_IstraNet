package repository

import "errors"

var ErrAppNotFound = errors.New("application not found")
var ErrForbidden = errors.New("forbidden")
var ErrAlreadyAssigned = errors.New("application already assigned")
var ErrInvalidStatusTransition = errors.New("invalid status transition")
var ErrInvalidStatusCode = errors.New("invalid status code")

