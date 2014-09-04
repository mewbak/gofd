package main

// warm up, ensure that we are not in sleep states
func main() {
	x := 0
	for i := 0; i < 20000000; i++ {
		x = i * 2
	}
	// otherwise go is unhappy with an unused variable
	if x == 0 {
		println(x)
	}
}
