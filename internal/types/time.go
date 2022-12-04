package types

import "time"

func ParseTime(l string) time.Time {
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339} {
		t, err := time.Parse(layout, l)
		if err == nil {
			return t
		}
	}
	return time.Time{}
}
