// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// SPARC64 TOS is STACK_BIAS bytes *above* from RSP (R14, %o6).
// Current frame is STACK_BIAS bytes *above* RFP (R30, %i6).
#define STACK_BIAS 0x7ff
