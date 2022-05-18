package notes

import (
	"sync"
	"time"
)

var jstLocation *time.Location
var jstOnce sync.Once

func JST() *time.Location {
	if jstLocation == nil {
		jstOnce.Do(func() {
			l, err := time.LoadLocation("Asia/Tokyo")
			if err != nil {
				l = time.FixedZone("JST2", +9*60*60)
			}
			jstLocation = l
		})
	}
	return jstLocation
}

type TimeStamp struct {
	value time.Time
}

func (t TimeStamp) Format() string {
	return t.value.Format(time.RFC3339)
}
