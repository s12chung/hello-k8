package routes

import "time"

type tClock struct {
	now time.Time
}

func (c *tClock) Now() time.Time {
	return c.now
}

var testClock = &tClock{time.Now()}
