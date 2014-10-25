package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"strings"
	"testing"
)

func fmtIntSlice(is []int) string {
	return "[" + strings.Join(core.IntSliceToStringSlice(is), ",") + "]"
}

func AlldifferentOffset_test(t *testing.T, xinit []int, yinit []int,
	zinit []int, qinit []int, offsets []int, expx []int, expy []int,
	expz []int, expq []int, expready bool) {
	msg := "AlldifferentOffset_test X:%s Y:%s Z:%s Q:%s offsets=%s"
	msg = fmt.Sprintf(msg, fmtIntSlice(xinit), fmtIntSlice(yinit), fmtIntSlice(zinit),
		fmtIntSlice(qinit), fmtIntSlice(offsets))
	log(msg)
	X := core.CreateIntVarExValues("X", store, xinit)
	Y := core.CreateIntVarExValues("Y", store, yinit)
	Z := core.CreateIntVarExValues("Z", store, zinit)
	Q := core.CreateIntVarExValues("Q", store, qinit)
	AlldifferentOffset := CreateAlldifferentOffset(
		[]core.VarId{X, Y, Z, Q}, offsets)
	store.AddPropagators(AlldifferentOffset)
	ready := store.IsConsistent()
	ready_test(t, "AlldifferentOffset", ready, expready)
	if expready {
		domainEquals_test(t, "AlldifferentOffset", X, expx)
		domainEquals_test(t, "AlldifferentOffset", Y, expy)
		domainEquals_test(t, "AlldifferentOffset", Z, expz)
		domainEquals_test(t, "AlldifferentOffset", Q, expq)
	}
}

func Test_AlldifferentOffseta(t *testing.T) {
	setup()
	defer teardown()
	offsets := []int{0, -1, -2, -3}
	AlldifferentOffset_test(t,
		[]int{0}, []int{0, 1}, []int{1, 2}, []int{2, 3}, offsets,
		[]int{}, []int{}, []int{}, []int{}, false)
}

func Test_AlldifferentOffsetb(t *testing.T) {
	setup()
	defer teardown()
	offsets := []int{0, -1, -2, -3}
	AlldifferentOffset_test(t,
		[]int{0}, []int{0, 1}, []int{1, 2}, []int{1, 2, 3}, offsets,
		[]int{}, []int{}, []int{}, []int{}, false)
}

func Test_AlldifferentOffsetc(t *testing.T) {
	setup()
	defer teardown()
	offsets := []int{0, -1, -2, -3}
	AlldifferentOffset_test(t,
		[]int{0}, []int{0, 1}, []int{0, 1, 2}, []int{0, 1, 2, 3}, offsets,
		[]int{0}, []int{0}, []int{0}, []int{0}, true)
}

func Test_AlldifferentOffsetd(t *testing.T) {
	setup()
	defer teardown()
	offsets := []int{0, -1, -2, -3}
	AlldifferentOffset_test(t,
		[]int{0, 1}, []int{1}, []int{1, 3}, []int{5}, offsets,
		[]int{1}, []int{1}, []int{1}, []int{5}, true)
}
