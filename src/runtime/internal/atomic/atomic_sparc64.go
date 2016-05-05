// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package atomic

import "unsafe"

//go:noescape
func Xadd(ptr *uint32, delta int32) uint32

//go:noescape
func Xadd64(ptr *uint64, delta int64) uint64

//go:noescape
func Xadduintptr(ptr *uintptr, delta uintptr) uintptr

//go:noescape
func Xchg(ptr *uint32, new uint32) uint32

//go:noescape
func Xchg64(ptr *uint64, new uint64) uint64

//go:noescape
func Xchguintptr(ptr *uintptr, new uintptr) uintptr

//go:noescape
func Load(ptr *uint32) uint32

//go:noescape
func Load64(ptr *uint64) uint64

//go:noescape
func Loadp(ptr unsafe.Pointer) unsafe.Pointer

//go:nosplit
func And8(addr *uint8, v uint8) {
	// TODO(aram) implement this in asm.
	// Align down to 4 bytes and use 32-bit CAS.
	uaddr := uintptr(unsafe.Pointer(addr))
	addr32 := (*uint32)(unsafe.Pointer(uaddr &^ 3))
	word := (uint32(v) << 24) >> ((uaddr & 3) * 8)    // big endian
	mask := (uint32(0xFF) << 24) >> ((uaddr & 3) * 8) // big endian
	word |= ^mask
	for {
		old := *addr32
		if Cas(addr32, old, old&word) {
			return
		}
	}
}

//go:nosplit
func Or8(addr *uint8, v uint8) {
	// TODO(aram) implement this in asm.
	// Align down to 4 bytes and use 32-bit CAS.
	uaddr := uintptr(unsafe.Pointer(addr))
	addr32 := (*uint32)(unsafe.Pointer(uaddr &^ 3))
	word := (uint32(v) << 24) >> ((uaddr & 3) * 8) // big endian
	for {
		old := *addr32
		if Cas(addr32, old, old|word) {
			return
		}
	}
}

// NOTE: Do not add atomicxor8 (XOR is not idempotent).

//go:noescape
func Cas64(ptr *uint64, old, new uint64) bool

//go:noescape
func Store(ptr *uint32, val uint32)

//go:noescape
func Store64(ptr *uint64, val uint64)

// NO go:noescape annotation; see atomic_pointer.go.
func Storep1(ptr unsafe.Pointer, val unsafe.Pointer)
