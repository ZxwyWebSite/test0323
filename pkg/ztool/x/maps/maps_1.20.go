//go:build !go1.21

package maps

// Clone returns a copy of m.  This is a shallow clone:
// the new keys and values are set using ordinary assignment.
func Clone[M ~map[K]V, K comparable, V any](m M) M {
	// Preserve nil in case it matters.
	if m == nil {
		return nil
	}
	n := make(M, len(m))
	Copy(n, m)
	return n
}
