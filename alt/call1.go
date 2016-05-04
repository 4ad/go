package main

import "unsafe"

var _ = unsafe.Sizeof(libc_write)

//go:cgo_import_dynamic libc_write write "libc.so"
//go:linkname libc_write libc_write
var libc_write uintptr

//go:linkname hello hello
func hello()

//go:linkname main main
//go:nosplit
func main() {
	foo()
}

//go:nosplit
func foo() {
	bar(1, 2, 3)
}

//go:nosplit
func bar(x, y, z int) {
	hello()
}
