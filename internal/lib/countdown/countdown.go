package countdown

import (
	"fmt"
	"time"
)

type Countdown struct {
	d int
	h int
	m int
	s int
}

func (c Countdown) Count(t time.Time) Countdown {
	diff := time.Since(t)

	total := int(diff.Seconds())

	days := int(total / (60 * 60 * 24))
	hours := int(total / (60 * 60) % 24)
	minutes := int(total/60) % 60
	seconds := int(total % 60)

	return Countdown{
		d: days,
		h: hours,
		m: minutes,
		s: seconds,
	}
}

func (c Countdown) String() string {

	return fmt.Sprintf("**Days**: %d\n**Hours**: %d\n**Minutes**: %d\n**Seconds**: %d", c.d, c.h, c.m, c.s)
}
