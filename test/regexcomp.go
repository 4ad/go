// run
package main

import (
	"fmt"
	"regexp/syntax"
)

// Check that one-pass cutoff does trigger.
func main() {
        re, err := syntax.Parse(`^x{1,1000}y{1,1000}$`, syntax.Perl)
        if err != nil {
                panic(fmt.Sprintf("parse: %v", err))
        }
        syntax.Compile(re.Simplify())
}
