// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build arm64

// Hashing algorithm inspired by ?

package runtime

import "unsafe"

func memhash(p unsafe.Pointer, seed, s uintptr) uintptr {
	return 0
}

func nilinterhash(p unsafe.Pointer, h uintptr) uintptr {
	return 0
}
