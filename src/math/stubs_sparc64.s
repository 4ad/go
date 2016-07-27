// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

TEXT ·Asin(SB),NOSPLIT|NOFRAME,$0
	JMP ·asin(SB)

TEXT ·Acos(SB),NOSPLIT|NOFRAME,$0
	JMP ·acos(SB)

TEXT ·Atan2(SB),NOSPLIT|NOFRAME,$0
	JMP ·atan2(SB)

TEXT ·Atan(SB),NOSPLIT|NOFRAME,$0
	JMP ·atan(SB)

TEXT ·Dim(SB),NOSPLIT|NOFRAME,$0
	JMP ·dim(SB)

TEXT ·Min(SB),NOSPLIT|NOFRAME,$0
	JMP ·min(SB)

TEXT ·Max(SB),NOSPLIT|NOFRAME,$0
	JMP ·max(SB)

TEXT ·Exp2(SB),NOSPLIT|NOFRAME,$0
	JMP ·exp2(SB)

TEXT ·Expm1(SB),NOSPLIT|NOFRAME,$0
	JMP ·expm1(SB)

TEXT ·Exp(SB),NOSPLIT|NOFRAME,$0
	JMP ·exp(SB)

TEXT ·Floor(SB),NOSPLIT|NOFRAME,$0
	JMP ·floor(SB)

TEXT ·Ceil(SB),NOSPLIT|NOFRAME,$0
	JMP ·ceil(SB)

TEXT ·Trunc(SB),NOSPLIT|NOFRAME,$0
	JMP ·trunc(SB)

TEXT ·Frexp(SB),NOSPLIT|NOFRAME,$0
	JMP ·frexp(SB)

TEXT ·Hypot(SB),NOSPLIT|NOFRAME,$0
	JMP ·hypot(SB)

TEXT ·Ldexp(SB),NOSPLIT|NOFRAME,$0
	JMP ·ldexp(SB)

TEXT ·Log10(SB),NOSPLIT|NOFRAME,$0
	JMP ·log10(SB)

TEXT ·Log2(SB),NOSPLIT|NOFRAME,$0
	JMP ·log2(SB)

TEXT ·Log1p(SB),NOSPLIT|NOFRAME,$0
	JMP ·log1p(SB)

TEXT ·Log(SB),NOSPLIT|NOFRAME,$0
	JMP ·log(SB)

TEXT ·Modf(SB),NOSPLIT|NOFRAME,$0
	JMP ·modf(SB)

TEXT ·Mod(SB),NOSPLIT|NOFRAME,$0
	JMP ·mod(SB)

TEXT ·Remainder(SB),NOSPLIT|NOFRAME,$0
	JMP ·remainder(SB)

TEXT ·Sincos(SB),NOSPLIT|NOFRAME,$0
	JMP ·sincos(SB)

TEXT ·Sin(SB),NOSPLIT|NOFRAME,$0
	JMP ·sin(SB)

TEXT ·Cos(SB),NOSPLIT|NOFRAME,$0
	JMP ·cos(SB)

TEXT ·Sqrt(SB),NOSPLIT|NOFRAME,$0
	JMP ·sqrt(SB)

TEXT ·Tan(SB),NOSPLIT|NOFRAME,$0
	JMP ·tan(SB)
