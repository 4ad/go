// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// SPARC64 TOS is STACK_BIAS bytes *above* from RSP (R14, %o6).
// Current frame is STACK_BIAS bytes *above* RFP (R30, %i6).
#define STACK_BIAS 0x7ff

// FIXED_FRAME defines the size of the fixed part of a stack frame. A stack
// frame looks like this:
//
// +---------------------+
// | local variable area |
// +---------------------+
// | argument area       |
// +---------------------+ <- BSP+FIXED_FRAME
// | fixed area          |
// +---------------------+ <- BSP
//
// So a function that sets up a stack frame at all uses as least FIXED_FRAME
// bytes of stack. This mostly affects assembly that calls other functions
// with arguments (the arguments should be stored at FIXED_FRAME+0(BSP),
// FIXED_FRAME+8(BSP) etc) and some other low-level places.

#define FIXED_FRAME 176
