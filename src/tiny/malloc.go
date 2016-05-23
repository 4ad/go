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

// implementation of make builtin for slices
func newarray(typ *_type, n uintptr) unsafe.Pointer {
	if int(n) < 0 || typ.size > 0 {
		panic("runtime: allocation size out of range")
	}
	return malloc(uintptr(typ.size)*n, 8)
}

// rawmem returns a chunk of pointerless memory. It is
// not zeroed.
func rawmem(size uintptr) unsafe.Pointer {
	return malloc(size, 8)
}

// base address for all 0-byte allocations
var zerobase uintptr

// Returns size of the memory block that mallocgc will allocate if you ask for the size.
func roundupsize(size uintptr) uintptr {
	return round(size, 8192)
}
