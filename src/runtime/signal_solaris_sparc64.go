// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"runtime/internal/sys"
	"unsafe"
)

type sigctxt struct {
	info *siginfo
	ctxt unsafe.Pointer
}

func (c *sigctxt) regs() *mcontext {
	return (*mcontext)(unsafe.Pointer(&(*ucontext)(c.ctxt).uc_mcontext))
}

func (c *sigctxt) r1() uint64  { return uint64(c.regs().gregs[_REG_G1]) }
func (c *sigctxt) r2() uint64  { return uint64(c.regs().gregs[_REG_G2]) }
func (c *sigctxt) r3() uint64  { return uint64(c.regs().gregs[_REG_G3]) }
func (c *sigctxt) r4() uint64  { return uint64(c.regs().gregs[_REG_G4]) }
func (c *sigctxt) r5() uint64  { return uint64(c.regs().gregs[_REG_G5]) }
func (c *sigctxt) r6() uint64  { return uint64(c.regs().gregs[_REG_G6]) }
func (c *sigctxt) r7() uint64  { return uint64(c.regs().gregs[_REG_G7]) }
func (c *sigctxt) r8() uint64  { return uint64(c.regs().gregs[_REG_O0]) }
func (c *sigctxt) r9() uint64  { return uint64(c.regs().gregs[_REG_O1]) }
func (c *sigctxt) r10() uint64 { return uint64(c.regs().gregs[_REG_O2]) }
func (c *sigctxt) r11() uint64 { return uint64(c.regs().gregs[_REG_O3]) }
func (c *sigctxt) r12() uint64 { return uint64(c.regs().gregs[_REG_O4]) }
func (c *sigctxt) r13() uint64 { return uint64(c.regs().gregs[_REG_O5]) }

func (c *sigctxt) sp() uint64 { return uint64(c.regs().gregs[_REG_O6]) }
func (c *sigctxt) lr() uint64 { return uint64(c.regs().gregs[_REG_O7]) }

func (c *sigctxt) r16() uint64 {
	if c.regs().gwins != nil {
		cwp := int(c.regs().gregs[_REG_CCR] & 0x1f)
		return uint64(c.regs().gwins.wbuf[cwp].local[0])
	}
	return *(*uint64)(unsafe.Pointer((uintptr)(c.regs().gregs[_REG_O6] + sys.StackBias + 0*8)))
}

func (c *sigctxt) r17() uint64 {
	if c.regs().gwins != nil {
		cwp := int(c.regs().gregs[_REG_CCR] & 0x1f)
		return uint64(c.regs().gwins.wbuf[cwp].local[1])
	}
	return *(*uint64)(unsafe.Pointer((uintptr)(c.regs().gregs[_REG_O6] + sys.StackBias + 1*8)))
}

func (c *sigctxt) r18() uint64 {
	if c.regs().gwins != nil {
		cwp := int(c.regs().gregs[_REG_CCR] & 0x1f)
		return uint64(c.regs().gwins.wbuf[cwp].local[2])
	}
	return *(*uint64)(unsafe.Pointer((uintptr)(c.regs().gregs[_REG_O6] + sys.StackBias + 2*8)))
}

func (c *sigctxt) r19() uint64 {
	if c.regs().gwins != nil {
		cwp := int(c.regs().gregs[_REG_CCR] & 0x1f)
		return uint64(c.regs().gwins.wbuf[cwp].local[3])
	}
	return *(*uint64)(unsafe.Pointer((uintptr)(c.regs().gregs[_REG_O6] + sys.StackBias + 3*8)))
}

func (c *sigctxt) r20() uint64 {
	if c.regs().gwins != nil {
		cwp := int(c.regs().gregs[_REG_CCR] & 0x1f)
		return uint64(c.regs().gwins.wbuf[cwp].local[4])
	}
	return *(*uint64)(unsafe.Pointer((uintptr)(c.regs().gregs[_REG_O6] + sys.StackBias + 4*8)))
}

func (c *sigctxt) r21() uint64 {
	if c.regs().gwins != nil {
		cwp := int(c.regs().gregs[_REG_CCR] & 0x1f)
		return uint64(c.regs().gwins.wbuf[cwp].local[5])
	}
	return *(*uint64)(unsafe.Pointer((uintptr)(c.regs().gregs[_REG_O6] + sys.StackBias + 5*8)))
}

func (c *sigctxt) r22() uint64 {
	if c.regs().gwins != nil {
		cwp := int(c.regs().gregs[_REG_CCR] & 0x1f)
		return uint64(c.regs().gwins.wbuf[cwp].local[6])
	}
	return *(*uint64)(unsafe.Pointer((uintptr)(c.regs().gregs[_REG_O6] + sys.StackBias + 6*8)))
}

