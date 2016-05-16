// +build arm64

package runtime

import "unsafe"

func unimpl(name string) {
	print("UNIMPLEMENTED ", name, "!\n")
	exit(42)
}

// NOTE: please do not trust the prototype defined in this file.
// Always lookup the real prototype in the original runtime.

func panicindex()                            { unimpl("panicindex") }
func panicslice()                            { unimpl("panicslice") }
func panicdivide()                           { unimpl("panicdivide") }
func throwreturn()                           { unimpl("throwreturn") }
func throwinit()                             { unimpl("throwinit") }
func panicwrap(_ string, _ string, _ string) { unimpl("panicwrap") }
func gopanic(_ interface{})                  { unimpl("gopanic") }
func gorecover(_ *int32)/*(_ interface{})*/ { unimpl("gorecover") }
func concatstring2(_ string, _ string)/*(_ string)*/ { unimpl("concatstring2") }
func concatstring3(_ string, _ string, _ string)/*(_ string)*/ { unimpl("concatstring3") }
func concatstring4(_ string, _ string, _ string, _ string)/*(_ string)*/ { unimpl("concatstring4") }
func concatstring5(_ string, _ string, _ string, _ string, _ string) /*(_ string)*/ {
	unimpl("concatstring5")
}
func concatstrings(_ []string)/*(_ string)*/ { unimpl("concatstrings") }
func cmpstring(_ string, _ string)/*(_ int)*/ { unimpl("cmpstring") }
func eqstring(_ string, _ string)/*(_ bool)*/ { unimpl("eqstring") }
func intstring(_ int64)/*(_ string)*/ { unimpl("intstring") }
func slicebytetostring(_ []byte)/*(_ string)*/ { unimpl("slicebytetostring") }
func slicebytetostringtmp(_ []byte)/*(_ string)*/ { unimpl("slicebytetostringtmp") }
func slicerunetostring(_ []rune)/*(_ string)*/ { unimpl("slicerunetostring") }
func stringtoslicebyte(_ string)/*(_ []byte)*/ { unimpl("stringtoslicebyte") }
func stringtoslicerune(_ string)/*(_ []rune)*/ { unimpl("stringtoslicerune") }
func stringiter(_ string, _ int)/*(_ int)*/ { unimpl("stringiter") }
func stringiter2(_ string, _ int)/*(retk int, retv rune)*/ { unimpl("stringiter2") }
func slicecopy(to sliceStruct, fr sliceStruct, wid uintptr)/*(_ int)*/ { unimpl("slicecopy") }
func slicestringcopy(to []byte, fr string)/*(_ int)*/ { unimpl("slicestringcopy") }
func typ2Itab(typ *byte, typ2 *byte, cache **byte)/*(ret *byte)*/ { unimpl("typ2Itab") }
func convI2E(elem fInterface)/*(ret interface{})*/ { unimpl("convI2E") }
func convI2I(typ *interfacetype, elem fInterface)/*(ret fInterface)*/ { unimpl("convI2I") }
func convT2E(typ *byte, elem unsafe.Pointer)/*(ret interface{})*/ { unimpl("convT2E") }
func convT2I(typ *byte, typ2 *byte, cache **byte, elem unsafe.Pointer) /*(ret fInterface)*/ {
	unimpl("convT2I")
}
func assertE2E(typ *byte, iface interface{}, ret *interface{}) { unimpl("assertE2E") }
func assertE2E2(typ *byte, iface interface{}, ret *interface{})/*(_ bool)*/ { unimpl("assertE2E2") }
func assertE2I(typ *byte, iface interface{}, ret *fInterface) { unimpl("assertE2I") }
func assertE2I2(typ *byte, iface interface{}, ret *fInterface)/*(_ bool)*/ { unimpl("assertE2I2") }
func assertE2T(typ *byte, iface interface{}, ret unsafe.Pointer) { unimpl("assertE2T") }
func assertE2T2(typ *byte, iface interface{}, ret unsafe.Pointer)/*(_ bool)*/ { unimpl("assertE2T2") }
func assertI2E(typ *byte, iface interface{}, ret unsafe.Pointer) { unimpl("assertI2E") }
func assertI2E2(typ *byte, iface interface{}, ret unsafe.Pointer)/*(_ bool)*/ { unimpl("assertI2E2") }
func assertI2I(typ *byte, iface interface{}, ret unsafe.Pointer) { unimpl("assertI2I") }
func assertI2I2(typ *byte, iface interface{}, ret unsafe.Pointer)/*(_ bool)*/ { unimpl("assertI2I2") }
func assertI2T(typ *byte, iface interface{}, ret unsafe.Pointer) { unimpl("assertI2T") }
func assertI2T2(typ *byte, iface interface{}, ret unsafe.Pointer)/*(_ bool)*/ { unimpl("assertI2T2") }
func ifaceeq(i1 fInterface, i2 fInterface)/*(ret bool)*/ { unimpl("ifaceeq") }
func efaceeq(i1 interface{}, i2 interface{})/*(ret bool)*/ { unimpl("efaceeq") }
func ifacethash(i1 fInterface)/*(ret uint32)*/ { unimpl("ifacethash") }
func efacethash(i1 interface{})/*(ret uint32)*/ { unimpl("efacethash") }
func makemap(mapType *byte, hint int64)/*(hmap unsafe.Pointer)*/ { unimpl("makemap") }
func mapaccess1(mapType *byte, hmap unsafe.Pointer, key unsafe.Pointer) /*(val unsafe.Pointer)*/ {
	unimpl("mapaccess1")
}
func mapaccess1_fast32(mapType *byte, hmap unsafe.Pointer, key uint32) /*(val unsafe.Pointer)*/ {
	unimpl("mapaccess1_fast32")
}
func mapaccess1_fast64(mapType *byte, hmap unsafe.Pointer, key uint64) /*(val unsafe.Pointer)*/ {
	unimpl("mapaccess1_fast64")
}
func mapaccess1_faststr(mapType *byte, hmap unsafe.Pointer, key string) /*(val unsafe.Pointer)*/ {
	unimpl("mapaccess1_faststr")
}
func mapaccess2(mapType *byte, hmap unsafe.Pointer, key unsafe.Pointer) /*(val unsafe.Pointer, pres bool)*/ {
	unimpl("mapaccess2")
}
func mapaccess2_fast32(mapType *byte, hmap unsafe.Pointer, key uint32) /*(val unsafe.Pointer, pres bool)*/ {
	unimpl("mapaccess2_fast32")
}
func mapaccess2_fast64(mapType *byte, hmap unsafe.Pointer, key uint64) /*(val unsafe.Pointer, pres bool)*/ {
	unimpl("mapaccess2_fast64")
}
func mapaccess2_faststr(mapType *byte, hmap unsafe.Pointer, key string) /*(val unsafe.Pointer, pres bool)*/ {
	unimpl("mapaccess2_faststr")
}
func mapassign1(mapType *byte, hmap unsafe.Pointer, key unsafe.Pointer, val unsafe.Pointer) {
	unimpl("mapassign1")
}
func mapiterinit(mapType *byte, hmap unsafe.Pointer, hiter unsafe.Pointer) { unimpl("mapiterinit") }
func mapdelete(mapType *byte, hmap unsafe.Pointer, key unsafe.Pointer)     { unimpl("mapdelete") }
func mapiternext(hiter unsafe.Pointer)                                     { unimpl("mapiternext") }
func makechan(chanType *byte, hint int64)/*(hchan unsafe.Pointer)*/ { unimpl("makechan") }
func chanrecv1(chanType *byte, hchan unsafe.Pointer, elem unsafe.Pointer) { unimpl("chanrecv1") }
func chanrecv2(chanType *byte, hchan unsafe.Pointer, elem unsafe.Pointer) /*(_ bool)*/ {
	unimpl("chanrecv2")
}
func chansend1(chanType *byte, hchan unsafe.Pointer, elem unsafe.Pointer) { unimpl("chansend1") }
func closechan(hchan unsafe.Pointer)                                      { unimpl("closechan") }

