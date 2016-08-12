// run

// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
	"regexp/syntax"
)

// Verify that stack growth works as expected for heavily-recursive
// functions such as writeRegexp, the one triggered below, when printing a
// regex to a buffer.
func main() {
	text := `^x{1,1000}y{1,1000}$`
	re, err := syntax.Parse(text, syntax.Perl)
	if err != nil {
		panic(fmt.Sprintf("parse: %v", err))
	}

	sre := re.Simplify()
	buf := fmt.Sprintf("	%+v\n", sre)
	if len(buf) != 11993 {
		panic(fmt.Sprintf("simplified regex not expected length: %d", len(buf)))
	}
}
