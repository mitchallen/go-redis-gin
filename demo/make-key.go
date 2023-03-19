/**
 * Author: Mitch Allen
 * File: make-key.go
 */

package demo

import (
	"fmt"
	"strings"
)

func MakeKey(namespace string, resource string) string {
	return fmt.Sprintf(
		"%s:%s",
		strings.ToLower(namespace),
		strings.ToLower(resource),
	)
}