func (c *sigctxt) r23() uint64 {
	if c.regs().gwins != nil {
		cwp := int(c.regs().gregs[_REG_CCR] & 0x1f)
		return uint64(c.regs().gwins.wbuf[cwp].local[7])
	}
	return *(*uint64)(unsafe.Pointer((uintptr)(c.regs().gregs[_REG_O6] + sys.StackBias + 7*8)))
}

func (c *sigctxt) r24() uint64 {
	if c.regs().gwins != nil {
		cwp := int(c.regs().gregs[_REG_CCR] & 0x1f)
		return uint64(c.regs().gwins.wbuf[cwp].in[0])
	}
	return *(*uint64)(unsafe.Pointer((uintptr)(c.regs().gregs[_REG_O6] + sys.StackBias + 8*8)))
}

func (c *sigctxt) r25() uint64 {
	if c.regs().gwins != nil {
		cwp := int(c.regs().gregs[_REG_CCR] & 0x1f)
		return uint64(c.regs().gwins.wbuf[cwp].in[1])
	}
	return *(*uint64)(unsafe.Pointer((uintptr)(c.regs().gregs[_REG_O6] + sys.StackBias + 9*8)))
}

func (c *sigctxt) r26() uint64 {
	if c.regs().gwins != nil {
		cwp := int(c.regs().gregs[_REG_CCR] & 0x1f)
		return uint64(c.regs().gwins.wbuf[cwp].in[2])
	}
	return *(*uint64)(unsafe.Pointer((uintptr)(c.regs().gregs[_REG_O6] + sys.StackBias + 10*8)))
}

func (c *sigctxt) r27() uint64 {
	if c.regs().gwins != nil {
		cwp := int(c.regs().gregs[_REG_CCR] & 0x1f)
		return uint64(c.regs().gwins.wbuf[cwp].in[3])
	}
	return *(*uint64)(unsafe.Pointer((uintptr)(c.regs().gregs[_REG_O6] + sys.StackBias + 11*8)))
}

func (c *sigctxt) r28() uint64 {
	if c.regs().gwins != nil {
		cwp := int(c.regs().gregs[_REG_CCR] & 0x1f)
		return uint64(c.regs().gwins.wbuf[cwp].in[4])
	}
	return *(*uint64)(unsafe.Pointer((uintptr)(c.regs().gregs[_REG_O6] + sys.StackBias + 12*8)))
}

func (c *sigctxt) r29() uint64 {
	if c.regs().gwins != nil {
		cwp := int(c.regs().gregs[_REG_CCR] & 0x1f)
		return uint64(c.regs().gwins.wbuf[cwp].in[5])
	}
	return *(*uint64)(unsafe.Pointer((uintptr)(c.regs().gregs[_REG_O6] + sys.StackBias + 13*8)))
}

func (c *sigctxt) r30() uint64 {
	if c.regs().gwins != nil {
		cwp := int(c.regs().gregs[_REG_CCR] & 0x1f)
		return uint64(c.regs().gwins.wbuf[cwp].in[6])
	}
	return *(*uint64)(unsafe.Pointer((uintptr)(c.regs().gregs[_REG_O6] + sys.StackBias + 14*8)))
}

func (c *sigctxt) r31() uint64 {
	if c.regs().gwins != nil {
		cwp := int(c.regs().gregs[_REG_CCR] & 0x1f)
		return uint64(c.regs().gwins.wbuf[cwp].in[7])
	}
	return *(*uint64)(unsafe.Pointer((uintptr)(c.regs().gregs[_REG_O6] + sys.StackBias + 15*8)))
}

func (c *sigctxt) pc() uint64     { return uint64(c.regs().gregs[_REG_PC]) }
func (c *sigctxt) tstate() uint64 { return uint64(c.regs().gregs[_REG_CCR]) }

func (c *sigctxt) sigcode() uint64 { return uint64(c.info.si_code) }
func (c *sigctxt) sigaddr() uint64 { return *(*uint64)(unsafe.Pointer(&c.info.__data[0])) }

func (c *sigctxt) set_pc(x uint64) { c.regs().gregs[_REG_PC] = int64(x) }
func (c *sigctxt) set_sp(x uint64) { c.regs().gregs[_REG_O6] = int64(x) }
func (c *sigctxt) set_lr(x uint64) { c.regs().gregs[_REG_O7] = int64(x) }

func (c *sigctxt) set_sigcode(x uint64) { c.info.si_code = int32(x) }
func (c *sigctxt) set_sigaddr(x uint64) {
	*(*uintptr)(unsafe.Pointer(&c.info.__data[0])) = uintptr(x)
}
