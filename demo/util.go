package demo

// sum_intarray computes the sum of all elements in an []int
func sum_intarray(a []int) int {
	res := 0
	for _, ele := range a {
		res += ele
	}
	return res
}
