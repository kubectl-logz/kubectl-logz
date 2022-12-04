package types

import "time"

func ParseTime(l string) time.Time {
	for _, layout := range []string{time.RFC3339Nano, time.RFC3339, "2006-01-02T15:04:05,999-0700"} {
		t, err := time.Parse(layout, l)
		if err == nil {
			return t
		}
	}
	return time.Time{}
}
