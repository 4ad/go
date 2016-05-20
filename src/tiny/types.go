package runtime

import "unsafe"

// Stack describes a Go execution stack.
// The bounds of the stack are exactly [lo, hi),
// with no implicit data structures on either side.
type stack struct {
	lo uintptr
	hi uintptr
}

type g struct {
	// Stack parameters.
	// stack describes the actual stack memory: [stack.lo, stack.hi).
	// stackguard0 is the stack pointer compared in the Go stack growth prologue.
	// It is stack.lo+StackGuard normally, but can be StackPreempt to trigger a preemption.
	// stackguard1 is the stack pointer compared in the C stack growth prologue.
	// It is stack.lo+StackGuard on g0 and gsignal stacks.
	// It is ~0 on other goroutine stacks, to trigger a call to morestackc (and crash).
	stack       stack   // offset known to runtime/cgo
	stackguard0 uintptr // offset known to liblink
	stackguard1 uintptr // offset known to liblink

	_panic       *_panic // innermost panic - offset known to liblink
	_defer       *_defer // innermost defer
	sched        gobuf
	syscallsp    uintptr        // if status==gsyscall, syscallsp = sched.sp to use during gc
	syscallpc    uintptr        // if status==gsyscall, syscallpc = sched.pc to use during gc
	param        unsafe.Pointer // passed parameter on wakeup
	atomicstatus uint32
	goid         int64
	waitsince    int64  // approx time when the g become blocked
	waitreason   string // if status==gwaiting
	schedlink    *g
	issystem     bool // do not output in stack dump, ignore in deadlock detector
	preempt      bool // preemption signal, duplicates stackguard0 = stackpreempt
	paniconfault bool // panic (instead of crash) on unexpected fault address
	preemptscan  bool // preempted g does scan for gc
	gcworkdone   bool // debug: cleared at begining of gc work phase cycle, set by gcphasework, tested at end of cycle
	throwsplit   bool // must not split stack
	raceignore   int8 // ignore race detection events
	m            *m   // for debuggers, but offset not hard-coded
	lockedm      *m
	sig          uint32
	writebuf     []byte
	sigcode0     uintptr
	sigcode1     uintptr
	sigpc        uintptr
	gopc         uintptr // pc of go statement that created this goroutine
	racectx      uintptr
	waiting      *sudog // sudog structures this g is waiting on (that have a valid elem ptr)
}

type gobuf struct {
	// The offsets of sp, pc, and g are known to (hard-coded in) libmach.
	sp   uintptr
	pc   uintptr
	g    guintptr
	ctxt unsafe.Pointer // this has to be a pointer so that gc scans it
	ret  uintreg
	lr   uintptr
}

// Known to compiler.
// Changes here must also be made in src/cmd/gc/select.c's selecttype.
type sudog struct {
	g           *g
	selectdone  *uint32
	next        *sudog
	prev        *sudog
	elem        unsafe.Pointer // data element
	releasetime int64
	nrelease    int32  // -1 for acquire
	waitlink    *sudog // g.waiting list
}

