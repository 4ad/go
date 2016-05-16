// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sys

// Declarations for runtime services implemented in C or assembly.

const PtrSize = 4 << (^uintptr(0) >> 63) // unsafe.Sizeof(uintptr(0)) but an ideal const
const RegSize = 4 << (^Uintreg(0) >> 63) // unsafe.Sizeof(uintreg(0)) but an ideal const
const StackBias = 0x7ff * GoarchSparc64  // 0 normally, 0x7ff for SPARC64
