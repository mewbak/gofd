package core

import (
	"fmt"
	"testing"
)

func domainAddTest(t *testing.T, a []int) {
	var d *ExDomain
	d = CreateExDomain()
	for _, ele := range a {
		if d.Contains(ele) {
			t.Errorf("Domain.Add: d.Contains(%d) = %v, want %v",
				ele, true, false)
		}
	}
	for _, ele := range a {
		d.Add(ele)
	}
	for _, ele := range a {
		if !d.Contains(ele) {
			t.Errorf("Domain.Add: d.Contains(%d) = %v, want %v",
				ele, false, true)
		}
	}
}

func Test_DomainAdd(t *testing.T) {
	setup()
	defer teardown()
	log("DomainAdd")
	a1 := []int{1, 2, 3, 4, 5, 6}
	domainAddTest(t, a1)
	a2 := []int{}
	domainAddTest(t, a2)
	a3 := []int{1, 3, 5, 7}
	domainAddTest(t, a3)
	a4 := []int{3, 7, 1, 5}
	domainAddTest(t, a4)
	a5 := []int{1, -7, 5, -3}
	domainAddTest(t, a5)
}

func domainAddsTest(t *testing.T, a []int) {
	var d *ExDomain
	d = CreateExDomain()
	d.Adds(a)
	for _, ele := range a {
		if !d.Contains(ele) {
			t.Errorf("Domain.Set(%v): d.Contains(%d) = %v, want %v",
				a, ele, false, true)
		}
	}
}

func Test_DomainAdds(t *testing.T) {
	setup()
	defer teardown()
	log("DomainAdds")
	a1 := []int{1, 2, 3, 4, 5, 6}
	domainAddsTest(t, a1[:])
	a2 := []int{}
	domainAddsTest(t, a2[:])
	a3 := []int{1, 3, 5, 7}
	domainAddsTest(t, a3[:])
	a4 := []int{3, 7, 1, 5}
	domainAddsTest(t, a4[:])
	a5 := []int{1, -7, 5, -3}
	domainAddsTest(t, a5[:])
}

func domainMinMaxTest(t *testing.T, a []int, min int, max int) {
	d := CreateExDomainAdds(a)
	cmin := d.GetMin() // should be called Min instead of getMin?
	if cmin != min {
		t.Errorf("Domain(%v): d.GetMin() = %v, want %v", a, cmin, min)
	}
	cmax := d.GetMax() // should be called Min instead of getMin?
	if cmax != max {
		t.Errorf("Domain(%v): d.GetMax() = %v, want %v", a, cmax, max)
	}
	cmin, cmax = d.GetMinAndMax()
	if cmin != min {
		t.Errorf("Domain(%v): d.GetMin() = %v, want %v", a, cmin, min)
	}
	if cmax != max {
		t.Errorf("Domain(%v): d.GetMax() = %v, want %v", a, cmax, max)
	}
}

func Test_DomainMinMax(t *testing.T) {
	setup()
	defer teardown()
	log("DomainMinMax")
	a1 := []int{1, 2, 3, 4, 5, 6}
	domainMinMaxTest(t, a1[:], 1, 6)
	a2 := []int{1, 3, 5, 7}
	domainMinMaxTest(t, a2[:], 1, 7)
	a3 := []int{3, 7, 1, 5}
	domainMinMaxTest(t, a3[:], 1, 7)
	a4 := []int{1, -7, 5, -3}
	domainMinMaxTest(t, a4[:], -7, 5)
	// on an empty domain we panic
}

func domainStringTest(t *testing.T, a []int, exp string) {
	d := CreateExDomainAdds(a)
	got := d.String()
	if got != exp {
		t.Errorf("Domain(%v): d.String() = %v, want %v", a, got, exp)
	}
}

func Test_DomainString(t *testing.T) {
	setup()
	defer teardown()
	log("DomainString")
	d := []int{1, 2, 3, 4, 5, 6}
	domainStringTest(t, d, "[1..6]")
	d = []int{1, 2, 6}
	domainStringTest(t, d, "[1..2,6]")
	d = []int{1, 2, 3, 4, 6}
	domainStringTest(t, d, "[1..4,6]")
	d = []int{1, 2, 3, 4, 7, 8, 9, 15, 17, 19, 20, 21}
	domainStringTest(t, d, "[1..4,7..9,15,17,19..21]")
	d = []int{0, 1, 3, 4, 8, 9, 15, 16, 17, 19, 20, 21, 23}
	domainStringTest(t, d, "[0..1,3..4,8..9,15..17,19..21,23]")
	d = []int{0, 1, 3, 4, 8, 9, 15, 16, 17, 19, 20, 21, 23, 25}
	domainStringTest(t, d, "[0..1,3..4,8..9,15..17,19..21,23,25]")
	d = []int{0, 3, 4}
	domainStringTest(t, d, "[0,3..4]")
	d = []int{0, 3, 5}
	domainStringTest(t, d, "[0,3,5]")
}

