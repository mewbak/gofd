package core

import (
	"fmt"
	"testing"
)

func Test_ValuesOfMapVarIdToIvDomain(t *testing.T) {
	setup()
	defer teardown()
	log("ValuesOfMapVarIdToIvDomain")
	varids := []VarId{1, 2, 3}
	fromTos := [][]int{{1, 10}, {20, 30}, {40, 50}}
	checkValuesOfMapVarIdToIvDomain(t, varids, fromTos)
}

func makeIvDomains(fromTos [][]int) []*IvDomain {
	doms := make([]*IvDomain, len(fromTos))
	for i, fromTo := range fromTos {
		doms[i] = CreateIvDomainFromTo(fromTo[0], fromTo[1])
	}
	return doms
}

func makeExDomains(fromTos [][]int) []*ExDomain {
	doms := make([]*ExDomain, len(fromTos))
	for i, fromTo := range fromTos {
		doms[i] = CreateExDomainFromTo(fromTo[0], fromTo[1])
	}
	return doms
}

func checkValuesOfMapVarIdToIvDomain(t *testing.T,
	varids []VarId, fromTos [][]int) {
	expDoms := makeIvDomains(fromTos)
	m := make(map[VarId]*IvDomain)
	m[varids[0]] = expDoms[0]
	m[varids[1]] = expDoms[1]
	m[varids[2]] = expDoms[2]
	doms := ValuesOfMapVarIdToIvDomain(varids, m)
	checkSameSliceResultIvDomain(t, doms, expDoms)
	expDomsEx := makeExDomains(fromTos)
	mEx := make(map[VarId]*ExDomain)
	mEx[varids[0]] = expDomsEx[0]
	mEx[varids[1]] = expDomsEx[1]
	mEx[varids[2]] = expDomsEx[2]
	doms = ValuesOfMapVarIdToExDomain(varids, mEx)
	checkSameSliceResultExDomain(t, doms, expDomsEx)
}

func Test_ScalarSlice(t *testing.T) {
	setup()
	defer teardown()
	log("ScalarSlice")
	result := ScalarSlice(5, []int{1, 2, 3, 4, 5})
	expResult := []int{5, 10, 15, 20, 25}
	checkSameSliceResult(t, result, expResult)
	result = ScalarSlice(10, []int{0, 2, 7, 10})
	expResult = []int{0, 20, 70, 100}
	checkSameSliceResult(t, result, expResult)
}

func checkSameDomain(t *testing.T, errorS string, got, want Domain) {
	if !got.Equals(want) {
		t.Errorf(errorS, got, want)
	}
}

func checkSameSliceResultExDomain(t *testing.T,
	got []Domain, want []*ExDomain) {
	for i := 0; i < len(got); i++ {
		msg := "checkSameSliceResultExDomain: got %s, want %s"
		checkSameDomain(t, msg, got[i], want[i])
	}
}

func checkSameSliceResultIvDomain(t *testing.T,
	got []Domain, want []*IvDomain) {
	for i := 0; i < len(got); i++ {
		msg := "checkSameSliceResultIvDomain: got %s, want=%s"
		checkSameDomain(t, msg, got[i], want[i])
	}
}

func checkSameSliceResult(t *testing.T, got, want []int) {
	for i := 0; i < len(got); i++ {
		if got[i] != want[i] {
			t.Errorf("checkSameSliceResult: got %v, want %v",
				got, want)
			return
		}
	}
}

func Test_AbsInt(t *testing.T) {
	setup()
	defer teardown()
	log("AbsInt")
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
	log("BeautifulOutput")
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

func Test_IntSliceToStringSlice(t *testing.T) {
	got := IntSliceToStringSlice([]int{1, 5, 10, 100, 1000})
	want := []string{"1", "5", "10", "100", "1000"}
	checkStrings(t, "IntSliceToStringSlice: got %s, want %s", got, want)
	got = IntSliceToStringSlice([]int{0})
	want = []string{"0"}
	checkStrings(t, "IntSliceToStringSlice: got %v, want %v", got, want)
}

func Test_IntSliceToStringSliceFormatted(t *testing.T) {
	got := IntSliceToStringSliceFormatted([]int{1, 5, 10, 100, 1000},
		"test_%d:test")
	want := []string{"test_1:test", "test_5:test", "test_10:test",
		"test_100:test", "test_1000:test"}
	checkStrings(t, "IntSliceToStringSlice: got %v, want %v", got, want)
}

func checkStrings(t *testing.T, msg string, got []string, want []string) {
	if len(got) != len(want) {
		t.Errorf(msg, got, want)
	}
	for i, _ := range got {
		if got[i] != want[i] {
			t.Errorf(msg, got, want)
		}
	}
}

func Test_Keys_MapVarIdsToBool(t *testing.T) {
	values := make(map[VarId]Domain)
	values[1] = CreateIvDomainFromTo(0, 10)
	values[2] = CreateIvDomainFromTo(50, 50)
	gotVarids := Keys_MapVarIdsToBool(values)
	wantVarids := []VarId{1, 2}
	checkVarids(t, "Keys_MapVarIdsToBool: got %v, want %v", gotVarids, wantVarids)
}

func Test_SortedKeys_MapVarIdToInt(t *testing.T) {
	values := make(map[VarId]int)
	values[1] = 5
	values[3] = 10
	values[2] = 50
	gotVarids := SortedKeys_MapVarIdToInt(values)
	wantVarids := []VarId{1, 2, 3}
	checkVarids(t, "SortedKeys_MapVarIdToInt: got %v, want %v",
		gotVarids, wantVarids)
	values = make(map[VarId]int)
	values[1] = 5
	gotVarids = SortedKeys_MapVarIdToInt(values)
	wantVarids = []VarId{1}
	checkVarids(t, "SortedKeys_MapVarIdToInt: got %v, want %v",
		gotVarids, wantVarids)
}

func checkVarids(t *testing.T, msg string, got []VarId, want []VarId) {
	if len(got) != len(want) {
		t.Errorf(msg, got, want)
	}
	for i, _ := range got {
		if got[i] != want[i] {
			t.Errorf(msg, got, want)
		}
	}
}

func Test_SliceToKeys_MapIntToBool(t *testing.T) {
	vals := []int{1, 4, 6, 2}
	m := SliceToKeys_MapIntToBool(vals)
	mapContainsCheck(t, "SliceToKeys_MapIntToBool: have ints %v, got map %v",
		vals, m)
	vals = []int{1}
	m = SliceToKeys_MapIntToBool(vals)
	mapContainsCheck(t, "SliceToKeys_MapIntToBool: have ints %v, got map %v",
		vals, m)
}

func mapContainsCheck(t *testing.T, msg string, vals []int, m map[int]bool) {
	for _, v := range vals {
		if _, ok := m[v]; !ok {
			t.Errorf(msg, vals, m)
		}
	}
}
