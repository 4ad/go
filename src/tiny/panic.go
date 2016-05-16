package runtime

import "unsafe"

func throw(s string) {
	println("THROW:", s)
	exit(41)
}

func newdefer(siz int32) *_defer {
	var d *_defer
	mp := acquirem()

	// Allocate new defer+args.
	const deferHeaderSize = unsafe.Sizeof(_defer{})
	total := round(deferHeaderSize, regSize) + uintptr(siz)
	d = (*_defer)(malloc(total, 0))

	d.siz = siz
	gp := mp.curg
	d.link = gp._defer
	gp._defer = d
	releasem(mp)
	return d
}

// The arguments associated with a deferred call are stored
// immediately after the _defer header in memory.
//go:nosplit
func deferArgs(d *_defer) unsafe.Pointer {
	return add(unsafe.Pointer(d), unsafe.Sizeof(*d))
}

// Create a new deferred function fn with siz bytes of arguments.
// The compiler turns a defer statement into a call to this.
//go:nosplit
func deferproc(siz int32, fn *funcval) { // arguments of fn follow fn
	if getg().m.curg != getg() {
		// go code on the system stack can't defer
		throw("defer on system stack")
	}

	// the arguments of fn are in a perilous state.  The stack map
	// for deferproc does not describe them.  So we can't let garbage
	// collection or stack copying trigger until we've copied them out
	// to somewhere safe.  The memmove below does that.
	// Until the copy completes, we can only call nosplit routines.
	sp := getcallersp(unsafe.Pointer(&siz))
	argp := uintptr(unsafe.Pointer(&fn)) + unsafe.Sizeof(fn)
	callerpc := getcallerpc(unsafe.Pointer(&siz))

	//systemstack(func() {
	d := newdefer(siz)
	if d._panic != nil {
		throw("deferproc: d.panic != nil after newdefer")
	}
	d.fn = fn
	d.pc = callerpc
	d.sp = sp
	memmove(add(unsafe.Pointer(d), unsafe.Sizeof(*d)), unsafe.Pointer(argp), uintptr(siz))
	//})

	// deferproc returns 0 normally.
	// a deferred func that stops a panic
	// makes the deferproc return 1.
	// the code the compiler generates always
	// checks the return value and jumps to the
	// end of the function if deferproc returns != 0.
	return0()
	// No code can go here - the C return register has
	// been set and must not be clobbered.
}

// Run a deferred function if there is one.
// The compiler inserts a call to this at the end of any
// function which calls defer.
// If there is a deferred function, this will call runtime·jmpdefer,
// which will jump to the deferred function such that it appears
// to have been called by the caller of deferreturn at the point
// just before deferreturn was called.  The effect is that deferreturn
// is called again and again until there are no more deferred functions.
// Cannot split the stack because we reuse the caller's frame to
// call the deferred function.

// The single argument isn't actually used - it just has its address
// taken so it can be matched against pending defers.
//go:nosplit
func deferreturn(arg0 uintptr) {
	gp := getg()
	d := gp._defer
	if d == nil {
		return
	}
	sp := getcallersp(unsafe.Pointer(&arg0))
	if d.sp != sp {
		return
	}

	// Moving arguments around.
	// Do not allow preemption here, because the garbage collector
	// won't know the form of the arguments until the jmpdefer can
	// flip the PC over to fn.
	mp := acquirem()
	memmove(unsafe.Pointer(&arg0), deferArgs(d), uintptr(d.siz))
	fn := d.fn
	d.fn = nil
	gp._defer = d.link
	//freedefer(d)
	releasem(mp)
	jmpdefer(fn, uintptr(unsafe.Pointer(&arg0)))
}
