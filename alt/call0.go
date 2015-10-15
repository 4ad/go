package p

func baz(a, b, c int)

func foo() {
	baz(1, 2, 3)
}

func bar(x, y, z int) {
	baz(x, y, z)
}
