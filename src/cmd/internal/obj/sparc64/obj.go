// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

import (
	"cmd/internal/obj"
	"encoding/binary"
)

var Linksparc64 = obj.LinkArch{
	ByteOrder: binary.BigEndian,
	Name:      "sparc64",
	Thechar:   'u',
	Minlc:     4,
	Ptrsize:   8,
	Regsize:   8,
}
