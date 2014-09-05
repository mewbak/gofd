package core

import (
	"fmt"
	"testing"
)

func Test_ScalarSlice(t *testing.T) {
	setup()
	defer teardown()
	log("util_ScalarSlice")
	result := ScalarSlice(5, []int{1, 2, 3, 4, 5})
	expResult := []int{5, 10, 15, 20, 25}
	checkSameSliceResult(t, result, expResult)
	result = ScalarSlice(10, []int{0, 2, 7, 10})
	expResult = []int{0, 20, 70, 100}
	checkSameSliceResult(t, result, expResult)
}

func checkSameSliceResult(t *testing.T, got, want []int) {
	for i := 0; i < len(got); i++ {
		if got[i] != want[i] {
			t.Errorf("util.ScalarSlice: result=%v, want %v", got, want)
			return
		}
	}
}

func Test_AbsInt(t *testing.T) {
	setup()
	defer teardown()
	log("util_AbsInt")
	msg := "AbsInt(%v)"
	checkSameInt(t, msg, -1, AbsInt(-1), 1)
	checkSameInt(t, msg, -3, AbsInt(-3), 3)
	checkSameInt(t, msg, 0, AbsInt(0), 0)
	checkSameInt(t, msg, 1, AbsInt(1), 1)
	checkSameInt(t, msg, 10000, AbsInt(10000), 10000)
}

func checkSameInt(t *testing.T, msg string, param, got, want int) {
	if got != want {
		msg = msg + ": got %d, want %d"
		t.Errorf(msg, param, got, want)
	}
}

func Test_BeautifulOutput(t *testing.T) {
	setup()
	defer teardown()
	log("util_BeautifulOutput")
	checkSameOutput(t, "", []int{2, 3, 5, 7, 11, 13}, "[2..3,5,7,11,13]")
	checkSameOutput(t, "", []int{}, "[]")
	checkSameOutput(t, "", []int{1}, "[1]")
	checkSameOutput(t, "", []int{1, 3}, "[1,3]")
	checkSameOutput(t, "", []int{1, 2}, "[1..2]")
	checkSameOutput(t, "", []int{1, 2, 3, 4}, "[1..4]")
	checkSameOutput(t, "", []int{1, 2, 4, 5}, "[1..2,4..5]")
	checkSameOutput(t, "", []int{1, 2, 4, 5, 7}, "[1..2,4..5,7]")
	checkSameOutput(t, "", []int{1, 2, 4, 6, 7}, "[1..2,4,6..7]")
	checkSameOutput(t, "", []int{1, 3, 4, 6, 7}, "[1,3..4,6..7]")
}

func checkSameOutput(t *testing.T, msg string, a []int, want string) {
	m := make(map[int]bool)
	for _, v := range a {
		m[v] = true
	}
	got := BeautifulOutput(msg, m)
	if got != want {
		message := "BeautifulOutput(%v): got %s, want %s"
		t.Errorf(message, m, got, want)
	} else {
		fmt.Printf("  BO(%v) = %s\n", a, want)
	}
}