type m struct {
	g0      *g    // goroutine with scheduling stack
	morebuf gobuf // gobuf arg to morestack

	// Fields not known to debuggers.
	procid  uint64 // for debuggers, but offset not hard-coded
	gsignal *g     // signal-handling g
	//tls           [4]uintptr     // thread-local storage (for x86 extern register)
	//mstartfn      unsafe.Pointer // todo go func()
	curg *g // current running goroutine
	//caughtsig     *g             // goroutine running during fatal signal
	//p             *p             // attached p for executing go code (nil if not executing go code)
	//nextp         *p
	//id            int32
	//mallocing     int32
	//throwing      int32
	//gcing         int32
	//locks         int32
	//softfloat     int32
	//dying         int32
	//profilehz     int32
	//helpgc        int32
	//spinning      bool // m is out of work and is actively looking for work
	//blocked       bool // m is blocked on a note
	//inwb          bool // m is executing a write barrier
	//printlock     int8
	fastrand uint32
	//ncgocall      uint64 // number of cgo calls in total
	//ncgo          int32  // number of cgo calls currently in progress
	//cgomal        *cgomal
	//park          note
	//alllink       *m // on allm
	//schedlink     *m
	//machport      uint32 // return address for mach ipc (os x)
	//mcache        *mcache
	//lockedg       *g
	//createstack   [32]uintptr // stack that created this thread.
	//freglo        [16]uint32  // d[i] lsb and f[i]
	//freghi        [16]uint32  // d[i] msb and f[i+16]
	//fflag         uint32      // floating point compare flags
	//locked        uint32      // tracking for lockosthread
	//nextwaitm     *m          // next m waiting for lock
	//waitsema      uintptr     // semaphore for parking on locks
	//waitsemacount uint32
	//waitsemalock  uint32
	//gcstats       gcstats
	//needextram    bool
	//traceback     uint8
	//waitunlockf   unsafe.Pointer // todo go func(*g, unsafe.pointer) bool
	//waitlock      unsafe.Pointer
	////#ifdef GOOS_windows
	//thread uintptr // thread handle
	//// these are here because they are too large to be on the stack
	//// of low-level NOSPLIT functions.
	//libcall   libcall
	//libcallpc uintptr // for cpu profiler
	//libcallsp uintptr
	//libcallg  *g
	////#endif
	////#ifdef GOOS_solaris
	//perrno *int32 // pointer to tls errno
	//// these are here because they are too large to be on the stack
	//// of low-level NOSPLIT functions.
	////LibCall	libcall;
	//ts      mts
	//scratch mscratch
	////#endif
	////#ifdef GOOS_plan9
	//notesig *int8
	//errstr  *byte
	////#endif
}

type _string struct {
	str *byte
	len int
}

type funcval struct {
	fn uintptr
	// variable-size, fn-specific data here
}

type iface struct {
	tab  *itab
	data unsafe.Pointer
}

type eface struct {
	_type *_type
	data  unsafe.Pointer
}

type slice struct {
	array *byte // actual data
	len   uint  // number of elements
	cap   uint  // allocated number of elements
}

// A guintptr holds a goroutine pointer, but typed as a uintptr
// to bypass write barriers. It is used in the Gobuf goroutine state.
//
// The Gobuf.g goroutine pointer is almost always updated by assembly code.
// In one of the few places it is updated by Go code - func save - it must be
// treated as a uintptr to avoid a write barrier being emitted at a bad time.
// Instead of figuring out how to emit the write barriers missing in the
// assembly manipulation, we change the type of the field to uintptr,
// so that it does not require write barriers at all.
//
// Goroutine structs are published in the allg list and never freed.
// That will keep the goroutine structs from being collected.
// There is never a time that Gobuf.g's contain the only references
// to a goroutine: the publishing of the goroutine in allg comes first.
// Goroutine pointers are also kept in non-GC-visible places like TLS,
// so I can't see them ever moving. If we did want to start moving data
// in the GC, we'd need to allocate the goroutine structs from an
// alternate arena. Using guintptr doesn't make that problem any worse.
type guintptr uintptr

func (gp guintptr) ptr() *g {
	return (*g)(unsafe.Pointer(gp))
}

// Layout of in-memory per-function information prepared by linker
// See http://golang.org/s/go12symtab.
// Keep in sync with linker and with ../../libmach/sym.c
// and with package debug/gosym and with symtab.go in package runtime.
type _func struct {
	entry   uintptr // start pc
	nameoff int32   // function name

	args  int32 // in/out args size
	frame int32 // legacy frame size; use pcsp if possible

	pcsp      int32
	pcfile    int32
	pcln      int32
	npcdata   int32
	nfuncdata int32
}

// layout of Itab known to compilers
// allocated in non-garbage-collected memory
type itab struct {
	inter  *interfacetype
	_type  *_type
	link   *itab
	bad    int32
	unused int32
	fun    [1]uintptr // variable sized
}

/*
 * deferred subroutine calls
 */
type _defer struct {
	siz     int32
	started bool
	sp      uintptr // sp at time of defer
	pc      uintptr
	fn      *funcval
	_panic  *_panic // panic that is running defer
	link    *_defer
}

/*
 * panics
 */
type _panic struct {
	argp      unsafe.Pointer // pointer to arguments of deferred call run during panic; cannot move - known to liblink
	arg       interface{}    // argument to panic
	link      *_panic        // link to earlier panic
	recovered bool           // whether this panic is over
	aborted   bool           // the panic was aborted
}

// Dummy types
type p struct{}

// A Func represents a Go function in the running binary.
type Func struct {
	opaque struct{} // unexported field to disallow conversions
}
