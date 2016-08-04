// run
package main

import (
	"fmt"
	"regexp/syntax"
)

// Check that one-pass cutoff does trigger.
func main() {
	text := `^x{1,1000}y{1,1000}$`
	println("text ", text)
	re, err := syntax.Parse(text, syntax.Perl)
	if err != nil {
		panic(fmt.Sprintf("parse: %v", err))
	}

	sre := re.Simplify()
	println("sre ", sre)
	fmt.Printf("	%+v\n", sre)
}