func domainEqualsTest(t *testing.T, a []int, b []int, exp bool) {
	da := CreateExDomainAdds(a)
	db := CreateExDomainAdds(b)
	got := da.Equals(db)
	if got != exp {
		t.Errorf("Domain(%v).Equals(Domain(%v)) = %v, want %v",
			da, db, got, exp)
	}
	dc := CreateIvDomainFromIntArr(b)
	got = da.Equals(dc)
	if got != exp {
		t.Errorf("Domain(%v).Equals(Domain(%v)) = %v, want %v",
			da, dc, got, exp)
	}
	dd := CreateIvDomainFromIntArr(a)
	got = db.Equals(dd)
	if got != exp {
		t.Errorf("Domain(%v).Equals(Domain(%v)) = %v, want %v",
			db, dd, got, exp)
	}
}

func Test_DomainEquals(t *testing.T) {
	setup()
	defer teardown()
	log("DomainEquals")
	a1 := []int{1, 2, 3, 4, 5, 6}
	a2 := []int{1, 2, 3, 4, 5, 6}
	domainEqualsTest(t, a1[:], a2[:], true)
	a1 = []int{1, 2, 3, 4, 5, 6}
	a2 = []int{1, 2, 3, 4, 5}
	domainEqualsTest(t, a1[:], a2[:], false)
	a1 = []int{1, 2, 3, 4, 5, 6}
	a2 = []int{1, 2, 3, 4, 5, -6}
	domainEqualsTest(t, a1[:], a2[:], false)
	a1 = []int{}
	a2 = []int{1}
	domainEqualsTest(t, a1[:], a2[:], false)
	a1 = []int{1}
	a2 = []int{}
	domainEqualsTest(t, a1[:], a2[:], false)
	a1 = []int{}
	a2 = []int{}
	domainEqualsTest(t, a1[:], a2[:], true)
	a1 = []int{1, -1}
	a2 = []int{-1, 1}
	domainEqualsTest(t, a1[:], a2[:], true)
}

func domainFromToTest(t *testing.T, da *ExDomain, from, to int) {
	db := CreateExDomainFromTo(from, to)
	if !da.Equals(db) {
		t.Errorf("CreateExDomainFromTo(%v, %v) = a, want b", from, to)
		t.Errorf("  a = %v\n", da)
		t.Errorf("  b = %v\n", db)
	}
}

func Test_DomainFromTo(t *testing.T) {
	setup()
	defer teardown()
	log("DomainFromTo")
	da := CreateExDomainAdds([]int{2, 3, 4, 5})
	domainFromToTest(t, da, 2, 5)
	da = CreateExDomainAdds([]int{2})
	domainFromToTest(t, da, 2, 2)
	da = CreateExDomainAdds([]int{})
	domainFromToTest(t, da, 2, 1)
	da = CreateExDomainAdds([]int{2, 3, 4, 5, 6, 7})
	domainFromToTest(t, da, 2, 7)
	da = CreateExDomainAdds([]int{0, 1, 2, 3, 4, 5, 6})
	domainFromToTest(t, da, 0, 6)
	da = CreateExDomainAdds([]int{})
	domainFromToTest(t, da, 1, -1)
	da = CreateExDomainAdds([]int{-2, -1, 0, 1, 2, 3})
	domainFromToTest(t, da, -2, 3)
}

func Test_DomainSortedValues(t *testing.T) {
	setup()
	defer teardown()
	log("DomainSortedValues")

	domainSortedValuesTest(t, []int{1, 2, 3, 4, 5, 6}, []int{1, 2, 3, 4, 5, 6})
	domainSortedValuesTest(t, []int{1}, []int{1})
	domainSortedValuesTest(t, []int{}, []int{})
	domainSortedValuesTest(t, []int{4, 1, 2, 6, 0, 8, 51}, []int{0, 1, 2, 4, 6, 8, 51})
}

func domainSortedValuesTest(t *testing.T, dvals1 []int, expdvalues []int) {

	if len(dvals1) != len(expdvalues) {
		t.Errorf("SortedValuesTest: Test wrongly called. " +
			" Length of Domain and of expected values must be identical.")
		return
	}

	d1 := CreateExDomainAdds(dvals1)
	d1valsSorted := d1.SortedValues()

	for i := 0; i < len(d1valsSorted); i++ {
		if d1valsSorted[i] != expdvalues[i] {
			t.Errorf("SortedValuesTest: result %v, want %v",
				d1valsSorted, expdvalues)
		}
	}
}

func DomainRemoveTest(t *testing.T, d Domain, eles Domain, dExp Domain) {
	d.Removes(eles)
	if !d.Equals(dExp) {
		t.Errorf("Domain.Removes: d = %s, want %s",
			d, dExp)
	}
}

