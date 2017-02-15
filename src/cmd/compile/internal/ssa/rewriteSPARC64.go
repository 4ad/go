// autogenerated from gen/SPARC64.rules: do not edit!
// generated with: cd gen; go run *.go

package ssa

import "math"

var _ = math.MinInt8 // in case not otherwise used
func rewriteValueSPARC64(v *Value, config *Config) bool {
	switch v.Op {
	case OpAdd16:
		return rewriteValueSPARC64_OpAdd16(v, config)
	case OpAdd32:
		return rewriteValueSPARC64_OpAdd32(v, config)
	case OpAdd32F:
		return rewriteValueSPARC64_OpAdd32F(v, config)
	case OpAdd64:
		return rewriteValueSPARC64_OpAdd64(v, config)
	case OpAdd64F:
		return rewriteValueSPARC64_OpAdd64F(v, config)
	case OpAdd8:
		return rewriteValueSPARC64_OpAdd8(v, config)
	case OpAddPtr:
		return rewriteValueSPARC64_OpAddPtr(v, config)
	case OpAnd16:
		return rewriteValueSPARC64_OpAnd16(v, config)
	case OpAnd32:
		return rewriteValueSPARC64_OpAnd32(v, config)
	case OpAnd64:
		return rewriteValueSPARC64_OpAnd64(v, config)
	case OpAnd8:
		return rewriteValueSPARC64_OpAnd8(v, config)
	case OpAndB:
		return rewriteValueSPARC64_OpAndB(v, config)
	case OpConst16:
		return rewriteValueSPARC64_OpConst16(v, config)
	case OpConst32:
		return rewriteValueSPARC64_OpConst32(v, config)
	case OpConst32F:
		return rewriteValueSPARC64_OpConst32F(v, config)
	case OpConst64:
		return rewriteValueSPARC64_OpConst64(v, config)
	case OpConst64F:
		return rewriteValueSPARC64_OpConst64F(v, config)
	case OpConst8:
		return rewriteValueSPARC64_OpConst8(v, config)
	case OpConstBool:
		return rewriteValueSPARC64_OpConstBool(v, config)
	case OpConstNil:
		return rewriteValueSPARC64_OpConstNil(v, config)
	case OpDiv16:
		return rewriteValueSPARC64_OpDiv16(v, config)
	case OpDiv16u:
		return rewriteValueSPARC64_OpDiv16u(v, config)
	case OpDiv32:
		return rewriteValueSPARC64_OpDiv32(v, config)
	case OpDiv32F:
		return rewriteValueSPARC64_OpDiv32F(v, config)
	case OpDiv32u:
		return rewriteValueSPARC64_OpDiv32u(v, config)
	case OpDiv64:
		return rewriteValueSPARC64_OpDiv64(v, config)
	case OpDiv64F:
		return rewriteValueSPARC64_OpDiv64F(v, config)
	case OpDiv64u:
		return rewriteValueSPARC64_OpDiv64u(v, config)
	case OpDiv8:
		return rewriteValueSPARC64_OpDiv8(v, config)
	case OpDiv8u:
		return rewriteValueSPARC64_OpDiv8u(v, config)
	case OpMod16:
		return rewriteValueSPARC64_OpMod16(v, config)
	case OpMod16u:
		return rewriteValueSPARC64_OpMod16u(v, config)
	case OpMod32:
		return rewriteValueSPARC64_OpMod32(v, config)
	case OpMod32u:
		return rewriteValueSPARC64_OpMod32u(v, config)
	case OpMod64:
		return rewriteValueSPARC64_OpMod64(v, config)
	case OpMod64u:
		return rewriteValueSPARC64_OpMod64u(v, config)
	case OpMod8:
		return rewriteValueSPARC64_OpMod8(v, config)
	case OpMod8u:
		return rewriteValueSPARC64_OpMod8u(v, config)
	case OpMul16:
		return rewriteValueSPARC64_OpMul16(v, config)
	case OpMul32:
		return rewriteValueSPARC64_OpMul32(v, config)
	case OpMul32F:
		return rewriteValueSPARC64_OpMul32F(v, config)
	case OpMul64:
		return rewriteValueSPARC64_OpMul64(v, config)
	case OpMul64F:
		return rewriteValueSPARC64_OpMul64F(v, config)
	case OpMul8:
		return rewriteValueSPARC64_OpMul8(v, config)
	case OpNeg16:
		return rewriteValueSPARC64_OpNeg16(v, config)
	case OpNeg32:
		return rewriteValueSPARC64_OpNeg32(v, config)
	case OpNeg32F:
		return rewriteValueSPARC64_OpNeg32F(v, config)
	case OpNeg64:
		return rewriteValueSPARC64_OpNeg64(v, config)
	case OpNeg64F:
		return rewriteValueSPARC64_OpNeg64F(v, config)
	case OpNeg8:
		return rewriteValueSPARC64_OpNeg8(v, config)
	case OpNot:
		return rewriteValueSPARC64_OpNot(v, config)
	case OpOr16:
		return rewriteValueSPARC64_OpOr16(v, config)
	case OpOr32:
		return rewriteValueSPARC64_OpOr32(v, config)
	case OpOr64:
		return rewriteValueSPARC64_OpOr64(v, config)
	case OpOr8:
		return rewriteValueSPARC64_OpOr8(v, config)
	case OpOrB:
		return rewriteValueSPARC64_OpOrB(v, config)
	case OpSignExt16to32:
		return rewriteValueSPARC64_OpSignExt16to32(v, config)
	case OpSignExt16to64:
		return rewriteValueSPARC64_OpSignExt16to64(v, config)
	case OpSignExt32to64:
		return rewriteValueSPARC64_OpSignExt32to64(v, config)
	case OpSignExt8to16:
		return rewriteValueSPARC64_OpSignExt8to16(v, config)
	case OpSignExt8to32:
		return rewriteValueSPARC64_OpSignExt8to32(v, config)
	case OpSignExt8to64:
		return rewriteValueSPARC64_OpSignExt8to64(v, config)
	case OpSqrt:
		return rewriteValueSPARC64_OpSqrt(v, config)
	case OpSub16:
		return rewriteValueSPARC64_OpSub16(v, config)
	case OpSub32:
		return rewriteValueSPARC64_OpSub32(v, config)
	case OpSub32F:
		return rewriteValueSPARC64_OpSub32F(v, config)
	case OpSub64:
		return rewriteValueSPARC64_OpSub64(v, config)
	case OpSub64F:
		return rewriteValueSPARC64_OpSub64F(v, config)
	case OpSub8:
		return rewriteValueSPARC64_OpSub8(v, config)
	case OpSubPtr:
		return rewriteValueSPARC64_OpSubPtr(v, config)
	case OpTrunc16to8:
		return rewriteValueSPARC64_OpTrunc16to8(v, config)
	case OpTrunc32to16:
		return rewriteValueSPARC64_OpTrunc32to16(v, config)
	case OpTrunc32to8:
		return rewriteValueSPARC64_OpTrunc32to8(v, config)
	case OpTrunc64to16:
		return rewriteValueSPARC64_OpTrunc64to16(v, config)
	case OpTrunc64to32:
		return rewriteValueSPARC64_OpTrunc64to32(v, config)
	case OpTrunc64to8:
		return rewriteValueSPARC64_OpTrunc64to8(v, config)
	case OpXor16:
		return rewriteValueSPARC64_OpXor16(v, config)
	case OpXor32:
		return rewriteValueSPARC64_OpXor32(v, config)
	case OpXor64:
		return rewriteValueSPARC64_OpXor64(v, config)
	case OpXor8:
		return rewriteValueSPARC64_OpXor8(v, config)
	case OpZeroExt16to32:
		return rewriteValueSPARC64_OpZeroExt16to32(v, config)
	case OpZeroExt16to64:
		return rewriteValueSPARC64_OpZeroExt16to64(v, config)
	case OpZeroExt32to64:
		return rewriteValueSPARC64_OpZeroExt32to64(v, config)
	case OpZeroExt8to16:
		return rewriteValueSPARC64_OpZeroExt8to16(v, config)
	case OpZeroExt8to32:
		return rewriteValueSPARC64_OpZeroExt8to32(v, config)
	case OpZeroExt8to64:
		return rewriteValueSPARC64_OpZeroExt8to64(v, config)
	}
	return false
}
func rewriteValueSPARC64_OpAdd16(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Add16  x y)
	// cond:
	// result: (ADD x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64ADD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpAdd32(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Add32  x y)
	// cond:
	// result: (ADD x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64ADD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpAdd32F(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Add32F x y)
	// cond:
	// result: (FADDS x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64FADDS)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpAdd64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Add64  x y)
	// cond:
	// result: (ADD  x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64ADD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpAdd64F(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Add64F x y)
	// cond:
	// result: (FADDD x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64FADDD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpAdd8(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Add8   x y)
	// cond:
	// result: (ADD x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64ADD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpAddPtr(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (AddPtr x y)
	// cond:
	// result: (ADD  x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64ADD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpAnd16(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (And16 x y)
	// cond:
	// result: (AND x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64AND)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpAnd32(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (And32 x y)
	// cond:
	// result: (AND x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64AND)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpAnd64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (And64 x y)
	// cond:
	// result: (AND x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64AND)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpAnd8(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (And8 x y)
	// cond:
	// result: (AND x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64AND)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpAndB(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (AndB x y)
	// cond:
	// result: (AND x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64AND)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpConst16(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Const16  [val])
	// cond:
	// result: (MOVWconst [val])
	for {
		val := v.AuxInt
		v.reset(OpSPARC64MOVWconst)
		v.AuxInt = val
		return true
	}
}
func rewriteValueSPARC64_OpConst32(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Const32  [val])
	// cond:
	// result: (MOVWconst [val])
	for {
		val := v.AuxInt
		v.reset(OpSPARC64MOVWconst)
		v.AuxInt = val
		return true
	}
}
func rewriteValueSPARC64_OpConst32F(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Const32F [val])
	// cond:
	// result: (FMOVSconst [val])
	for {
		val := v.AuxInt
		v.reset(OpSPARC64FMOVSconst)
		v.AuxInt = val
		return true
	}
}
func rewriteValueSPARC64_OpConst64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Const64  [val])
	// cond:
	// result: (MOVDconst [val])
	for {
		val := v.AuxInt
		v.reset(OpSPARC64MOVDconst)
		v.AuxInt = val
		return true
	}
}
func rewriteValueSPARC64_OpConst64F(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Const64F [val])
	// cond:
	// result: (FMOVDconst [val])
	for {
		val := v.AuxInt
		v.reset(OpSPARC64FMOVDconst)
		v.AuxInt = val
		return true
	}
}
func rewriteValueSPARC64_OpConst8(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Const8   [val])
	// cond:
	// result: (MOVWconst [val])
	for {
		val := v.AuxInt
		v.reset(OpSPARC64MOVWconst)
		v.AuxInt = val
		return true
	}
}
func rewriteValueSPARC64_OpConstBool(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (ConstBool [b])
	// cond:
	// result: (MOVWconst [b])
	for {
		b := v.AuxInt
		v.reset(OpSPARC64MOVWconst)
		v.AuxInt = b
		return true
	}
}
func rewriteValueSPARC64_OpConstNil(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (ConstNil)
	// cond:
	// result: (MOVDconst [0])
	for {
		v.reset(OpSPARC64MOVDconst)
		v.AuxInt = 0
		return true
	}
}
func rewriteValueSPARC64_OpDiv16(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Div16 x y)
	// cond:
	// result: (SDIVD x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64SDIVD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpDiv16u(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Div16u x y)
	// cond:
	// result: (UDIVD x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64UDIVD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpDiv32(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Div32 x y)
	// cond:
	// result: (SDIVD x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64SDIVD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpDiv32F(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Div32F x y)
	// cond:
	// result: (FDIVS x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64FDIVS)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpDiv32u(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Div32u x y)
	// cond:
	// result: (UDIVD x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64UDIVD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpDiv64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Div64 x y)
	// cond:
	// result: (SDIVD x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64SDIVD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpDiv64F(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Div64F x y)
	// cond:
	// result: (FDIVD x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64FDIVD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpDiv64u(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Div64u x y)
	// cond:
	// result: (UDIVD x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64UDIVD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpDiv8(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Div8 x y)
	// cond:
	// result: (SDIVD x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64SDIVD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpDiv8u(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Div8u x y)
	// cond:
	// result: (UDIVD x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64UDIVD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpMod16(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Mod16 x y)
	// cond:
	// result: (Mod64 (SignExt16to64 x) (SignExt16to64 y))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpMod64)
		v0 := b.NewValue0(v.Line, OpSignExt16to64, config.fe.TypeInt64())
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Line, OpSignExt16to64, config.fe.TypeInt64())
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueSPARC64_OpMod16u(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Mod16u x y)
	// cond:
	// result: (Mod64u (ZeroExt16to64 x) (ZeroExt16to64 y))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpMod64u)
		v0 := b.NewValue0(v.Line, OpZeroExt16to64, config.fe.TypeUInt64())
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Line, OpZeroExt16to64, config.fe.TypeUInt64())
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueSPARC64_OpMod32(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Mod32 x y)
	// cond:
	// result: (Mod64 (SignExt32to64 x) (SignExt32to64 y))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpMod64)
		v0 := b.NewValue0(v.Line, OpSignExt32to64, config.fe.TypeInt64())
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Line, OpSignExt32to64, config.fe.TypeInt64())
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueSPARC64_OpMod32u(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Mod32u x y)
	// cond:
	// result: (Mod64u (ZeroExt32to64 x) (ZeroExt32to64 y))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpMod64u)
		v0 := b.NewValue0(v.Line, OpZeroExt32to64, config.fe.TypeUInt64())
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Line, OpZeroExt32to64, config.fe.TypeUInt64())
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueSPARC64_OpMod64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Mod64 x y)
	// cond:
	// result: (SUB x (MULD y (SDIVD x y)))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64SUB)
		v.AddArg(x)
		v0 := b.NewValue0(v.Line, OpSPARC64MULD, config.fe.TypeInt64())
		v0.AddArg(y)
		v1 := b.NewValue0(v.Line, OpSPARC64SDIVD, config.fe.TypeInt64())
		v1.AddArg(x)
		v1.AddArg(y)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueSPARC64_OpMod64u(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Mod64u x y)
	// cond:
	// result: (SUB x (MULD y (UDIVD x y)))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64SUB)
		v.AddArg(x)
		v0 := b.NewValue0(v.Line, OpSPARC64MULD, config.fe.TypeInt64())
		v0.AddArg(y)
		v1 := b.NewValue0(v.Line, OpSPARC64UDIVD, config.fe.TypeUInt64())
		v1.AddArg(x)
		v1.AddArg(y)
		v0.AddArg(v1)
		v.AddArg(v0)
		return true
	}
}
func rewriteValueSPARC64_OpMod8(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Mod8 x y)
	// cond:
	// result: (Mod64 (SignExt8to64 x) (SignExt8to64 y))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpMod64)
		v0 := b.NewValue0(v.Line, OpSignExt8to64, config.fe.TypeInt64())
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Line, OpSignExt8to64, config.fe.TypeInt64())
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueSPARC64_OpMod8u(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Mod8u x y)
	// cond:
	// result: (Mod64u (ZeroExt8to64 x) (ZeroExt8to64 y))
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpMod64u)
		v0 := b.NewValue0(v.Line, OpZeroExt8to64, config.fe.TypeUInt64())
		v0.AddArg(x)
		v.AddArg(v0)
		v1 := b.NewValue0(v.Line, OpZeroExt8to64, config.fe.TypeUInt64())
		v1.AddArg(y)
		v.AddArg(v1)
		return true
	}
}
func rewriteValueSPARC64_OpMul16(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Mul16 x y)
	// cond:
	// result: (MULD x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64MULD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpMul32(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Mul32 x y)
	// cond:
	// result: (MULD x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64MULD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpMul32F(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Mul32F x y)
	// cond:
	// result: (FMULS x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64FMULS)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpMul64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Mul64 x y)
	// cond:
	// result: (MULD x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64MULD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpMul64F(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Mul64F x y)
	// cond:
	// result: (FMULD x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64FMULD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpMul8(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Mul8 x y)
	// cond:
	// result: (MULD x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64MULD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpNeg16(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Neg16 x)
	// cond:
	// result: (NEG x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64NEG)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpNeg32(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Neg32 x)
	// cond:
	// result: (NEG x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64NEG)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpNeg32F(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Neg32F x)
	// cond:
	// result: (FNEGS x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64FNEGS)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpNeg64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Neg64 x)
	// cond:
	// result: (NEG x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64NEG)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpNeg64F(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Neg64F x)
	// cond:
	// result: (FNEGD x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64FNEGD)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpNeg8(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Neg8 x)
	// cond:
	// result: (NEG x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64NEG)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpNot(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Not x)
	// cond:
	// result: (XORconst [1] x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64XORconst)
		v.AuxInt = 1
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpOr16(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Or16 x y)
	// cond:
	// result: (OR x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64OR)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpOr32(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Or32 x y)
	// cond:
	// result: (OR x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64OR)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpOr64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Or64 x y)
	// cond:
	// result: (OR x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64OR)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpOr8(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Or8 x y)
	// cond:
	// result: (OR x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64OR)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpOrB(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (OrB x y)
	// cond:
	// result: (OR x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64OR)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpSignExt16to32(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (SignExt16to32 x)
	// cond:
	// result: (MOVHreg x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64MOVHreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpSignExt16to64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (SignExt16to64 x)
	// cond:
	// result: (MOVHreg x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64MOVHreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpSignExt32to64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (SignExt32to64 x)
	// cond:
	// result: (MOVWreg x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64MOVWreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpSignExt8to16(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (SignExt8to16  x)
	// cond:
	// result: (MOVBreg x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64MOVBreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpSignExt8to32(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (SignExt8to32  x)
	// cond:
	// result: (MOVBreg x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64MOVBreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpSignExt8to64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (SignExt8to64  x)
	// cond:
	// result: (MOVBreg x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64MOVBreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpSqrt(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Sqrt x)
	// cond:
	// result: (FSQRTD x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64FSQRTD)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpSub16(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Sub16 x y)
	// cond:
	// result: (SUB x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64SUB)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpSub32(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Sub32 x y)
	// cond:
	// result: (SUB x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64SUB)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpSub32F(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Sub32F x y)
	// cond:
	// result: (FSUBS x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64FSUBS)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpSub64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Sub64 x y)
	// cond:
	// result: (SUB x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64SUB)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpSub64F(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Sub64F x y)
	// cond:
	// result: (FSUBD x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64FSUBD)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpSub8(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Sub8 x y)
	// cond:
	// result: (SUB x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64SUB)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpSubPtr(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (SubPtr x y)
	// cond:
	// result: (SUB x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64SUB)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpTrunc16to8(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Trunc16to8  x)
	// cond:
	// result: (MOVBreg x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64MOVBreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpTrunc32to16(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Trunc32to16 x)
	// cond:
	// result: (MOVHreg x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64MOVHreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpTrunc32to8(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Trunc32to8  x)
	// cond:
	// result: (MOVBreg x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64MOVBreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpTrunc64to16(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Trunc64to16 x)
	// cond:
	// result: (MOVHreg x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64MOVHreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpTrunc64to32(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Trunc64to32 x)
	// cond:
	// result: (MOVWreg x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64MOVWreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpTrunc64to8(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Trunc64to8  x)
	// cond:
	// result: (MOVBreg x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64MOVBreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpXor16(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Xor16 x y)
	// cond:
	// result: (XOR x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64XOR)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpXor32(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Xor32 x y)
	// cond:
	// result: (XOR x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64XOR)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpXor64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Xor64 x y)
	// cond:
	// result: (XOR x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64XOR)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpXor8(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (Xor8 x y)
	// cond:
	// result: (XOR x y)
	for {
		x := v.Args[0]
		y := v.Args[1]
		v.reset(OpSPARC64XOR)
		v.AddArg(x)
		v.AddArg(y)
		return true
	}
}
func rewriteValueSPARC64_OpZeroExt16to32(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (ZeroExt16to32 x)
	// cond:
	// result: (MOVUHreg x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64MOVUHreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpZeroExt16to64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (ZeroExt16to64 x)
	// cond:
	// result: (MOVUHreg x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64MOVUHreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpZeroExt32to64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (ZeroExt32to64 x)
	// cond:
	// result: (MOVUWreg x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64MOVUWreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpZeroExt8to16(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (ZeroExt8to16  x)
	// cond:
	// result: (MOVUBreg x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64MOVUBreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpZeroExt8to32(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (ZeroExt8to32  x)
	// cond:
	// result: (MOVUBreg x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64MOVUBreg)
		v.AddArg(x)
		return true
	}
}
func rewriteValueSPARC64_OpZeroExt8to64(v *Value, config *Config) bool {
	b := v.Block
	_ = b
	// match: (ZeroExt8to64  x)
	// cond:
	// result: (MOVUBreg x)
	for {
		x := v.Args[0]
		v.reset(OpSPARC64MOVUBreg)
		v.AddArg(x)
		return true
	}
}
func rewriteBlockSPARC64(b *Block, config *Config) bool {
	switch b.Kind {
	}
	return false
}