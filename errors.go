package main

import (
	"errors"
	"fmt"
)

type ValidationError struct {
	Err error
}

func (v *ValidationError) Error() string {
	return fmt.Sprintf("⚠ %s", v.Err)
}

func NewValidationError(err error) *ValidationError {
	return &ValidationError{err}
}

var (
	errNoUsername           = NewValidationError(errors.New("You must supply a username"))
	errNoEmail              = NewValidationError(errors.New("You must supply an email"))
	errNoPassword           = NewValidationError(errors.New("You must supply a password"))
	errPasswordTooShort     = NewValidationError(errors.New("Your password is too short"))
	errUsernameExists       = NewValidationError(errors.New("That username is taken"))
	errEmailExists          = NewValidationError(errors.New("That email address has an account"))
	errCredentialsIncorrect = NewValidationError(errors.New("We couldn't find a user with the supplied username and password combination"))
	errPasswordIncorrect    = NewValidationError(errors.New("Password did not match"))
	errInvalidImageType     = NewValidationError(errors.New("Please upload only jpeg, gif or png images"))
	errNoImage              = NewValidationError(errors.New("Please select an image to upload"))
	errImageUrlInvalid      = NewValidationError(errors.New("Couldn't download image from the Url you provided"))
)

func IsValidationError(err error) bool {
	_, ok := err.(*ValidationError)
	return ok
}
