package core

import (
	"testing"
)

func ivDomPartCopyTest(t *testing.T, a *IvDomPart) {
	b := a.Copy()
	if !a.Equals(b) {
		t.Errorf("IvDomPart(%v).Copy() = %v",
			a, b)
	}
}

func Test_IvDomPartCopy(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomPartCopy")
	dp := CreateIvDomPart(1, 100)
	ivDomPartCopyTest(t, dp)
	dp = CreateIvDomPart(1, 1)
	ivDomPartCopyTest(t, dp)
}

func ivDomPartEqualsTest(t *testing.T, a *IvDomPart, b *IvDomPart) {
	if !a.Equals(b) {
		t.Errorf("IvDomPart(%v).Equals() = %v",
			a, b)
	}
}

func Test_IvDomPartEquals(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomPartEquals")
	a := CreateIvDomPart(1, 100)
	b := CreateIvDomPart(1, 100)
	ivDomPartEqualsTest(t, a, b)
	a = CreateIvDomPart(1, 1)
	b = CreateIvDomPart(1, 1)
	ivDomPartEqualsTest(t, a, b)
}

func ivDomPartRemoveWithTwoSplitsEleTest(t *testing.T, a *IvDomPart, e int,
	expA []*IvDomPart) {
	a1, a2 := a.removeWithTwoSplits(e)
	expA1 := expA[0]
	expA2 := expA[1]
	if !a1.Equals(expA1) || !a2.Equals(expA2) {
		t.Errorf("IvDomPart(%v) RemoveWithTwoSplitsEle(%v,%v) want (%v,%v)",
			a, a1, a2, expA1, expA2)
	}
}

func Test_IvDomPartRemoveWithTwoSplitsEle(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomPartRemoveWithTwoSplitsEle")
	a := CreateIvDomPart(1, 100)
	e := 5
	expa1 := CreateIvDomPart(1, 4)
	expa2 := CreateIvDomPart(6, 100)
	ivDomPartRemoveWithTwoSplitsEleTest(t, a, e, []*IvDomPart{expa1, expa2})
	a = CreateIvDomPart(1, 100)
	e = 50
	expa1 = CreateIvDomPart(1, 49)
	expa2 = CreateIvDomPart(51, 100)
	ivDomPartRemoveWithTwoSplitsEleTest(t, a, e, []*IvDomPart{expa1, expa2})
}

func ivDomPartRemoveWithTwoSplitsIvDomTest(t *testing.T, a, b *IvDomPart,
	expA []*IvDomPart) {
	a1, a2 := a.removeWithTwoSplitsIvDom(b)
	expA1 := expA[0]
	expA2 := expA[1]
	if !a1.Equals(expA1) || !a2.Equals(expA2) {
		t.Errorf("IvDomPart(%v) RemoveWithTwoSplitsIvDom(%v,%v) want (%v,%v)",
			a, a1, a2, expA1, expA2)
	}
}

func Test_IvDomPartRemoveWithTwoSplitsIvDom(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomPartRemoveWithTwoSplitsIvDom")
	a := CreateIvDomPart(1, 100)
	b := CreateIvDomPart(5, 5)
	expa1 := CreateIvDomPart(1, 4)
	expa2 := CreateIvDomPart(6, 100)
	ivDomPartRemoveWithTwoSplitsIvDomTest(t, a, b, []*IvDomPart{expa1, expa2})
	a = CreateIvDomPart(1, 100)
	b = CreateIvDomPart(5, 50)
	expa1 = CreateIvDomPart(1, 4)
	expa2 = CreateIvDomPart(51, 100)
	ivDomPartRemoveWithTwoSplitsIvDomTest(t, a, b, []*IvDomPart{expa1, expa2})
}

func Test_IvDomPartINTest(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomPartIN")
	a := CreateIvDomPart(1, 100)
	if !a.ContainsInt(5) ||
		!a.ContainsInt(100) ||
		!a.ContainsInt(1) ||
		a.ContainsInt(105) ||
		a.ContainsInt(0) {
		t.Errorf("IvDomPart(%v) failed", a)
	}
}

