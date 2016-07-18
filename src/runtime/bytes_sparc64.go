// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import "unsafe"

var _ = unsafe.Pointer(uintptr(42))

//go:linkname indexByte_bytes bytes.IndexByte
//go:nosplit
func indexByte_bytes(s []byte, c byte) int {
	for i, b := range s {
		if b == c {
			return i
		}
	}
	return -1
}

//go:linkname indexByte_strings strings.IndexByte
//go:nosplit
func indexByte_strings(s string, c byte) int {
	for i, b := range []byte(s) {
		if b == c {
			return i
		}
	}
	return -1
}

//go:linkname eqstring_go runtime.eqstring
//go:nosplit
func eqstring_go(s1, s2 string) bool {
	if len(s1) != len(s2) {
		return false
	}
	// optimization in assembly versions:
	// if s1.str == s2.str { return true }
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

//go:linkname equal_bytes bytes.Equal
//go:nosplit
func equal_bytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, c := range a {
		if c != b[i] {
			return false
		}
	}
	return true
}
