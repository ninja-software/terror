package terror

import "fmt"

// A list of errors that is pre-made

// ErrDataloader when could not retrieve data
var ErrDataloader = fmt.Errorf("could not fetch data")

// ErrDataBlank when data is blank/null, when it shouldnt
var ErrDataBlank = fmt.Errorf("blank data")

// ErrNoKeys error when no keys provided to GetMany
var ErrNoKeys = fmt.Errorf("no keys provided")

// ErrDataArchived when data is archived and cannot be changed
var ErrDataArchived = fmt.Errorf("data already archived")

// ErrWrongLength error when GetMany return wrong number of data
var ErrWrongLength = fmt.Errorf("get many return wrong length")

// ErrInvalidInput when data is invalid
var ErrInvalidInput = fmt.Errorf("invalid input data")

// ErrBadClaims when JWT could not be read
var ErrBadClaims = fmt.Errorf("could not read credentials from JWT")

// ErrBlacklisted when a typecast fails
var ErrBlacklisted = fmt.Errorf("token has been blacklisted")

// ErrTypeCast when a typecast fails
var ErrTypeCast = fmt.Errorf("could not cast interface to type")

// ErrParse when a parse fails
var ErrParse = fmt.Errorf("could not parse input")

// ErrBadCredentials when a bad username or password is passed in
var ErrBadCredentials = fmt.Errorf("bad credentials")

// ErrNotImplemented for non-implemented funcs
var ErrNotImplemented = fmt.Errorf("not implemented")

// ErrUnauthorised is for when unauthenticated
var ErrUnauthorised = fmt.Errorf("unauthorised")

// ErrNilUUID for nil uuid
var ErrNilUUID = fmt.Errorf("uuid is nil")

// ErrBadContext for missing context values
var ErrBadContext = fmt.Errorf("bad context")

// ErrAuthNoEmail during authentication when login failed due to non-existant user email address
var ErrAuthNoEmail = fmt.Errorf("user not found")

// ErrAuthWrongPassword during authentication when login failed due to incorrect incorrect password
var ErrAuthWrongPassword = fmt.Errorf("wrong password")

// ErrForbidden when user access is for not allowed access
var ErrForbidden = fmt.Errorf("access forbidden")
