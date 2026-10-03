package system

import (
	"time"

	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Clock is the wall-clock implementation of out.Clock.
type Clock struct{}

var _ out.Clock = Clock{}

// Now returns the current time.
func (Clock) Now() time.Time { return time.Now() }
