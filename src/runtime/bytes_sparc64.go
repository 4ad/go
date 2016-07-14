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
