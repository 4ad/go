// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"
#include "funcdata.h"
#include "asm_sparc64.h"

// makeFuncStub is the code half of the function returned by MakeFunc.
// See the comment on the declaration of makeFuncStub in makefunc.go
// for more details.
// No arg size here, runtime pulls arg map out of the func value.
TEXT ·makeFuncStub(SB),(NOSPLIT|WRAPPER),$16
	NO_LOCAL_POINTERS
	MOVD	CTXT, FIXED_FRAME+0(BSP)
	MOVD	$argframe+0(FP), RT1
	MOVD	RT1, FIXED_FRAME+8(BSP)
	CALL	·callReflect(SB)
	RET

// methodValueCall is the code half of the function returned by makeMethodValue.
// See the comment on the declaration of methodValueCall in makefunc.go
// for more details.
// No arg size here; runtime pulls arg map out of the func value.
TEXT ·methodValueCall(SB),(NOSPLIT|WRAPPER),$16
	NO_LOCAL_POINTERS
	MOVD	CTXT, FIXED_FRAME+0(BSP)
	MOVD	$argframe+0(FP), RT1
	MOVD	RT1, FIXED_FRAME+8(BSP)
	CALL	·callMethod(SB)
	RET
