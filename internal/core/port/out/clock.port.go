package out

import "time"

// Clock abstracts time for TTL testing.
// Implemented by adapter/out/clock/*.
type Clock interface {
	Now() time.Time
}
