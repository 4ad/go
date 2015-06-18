package sparc64

func rclass(r int16) int {
	switch {
	case r == RegZero:
		return ClassZero
	case REG_R1 <= r && r <= REG_R31:
		return ClassReg
	case REG_F0 <= r && r <= REG_F31:
		return ClassFloatReg
	case r == REG_BSP || r == REG_BFP:
		return ClassBiased
	}
	return ClassUnknown
}
