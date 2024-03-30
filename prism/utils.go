package prism

import "unsafe"

func stringToBytes(s string) []byte {
	p := unsafe.StringData(s)
	b := unsafe.Slice(p, len(s))
	return b
}

func bytesToString(b []byte) string {
	p := unsafe.SliceData(b)
	return unsafe.String(p, len(b))
}
