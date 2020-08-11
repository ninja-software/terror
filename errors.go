package terror

import (
	"errors"
)

// A list of errors that is pre-made

// ErrDataloader when could not retrieve data
var ErrDataloader = errors.New("could not fetch data")

// ErrDataBlank when data is blank/null, when it shouldnt
var ErrDataBlank = errors.New("blank data")

// ErrNoKeys error when no keys provided to GetMany
var ErrNoKeys = errors.New("no keys provided")

// ErrDataArchived when data is archived and cannot be changed
var ErrDataArchived = errors.New("data already archived")

// ErrWrongLength error when GetMany return wrong number of data
var ErrWrongLength = errors.New("get many return wrong length")

// ErrInvalidInput when data is invalid
var ErrInvalidInput = errors.New("invalid input data")

// ErrBadClaims when JWT could not be read
var ErrBadClaims = errors.New("could not read credentials from JWT")

// ErrBlacklisted when a typecast fails
var ErrBlacklisted = errors.New("token has been blacklisted")

// ErrTypeCast when a typecast fails
var ErrTypeCast = errors.New("could not cast interface to type")

// ErrParse when a parse fails
var ErrParse = errors.New("could not parse input")

// ErrBadCredentials when a bad username or password is passed in
var ErrBadCredentials = errors.New("bad credentials")

// ErrNotImplemented for non-implemented funcs
var ErrNotImplemented = errors.New("not implemented")

// ErrUnauthorized for bad permissions.
// Deprecated: Use terror.ErrUnauthorised instead
var ErrUnauthorized = errors.New("unauthorized")

// ErrUnauthorised for bad permissions
var ErrUnauthorised = errors.New("unauthorised")

// ErrNilUUID for nil uuid
var ErrNilUUID = errors.New("uuid is nil")

// ErrBadContext for missing context values
var ErrBadContext = errors.New("bad context")

// ErrAuthNoEmail during authentication when login failed due to non-existant user email address
var ErrAuthNoEmail = errors.New("user not found")

// ErrAuthWrongPassword during authentication when login failed due to incorrect incorrect password
var ErrAuthWrongPassword = errors.New("wrong password")