func ivDomPartRelationCheckDomPartTest(t *testing.T,
	base *IvDomPart, a *IvDomPart, expRes int) {
	res := base.RelationCheck(a)
	if res != expRes {
		t.Errorf("IvDomPart(%s).RelationCheck(%s) want %v got %v",
			base, a, expRes, res)
	}
}

func Test_IvDomPartRelationCheckDomPart(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomPartRelationCheckDomPart")
	base := CreateIvDomPart(1, 100)
	b := CreateIvDomPart(10, 20)
	ivDomPartRelationCheckDomPartTest(t, base, b, CONTAINS_NOT_SAME_FROM_OR_TO)
	b = CreateIvDomPart(1, 20)
	ivDomPartRelationCheckDomPartTest(t, base, b, CONTAINS_SAME_FROM)
	b = CreateIvDomPart(50, 100)
	ivDomPartRelationCheckDomPartTest(t, base, b, CONTAINS_SAME_TO)
	b = CreateIvDomPart(102, 150)
	ivDomPartRelationCheckDomPartTest(t, base, b, NOT_IN_LOWER)
	b = CreateIvDomPart(0, 0)
	ivDomPartRelationCheckDomPartTest(t, base, b, NOT_IN_HIGHER)
	b = CreateIvDomPart(1, 100)
	ivDomPartRelationCheckDomPartTest(t, base, b, SAME)
	b = CreateIvDomPart(0, 200)
	ivDomPartRelationCheckDomPartTest(t, base, b,
		IS_CONTAINED_NOT_SAME_FROM_OR_TO)
	b = CreateIvDomPart(1, 200)
	ivDomPartRelationCheckDomPartTest(t, base, b, IS_CONTAINED_SAME_FROM)
	b = CreateIvDomPart(0, 100)
	ivDomPartRelationCheckDomPartTest(t, base, b, IS_CONTAINED_SAME_TO)
	b = CreateIvDomPart(1, 100)
	ivDomPartRelationCheckDomPartTest(t, base, b, SAME)
}

func ivDomPartRelationCheckIntTest(t *testing.T, base *IvDomPart,
	a int, expRes int) {
	res := base.RelationCheckInt(a)
	if res != expRes {
		t.Errorf("IvDomPart(%s).RelationCheckEle(%v) want %v got %v",
			base, a, expRes, res)
	}
}

func Test_IvDomPartRelationCheckInt(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomPartRelationCheckInt")
	base := CreateIvDomPart(1, 100)
	b := 20
	ivDomPartRelationCheckIntTest(t, base, b, CONTAINS_NOT_SAME_FROM_OR_TO)
	b = 1
	ivDomPartRelationCheckIntTest(t, base, b, CONTAINS_SAME_FROM)
	b = 100
	ivDomPartRelationCheckIntTest(t, base, b, CONTAINS_SAME_TO)
	b = 105
	ivDomPartRelationCheckIntTest(t, base, b, NOT_IN_LOWER)
	base = CreateIvDomPart(50, 50)
	b = 50
	ivDomPartRelationCheckIntTest(t, base, b, SAME)
}

func ivDomPartDIFFERENCE_IntsTest(t *testing.T, baseI, a []int,
	expVals [][]int) {
	base := CreateIvDomPart(baseI[0], baseI[1])
	exp := CreateIvDomParts(expVals)
	res := DIFFERENCE_Ints(base, a...)
	if len(res) != len(exp) {
		msg := "Split of %s with %v result %s: wrong count. "
		msg += "Expected %v, got %v."
		t.Errorf(msg, base, a, res, len(exp), len(res))
		return
	}
	for i := 0; i < len(res); i++ {
		if !res[i].Equals(exp[i]) {
			t.Errorf("Split of %s - at least one part wrong : %s want %s",
				base, res[i], exp[i])
		}
	}
}

