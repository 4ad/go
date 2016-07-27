// run

package main

func main() {
	const (
		I = 103
		N = 3
	)
	type T [N]int
	m := make(map[T]int)
	for i := 0; i < I; i++ {
		var v T
		for j := 0; j < N; j++ {
			v[j] = i + j
		}
		m[v] = i
	}

	m[T{104, 105, 106}] = 104

	t0 := T{100, 101, 102}
	v := m[t0]
	if v != 100 {
		panic("map lookup failure")
	}

	m[T{105, 106, 107}] = 105

	v = m[t0]
	if v != 100 {
		panic("map lookup failure")
	}
}
