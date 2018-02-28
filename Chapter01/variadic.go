package main

func simpleVariadicToSlice(numbers ...int) []int {
	return numbers
}

func mixedVariadicToSlice(name string, numbers ...int) (string, []int) {
	return name, numbers
}

// Does not work.
// func badVariadic(name ...string, numbers ...int) {}