func Test_IvDomPartDIFFERENCE_Ints(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomPartDIFFERENCE_Ints")
	ivDomPartDIFFERENCE_IntsTest(t, []int{1, 100}, []int{10},
		[][]int{{1, 9}, {11, 100}})
	ivDomPartDIFFERENCE_IntsTest(t, []int{1, 100}, []int{10, 15, 60},
		[][]int{{1, 9}, {11, 14}, {16, 59}, {61, 100}})
	ivDomPartDIFFERENCE_IntsTest(t, []int{1, 100}, []int{1},
		[][]int{{2, 100}})
	ivDomPartDIFFERENCE_IntsTest(t, []int{1, 100}, []int{100},
		[][]int{{1, 99}})
	ivDomPartDIFFERENCE_IntsTest(t, []int{1, 100}, []int{1, 5, 60},
		[][]int{{2, 4}, {6, 59}, {61, 100}})
	ivDomPartDIFFERENCE_IntsTest(t, []int{1, 100}, []int{5, 60, 100},
		[][]int{{1, 4}, {6, 59}, {61, 99}})
	ivDomPartDIFFERENCE_IntsTest(t, []int{1, 100}, []int{1, 5, 60, 100},
		[][]int{{2, 4}, {6, 59}, {61, 99}})
	ivDomPartDIFFERENCE_IntsTest(t, []int{1, 100}, []int{1, 5, 60, 100, 102},
		[][]int{{2, 4}, {6, 59}, {61, 99}})
	ivDomPartDIFFERENCE_IntsTest(t, []int{1, 100}, []int{102},
		[][]int{{1, 100}})
	ivDomPartDIFFERENCE_IntsTest(t, []int{50, 100}, []int{50, 99, 100},
		[][]int{{51, 98}})
	ivDomPartDIFFERENCE_IntsTest(t, []int{1, 5}, []int{1, 2, 3, 4, 5},
		nil)
	ivDomPartDIFFERENCE_IntsTest(t, []int{1, 100}, []int{0, 1, 2, 3},
		[][]int{{4, 100}})
}

func ivDomPartDIFFERENCE_DomPartsTest(t *testing.T, baseI []int,
	splitDpsI [][]int, expVals [][]int) {
	base := CreateIvDomPart(baseI[0], baseI[1])
	splitPart := CreateIvDomParts(splitDpsI)
	exp := CreateIvDomParts(expVals)
	res := DIFFERENCE_DomParts(base, splitPart...)
	if len(res) != len(exp) {
		t.Errorf("Split of %s - wrong splits-count. Expected %v, got %v",
			base, len(exp), len(res))
		return
	}
	for i := 0; i < len(res); i++ {
		if !res[i].Equals(exp[i]) {
			t.Errorf("Split of %s - at least one part wrong : %s want %s",
				base, res[i], exp[i])
		}
	}
}

func Test_IvDomPartDIFFERENCE_DomParts(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomPartDIFFERENCE_DomParts")
	ivDomPartDIFFERENCE_DomPartsTest(t, []int{1, 100},
		[][]int{{10, 10}}, [][]int{{1, 9}, {11, 100}})
	ivDomPartDIFFERENCE_DomPartsTest(t, []int{1, 100},
		[][]int{{10}, {15}, {60}},
		[][]int{{1, 9}, {11, 14}, {16, 59}, {61, 100}})
	ivDomPartDIFFERENCE_DomPartsTest(t, []int{1, 100},
		[][]int{{1, 5}}, [][]int{{6, 100}})
	ivDomPartDIFFERENCE_DomPartsTest(t, []int{1, 100},
		[][]int{{98, 100}}, [][]int{{1, 97}})
	ivDomPartDIFFERENCE_DomPartsTest(t, []int{1, 100},
		[][]int{{1}, {5}, {60}}, [][]int{{2, 4}, {6, 59}, {61, 100}})
	ivDomPartDIFFERENCE_DomPartsTest(t, []int{1, 100},
		[][]int{{5}, {60}, {100}}, [][]int{{1, 4}, {6, 59}, {61, 99}})
	ivDomPartDIFFERENCE_DomPartsTest(t, []int{1, 100},
		[][]int{{1}, {5}, {60}, {100}}, [][]int{{2, 4}, {6, 59}, {61, 99}})
	ivDomPartDIFFERENCE_DomPartsTest(t, []int{1, 100},
		[][]int{{1}, {5}, {60}, {100}, {102}},
		[][]int{{2, 4}, {6, 59}, {61, 99}})
	ivDomPartDIFFERENCE_DomPartsTest(t, []int{1, 100},
		[][]int{{1, 5}, {20, 30}, {40, 60}},
		[][]int{{6, 19}, {31, 39}, {61, 100}})
	ivDomPartDIFFERENCE_DomPartsTest(t, []int{1, 100},
		[][]int{{1, 5}, {20, 30}, {40, 60}, {70, 100}},
		[][]int{{6, 19}, {31, 39}, {61, 69}})
	ivDomPartDIFFERENCE_DomPartsTest(t, []int{1, 100},
		[][]int{{102, 110}}, [][]int{{1, 100}})
	ivDomPartDIFFERENCE_DomPartsTest(t, []int{1, 5},
		[][]int{{1, 5}}, nil)
	ivDomPartDIFFERENCE_DomPartsTest(t, []int{1, 10},
		[][]int{{1, 5}, {7, 10}}, [][]int{{6, 6}})
	ivDomPartDIFFERENCE_DomPartsTest(t, []int{5, 95},
		[][]int{{0, 4}, {10, 20}, {96, 100}}, [][]int{{5, 9}, {21, 95}})
}

