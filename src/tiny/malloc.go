// +build arm64

package runtime

import "unsafe"

var oldBreak uintptr

const mallocVerbose = !true

func newobject(typ *_type) unsafe.Pointer {
	return malloc(typ.size, uintptr(typ.align))
}

// malloc allocates memory for at least size bytes, and aligned to align bytes.
// If align is zero, a default alignment of regSize is used.
func malloc(size, align uintptr) unsafe.Pointer {
	if oldBreak == 0 {
		oldBreak = brk(0)
	}
	if mallocVerbose {
		print("Allocating ", size, " bytes, align = ", align, ", brk = ", hex(oldBreak))
	}
	if align == 0 {
		align = regSize
	}
	addr := round(oldBreak, align)
	oldBreak = addr + size
	if newBrk := brk(oldBreak); newBrk != oldBreak {
		panic("brk failed")
	}
	if mallocVerbose {
		println(", newbrk =", hex(oldBreak))
	}
	return unsafe.Pointer(addr)
}

// round n up to a multiple of a.  a must be a power of 2.
func round(n, a uintptr) uintptr {
	return (n + a - 1) &^ (a - 1)
}
