// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"unsafe"
)

func asmcgocall2(fn, arg unsafe.Pointer) int32
func save_g()

// Call fn(arg) on the scheduler stack,
// aligned appropriately for the gcc ABI.
// See cgocall.go for more details.
//go:linkname asmcgocall_go runtime.asmcgocall
//go:nosplit
func asmcgocall_go(fn, arg unsafe.Pointer) (r int32) {
	systemstack(func() {
		save_g()
		r = asmcgocall2(fn, arg)
	})
	return r
}
