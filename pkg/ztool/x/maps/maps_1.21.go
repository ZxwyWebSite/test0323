//go:build go1.21

package maps

import (
	_ "maps"
	_ "unsafe"
)

// clone is implemented in the runtime package.
//
//go:linkname clone maps.clone
func clone(m any) any

// Clone returns a copy of m.  This is a shallow clone:
// the new keys and values are set using ordinary assignment.
func Clone[M ~map[K]V, K comparable, V any](m M) M {
	// Preserve nil in case it matters.
	if m == nil {
		return nil
	}
	return clone(m).(M)
}