func intersection_test(t *testing.T, p1vals, p2vals, expIntervals []int) {
	p1 := CreateIvDomPart(p1vals[0], p1vals[1])
	p2 := CreateIvDomPart(p2vals[0], p2vals[1])
	var expInter *IvDomPart
	if len(expIntervals) != 0 {
		expInter = CreateIvDomPart(expIntervals[0], expIntervals[1])
	}
	intersectionCheck(t, p1, p2, expInter)
	intersectionCheck(t, p2, p1, expInter)
}

func intersectionCheck(t *testing.T, p1, p2, expInter *IvDomPart) {
	calcInter := p1.INTERSECTION(p2)
	msg := "intersectionCheck : %s intersected with %s results"
	msg += "in %s, expected %s"
	if calcInter != nil {
		if !calcInter.Equals(expInter) {
			t.Errorf(msg, p1, p2, calcInter, expInter)
		}
	} else {
		if expInter != nil {
			t.Errorf(msg, p2, p1, calcInter, nil)
		}
	}
}

func Test_IvDomPartINTERSECTION(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomPartINTERSECTION")
	intersection_test(t, []int{1, 100}, []int{10, 10}, []int{10, 10})
	intersection_test(t, []int{1, 100}, []int{1, 10}, []int{1, 10})
	intersection_test(t, []int{1, 100}, []int{1, 10}, []int{1, 10})
	intersection_test(t, []int{5, 100}, []int{1, 10}, []int{5, 10})
	intersection_test(t, []int{5, 100}, []int{99, 105}, []int{99, 100})
	intersection_test(t, []int{5, 100}, []int{102, 106}, []int{})
	intersection_test(t, []int{5, 100}, []int{1, 2}, []int{})
	intersection_test(t, []int{10, 10}, []int{1, 100}, []int{10, 10})
	intersection_test(t, []int{1, 100}, []int{0, 1}, []int{1, 1})
	intersection_test(t, []int{1, 100}, []int{100, 105}, []int{100, 100})
	intersection_test(t, []int{5, 5}, []int{3, 4}, []int{})
	intersection_test(t, []int{5, 5}, []int{3, 6}, []int{5, 5})
	intersection_test(t, []int{5, 5}, []int{3, 5}, []int{5, 5})
}

