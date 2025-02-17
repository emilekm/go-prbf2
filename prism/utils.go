package prism

import "unsafe"

func stringToBytes(s string) []byte {
	p := unsafe.StringData(s)
	b := unsafe.Slice(p, len(s))
	return b
}

func fastBoolConv(b bool) int {
	return int(*(*byte)(unsafe.Pointer(&b)))
}
