package runtime

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
