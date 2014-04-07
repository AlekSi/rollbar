package rollbar

import (
	"net"
)

type FatalError struct {
	error
}

// check interfaces
var (
	_ net.Error = FatalError{}
)

func (f FatalError) Timeout() bool   { return false }
func (f FatalError) Temporary() bool { return false }