func difference_test(t *testing.T, p1vals, p2vals []int,
	expDifference [][]int, reverseExpDifference [][]int) {
	p1 := CreateIvDomPart(p1vals[0], p1vals[1])
	p2 := CreateIvDomPart(p2vals[0], p2vals[1])
	var expDiff []*IvDomPart
	var expRevDiff []*IvDomPart
	if len(expDifference) == 1 {
		expDiff = make([]*IvDomPart, 1)
		expDiff[0] = CreateIvDomPart(expDifference[0][0], expDifference[0][1])
	} else if len(expDifference) == 2 {
		expDiff = make([]*IvDomPart, 2)
		expDiff[0] = CreateIvDomPart(expDifference[0][0], expDifference[0][1])
		expDiff[1] = CreateIvDomPart(expDifference[1][0], expDifference[1][1])
	}
	if len(reverseExpDifference) == 1 {
		expRevDiff = make([]*IvDomPart, 1)
		expRevDiff[0] = CreateIvDomPart(
			reverseExpDifference[0][0], reverseExpDifference[0][1])
	} else if len(reverseExpDifference) == 2 {
		expRevDiff = make([]*IvDomPart, 2)
		expRevDiff[0] = CreateIvDomPart(
			reverseExpDifference[0][0], reverseExpDifference[0][1])
		expRevDiff[1] = CreateIvDomPart(
			reverseExpDifference[1][0], reverseExpDifference[1][1])
	}
	_, calcDiff := p1.DIFFERENCE_NEW(p2)
	_, calcReverseDiff := p2.DIFFERENCE_NEW(p1)
	if calcDiff != nil {
		if len(calcDiff) != len(expDiff) {
			msg := "DifferenceTest: %s differenced with %s wrong "
			msg += "result count. got %v, want %v"
			t.Errorf(msg, p1, p2, len(calcDiff), len(expDiff))
		} else {
			for i, cDiff := range calcDiff {
				if !cDiff.Equals(expDiff[i]) {
					msg := "DifferenceTest: %s differenced with %s results "
					msg += "in %s, expected %s"
					t.Errorf(msg, p1, p2, cDiff, expDiff[i])
				}
			}
		}
	} else {
		if expDiff != nil {
			msg := "DifferenceTest: %s differenced with %s results "
			msg += "in %s, expected not nil"
			t.Errorf(msg, p1, p2, calcDiff)
		}
	}
	if calcReverseDiff != nil {
		if len(calcReverseDiff) != len(expRevDiff) {
			msg := "DifferenceTest: %s differenced with %s wrong "
			msg += "result count. Got %v, want %v"
			t.Errorf(msg, p2, p2, len(calcReverseDiff), len(expRevDiff))
		} else {
			for i, cDiff := range calcReverseDiff {
				if !cDiff.Equals(expRevDiff[i]) {
					msg := "DifferenceTest: %s differenced with %s results "
					msg += "in %s, expected %s"
					t.Errorf(msg, p2, p1, cDiff, expRevDiff[i])
				}
			}
		}
	} else {
		if expRevDiff != nil {
			msg := "DifferenceTest: %s differenced with %s results "
			msg += "in %s, expected not nil"
			t.Errorf(msg, p2, p1, calcReverseDiff)
		}
	}
}

func Test_IvDomPartDIFFERENCE_Parts(t *testing.T) {
	setup()
	defer teardown()
	log("Test_IvDomPartDIFFERENCE_Parts")
	difference_test(t, []int{1, 6}, []int{5, 9},
		[][]int{{1, 4}}, [][]int{{7, 9}})
	difference_test(t, []int{1, 6}, []int{0, 2},
		[][]int{{3, 6}}, [][]int{{0, 0}})
	difference_test(t, []int{1, 6}, []int{7, 9},
		[][]int{{1, 6}}, [][]int{{7, 9}})
	difference_test(t, []int{1, 6}, []int{0, 0},
		[][]int{{1, 6}}, [][]int{{0, 0}})
	difference_test(t, []int{1, 6}, []int{1, 6},
		[][]int{}, [][]int{})
	difference_test(t, []int{1, 6}, []int{3, 4},
		[][]int{{1, 2}, {5, 6}}, [][]int{})
	difference_test(t, []int{1, 6}, []int{2, 2},
		[][]int{{1, 1}, {3, 6}}, [][]int{})
}

func addTest(t *testing.T, p1i, p2i, expPi []int) {
	p1 := CreateIvDomPart(p1i[0], p1i[1])
	p2 := CreateIvDomPart(p2i[0], p2i[1])
	expP := CreateIvDomPart(expPi[0], expPi[1])
	addCheck(t, p1, p2, expP)
	addCheck(t, p2, p1, expP)
}

