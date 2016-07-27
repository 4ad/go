// run

package main

func main() {
	defer func() {
		p := recover()
		if p != "test panic" {
			panic("recover failed")
		}
	}()
	panic("test panic")
}
