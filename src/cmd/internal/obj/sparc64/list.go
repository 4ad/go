// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sparc64

import (
	"cmd/internal/obj"
	"fmt"
)

func init() {
	obj.RegisterRegister(obj.RBaseSPARC64, REG_LAST, Rconv)
	obj.RegisterOpcode(obj.ABaseSPARC64, Anames)
}

func Rconv(r int) string {
	switch {
	case r == RegRFP:
		return "RFP"
	case r == RegLink:
		return "LR"
	case r == RegTMP:
		return "TMP"
	case r == RegGOT:
		return "GOT"
	case r == RegRSP:
		return "RSP"
	case r == RegZero:
		return "ZR"
	case r == REG_BSP:
		return "BSP"
	case r == REG_BFP:
		return "BFP"
	case r == REG_ICC:
		return "ICC"
	case r == REG_XCC:
		return "XCC"
	case r == REG_CCR:
		return "CCR"
	case r == REG_TICK:
		return "TICK"
	case r == REG_RPC:
		return "RPC"
	}
	switch {
	case REG_R0 <= r && r <= REG_R31:
		return fmt.Sprintf("R%d", r-REG_R0)
	case REG_F0 <= r && r <= REG_F31:
		return fmt.Sprintf("F%d", r-REG_F0)
	case REG_D0 <= r && r <= REG_D30 && r%2 == 0:
		return fmt.Sprintf("D%d", r-REG_D0)
	case REG_D32 <= r && r <= REG_D62 && r%2 == 1:
		return fmt.Sprintf("D%d", r-REG_D0+31)
	case REG_FCC0 <= r && r <= REG_FCC3:
		return fmt.Sprintf("FCC%d", r-REG_FCC0)
	}
	return fmt.Sprintf("badreg(%d)", r)
}

func DRconv(a int8) string {
	if a >= ClassUnknown && a <= ClassNone {
		return cnames[a]
	}
	return "C_??"
}
