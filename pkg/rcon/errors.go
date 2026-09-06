package rcon

import "errors"

// Errors define the possible errors that can occur when connecting,
// authenticating, or communicating with an RCON server.
var (
	// ErrConnectionFailed is returned when the TCP connection to the
	// server could not be established (e.g. server offline, port closed).
	ErrConnectionFailed = errors.New("connection failed")
	// ErrAuthenticationFailed is returned when the server rejects the
	// provided password.
	ErrAuthenticationFailed = errors.New("authentication failed")
)
