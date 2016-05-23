// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import "unsafe"

// fInterface is our standard non-empty interface.  We use it instead
// of interface{f()} in function prototypes because gofmt insists on
// putting lots of newlines in the otherwise concise interface{f()}.
type fInterface interface {
	f()
}

func convT2E(t *_type, elem unsafe.Pointer, x unsafe.Pointer) (e eface) {
	if isDirectIface(t) {
		e._type = t
		typedmemmove(t, unsafe.Pointer(&e.data), elem)
	} else {
		if x == nil {
			x = newobject(t)
		}
		// TODO: We allocate a zeroed object only to overwrite it with
		// actual data. Figure out how to avoid zeroing. Also below in convT2I.
		typedmemmove(t, x, elem)
		e._type = t
		e.data = x
	}
	return
}
