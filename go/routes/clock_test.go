package routes

import (
	"testing"
	"time"
)

type tClock struct {
	now time.Time
}

func (c *tClock) Now() time.Time {
	return c.now
}

var testClock = &tClock{time.Now()}

func TestRealClock_Now(t *testing.T) {
	now := realClock{}.Now()
	timeDiff := time.Since(now)

	if timeDiff > time.Second {
		t.Error("timeDiff > time.Second")
	}
	if timeDiff < 0 {
		t.Error("timeDiff < 0")
	}
}
