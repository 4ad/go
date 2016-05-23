// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import "unsafe"

var _ = unsafe.Sizeof(0)

//go:nosplit
//go:linkname time_now_sparc64 time.now
func time_now_sparc64() (sec int64, nsec int32) {
	ns := nanotime()
	return ns / 1000000000, int32(ns % 100000000)
}
