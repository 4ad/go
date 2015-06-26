// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

func DRconv(a int8) string {
	if a >= ClassUnknown && a <= ClassNone {
		return cnames[a]
	}
	return "C_??"
}
