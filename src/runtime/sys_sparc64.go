// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import "unsafe"

// adjust Gobuf as if it executed a call to fn with context ctxt
// and then did an immediate Gosave.
func gostartcall(buf *gobuf, fn, ctxt unsafe.Pointer) {
	if buf.lr != 0 {
		throw("invalid use of gostartcall")
	}
	buf.lr = buf.pc
	buf.pc = uintptr(fn)
	buf.ctxt = ctxt
}

// Called to rewind context saved during morestack back to beginning of function.
// To help us, the linker emits a call back to the beginning 8 bytes after the
// call to morestack. We just have to decode and apply that jump.
func rewindmorestack(buf *gobuf) {
	var inst uint32
	if buf.pc&3 == 0 && buf.pc != 0 {
		// For sparc, the pc register holds the address of the
		// *current* instruction, rather than the next instruction
		// to execute, and CTIs are padded with a nop to avoid DCTI
		// coupling, so we need to skip ahead to get the jump.
		buf.pc += 8
		inst = *(*uint32)(unsafe.Pointer(buf.pc))
		// Extract opcode
		op1 := inst >> 30
		// call and link
		mop1 := 1 << 30

		call := uint32(mop1 >> 30)
		if op1 == call {
			// Extract pc-relative address (4*sign_ext(disp30))
			idisp30 := 4 * (int32(inst<<2) >> 2)

			//ipc := uintptr(unsafe.Pointer(buf.pc))
			buf.pc += uintptr(idisp30)

			//print("runtime: rewind pc=", hex(ipc), " to pc=", hex(buf.pc), "\n");
			return
		}
	}
	print("runtime: pc=", hex(buf.pc), " ", hex(inst), "\n")
	throw("runtime: misuse of rewindmorestack")
}

func usleep2(us uint32)

//go:linkname usleep1_go runtime.usleep1
//go:nosplit
func usleep1_go(µs uint32) {
	_g_ := getg()

	// Check the validity of m because we might be called in cgo callback
	// path early enough where there isn't a m available yet.
	if _g_ != nil && _g_.m != nil {
		sysvicall1(&libc_usleep, uintptr(µs))
		return
	}
	usleep2(µs)
}
