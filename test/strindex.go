// run

package main

import "strings"

func main() {
	str := "hellolllo"
	if strings.IndexByte(str, byte('l')) != 2 {
		panic("strings.Index broken")
	}
}