func selectnbsend(chanType *byte, hchan unsafe.Pointer, elem unsafe.Pointer) /*(_ bool)*/ {
	unimpl("selectnbsend")
}
func selectnbrecv(chanType *byte, elem unsafe.Pointer, hchan unsafe.Pointer) /*(_ bool)*/ {
	unimpl("selectnbrecv")
}
func selectnbrecv2(chanType *byte, elem unsafe.Pointer, received *bool, hchan unsafe.Pointer) /*(_ bool)*/ {
	unimpl("selectnbrecv2")
}
func newselect(sel *byte, selsize int64, size int32) { unimpl("newselect") }
func selectsend(sel *byte, hchan unsafe.Pointer, elem unsafe.Pointer) /*(selected bool)*/ {
	unimpl("selectsend")
}
func selectrecv(sel *byte, hchan unsafe.Pointer, elem unsafe.Pointer) /*(selected bool)*/ {
	unimpl("selectrecv")
}
func selectrecv2(sel *byte, hchan unsafe.Pointer, elem unsafe.Pointer, received *bool) /*(selected bool)*/ {
	unimpl("selectrecv2")
}
func selectdefault(sel *byte)/*(selected bool)*/ { unimpl("selectdefault") }
func selectgo(sel *byte) { unimpl("selectgo") }
func block()             { unimpl("block") }
func makeslice(typ *byte, nel int64, cap int64)/*(ary sliceStruct)*/ { unimpl("makeslice") }
func growslice(typ *byte, old sliceStruct, n int64)/*(ary sliceStruct)*/ { unimpl("growslice") }
func memequal(x unsafe.Pointer, y unsafe.Pointer, size uintptr)/*(_ bool)*/ { unimpl("memequal") }
func memequal8(x unsafe.Pointer, y unsafe.Pointer)/*(_ bool)*/ { unimpl("memequal8") }
func memequal16(x unsafe.Pointer, y unsafe.Pointer)/*(_ bool)*/ { unimpl("memequal16") }
func memequal32(x unsafe.Pointer, y unsafe.Pointer)/*(_ bool)*/ { unimpl("memequal32") }
func memequal64(x unsafe.Pointer, y unsafe.Pointer)/*(_ bool)*/ { unimpl("memequal64") }
func memequal128(x unsafe.Pointer, y unsafe.Pointer)/*(_ bool)*/ { unimpl("memequal128") }
func int64div(_ int64, _ int64)/*(_ int64)*/ { unimpl("int64div") }
func uint64div(_ uint64, _ uint64)/*(_ uint64)*/ { unimpl("uint64div") }
func int64mod(_ int64, _ int64)/*(_ int64)*/ { unimpl("int64mod") }
func uint64mod(_ uint64, _ uint64)/*(_ uint64)*/ { unimpl("uint64mod") }
func float64toint64(_ float64)/*(_ int64)*/ { unimpl("float64toint64") }
func float64touint64(_ float64)/*(_ uint64)*/ { unimpl("float64touint64") }
func int64tofloat64(_ int64)/*(_ float64)*/ { unimpl("int64tofloat64") }
func uint64tofloat64(_ uint64)/*(_ float64)*/ { unimpl("uint64tofloat64") }
func complex128div(num complex128, den complex128)/*(quo complex128)*/ { unimpl("complex128div") }
func racefuncenter(_ uintptr)                   { unimpl("racefuncenter") }
func racefuncexit()                             { unimpl("racefuncexit") }
func raceread(_ uintptr)                        { unimpl("raceread") }
func racewrite(_ uintptr)                       { unimpl("racewrite") }
func racereadrange(addr uintptr, size uintptr)  { unimpl("racereadrange") }
func racewriterange(addr uintptr, size uintptr) { unimpl("racewriterange") }

func newproc() { unimpl("newproc") }
