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
// FIXED_FRAME+8(BSP) etc.) and some other low-level places.
//
// Register usage conventions keeping the above in mind:
//
// * FP will make use of RFP, which is not always appropiate, e.g. in NOFRAME
//   functions.
// * Generally, push arguments to off(BFP); in NOFRAME functions use
//   use off(BSP).
// * Example: setcallerpc does not return any arguments, it modifies the
//   link register saved on the caller's stack, so it has to use BFP.
// * Generally, use ret+off(FP) so go vet can actually check it for
//   corectness if there's a go declaration.  Although this cannot be used in
//   NOFRAME functions, etc.
//
// See runtime/stack.go for details.
#define ARG_PUSH_SIZE 6*8
#define WINDOW_SIZE 16*8
#define FIXED_FRAME WINDOW_SIZE+ARG_PUSH_SIZE
