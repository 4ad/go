package main

func main() {
	f32 := float32(10)
	i64 := int64(f32)
	u32 := uint32(10)
	f64 := float64(u32)
	i8 := int8(42)
	f64 = float64(i8)

	_ = i8
	_ = f64
	_ = i64
}
