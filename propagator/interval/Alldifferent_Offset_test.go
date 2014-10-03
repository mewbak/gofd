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

func Alldifferent_Offset_test(t *testing.T, inits [][]int, names []string,
	offsets []int, exps [][]int, expready bool) {
	msg := "Alldifferent_Offset_test "
	sForm := make([]string, len(inits)+1)

	s := make([]string, len(inits)+1)
	i := 0
	for _, init := range inits {
		s[i] = fmtIntSliceName(names[i], init)
		sForm[i] = "%s"
		i = i + 1
	}
	s[i] = fmtIntSlice(offsets)
	sForm[i] = "offsets=%s"

	msg = fmt.Sprintf(msg+strings.Join(sForm, " "), strings.Join(s, " "))
	log(msg)

	vars := createTestVars(inits, names)

	Alldifferent_Offset := CreateAlldifferent_Offset(
		vars, offsets)
	store.AddPropagators(Alldifferent_Offset)
	ready := store.IsConsistent()
	ready_test(t, "Alldifferent_Offset", ready, expready)
	if expready {
		for i, exp := range exps {
			domainEquals_test(t, "Alldifferent_Offset", vars[i], exp)
		}
	}
}

func Test_Alldifferent_Offseta(t *testing.T) {
	setup()
	defer teardown()
	offsets := []int{0, -1, -2, -3}
	Alldifferent_Offset_test(t,
		[][]int{[]int{0}, []int{0, 1}, []int{1, 2}, []int{2, 3}},
		[]string{"X", "Y", "Z", "Q"},
		offsets,
		[][]int{[]int{}, []int{}, []int{}, []int{}}, false)
}

func Test_Alldifferent_Offsetb(t *testing.T) {
	setup()
	defer teardown()
	offsets := []int{0, -1, -2, -3}
	Alldifferent_Offset_test(t,
		[][]int{[]int{0}, []int{0, 1}, []int{1, 2}, []int{1, 2, 3}},
		[]string{"X", "Y", "Z", "Q"},
		offsets,
		[][]int{[]int{}, []int{}, []int{}, []int{}}, false)
}

func Test_Alldifferent_Offsetc(t *testing.T) {
	setup()
	defer teardown()
	offsets := []int{0, -1, -2, -3}
	Alldifferent_Offset_test(t,
		[][]int{[]int{0}, []int{0, 1}, []int{0, 1, 2}, []int{0, 1, 2, 3}},
		[]string{"X", "Y", "Z", "Q"},
		offsets,
		[][]int{[]int{0}, []int{0}, []int{0}, []int{0}}, true)
}

func Test_Alldifferent_Offsetd(t *testing.T) {
	setup()
	defer teardown()
	offsets := []int{0, -1, -2, -3}
	Alldifferent_Offset_test(t,
		[][]int{[]int{0, 1}, []int{1}, []int{1, 3}, []int{5}},
		[]string{"X", "Y", "Z", "Q"},
		offsets,
		[][]int{[]int{1}, []int{1}, []int{1}, []int{5}}, true)
}

func Test_Alldifferent_Offset_clone(t *testing.T) {
	setup()
	defer teardown()
	log("Alldifferent_Offset_clone")

	inits := [][]int{[]int{0, 1}, []int{1}, []int{1, 3}, []int{5}}
	offsets := []int{0, -1, -2, -3}
	names := []string{"X", "Y", "Z", "Q"}

	vars := createTestVars(inits, names)
	c := CreateAlldifferent_Offset(vars, offsets)

	clone_test(t, store, c)
}
