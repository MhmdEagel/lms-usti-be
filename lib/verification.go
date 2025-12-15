package lib

import "time"

func IsExpired(expires *time.Time) bool {
	val := expires.Compare(time.Now())
	return val == -1
}