// run

// Copyright 2011 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test divide corner cases.

package main

func f8(x, y, q, r int8) {
	if t := x / y; t != q {
		println(x, "/", y, "=", t, "want", q)
		panic("divide")
	}
	if t := x % y; t != r {
		println(x, "/", y, "=", t, "want", r)
		panic("divide")
	}
}

func f16(x, y, q, r int16) {
	if t := x / y; t != q {
		println(x, "/", y, "=", t, "want", q)
		panic("divide")
	}
	if t := x % y; t != r {
		println(x, "/", y, "=", t, "want", r)
		panic("divide")
	}
}

func f32(x, y, q, r int32) {
	if t := x / y; t != q {
		println(x, "/", y, "=", t, "want", q)
		panic("divide")
	}
	if t := x % y; t != r {
		println(x, "/", y, "=", t, "want", r)
		panic("divide")
	}
}

func f64(x, y, q, r int64) {
	if t := x / y; t != q {
		println(x, "/", y, "=", t, "want", q)
		panic("divide")
	}
	if t := x % y; t != r {
		println(x, "/", y, "=", t, "want", r)
		panic("divide")
	}
}

func main() {
	f8(-1<<7, -1, -1<<7, 0)
	f16(-1<<15, -1, -1<<15, 0)
	f32(-1<<31, -1, -1<<31, 0)
	f64(-1<<63, -1, -1<<63, 0)
}
