package main

var i8 int8
var u8 uint8
var i16 int16
var u16 uint16
var i32 int32
var u32 uint32
var i64 int64
var u64 uint64
var f32 float32
var f64 float64

var w float64

func main() {
	f64 = 16717361816799281152
	u64 = uint64(f64)
	w = float64(u64)
	println(w)
}
