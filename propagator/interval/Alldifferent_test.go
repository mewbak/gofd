package interval

import (
	"testing"
)

func alldiffinterval_test(t *testing.T, inits [][]int, names []string,
	exps [][]int, expready bool) {

	vars := createTestVars(inits, names)

	store.AddPropagators(CreateAlldifferent(vars...))
	ready := store.IsConsistent()
	ready_test(t, "Alldifferent_intervals", ready, expready)
	if expready {
		for i, exp := range exps {
			domainEquals_test(t, "Alldifferent_intervals", vars[i], exp)
		}
	}
}

func Test_Alldifferenta(t *testing.T) {
	setup()
	defer teardown()
	log("Alldifferenta_interval: X:0, Y:0..1, Z:1..2, Q:2..3")
	alldiffinterval_test(t,
		[][]int{[]int{0}, []int{0, 1}, []int{1, 2}, []int{2, 3}},
		[]string{"X", "Y", "Z", "Q"},
		[][]int{[]int{0}, []int{1}, []int{2}, []int{3}}, true)
}

func Test_Alldifferentb(t *testing.T) {
	setup()
	defer teardown()
	log("Alldifferentb_interval: X:0..1, Y:1, Z:2..3, Q:3")
	alldiffinterval_test(t,
		[][]int{[]int{0, 1}, []int{1}, []int{2, 3}, []int{3}},
		[]string{"X", "Y", "Z", "Q"},
		[][]int{[]int{0}, []int{1}, []int{2}, []int{3}}, true)
}

func Test_Alldifferentc(t *testing.T) {
	setup()
	defer teardown()
	log("Alldifferentc_interval: X:0, Y:1, Z:2, Q:3")
	alldiffinterval_test(t,
		[][]int{[]int{0}, []int{1}, []int{2}, []int{3}},
		[]string{"X", "Y", "Z", "Q"},
		[][]int{[]int{0}, []int{1}, []int{2}, []int{3}}, true)
}

func Test_Alldifferentd(t *testing.T) {
	setup()
	defer teardown()
	log("Alldifferentd_interval: X:0, Y:0, Z:0, Q:0")
	alldiffinterval_test(t,
		[][]int{[]int{0}, []int{0}, []int{0}, []int{0}},
		[]string{"X", "Y", "Z", "Q"},
		[][]int{[]int{}, []int{}, []int{}, []int{}}, false)
}

func Test_Alldifferent_clone(t *testing.T) {
	setup()
	defer teardown()
	log("Alldifferent_clone")

	inits := [][]int{[]int{0, 1}, []int{1}, []int{2, 3}, []int{3}}
	names := []string{"X", "Y", "Z", "Q"}

	vars := createTestVars(inits, names)
	c := CreateAlldifferent(vars...)

	clone_test(t, store, c)
}