func convertDomainRemoveTest(t *testing.T,
	from, to int, eles [][]int, expEles [][]int) {
	msg := "Removing "
	for _, ft := range eles {
		msg += fmt.Sprintf("%d..%d ", ft[0], ft[1])
	}
	msg += fmt.Sprintf("from %d..%d", from, to)
	log(msg)
	d := CreateExDomainFromTo(from, to)
	xeles := CreateExDomainAdds(makeTwoDim_OneDim(eles))
	dExp := CreateExDomainAdds(makeTwoDim_OneDim(expEles))
	DomainRemoveTest(t, d, xeles, dExp)
}

func Test_DomainPerformance1(t *testing.T) {
	setup()
	defer teardown()
	log("DomainPerformance1")
	from, to := 1, 50000
	eles := [][]int{{100, 1500}, {10000, 15000}}
	expEles := [][]int{{1, 99}, {1501, 9999}, {15001, 50000}}
	convertDomainRemoveTest(t, from, to, eles, expEles)
}

func Test_DomainPerformance2(t *testing.T) {
	setup()
	defer teardown()
	log("DomainPerformance2")
	from, to := 1, 50000
	eles := [][]int{{100, 1500}, {10000, 15000}}
	expEles := [][]int{{1, 99}, {1501, 9999}, {15001, 50000}}
	convertDomainRemoveTest(t, from, to, eles, expEles)
}

// takes long
func NOTest_DomainPerformance3(t *testing.T) {
	setup()
	defer teardown()
	log("DomainPerformance3")
	from, to := 1, 500000
	eles := [][]int{{100, 1500}, {10000, 15000}, {200000, 200100}}
	expEles := [][]int{{1, 99}, {1501, 9999}, {15001, 199999},
		{200101, 500000}}
	convertDomainRemoveTest(t, from, to, eles, expEles)
}

// takes long
func NOTest_DomainPerformance4(t *testing.T) {
	setup()
	defer teardown()
	log("DomainPerformance4")
	from, to := 1, 500000
	eles := [][]int{{100, 1500}, {10000, 15000}, {200000, 200100}}
	expEles := [][]int{{1, 99}, {1501, 9999}, {15001, 199999},
		{200101, 500000}}
	convertDomainRemoveTest(t, from, to, eles, expEles)
}

func Test_DomainPerformance5(t *testing.T) {
	setup()
	defer teardown()
	log("DomainPerformance5")
	log("Removing ele 9000 from 0..0, 2..2, ..., 10000..10000")
	vals := makeIvDomainWorstCase(10000, -1)
	d := CreateExDomainAdds(vals)
	eles := CreateExDomainAdds([]int{9000, 9000})
	valsExp := makeIvDomainWorstCase(10000, 9000)
	dExp := CreateExDomainAdds(valsExp)
	DomainRemoveTest(t, d, eles, dExp)
}

func Test_GetDomainOutOfBounds(t *testing.T) {
	setup()
	defer teardown()
	log("GetDomainOutOfBounds")

	d := CreateExDomainAdds([]int{1, 2, 3, 4, 5, 10, 11, 12})
	min := 3
	max := 11
	expD := CreateExDomainAdds([]int{1, 2, 12})
	getDomainOutOfBounds_test(t, d, min, max, expD)

	d = CreateExDomainAdds([]int{1, 2, 3, 4, 5, 10, 11, 12})
	min = 0
	max = 13
	expD = CreateExDomain()
	getDomainOutOfBounds_test(t, d, min, max, expD)

	d = CreateExDomainAdds([]int{1, 2, 3, 4, 5, 10, 11, 12})
	min = 1
	max = 5
	expD = CreateExDomainAdds([]int{10, 11, 12})
	getDomainOutOfBounds_test(t, d, min, max, expD)

	d = CreateExDomainAdds([]int{1, 2, 3, 4, 5, 10, 11, 12})
	min = 0
	max = 5
	expD = CreateExDomainAdds([]int{10, 11, 12})
	getDomainOutOfBounds_test(t, d, min, max, expD)

	d = CreateExDomainAdds([]int{1, 2, 3, 4, 5, 10, 11, 12})
	min = 10
	max = 12
	expD = CreateExDomainAdds([]int{1, 2, 3, 4, 5})
	getDomainOutOfBounds_test(t, d, min, max, expD)

	d = CreateExDomainAdds([]int{1, 2, 3, 4, 5, 10, 11, 12})
	min = 10
	max = 13
	expD = CreateExDomainAdds([]int{1, 2, 3, 4, 5})
	getDomainOutOfBounds_test(t, d, min, max, expD)

	d = CreateExDomainAdds([]int{1, 2, 3, 4, 5, 10, 11, 12})
	min = 5
	max = 5
	expD = CreateExDomainAdds([]int{1, 2, 3, 4, 10, 11, 12})
	getDomainOutOfBounds_test(t, d, min, max, expD)
}
