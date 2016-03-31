package foo

func f1() int {
	var x = 7

	if x < 3 {
		return 1
	}
	return 2
}

func f() int {
	var x, y = 5, 6

	if x < y {
		return 1
	}
	return 2
}
