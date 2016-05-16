// Copyright 2015 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Garbage collector: write barriers.
//
// For the concurrent garbage collector, the Go compiler implements
// updates to pointer-valued fields that may be in heap objects by
// emitting calls to write barriers. This file contains the actual write barrier
// implementation, markwb, and the various wrappers called by the
// compiler to implement pointer assignment, slice assignment,
// typed memmove, and so on.
//
// To check for missed write barriers, the GODEBUG=wbshadow debugging
// mode allocates a second copy of the heap. Write barrier-based pointer
// updates make changes to both the real heap and the shadow, and both
// the pointer updates and the GC look for inconsistencies between the two,
// indicating pointer writes that bypassed the barrier.

package runtime

import "unsafe"

// NOTE: Really dst *unsafe.Pointer, src unsafe.Pointer,
// but if we do that, Go inserts a write barrier on *dst = src.
//go:nosplit
func writebarrierptr(dst *uintptr, src uintptr) {
	*dst = src
	return
}

//go:nosplit
func writebarrierstring(dst *[2]uintptr, src [2]uintptr) {
	writebarrierptr(&dst[0], src[0])
	dst[1] = src[1]
}

//go:nosplit
func writebarrierslice(dst *[3]uintptr, src [3]uintptr) {
	writebarrierptr(&dst[0], src[0])
	dst[1] = src[1]
	dst[2] = src[2]
}

//go:nosplit
func writebarrieriface(dst *[2]uintptr, src [2]uintptr) {
	writebarrierptr(&dst[0], src[0])
	writebarrierptr(&dst[1], src[1])
}

//go:generate go run wbfat_gen.go -- wbfat.go
//
// The above line generates multiword write barriers for
// all the combinations of ptr+scalar up to four words.
// The implementations are written to wbfat.go.

// typedmemmove copies a value of type t to dst from src.
//go:nosplit
func typedmemmove(typ *_type, dst, src unsafe.Pointer) {
	memmove(dst, src, typ.size)
	return
}

//go:nosplit
func typedslicecopy(typ *_type, dst, src slice) int {
	n := dst.len
	if n > src.len {
		n = src.len
	}
	if n == 0 {
		return 0
	}
	dstp := unsafe.Pointer(dst.array)
	srcp := unsafe.Pointer(src.array)

	memmove(dstp, srcp, uintptr(n)*typ.size)
	return int(n)
}

//go:linkname reflect_typedslicecopy reflect.typedslicecopy
func reflect_typedslicecopy(elemType *_type, dst, src slice) int {
	return typedslicecopy(elemType, dst, src)
}
