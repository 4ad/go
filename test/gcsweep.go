// run
package main

import (
	. "reflect"
	"runtime"
)

func main() {
        type T int
        st := SliceOf(TypeOf(T(1)))
        v := MakeSlice(st, 1, 1)
        runtime.GC()
        for i := 0; i < v.Len(); i++ {
                v.Index(i).Set(ValueOf(T(i)))
                runtime.GC()
        }
}
