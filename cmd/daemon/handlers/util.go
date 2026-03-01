package handlers

import (
	"time"
)

func toPtr[T any](v T) *T {
	return &v
}

func parseTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	return &t
}

func parseTimeString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
