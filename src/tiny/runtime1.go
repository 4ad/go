package runtime

import (
	"runtime/internal/sys"
	"unsafe"
)

type libcall struct {
	fn   uintptr
	n    uintptr // number of parameters
	args uintptr // parameters
	r1   uintptr // return values
	r2   uintptr
	err  uintptr // error number
}

// Helpers for Go. Must be NOSPLIT, must only call NOSPLIT functions, and must not block.

//go:nosplit
func acquirem() *m {
	_g_ := getg()
	//_g_.m.locks++
	return _g_.m
}

//go:nosplit
func releasem(mp *m) {
	//_g_ := getg()
	//mp.locks--
	//if mp.locks == 0 && _g_.preempt {
	//	// restore the preemption request in case we've cleared it in newstack
	//	_g_.stackguard0 = stackPreempt
	//}
}

// funcPC returns the entry PC of the function f.
// It assumes that f is a func value. Otherwise the behavior is undefined.
//go:nosplit
func funcPC(f interface{}) uintptr {
	return **(**uintptr)(add(unsafe.Pointer(&f), sys.PtrSize))
}
