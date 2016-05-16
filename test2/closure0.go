// run

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Test the behavior of closures.

package main

var fail bool

func newfunc() func(int) int       { return func(x int) int { return x } }
func newfunc2(x int) func(int) int { return func(int) int { return x } }

func main() {
	x, y := newfunc(), newfunc()
	if x(1) != 1 || y(2) != 2 {
		println("newfunc returned broken funcs")
		fail = true
	}
	x, y = newfunc2(2), newfunc2(1)
	if x(1) != 2 || y(2) != 1 {
		println("newfunc2 returned broken funcs")
		fail = true
	}

	if fail {
		panic("fail")
	}
}

func ff(x int) {
	call(func() {
		_ = x
	})
}

func call(func()) {
}
