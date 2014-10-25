package interval

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"strings"
	"testing"
)

func fmtIntSliceName(name string, is []int) string {
	return name + ":" + fmtIntSlice(is)
}

func fmtIntSlice(is []int) string {
	return "[" + strings.Join(core.IntSliceToStringSlice(is), ",") + "]"
}

func AlldifferentOffset_test(t *testing.T, inits [][]int, names []string,
	offsets []int, exps [][]int, expready bool) {
	msg := "AlldifferentOffset_test %s"

	s := make([]string, len(inits)+1)
	i := 0
	for _, init := range inits {
		s[i] = fmtIntSliceName(names[i], init)
		i = i + 1
	}
	s[i] = fmtIntSliceName("offsets", offsets)

	//msg = fmt.Sprintf(msg+strings.Join(sForm, " "), strings.Join(s, " "))
	msg = fmt.Sprintf(msg, strings.Join(s, " "))
	log(msg)

	vars := createTestVars(inits, names)

	AlldifferentOffset := CreateAlldifferentOffset(
		vars, offsets)
	store.AddPropagators(AlldifferentOffset)
	ready := store.IsConsistent()
	ready_test(t, "AlldifferentOffset", ready, expready)
	if expready {
		for i, exp := range exps {
			domainEquals_test(t, "AlldifferentOffset", vars[i], exp)
		}
	}
}

func Test_AlldifferentOffseta(t *testing.T) {
	setup()
	defer teardown()
	offsets := []int{0, -1, -2, -3}
	AlldifferentOffset_test(t,
		[][]int{[]int{0}, []int{0, 1}, []int{1, 2}, []int{2, 3}},
		[]string{"X", "Y", "Z", "Q"},
		offsets,
		[][]int{[]int{}, []int{}, []int{}, []int{}}, false)
}

func Test_AlldifferentOffsetb(t *testing.T) {
	setup()
	defer teardown()
	offsets := []int{0, -1, -2, -3}
	AlldifferentOffset_test(t,
		[][]int{[]int{0}, []int{0, 1}, []int{1, 2}, []int{1, 2, 3}},
		[]string{"X", "Y", "Z", "Q"},
		offsets,
		[][]int{[]int{}, []int{}, []int{}, []int{}}, false)
}

func Test_AlldifferentOffsetc(t *testing.T) {
	setup()
	defer teardown()
	offsets := []int{0, -1, -2, -3}
	AlldifferentOffset_test(t,
		[][]int{[]int{0}, []int{0, 1}, []int{0, 1, 2}, []int{0, 1, 2, 3}},
		[]string{"X", "Y", "Z", "Q"},
		offsets,
		[][]int{[]int{0}, []int{0}, []int{0}, []int{0}}, true)
}

func Test_AlldifferentOffsetd(t *testing.T) {
	setup()
	defer teardown()
	offsets := []int{0, -1, -2, -3}
	AlldifferentOffset_test(t,
		[][]int{[]int{0, 1}, []int{1}, []int{1, 3}, []int{5}},
		[]string{"X", "Y", "Z", "Q"},
		offsets,
		[][]int{[]int{1}, []int{1}, []int{1}, []int{5}}, true)
}

func Test_AlldifferentOffset_clone(t *testing.T) {
	setup()
	defer teardown()
	log("AlldifferentOffset_clone")

	inits := [][]int{[]int{0, 1}, []int{1}, []int{1, 3}, []int{5}}
	offsets := []int{0, -1, -2, -3}
	names := []string{"X", "Y", "Z", "Q"}

	vars := createTestVars(inits, names)
	c := CreateAlldifferentOffset(vars, offsets)

	clone_test(t, store, c)
}
