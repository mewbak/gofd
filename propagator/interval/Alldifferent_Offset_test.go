package interval

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"strings"
	"testing"
)

func fmtIntSlice(is []int) string {
	return "[" + strings.Join(core.IntSliceToStringSlice(is), ",") + "]"
}

func Alldifferent_Offset_test(t *testing.T, xinit []int, yinit []int,
	zinit []int, qinit []int, offsets []int, expx []int, expy []int,
	expz []int, expq []int, expready bool) {
	msg := "Alldifferent_Offset_test X:%s Y:%s Z:%s Q:%s offsets=%s"
	msg = fmt.Sprintf(msg, fmtIntSlice(xinit), fmtIntSlice(yinit), fmtIntSlice(zinit),
		fmtIntSlice(qinit), fmtIntSlice(offsets))
	log(msg)
	X := core.CreateIntVarIvValues("X", store, xinit)
	Y := core.CreateIntVarIvValues("Y", store, yinit)
	Z := core.CreateIntVarIvValues("Z", store, zinit)
	Q := core.CreateIntVarIvValues("Q", store, qinit)
	Alldifferent_Offset := CreateAlldifferent_Offset(
		[]core.VarId{X, Y, Z, Q}, offsets)
	store.AddPropagators(Alldifferent_Offset)
	ready := store.IsConsistent()
	ready_test(t, "Alldifferent_Offset", ready, expready)
	if expready {
		domainEquals_test(t, "Alldifferent_Offset", X, expx)
		domainEquals_test(t, "Alldifferent_Offset", Y, expy)
		domainEquals_test(t, "Alldifferent_Offset", Z, expz)
		domainEquals_test(t, "Alldifferent_Offset", Q, expq)
	}
}

func Test_Alldifferent_Offseta(t *testing.T) {
	setup()
	defer teardown()
	offsets := []int{0, -1, -2, -3}
	Alldifferent_Offset_test(t,
		[]int{0}, []int{0, 1}, []int{1, 2}, []int{2, 3}, offsets,
		[]int{}, []int{}, []int{}, []int{}, false)
}

func Test_Alldifferent_Offsetb(t *testing.T) {
	setup()
	defer teardown()
	offsets := []int{0, -1, -2, -3}
	Alldifferent_Offset_test(t,
		[]int{0}, []int{0, 1}, []int{1, 2}, []int{1, 2, 3}, offsets,
		[]int{}, []int{}, []int{}, []int{}, false)
}

func Test_Alldifferent_Offsetc(t *testing.T) {
	setup()
	defer teardown()
	offsets := []int{0, -1, -2, -3}
	Alldifferent_Offset_test(t,
		[]int{0}, []int{0, 1}, []int{0, 1, 2}, []int{0, 1, 2, 3}, offsets,
		[]int{0}, []int{0}, []int{0}, []int{0}, true)
}

func Test_Alldifferent_Offsetd(t *testing.T) {
	setup()
	defer teardown()
	offsets := []int{0, -1, -2, -3}
	Alldifferent_Offset_test(t,
		[]int{0, 1}, []int{1}, []int{1, 3}, []int{5}, offsets,
		[]int{1}, []int{1}, []int{1}, []int{5}, true)
}
