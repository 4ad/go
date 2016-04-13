package main

import "unsafe"

//go:cgo_import_dynamic libc_write write "libc.so"
//go:linkname libc_write libc_write
var libc_write uintptr

//go:nosplit
func foo(v uintptr) uintptr {
	return v
}

//go:linkname main main
//go:nosplit
func main() {
	_ = unsafe.Sizeof(libc_write)
	foo(libc_write)
}
