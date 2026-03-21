package prism

import "unsafe"

// stringToBytes returns the underlying bytes of s without copying.
// The returned slice must not be mutated — it shares memory with the string.
func stringToBytes(s string) []byte {
	p := unsafe.StringData(s)
	b := unsafe.Slice(p, len(s))
	return b
}

// fastBoolConv converts a bool to 0 or 1.
// Safe because the Go spec guarantees false=0, true=1 in memory.
func fastBoolConv(b bool) int {
	return int(*(*byte)(unsafe.Pointer(&b)))
}
