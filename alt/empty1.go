package main

import "unsafe"

//go:linkname main main
//go:nosplit
func main() {
	x := 42
	_ = unsafe.Sizeof(x)
}