func addCheck(t *testing.T, p1, p2, expP *IvDomPart) {
	calcP := p1.ADD(p2)
	if !calcP.Equals(expP) {
		msg := "addCheck: %s + %s results in %s, expected %s"
		t.Errorf(msg, p1, p2, calcP, expP)
	}
}

func Test_IvDomPartADD(t *testing.T) {
	setup()
	defer teardown()
	log("Test_IvDomPartADD")
	addTest(t, []int{1, 6}, []int{1, 6}, []int{2, 12})
	addTest(t, []int{0, 1}, []int{0, 1}, []int{0, 2})
	addTest(t, []int{0, 0}, []int{0, 0}, []int{0, 0})
	addTest(t, []int{0, 2}, []int{4, 5}, []int{4, 7})
}

func subtractTest(t *testing.T, p1i, p2i, expPi, expPi2 []int) {
	p1 := CreateIvDomPart(p1i[0], p1i[1])
	p2 := CreateIvDomPart(p2i[0], p2i[1])
	expP := CreateIvDomPart(expPi[0], expPi[1])
	subtractCheck(t, p1, p2, expP)
	expP = CreateIvDomPart(expPi2[0], expPi2[1])
	subtractCheck(t, p2, p1, expP)
}

func subtractCheck(t *testing.T, p1, p2, expP *IvDomPart) {
	calcP := p1.SUBTRACT(p2)
	if !calcP.Equals(expP) {
		msg := "subtractCheck: %s - %s results in %s, expected %s"
		t.Errorf(msg, p2, p1, calcP, expP)
	}
}

func Test_IvDomPartSUBTRACT(t *testing.T) {
	setup()
	defer teardown()
	log("Test_IvDomPartSUBTRACT")
	subtractTest(t, []int{1, 6}, []int{1, 6}, []int{-5, 5}, []int{-5, 5})
	subtractTest(t, []int{0, 1}, []int{0, 1}, []int{-1, 1}, []int{-1, 1})
	subtractTest(t, []int{0, 0}, []int{0, 0}, []int{0, 0}, []int{0, 0})
	subtractTest(t, []int{1, 1}, []int{1, 1}, []int{0, 0}, []int{0, 0})
	subtractTest(t, []int{0, 2}, []int{4, 5}, []int{-5, -2}, []int{2, 5})
}

func neg_test(t *testing.T, from, to, expFrom, expTo int) {
	p := CreateIvDomPart(from, to)
	expP := CreateIvDomPart(expFrom, expTo)
	calcP := p.NEG()
	if !calcP.Equals(expP) {
		msg := "neg_test: neg(%s) results in %s, expected %s"
		t.Errorf(msg, p, calcP, expP)
	}
}

func Test_IvDomPartNEG(t *testing.T) {
	setup()
	defer teardown()
	log("Test_IvDomPartNEG")
	neg_test(t, -6, -5, 5, 6)
	neg_test(t, -5, -5, 5, 5)
	neg_test(t, 0, 1, -1, 0)
	neg_test(t, 5, 5, -5, -5)
	neg_test(t, 5, 6, -6, -5)
	neg_test(t, -5, 5, -5, 5)
	neg_test(t, -4, 5, -5, 4)
	neg_test(t, -6, 5, -5, 6)
}

func abs_test(t *testing.T, from, to, expFrom, expTo int) {
	p := CreateIvDomPart(from, to)
	expP := CreateIvDomPart(expFrom, expTo)
	calcP := p.ABS()
	if !calcP.Equals(expP) {
		msg := "abs_test: abs(%s) results in %s, expected %s"
		t.Errorf(msg, p, calcP, expP)
	}
}

func Test_IvDomPartABS(t *testing.T) {
	setup()
	defer teardown()
	log("Test_IvDomPartABS")
	abs_test(t, -6, -3, 3, 6)
	abs_test(t, -5, 5, 0, 5)
	abs_test(t, -4, 5, 0, 5)
	abs_test(t, -6, 5, 0, 6)
	abs_test(t, 3, 5, 3, 5)
	abs_test(t, 0, 0, 0, 0)
	abs_test(t, -1, 0, 0, 1)
	abs_test(t, 0, 1, 0, 1)
}
