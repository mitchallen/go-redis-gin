/**
 * Author: Mitch Allen
 * File: lock.go
 */

package demo

type Lock struct {
	Resource string `json:"resource"`
	UserID   string `json:"userId"`
	Duration string `json:"duration"`
}

const LOCK_NAMESPACE = "lock"

func MakeLockKey(resource string) string {

	return MakeKey(LOCK_NAMESPACE, resource)
}
