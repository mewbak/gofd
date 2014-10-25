package core

import (
	"fmt"
	"testing"
)

func Test_CreateIvDomainFromIntArr(t *testing.T) {
	setup()
	defer teardown()
	log("CreateIvDomainFromIntArr")
	vals := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	expVals := [][]int{{1, 9}}
	da := CreateIvDomainFromIntArr(vals)
	partsTest(t, "CreateIvDomainFromIntArr-", da, expVals, 1, 9)
	vals = []int{6, 1, 2, 8, 4, 7}
	expVals = [][]int{{1, 2}, {4, 4}, {6, 8}}
	da = CreateIvDomainFromIntArr(vals)
	partsTest(t, "CreateIvDomainFromIntArr-", da, expVals, 1, 8)
}

func partsTest(t *testing.T, test string, da *IvDomain,
	vals [][]int, min int, max int) {
	parts := da.GetParts()
	if len(vals) != len(parts) {
		t.Errorf(test+"IvDomain %s don't match with param %v\n", da, vals)
	}
	for i := 0; i < len(vals); i++ {
		if parts[i].From != vals[i][0] {
			t.Errorf(test+"IvDomain %s wrong. Part%v.From = %v, want %v \n",
				da, i, parts[i].From, vals[i][0])
		}
		if parts[i].To != vals[i][1] {
			t.Errorf(test+"IvDomain %s wrong. Part%v.To = %v, want %v \n",
				da, i, parts[i].To, vals[i][1])
		}
	}
	minTest(t, test, da, min)
	maxTest(t, test, da, max)
}

func Test_CreateIvDomainFromTos(t *testing.T) {
	setup()
	defer teardown()
	log("CreateIvDomainFromTos")
	vals := [][]int{{1, 100}, {200, 300}}
	da := CreateIvDomainFromTos(vals)
	partsTest(t, "CreateIvDomainFromTos-", da, vals, 1, 300)
	vals = [][]int{{5, 50}}
	da = CreateIvDomainFromTos(vals)
	partsTest(t, "CreateIvDomainFromTos-", da, vals, 5, 50)
}

func CreateIvDomainTest(t *testing.T, da *IvDomain) {
	if !da.IsEmpty() {
		t.Errorf("IvDomain.IsEmpty()= %v want %v",
			da.IsEmpty(), true)
	}
}

func Test_CreateIvDomain(t *testing.T) {
	setup()
	defer teardown()
	log("CreateIvDomain")
	da := CreateIvDomain()
	CreateIvDomainTest(t, da)
	da = CreateIvDomain()
	CreateIvDomainTest(t, da)
}

func CreateIvDomainFromToTest(t *testing.T, da *IvDomain, from, to int) {
	db := CreateIvDomainFromTo(from, to)
	if !da.Equals(db) {
		t.Errorf("CreateIvDomainFromTo(%v, %v) = a, want b", from, to)
		t.Errorf("  a = %v\n", da)
		t.Errorf("  b = %v\n", db)
	}
	parts := da.GetParts()
	if parts[0].From != from {
		t.Errorf("  IvDomain %s wrong. Part0.From = %v, want %v \n", da,
			parts[0].From, from)
	}
	if parts[0].To != to {
		t.Errorf("  IvDomain %s wrong. Part0.To = %v, want %v \n", da,
			parts[0].To, to)
	}
}

func Test_CreateIvDomainFromTo(t *testing.T) {
	setup()
	defer teardown()
	log("CreateIvDomainFromTo")
	da := CreateIvDomainFromTo(1, 100)
	CreateIvDomainFromToTest(t, da, 1, 100)
	da = CreateIvDomainFromTo(5, 50)
	CreateIvDomainFromToTest(t, da, 5, 50)
}

func IvDomainMinMaxTest(t *testing.T, min int, max int) {
	d := CreateIvDomainFromTo(min, max)
	minTest(t, "IvDomainMinMax", d, min)
	maxTest(t, "IvDomainMinMax", d, max)
}

func Test_IvDomainMinMax(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainMinMax")
	IvDomainMinMaxTest(t, 1, 6)
	IvDomainMinMaxTest(t, 1, 50)
	IvDomainMinMaxTest(t, 50, 150)
	IvDomainMinMaxTest(t, -7, 5)
}

func IvDomainRemovesTest(t *testing.T, dom *IvDomain, eles []int,
	expDom *IvDomain) {
	elesDom := CreateIvDomainFromIntArr(eles)
	dom.Removes(elesDom)
	if !dom.Equals(expDom) {
		t.Errorf("IvDomain: d = %v, want %v", dom, expDom)
	}
}

// dom and expDom must have the same type
func IvDomainRemovesTestDomain(t *testing.T, dom Domain,
	removingDom Domain, expDom Domain) {
	dom.Removes(removingDom)
	if !dom.Equals(expDom) {
		t.Errorf("IvDomain: d = %v, want %v", dom, expDom)
	}
}

func intersectionIvD_DTest(t *testing.T, kind string, da *IvDomain,
	db *ExDomain, expIntersectionPerPart map[int][]int) {
	intersectionPerPart := da.IntersectionWithDomainThisAssoc(db)
	errmsg := "%s-IntersectionCheck: calculated intersection between %s"
	errmsg += " and %s: %v, want %v"
	if len(intersectionPerPart) != len(expIntersectionPerPart) {
		t.Errorf(errmsg, kind, da, db,
			intersectionPerPart, expIntersectionPerPart)
	}
	for i, expIntersectionPerPartI := range expIntersectionPerPart {
		if len(intersectionPerPart[i]) != len(expIntersectionPerPart[i]) {
			t.Errorf(errmsg, kind, da, db,
				intersectionPerPart, expIntersectionPerPart)
		}
		for j, _ := range intersectionPerPart[i] {
			if expIntersectionPerPartI[j] != intersectionPerPart[i][j] {
				t.Errorf(errmsg, kind, da, db,
					intersectionPerPart, expIntersectionPerPart)
			}
		}
	}
}

func intersectionIvD_IvDTest(t *testing.T, kind string,
	da *IvDomain, db *IvDomain, expIntersectionPerPart map[int][]*IvDomPart) {
	intersectionPerPart := da.IntersectionWithIvDomainThisAssoc(db)
	errmsg := "%s-IntersectionCheck: calculated intersection between %s"
	errmsg += " and %s: %v, want %v"
	if len(intersectionPerPart) != len(expIntersectionPerPart) {
		t.Errorf(errmsg, kind, da, db,
			intersectionPerPart, expIntersectionPerPart)
	}
	for i, expIntersectionPerPartI := range expIntersectionPerPart {
		if len(intersectionPerPart[i]) != len(expIntersectionPerPart[i]) {
			t.Errorf(errmsg, kind, da, db,
				intersectionPerPart, expIntersectionPerPart)
		}
		for j, _ := range intersectionPerPart[i] {
			if !expIntersectionPerPartI[j].Equals(intersectionPerPart[i][j]) {
				t.Errorf(errmsg, kind, da, db,
					intersectionPerPart, expIntersectionPerPart)
			}
		}
	}
	intersectionPerPartInverse := db.IntersectionWithIvDomainThisAssoc(da)
	// check, if there are the same values in intersectionPerPart and
	// intersectionPerPartInverse (only index should be different)
	parts1 := getIntersectionAsSortedDomPartList(intersectionPerPart)
	parts2 := getIntersectionAsSortedDomPartList(intersectionPerPartInverse)
	if len(parts1) != len(parts2) {
		t.Errorf("%s-IntersectionCheck: intersection are not the same!", kind)
	}
	for i := 0; i < len(parts1); i++ {
		if !parts1[i].Equals(parts2[i]) {
			t.Errorf("%s-IntersectionCheck: intersection are not the same!",
				kind)
		}
	}
}

func Test_IntersectionWithDomain(t *testing.T) {
	setup()
	defer teardown()
	log("IntersectWithDomain: ")

	da := CreateIvDomainFromTo(5, 90)
	dbd := CreateExDomainFromTo(1, 100)
	expIntersectionPerPartD := make(map[int][]int)
	expIntersectionPerPartD[0] = makeSlice(5, 90)
	intersectionIvD_DTest(t, "IntersectionWithDomain", da, dbd,
		expIntersectionPerPartD)

	da = CreateIvDomainFromTo(1, 90)
	dbd = CreateExDomainFromTo(1, 100)
	expIntersectionPerPartD = make(map[int][]int)
	expIntersectionPerPartD[0] = makeSlice(1, 90)
	intersectionIvD_DTest(t, "IntersectionWithDomain", da, dbd,
		expIntersectionPerPartD)

	da = CreateIvDomainFromTo(5, 100)
	dbd = CreateExDomainFromTo(1, 100)
	expIntersectionPerPartD = make(map[int][]int)
	expIntersectionPerPartD[0] = makeSlice(5, 100)
	intersectionIvD_DTest(t, "IntersectionWithDomain", da, dbd,
		expIntersectionPerPartD)

	da = CreateIvDomainFromTo(5, 150)
	dbd = CreateExDomainFromTo(1, 100)
	expIntersectionPerPartD = make(map[int][]int)
	expIntersectionPerPartD[0] = makeSlice(5, 100)
	intersectionIvD_DTest(t, "IntersectionWithDomain", da, dbd,
		expIntersectionPerPartD)

	da = CreateIvDomainFromTo(1, 100)
	dbd = CreateExDomainFromTo(5, 150)
	expIntersectionPerPartD = make(map[int][]int)
	expIntersectionPerPartD[0] = makeSlice(5, 100)
	intersectionIvD_DTest(t, "IntersectionWithDomain", da, dbd,
		expIntersectionPerPartD)

	da = CreateIvDomainFromTo(1, 100)
	dbd = CreateExDomainFromTo(5, 90)
	expIntersectionPerPartD = make(map[int][]int)
	expIntersectionPerPartD[0] = makeSlice(5, 90)
	intersectionIvD_DTest(t, "IntersectionWithDomain", da, dbd,
		expIntersectionPerPartD)

	da = CreateIvDomainFromTo(1, 100)
	dbd = CreateExDomainFromTo(1, 90)
	expIntersectionPerPartD = make(map[int][]int)
	expIntersectionPerPartD[0] = makeSlice(1, 90)
	intersectionIvD_DTest(t, "IntersectionWithDomain", da, dbd,
		expIntersectionPerPartD)

	da = CreateIvDomainFromTo(1, 100)
	dbd = CreateExDomainFromTo(10, 100)
	expIntersectionPerPartD = make(map[int][]int)
	expIntersectionPerPartD[0] = makeSlice(10, 100)
	intersectionIvD_DTest(t, "IntersectionWithDomain", da, dbd,
		expIntersectionPerPartD)

	da = CreateIvDomainFromTos([][]int{{1, 5}, {7, 7}, {9, 12}, {15, 16},
		{19, 23}, {99, 120}})
	tmp := makeSlice(1, 10)
	tmp = append(tmp, makeSlice(20, 30)...)
	tmp = append(tmp, makeSlice(50, 80)...)
	tmp = append(tmp, makeSlice(100, 200)...)
	dbd = CreateExDomainAdds(tmp)
	expIntersectionPerPartD = make(map[int][]int)
	expIntersectionPerPartD[0] = makeSlice(1, 5)
	expIntersectionPerPartD[1] = makeSlice(7, 7)
	expIntersectionPerPartD[2] = makeSlice(9, 10)
	expIntersectionPerPartD[4] = makeSlice(20, 23)
	expIntersectionPerPartD[5] = makeSlice(100, 120)
	intersectionIvD_DTest(t, "IntersectionWithDomain", da, dbd,
		expIntersectionPerPartD)
	// inverse
	da = CreateIvDomainFromTos([][]int{{1, 10}, {20, 30}, {50, 80},
		{100, 200}})
	tmp = makeSlice(1, 5)
	tmp = append(tmp, 7)
	tmp = append(tmp, makeSlice(9, 12)...)
	tmp = append(tmp, makeSlice(15, 16)...)
	tmp = append(tmp, makeSlice(19, 23)...)
	tmp = append(tmp, makeSlice(99, 120)...)
	dbd = CreateExDomainAdds(tmp)
	expIntersectionPerPartD = make(map[int][]int)
	tmp = makeSlice(1, 5)
	tmp = append(tmp, 7)
	tmp = append(tmp, makeSlice(9, 10)...)
	expIntersectionPerPartD[0] = tmp
	expIntersectionPerPartD[1] = makeSlice(20, 23)
	expIntersectionPerPartD[3] = makeSlice(100, 120)
	intersectionIvD_DTest(t, "IntersectionWithDomain", da, dbd,
		expIntersectionPerPartD)
}

func Test_IntersectionWithIvDomain(t *testing.T) {
	setup()
	defer teardown()
	log("IntersectWithIvDomain: ")

	da := CreateIvDomainFromTo(1, 100)
	db := CreateIvDomainFromTo(5, 90)
	expIntersectionPerPart := make(map[int][]*IvDomPart)
	expIntersectionPerPart[0] = []*IvDomPart{CreateIvDomPart(5, 90)}
	intersectionIvD_IvDTest(t, "IntersectWithIvDomain", da, db,
		expIntersectionPerPart)

	da = CreateIvDomainFromTo(1, 100)
	db = CreateIvDomainFromTo(1, 90)
	expIntersectionPerPart = make(map[int][]*IvDomPart)
	expIntersectionPerPart[0] = []*IvDomPart{CreateIvDomPart(1, 90)}
	intersectionIvD_IvDTest(t, "IntersectWithIvDomain", da, db,
		expIntersectionPerPart)

	da = CreateIvDomainFromTo(1, 100)
	db = CreateIvDomainFromTo(5, 100)
	expIntersectionPerPart = make(map[int][]*IvDomPart)
	expIntersectionPerPart[0] = []*IvDomPart{CreateIvDomPart(5, 100)}
	intersectionIvD_IvDTest(t, "IntersectWithIvDomain", da, db,
		expIntersectionPerPart)

	da = CreateIvDomainFromTo(1, 100)
	db = CreateIvDomainFromTo(5, 150)
	expIntersectionPerPart = make(map[int][]*IvDomPart)
	expIntersectionPerPart[0] = []*IvDomPart{CreateIvDomPart(5, 100)}
	intersectionIvD_IvDTest(t, "IntersectWithIvDomain", da, db,
		expIntersectionPerPart)

	da = CreateIvDomainFromTo(5, 150)
	db = CreateIvDomainFromTo(1, 100)
	expIntersectionPerPart = make(map[int][]*IvDomPart)
	expIntersectionPerPart[0] = []*IvDomPart{CreateIvDomPart(5, 100)}
	intersectionIvD_IvDTest(t, "IntersectWithIvDomain", da, db,
		expIntersectionPerPart)

	da = CreateIvDomainFromTo(5, 90)
	db = CreateIvDomainFromTo(1, 100)
	expIntersectionPerPart = make(map[int][]*IvDomPart)
	expIntersectionPerPart[0] = []*IvDomPart{CreateIvDomPart(5, 90)}
	intersectionIvD_IvDTest(t, "IntersectWithIvDomain", da, db,
		expIntersectionPerPart)

	da = CreateIvDomainFromTo(1, 90)
	db = CreateIvDomainFromTo(1, 100)
	expIntersectionPerPart = make(map[int][]*IvDomPart)
	expIntersectionPerPart[0] = []*IvDomPart{CreateIvDomPart(1, 90)}
	intersectionIvD_IvDTest(t, "IntersectWithIvDomain", da, db,
		expIntersectionPerPart)

	da = CreateIvDomainFromTo(10, 100)
	db = CreateIvDomainFromTo(1, 100)
	expIntersectionPerPart = make(map[int][]*IvDomPart)
	expIntersectionPerPart[0] = []*IvDomPart{CreateIvDomPart(10, 100)}
	intersectionIvD_IvDTest(t, "IntersectWithIvDomain", da, db,
		expIntersectionPerPart)

	da = CreateIvDomainFromTos([][]int{{1, 10}, {20, 30}, {50, 80},
		{100, 200}})
	db = CreateIvDomainFromTos([][]int{{1, 5}, {7, 7}, {9, 12}, {15, 16},
		{19, 23}, {99, 120}})
	expIntersectionPerPart = make(map[int][]*IvDomPart)
	expIntersectionPerPart[0] = []*IvDomPart{CreateIvDomPart(1, 5),
		CreateIvDomPart(7, 7), CreateIvDomPart(9, 10)}
	expIntersectionPerPart[1] = []*IvDomPart{CreateIvDomPart(20, 23)}
	expIntersectionPerPart[3] = []*IvDomPart{CreateIvDomPart(100, 120)}
	intersectionIvD_IvDTest(t, "IntersectWithIvDomain", da, db,
		expIntersectionPerPart)

	// inverse
	da = CreateIvDomainFromTos([][]int{{1, 5}, {7, 7}, {9, 12}, {15, 16},
		{19, 23}, {99, 120}})
	db = CreateIvDomainFromTos([][]int{{1, 10}, {20, 30}, {50, 80},
		{100, 200}})
	expIntersectionPerPart = make(map[int][]*IvDomPart)
	expIntersectionPerPart[0] = []*IvDomPart{CreateIvDomPart(1, 5)}
	expIntersectionPerPart[1] = []*IvDomPart{CreateIvDomPart(7, 7)}
	expIntersectionPerPart[2] = []*IvDomPart{CreateIvDomPart(9, 10)}
	expIntersectionPerPart[4] = []*IvDomPart{CreateIvDomPart(20, 23)}
	expIntersectionPerPart[5] = []*IvDomPart{CreateIvDomPart(100, 120)}
	intersectionIvD_IvDTest(t, "IntersectWithIvDomain", da, db,
		expIntersectionPerPart)
}

func differenceWithDomainTest(t *testing.T, kind string,
	da *IvDomain, db *ExDomain,
	expDifferenceBA *IvDomain, expDifferenceAB *IvDomain) {
	diffDom := da.DifferenceWithDomain(db)
	if !diffDom.Equals(expDifferenceAB) {
		t.Errorf("DifferenceWithDomain: a is %s, b is %s, got %s, want %s",
			da, db, diffDom, expDifferenceAB)
	}
}

func differenceWithIvDomainTest(t *testing.T,
	da *IvDomain, db *IvDomain,
	expDifferenceBA *IvDomain, expDifferenceAB *IvDomain) {
	daC := da.Copy().(*IvDomain)
	dbC := db.Copy().(*IvDomain)
	diffDom := dbC.DifferenceWithIvDomain(daC)
	if !diffDom.Equals(expDifferenceBA) {
		t.Errorf("DifferenceWithIvDomain-BA: b %s, a %s, got %s, want %s",
			dbC, daC, diffDom, expDifferenceBA)
	}
	daC = da.Copy().(*IvDomain)
	dbC = db.Copy().(*IvDomain)
	diffDom = daC.DifferenceWithIvDomain(dbC)
	if !diffDom.Equals(expDifferenceAB) {
		t.Errorf("DifferenceWithIvDomain-AB: a %s, b %s, got %s, want %s",
			daC, dbC, diffDom, expDifferenceAB)
	}
}

func Test_DifferenceWithIvDomain(t *testing.T) {
	setup()
	defer teardown()
	log("DifferenceWithIvDomain: ")

	da := CreateIvDomainFromTo(1, 100)
	db := CreateIvDomainFromTo(5, 90)
	expDiffBA := CreateIvDomain()
	expDiffAB := CreateIvDomainFromTos([][]int{{1, 4}, {91, 100}})
	differenceWithIvDomainTest(t, da, db, expDiffBA, expDiffAB)

	da = CreateIvDomainFromTo(1, 10)
	db = CreateIvDomainFromTo(0, 0)
	expDiffBA = CreateIvDomainFromTo(0, 0)
	expDiffAB = CreateIvDomainFromTo(1, 10)
	differenceWithIvDomainTest(t, da, db, expDiffBA, expDiffAB)

	da = CreateIvDomainFromTo(1, 100)
	db = CreateIvDomainFromTo(1, 90)
	expDiffBA = CreateIvDomain()
	expDiffAB = CreateIvDomainFromTos([][]int{{91, 100}})
	differenceWithIvDomainTest(t, da, db, expDiffBA, expDiffAB)

	da = CreateIvDomainFromTo(1, 100)
	db = CreateIvDomainFromTo(5, 100)
	expDiffBA = CreateIvDomain()
	expDiffAB = CreateIvDomainFromTos([][]int{{1, 4}})
	differenceWithIvDomainTest(t, da, db, expDiffBA, expDiffAB)

	da = CreateIvDomainFromTo(1, 100)
	db = CreateIvDomainFromTo(5, 150)
	expDiffBA = CreateIvDomainFromTo(101, 150)
	expDiffAB = CreateIvDomainFromTo(1, 4)
	differenceWithIvDomainTest(t, da, db, expDiffBA, expDiffAB)

	da = CreateIvDomainFromTo(5, 150)
	db = CreateIvDomainFromTo(1, 100)
	expDiffBA = CreateIvDomainFromTo(1, 4)
	expDiffAB = CreateIvDomainFromTo(101, 150)
	differenceWithIvDomainTest(t, da, db, expDiffBA, expDiffAB)

	da = CreateIvDomainFromTo(5, 90)
	db = CreateIvDomainFromTo(1, 100)
	expDiffBA = CreateIvDomainFromTos([][]int{{1, 4}, {91, 100}})
	expDiffAB = CreateIvDomain()
	differenceWithIvDomainTest(t, da, db, expDiffBA, expDiffAB)

	da = CreateIvDomainFromTo(1, 90)
	db = CreateIvDomainFromTo(1, 100)
	expDiffBA = CreateIvDomainFromTo(91, 100)
	expDiffAB = CreateIvDomain()
	differenceWithIvDomainTest(t, da, db, expDiffBA, expDiffAB)

	da = CreateIvDomainFromTo(10, 100)
	db = CreateIvDomainFromTo(1, 100)
	expDiffBA = CreateIvDomainFromTo(1, 9)
	expDiffAB = CreateIvDomain()
	differenceWithIvDomainTest(t, da, db, expDiffBA, expDiffAB)

	da = CreateIvDomainFromTo(10, 100)
	db = CreateIvDomainFromTo(1, 5)
	expDiffBA = CreateIvDomainFromTo(1, 5)
	expDiffAB = CreateIvDomainFromTo(10, 100)
	differenceWithIvDomainTest(t, da, db, expDiffBA, expDiffAB)

	da = CreateIvDomainFromTos([][]int{{1, 10}, {20, 30}, {50, 80},
		{100, 200}})
	db = CreateIvDomainFromTos([][]int{{1, 5}, {7, 7}, {9, 12}, {15, 16},
		{19, 23}, {99, 120}})
	expDiffBA = CreateIvDomainFromTos([][]int{{11, 12}, {15, 16}, {19, 19},
		{99, 99}})
	expDiffAB = CreateIvDomainFromTos([][]int{{6, 6}, {8, 8}, {24, 30},
		{50, 80}, {121, 200}})
	differenceWithIvDomainTest(t, da, db, expDiffBA, expDiffAB)

	da = CreateIvDomainFromTos([][]int{{0, 8}})
	db = CreateIvDomainFromTos([][]int{{1, 5}, {7, 7}})
	expDiffBA = CreateIvDomain()
	expDiffAB = CreateIvDomainFromTos([][]int{{0, 0}, {6, 6}, {8, 8}})
	differenceWithIvDomainTest(t, da, db, expDiffBA, expDiffAB)
}

func IvDomainRemoveTest(t *testing.T, domVals [][]int, ele int,
	expVals [][]int) {
	dom := CreateIvDomainFromTos(domVals)
	expDom := CreateIvDomainFromTos(expVals)
	dom.Remove(ele)
	if !dom.Equals(expDom) {
		t.Errorf("IvDomain-Remove failed: d = %v, want %v", dom, expDom)
	}
}

func Test_IvDomainRemove_1(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainRemove_removing ele: 5 from 1..100")
	IvDomainRemoveTest(t, [][]int{{1, 100}}, 5, [][]int{{1, 4}, {6, 100}})
	IvDomainRemoveTest(t, [][]int{{1, 100}}, 1, [][]int{{2, 100}})
	IvDomainRemoveTest(t, [][]int{{1, 100}}, 100, [][]int{{1, 99}})
	IvDomainRemoveTest(t, [][]int{{1, 100}}, 101, [][]int{{1, 100}})
}

func Test_IvDomainRemoves_1(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainRemoves_removing eles: 5 from 1..100")
	d := CreateIvDomainFromTo(1, 100)
	eles := make([]int, 1)
	eles[0] = 5
	vals := [][]int{{1, 4}, {6, 100}}
	expD := CreateIvDomainFromTos(vals)
	IvDomainRemovesTest(t, d, eles, expD)
}

func Test_IvDomainRemoves_2(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainRemoves_removing eles: 5, 6, 7, 8 from 1..100")
	d := CreateIvDomainFromTo(1, 100)
	eles := make([]int, 4)
	eles[0] = 5
	eles[1] = 6
	eles[2] = 7
	eles[3] = 8
	vals := [][]int{{1, 4}, {9, 100}}
	expD := CreateIvDomainFromTos(vals)
	IvDomainRemovesTest(t, d, eles, expD)
}

func Test_IvDomainRemoves_3(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainRemoves_removing eles: 5, 8, 20, 50 from 1..100")
	d := CreateIvDomainFromTo(1, 100)
	eles := make([]int, 4)
	eles[0] = 5
	eles[1] = 8
	eles[2] = 20
	eles[3] = 50
	vals := [][]int{{1, 4}, {6, 7}, {9, 19}, {21, 49}, {51, 100}}
	expD := CreateIvDomainFromTos(vals)
	IvDomainRemovesTest(t, d, eles, expD)
}

func Test_IvDomainRemoves_4(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainRemoves_removing eles 5, 8, 20, 50 from 1..100" +
		"(with other functions)")
	d := CreateIvDomainFromTo(1, 100)
	eles := makeTwoDim_OneDim([][]int{{5, 5}, {8, 8}, {20, 20}, {50, 50}})
	expD := CreateIvDomainFromTos([][]int{{1, 4}, {6, 7}, {9, 19},
		{21, 49}, {51, 100}})
	IvDomainRemovesTest(t, d, eles, expD)
}

func Test_IvDomainRemoves_5(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainRemoves_removing eles 0 from 1..100")
	d := CreateIvDomainFromTo(1, 100)
	eles := makeTwoDim_OneDim([][]int{{0, 0}})
	expD := CreateIvDomainFromTos([][]int{{1, 100}})
	IvDomainRemovesTest(t, d, eles, expD)
}

func Test_IvDomainRemoves_6(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainRemoves_removing eles 1..100, 1000..2000 from 1..10000")
	d := CreateIvDomainFromTo(1, 10000)
	eles := makeTwoDim_OneDim([][]int{{1, 100}, {1000, 2000}})
	expD := CreateIvDomainFromTos([][]int{{101, 999}, {2001, 10000}})
	IvDomainRemovesTest(t, d, eles, expD)
}

func Test_IvDomainRemoves_7(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainRemoves_removing eles 11000 from 1..10000")
	d := CreateIvDomainFromTo(1, 10000)
	eles := makeTwoDim_OneDim([][]int{{11000, 11000}})
	expD := CreateIvDomainFromTo(1, 10000)
	IvDomainRemovesTest(t, d, eles, expD)
}

func Test_IvDomainRemoves_8(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainRemoves_removing eles 99 from 100..10000")
	d := CreateIvDomainFromTo(100, 10000)
	eles := makeTwoDim_OneDim([][]int{{99, 99}})
	expD := CreateIvDomainFromTo(100, 10000)
	IvDomainRemovesTest(t, d, eles, expD)
}

func createRemovingRange(from, to int) map[int]bool {
	size := (to + 1) - from
	r := make(map[int]bool, size)
	for i := from; i < to; i++ {
		r[i] = true
	}
	return r
}

func IvDomainStringTest(t *testing.T, a [][]int, exp string) {
	ivdom := CreateIvDomainFromTos(a)
	got := ivdom.String()
	if got != exp {
		t.Errorf("Domain(%v): d.String() = %v, want %v", a, got, exp)
	}
}

func Test_IvDomainString(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainString")
	d := [][]int{{1, 6}}
	IvDomainStringTest(t, d, "[1..6]")
	d = [][]int{{1, 2}, {6}}
	IvDomainStringTest(t, d, "[1..2,6]")
	d = [][]int{{1, 4}, {6}}
	IvDomainStringTest(t, d, "[1..4,6]")
	d = [][]int{{1, 4}, {7, 9}, {15}, {17}, {19, 21}}
	IvDomainStringTest(t, d, "[1..4,7..9,15,17,19..21]")
	d = [][]int{{0, 1}, {3, 4}, {8, 9}, {15, 17}, {19, 21}, {23}}
	IvDomainStringTest(t, d, "[0..1,3..4,8..9,15..17,19..21,23]")
	d = [][]int{{0, 1}, {3, 4}, {8, 9}, {15, 17}, {19, 21}, {23}, {25}}
	IvDomainStringTest(t, d, "[0..1,3..4,8..9,15..17,19..21,23,25]")
	d = [][]int{{0}, {3, 4}}
	IvDomainStringTest(t, d, "[0,3..4]")
	d = [][]int{{0}, {3}, {5}}
	IvDomainStringTest(t, d, "[0,3,5]")
	d = [][]int{{}}
	IvDomainStringTest(t, d, "[]")
}

func IvDomainEqualsTest(t *testing.T, a [][]int, b [][]int, exp bool) {
	da := CreateIvDomainFromTos(a)
	db := CreateIvDomainFromTos(b)
	got := da.Equals(db)
	if got != exp {
		t.Errorf("Domain(%v).Equals(Domain(%v)) = %v, want %v",
			da, db, got, exp)
	}

	bs := makeTwoDim_OneDim(b)
	dc := CreateExDomainAdds(bs)
	got = da.Equals(dc)
	if got != exp {
		t.Errorf("Domain(%v).Equals(Domain(%v)) = %v, want %v",
			da, dc, got, exp)
	}

	as := makeTwoDim_OneDim(a)
	dd := CreateExDomainAdds(as)
	got = db.Equals(dd)
	if got != exp {
		t.Errorf("Domain(%v).Equals(Domain(%v)) = %v, want %v",
			db, dd, got, exp)
	}
}

func Test_IvDomainEquals(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainEquals")
	a1 := [][]int{{1, 6}}
	a2 := [][]int{{1, 6}}
	IvDomainEqualsTest(t, a1, a2, true)
	a1 = [][]int{{1, 6}}
	a2 = [][]int{{1, 5}}
	IvDomainEqualsTest(t, a1, a2, false)
	a1 = [][]int{{1, 6}}
	a2 = [][]int{{-6}, {1, 5}}
	IvDomainEqualsTest(t, a1, a2, false)
	a1 = [][]int{{}}
	a2 = [][]int{{1}}
	IvDomainEqualsTest(t, a1, a2, false)
	a1 = [][]int{{1}}
	a2 = [][]int{{}}
	IvDomainEqualsTest(t, a1, a2, false)
	a1 = [][]int{{}}
	a2 = [][]int{{}}
	IvDomainEqualsTest(t, a1, a2, true)
}

func IvDomainCopyTest(t *testing.T, a [][]int) {
	da := CreateIvDomainFromTos(a)
	db := da.Copy()
	if !da.Equals(db) {
		t.Errorf("Domain(%v).Copy() = %v",
			da, db)
	}
}

func Test_IvDomainCopy(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainCopy")
	d := [][]int{{1, 6}}
	IvDomainCopyTest(t, d)
	d = [][]int{{-6}, {1, 5}}
	IvDomainCopyTest(t, d)
	d = [][]int{{}}
	IvDomainCopyTest(t, d)
	d = [][]int{{1}}
	IvDomainCopyTest(t, d)
	d = [][]int{{1, 5}, {10, 100}, {105}, {110, 120}, {150}}
	IvDomainCopyTest(t, d)
}

func IvDomainIsEmptyTest(t *testing.T, dom *IvDomain, isempty bool,
	min int, max int) {
	if dom.IsEmpty() != isempty {
		t.Errorf("Domain(%v).IsEmpty() = %v, want %v",
			dom, dom.IsEmpty(), isempty)
	}
	minTest(t, "IsEmpty", dom, min)
	maxTest(t, "IsEmpty", dom, max)
}

func Test_IsEmpty(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainIsEmpty")
	d := CreateIvDomainFromTos([][]int{{1, 3}, {6, 8}})
	d2 := CreateIvDomainFromTos([][]int{{1, 3}, {6, 8}})
	d.Removes(d2)
	IvDomainIsEmptyTest(t, d, true, 1, 0)
}

func IvDomainIsGroundTest(t *testing.T, dom *IvDomain, isground bool) {
	if dom.IsGround() != isground {
		t.Errorf("Domain(%v).IsGround() = %v, want %v",
			dom, dom.IsGround(), isground)
	}
	if isground && (dom.Size() != 1) {
		t.Errorf("Domain(%v).Size() = %v, want %v",
			dom, dom.Size(), 1)
	}
}

func Test_IsGround(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainIsGround")
	d := CreateIvDomainFromTos([][]int{{}})
	IvDomainIsGroundTest(t, d, false)

	d = CreateIvDomainFromTos([][]int{{1, 3}})
	d.Remove(1)
	d.Remove(2)
	d.Remove(3)
	IvDomainIsEmptyTest(t, d, true, 1, 0)

	d = CreateIvDomainFromTos([][]int{{1, 3}})
	d.Remove(1)
	d.Remove(2)
	IvDomainIsEmptyTest(t, d, false, 3, 3)
	IvDomainIsGroundTest(t, d, true)

	d = CreateIvDomainFromTos([][]int{{1, 3}, {4, 5}})
	d2 := CreateIvDomainFromTos([][]int{{1, 3}, {4, 4}})
	d.Removes(d2)
	IvDomainIsEmptyTest(t, d, false, 5, 5)
	IvDomainIsGroundTest(t, d, true)
	d.Remove(5)
	IvDomainIsEmptyTest(t, d, true, 1, 0)
}

func IvDomainContainsTest(t *testing.T, dom *IvDomain, ele int,
	contains bool) {
	if dom.Contains(ele) != contains {
		t.Errorf("Domain(%v).Contains(%v) = %v, want %v",
			dom, ele, dom.Contains(ele), contains)
	}
}

func Test_Contains(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainIsEmpty")
	d := CreateIvDomainFromTos([][]int{{}})
	IvDomainContainsTest(t, d, 1, false)

	d = CreateIvDomainFromTos([][]int{{1, 3}})
	d.Remove(1)
	d.Remove(2)
	IvDomainContainsTest(t, d, 3, true)
	IvDomainContainsTest(t, d, 1, false)

	d = CreateIvDomainFromTos([][]int{{1, 3}})
	IvDomainContainsTest(t, d, 2, true)

	d = CreateIvDomainFromTos([][]int{{1, 3}, {4, 6}, {7}})
	IvDomainContainsTest(t, d, 5, true)
	IvDomainContainsTest(t, d, 6, true)
	IvDomainContainsTest(t, d, 7, true)
}

func minTest(t *testing.T, test string, dom Domain, min int) {
	if !dom.IsEmpty() {
		if dom.GetMin() != min {
			t.Errorf(test+"\r\n"+"IvDomain-Min: min(%s) is %v, want %v",
				dom, dom.GetMin(), min)
		}
	} else {
		if 1 != min {
			t.Errorf(test+"\r\n"+"IvDomain-Min: min(%s) is %v, want %v",
				dom, 1, min)
		}
	}
}

func maxTest(t *testing.T, test string, dom Domain, max int) {
	if !dom.IsEmpty() {
		if dom.GetMax() != max {
			t.Errorf(test+"\r\n"+"IvDomain-Max: max(%s) is %v, want %v",
				dom, dom.GetMax(), max)
		}
	} else {
		if 0 != max {
			t.Errorf(test+"\r\n"+"IvDomain-Min: max(%s) is %v, want %v",
				dom, 0, max)
		}
	}
}

func getMaxTest(t *testing.T, test string, a [][]int, max int) {
	da := CreateIvDomainFromTos(a)
	maxTest(t, test, da, max)
}

func getMinTest(t *testing.T, test string, a [][]int, min int) {
	da := CreateIvDomainFromTos(a)
	minTest(t, test, da, min)
}

func Test_IvDomainGetMaxAndMin(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainGetMaxAndMin")

	getMaxTest(t, "GetMaxandMinTest", [][]int{{1, 6}, {10, 15}, {20, 28}}, 28)
	getMaxTest(t, "GetMaxandMinTest", [][]int{{1, 6}}, 6)
	getMaxTest(t, "GetMaxandMinTest", [][]int{{6}}, 6)

	getMinTest(t, "GetMaxandMinTest", [][]int{{1, 6}, {10, 15}, {20, 28}}, 1)
	getMinTest(t, "GetMaxandMinTest", [][]int{{1, 6}}, 1)
	getMinTest(t, "GetMaxandMinTest", [][]int{{6}}, 6)
}

func removingTestWithRemovingOtherCheck(t *testing.T, kind string,
	a Domain, b Domain, expA Domain, expB Domain,
	expAMin int, expAMax int, expBMin int, expBMax int) {
	beforeDomA := a.Copy()
	beforeDomB := b.Copy()
	a.RemovesWithOther(b)
	if !a.Equals(expA) || !b.Equals(expB) {
		t.Errorf("IvDomain-Removing-%s: DomA %s removing with DomB %s,\r\n"+
			" got DomA %s and DomB %s, want DomA %s and DomB %s",
			kind, beforeDomA, beforeDomB, a, b, expA, expB)
	}
	minTest(t, "IvDomain-Removing-"+kind, a, expAMin)
	maxTest(t, "IvDomain-Removing-"+kind, a, expAMax)
	minTest(t, "IvDomain-Removing-"+kind, b, expBMin)
	maxTest(t, "IvDomain-Removing-"+kind, b, expBMax)
}

// removingTestWithRemovingOther tests will removing-flag == true (also
// removing from other Domain b, if it is not in Domain a.
func removingTestWithRemovingOther(t *testing.T, a [][]int, b [][]int,
	expA [][]int, expB [][]int,
	expAMin int, expAMax int, expBMin int, expBMax int) {
	// IvD_IvD
	aIvDom := CreateIvDomainFromTos(a)
	expaIvDom := CreateIvDomainFromTos(expA)
	bIvDom := CreateIvDomainFromTos(b)
	expbIvDom := CreateIvDomainFromTos(expB)
	removingTestWithRemovingOtherCheck(t, "IvD_IvD", aIvDom, bIvDom,
		expaIvDom, expbIvDom, expAMin, expAMax, expBMin, expBMax)
	// IvD_D
	aIvDom = CreateIvDomainFromTos(a)
	expaIvDom = CreateIvDomainFromTos(expA)
	bDom := CreateExDomainAdds(makeTwoDim_OneDim(b))
	expbDom := CreateExDomainAdds(makeTwoDim_OneDim(expB))
	removingTestWithRemovingOtherCheck(t, "IvD_D", aIvDom, bDom,
		expaIvDom, expbDom, expAMin, expAMax, expBMin, expBMax)
	// D_IvD
	adDom := CreateExDomainAdds(makeTwoDim_OneDim(a))
	expaDom := CreateExDomainAdds(makeTwoDim_OneDim(expA))
	bIvDom = CreateIvDomainFromTos(b)
	expbIvDom = CreateIvDomainFromTos(expB)
	removingTestWithRemovingOtherCheck(t, "D_IvD", adDom, bIvDom,
		expaDom, expbIvDom, expAMin, expAMax, expBMin, expBMax)
	// D_D
	adDom = CreateExDomainAdds(makeTwoDim_OneDim(a))
	expaDom = CreateExDomainAdds(makeTwoDim_OneDim(expA))
	bDom = CreateExDomainAdds(makeTwoDim_OneDim(b))
	expbDom = CreateExDomainAdds(makeTwoDim_OneDim(expB))
	removingTestWithRemovingOtherCheck(t, "D_D", adDom, bDom,
		expaDom, expbDom, expAMin, expAMax, expBMin, expBMax)
}

func Test_IvDomainRemovesAndOtherRemoves(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainRemovesAndOtherRemoves")

	//a, b, expA, expB, expAMin, expAMax, expBMin, expBMax
	removingTestWithRemovingOther(t, [][]int{{1, 10}}, [][]int{{5, 7}},
		[][]int{{1, 4}, {8, 10}}, [][]int{{5, 7}}, 1, 10, 5, 7)
	removingTestWithRemovingOther(t, [][]int{{1, 10}}, [][]int{{1}},
		[][]int{{2, 10}}, [][]int{{1}}, 2, 10, 1, 1)
	removingTestWithRemovingOther(t, [][]int{{1, 10}}, [][]int{{10}},
		[][]int{{1, 9}}, [][]int{{10}}, 1, 9, 10, 10)
	removingTestWithRemovingOther(t, [][]int{{1, 10}}, [][]int{{0}},
		[][]int{{1, 10}}, [][]int{{}}, 1, 10, 1, 0)
	removingTestWithRemovingOther(t, [][]int{{1, 10}, {50, 100}},
		[][]int{{1}, {5, 7}, {50}, {99, 100}},
		[][]int{{2, 4}, {8, 10}, {51, 98}},
		[][]int{{1}, {5, 7}, {50}, {99, 100}}, 2, 98, 1, 100)
	removingTestWithRemovingOther(t, [][]int{{1, 10}, {50, 100}},
		[][]int{{1}, {20, 30}, {50}}, [][]int{{2, 10}, {51, 100}},
		[][]int{{1}, {50}}, 2, 100, 1, 50)
	removingTestWithRemovingOther(t, [][]int{{1, 4}}, [][]int{{1, 3}},
		[][]int{{4, 4}}, [][]int{{1, 3}}, 4, 4, 1, 3)
	removingTestWithRemovingOther(t, [][]int{{0, 10}}, [][]int{{0}},
		[][]int{{1, 10}}, [][]int{{0, 0}}, 1, 10, 0, 0)
	removingTestWithRemovingOther(t, [][]int{{0, 10}},
		[][]int{{0, 0}, {2, 2}, {5, 5}},
		[][]int{{1, 1}, {3, 4}, {6, 10}},
		[][]int{{0, 0}, {2, 2}, {5, 5}}, 1, 10, 0, 5)
}

func removingTestCheck(t *testing.T, kind string,
	dom Domain, removingDom Domain, expDom Domain,
	expMin int, expMax int) {
	beforeDom := dom.Copy()
	dom.Removes(removingDom)
	if !dom.Equals(expDom) {
		t.Errorf("IvDomain-Removing-%s: Dom %s, removing %s, got %s, want %s",
			kind, beforeDom, removingDom, dom, expDom)
	}
	minTest(t, "IvDomain-Removing-"+kind, dom, expMin)
	maxTest(t, "IvDomain-Removing-"+kind, dom, expMax)
}

func removingTest(t *testing.T, a [][]int, values []int, expA [][]int,
	expMin int, expMax int) {
	// IvD_IvD
	aIvDom := CreateIvDomainFromTos(a)
	expaIvDom := CreateIvDomainFromTos(expA)
	valuesIvDom := CreateIvDomainFromIntArr(values)
	removingTestCheck(t, "IvD_IvD", aIvDom, valuesIvDom,
		expaIvDom, expMin, expMax)
	// IvD_D
	aIvDom = CreateIvDomainFromTos(a)
	expaIvDom = CreateIvDomainFromTos(expA)
	valuesDom := CreateExDomainAdds(values)
	removingTestCheck(t, "IvD_D", aIvDom, valuesDom,
		expaIvDom, expMin, expMax)
	// D_IvD
	adDom := CreateExDomainAdds(makeTwoDim_OneDim(a))
	expaDom := CreateExDomainAdds(makeTwoDim_OneDim(expA))
	valuesIvDom = CreateIvDomainFromIntArr(values)
	removingTestCheck(t, "D_IvD", adDom, valuesIvDom, expaDom, expMin, expMax)

	// D_D
	adDom = CreateExDomainAdds(makeTwoDim_OneDim(a))
	expaDom = CreateExDomainAdds(makeTwoDim_OneDim(expA))
	valuesDom = CreateExDomainAdds(values)
	removingTestCheck(t, "D_D", adDom, valuesDom, expaDom, expMin, expMax)
}

func removingTest2(t *testing.T, a [][]int, values [][]int, expA [][]int,
	expMin int, expMax int) {
	// IvD_IvD
	aIvDom := CreateIvDomainFromTos(a)
	expaIvDom := CreateIvDomainFromTos(expA)
	valuesIvDom := CreateIvDomainFromTos(values)
	removingTestCheck(t, "IvD_IvD", aIvDom, valuesIvDom,
		expaIvDom, expMin, expMax)
}

func Test_IvDomainRemoves_normal(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainRemoves")
	removingTest(t, [][]int{{1, 10}}, []int{5, 6, 7},
		[][]int{{1, 4}, {8, 10}}, 1, 10)
	removingTest(t, [][]int{{1, 10}}, []int{1}, [][]int{{2, 10}}, 2, 10)
	removingTest(t, [][]int{{1, 10}}, []int{10}, [][]int{{1, 9}}, 1, 9)
	removingTest(t, [][]int{{1, 10}}, []int{0}, [][]int{{1, 10}}, 1, 10)
	removingTest(t, [][]int{{1, 10}, {50, 100}},
		[]int{1, 5, 6, 7, 50, 99, 100},
		[][]int{{2, 4}, {8, 10}, {51, 98}}, 2, 98)
	removingTest(t, [][]int{{1, 4}}, []int{1, 2, 3}, [][]int{{4}}, 4, 4)
	removingTest(t, [][]int{{1}, {4}, {6}}, []int{1, 4}, [][]int{{6}}, 6, 6)
	removingTest(t, [][]int{{1}, {4}, {6}}, []int{1, 6}, [][]int{{4}}, 4, 4)
	removingTest(t, [][]int{{1}, {4}, {6}, {8, 10}}, []int{1, 6},
		[][]int{{4}, {8, 10}}, 4, 10)
	removingTest(t, [][]int{{1}, {4}, {6}, {8, 10}},
		[]int{6, 8, 9, 10}, [][]int{{1}, {4}}, 1, 4)
	removingTest(t, [][]int{{1}, {4}, {6}, {8, 10}},
		[]int{6, 8, 10}, [][]int{{1}, {4}, {9}}, 1, 9)
	removingTest(t, [][]int{{1}, {4}, {6}, {8, 10}},
		[]int{1, 8, 9, 10}, [][]int{{4}, {6}}, 4, 6)
	removingTest(t, [][]int{{0, 10}}, []int{0}, [][]int{{1, 10}}, 1, 10)
	removingTest(t, [][]int{{0, 10}}, []int{0, 2, 5},
		[][]int{{1, 1}, {3, 4}, {6, 10}}, 1, 10)
	// fix-Tests (if variable is fixed to one value, remove all others)...
	removingTest2(t, [][]int{{0, 100}}, [][]int{{1, 100}},
		[][]int{{0, 0}}, 0, 0)
	removingTest2(t, [][]int{{0, 100}}, [][]int{{0, 99}},
		[][]int{{100, 100}}, 100, 100)
	removingTest2(t,
		[][]int{{0, 10}, {20, 30}, {31, 50}, {60, 70}, {80, 90}},
		[][]int{{1, 90}}, [][]int{{0, 0}}, 0, 0)
	removingTest2(t,
		[][]int{{0, 10}, {20, 30}, {31, 50}, {60, 70}, {80, 90}},
		[][]int{{0, 89}}, [][]int{{90, 90}}, 90, 90)
	removingTest2(t,
		[][]int{{0, 10}, {20, 30}, {31, 50}, {60, 70}, {80, 90}},
		[][]int{{0, 48}, {50, 90}}, [][]int{{49, 49}}, 49, 49)
}

func ivsortedAsDomPartsTest(t *testing.T, a [][]int) {
	da := CreateIvDomainDomPartsWithSort(CreateIvDomParts(a))
	parts := da.GetParts()
	for i := 0; i < len(parts)-1; i++ {
		if parts[i].From >= parts[i+1].From {
			t.Errorf("IvDomain-Sorting failure (DomPart)")
			return
		}
	}
}

func ivsortedIntTest(t *testing.T, a [][]int) {
	da := CreateIvDomainDomPartsWithSort(CreateIvDomParts(a))
	vals := da.SortedValues()
	for i := 0; i < len(vals)-1; i++ {
		if vals[i] >= vals[i+1] {
			t.Errorf("IvDomain-Sorting failure (int)")
			return
		}
	}
}

func Test_IvDomainSorted(t *testing.T) {
	ivsortedAsDomPartsTest(t, [][]int{{1}, {4}, {6}, {8, 10}})
	ivsortedAsDomPartsTest(t, [][]int{{1}, {8, 10}, {6}})
	ivsortedAsDomPartsTest(t, [][]int{{6}, {8, 10}, {1}})
	ivsortedAsDomPartsTest(t,
		[][]int{{6, 100}, {3, 4}, {1}, {200, 201}, {150, 156}})
	ivsortedIntTest(t, [][]int{{1}, {4}, {6}, {8, 10}})
	ivsortedIntTest(t, [][]int{{1}, {8, 10}, {6}})
	ivsortedIntTest(t, [][]int{{6}, {8, 10}, {1}})
	ivsortedIntTest(t, [][]int{{6, 100}, {3, 4}, {1}, {200, 201}, {150, 156}})
}

func makeIvDomainWorstCase(n int, except int) []int {
	vals := make([]int, n)
	v := 0
	for i := 0; v <= n; i++ {
		v = i * 2
		if except != v {
			vals = append(vals, v)
		}
	}
	return vals
}

func convertIvDomainRemoveTest(t *testing.T,
	from, to int, eles [][]int, expEles [][]int) {
	msg := "Removing "
	for _, ft := range eles {
		msg += fmt.Sprintf("%d..%d ", ft[0], ft[1])
	}
	msg += fmt.Sprintf("from %d..%d", from, to)
	log(msg)
	d := CreateIvDomainFromTo(from, to)
	xeles := CreateIvDomainFromIntArr(makeTwoDim_OneDim(eles))
	dExp := CreateIvDomainFromIntArr(makeTwoDim_OneDim(expEles))
	DomainRemoveTest(t, d, xeles, dExp)
}

func Test_IvDomainPerformance_1(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainPerformance_1")
	from, to := 1, 50000
	eles := [][]int{{100, 1500}, {10000, 15000}}
	expEles := [][]int{{1, 99}, {1501, 9999}, {15001, 50000}}
	convertIvDomainRemoveTest(t, from, to, eles, expEles)
}

func Test_IvDomainPerformance_2(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainPerformance_2")
	from, to := 1, 50000
	eles := [][]int{{100, 1500}, {10000, 15000}}
	expEles := [][]int{{1, 99}, {1501, 9999}, {15001, 50000}}
	convertIvDomainRemoveTest(t, from, to, eles, expEles)
}

func Test_IvDomainPerformance_3(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainPerformance_3")
	from, to := 1, 500000
	eles := [][]int{{100, 1500}, {10000, 15000}, {200000, 200100}}
	expEles := [][]int{{1, 99}, {1501, 9999}, {15001, 199999},
		{200101, 500000}}
	convertIvDomainRemoveTest(t, from, to, eles, expEles)
}

func Test_IvDomainPerformance_4(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainPerformance_4")
	from, to := 1, 500000
	eles := [][]int{{100, 1500}, {10000, 15000}, {200000, 200100}}
	expEles := [][]int{{1, 99}, {1501, 9999}, {15001, 199999},
		{200101, 500000}}
	convertIvDomainRemoveTest(t, from, to, eles, expEles)
}

func Test_IvDomainPerformance_5(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomainPerformance_5")
	log("Removing ele 9000 from 0..0, 2..2, ..., 10000..10000")
	vals := makeIvDomainWorstCase(10000, -1)
	d := CreateIvDomainFromIntArr(vals)
	eles := CreateIvDomainFromIntArr(makeTwoDim_OneDim([][]int{{9000, 9000}}))
	valsExp := makeIvDomainWorstCase(10000, 9000)
	dExp := CreateIvDomainFromIntArr(valsExp)
	DomainRemoveTest(t, d, eles, dExp)
}

func ivdompartAddTest(t *testing.T, dvals [][]int, val int, dExpVals [][]int) {
	var d *IvDomain
	if len(dvals) != 0 {
		d = CreateIvDomainFromTos(dvals)
	} else {
		d = CreateIvDomain()
	}
	dExp := CreateIvDomainFromTos(dExpVals)
	d.Add(val)
	if !d.Equals(dExp) {
		errmsg := "IvDomain-Add failure, "
		errmsg += "wanted Domain %s with min ==%v and max ==%v, "
		errmsg += "got %s with min ==%v and max ==%v"
		t.Errorf(errmsg,
			dExp, dExp.GetMin(), dExp.GetMax(), d, d.GetMin(), d.GetMax())
	}
}

func Test_IvDomPartAdd(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomain-IvDomPartAdd")

	ivdompartAddTest(t, [][]int{{}}, 10, [][]int{{10}})
	ivdompartAddTest(t, [][]int{{10}}, 9, [][]int{{9, 10}})
	ivdompartAddTest(t, [][]int{{9, 10}}, 11, [][]int{{9, 11}})
	ivdompartAddTest(t, [][]int{{9, 11}}, 11, [][]int{{9, 11}})
	ivdompartAddTest(t, [][]int{{9, 11}}, 5, [][]int{{5}, {9, 11}})
	ivdompartAddTest(t, [][]int{{5}, {9, 11}}, 7, [][]int{{5}, {7}, {9, 11}})
	ivdompartAddTest(t, [][]int{{5}, {7}, {9, 11}}, 6, [][]int{{5, 7}, {9, 11}})
	ivdompartAddTest(t, [][]int{{5, 7}, {9, 11}}, 20,
		[][]int{{5, 7}, {9, 11}, {20}})
	ivdompartAddTest(t, [][]int{{5, 7}, {9, 11}, {20}}, 15,
		[][]int{{5, 7}, {9, 11}, {15}, {20}})
}

func Test_IvDomPartAdd2(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomain-IvDomPartAdd2")
	ivdompartAddTest(t, [][]int{{}}, 10, [][]int{{10}})
	ivdompartAddTest(t, [][]int{{10}}, 7, [][]int{{7}, {10}})
	ivdompartAddTest(t, [][]int{{7}, {10}}, 8, [][]int{{7, 8}, {10}})
}

func ivdompartAddPartTest(t *testing.T, dvals [][]int, partVal []int,
	dExpVals [][]int) {
	var d *IvDomain
	if len(dvals) != 0 {
		d = CreateIvDomainFromTos(dvals)
	} else {
		d = CreateIvDomain()
	}
	dExp := CreateIvDomainFromTos(dExpVals)
	from := partVal[0]
	to := partVal[0]
	if len(partVal) == 2 {
		to = partVal[1]
	}
	p := CreateIvDomPart(from, to)
	d.partList.AddAnyPart(p)
	if !d.Equals(dExp) {
		errmsg := "IvDomain-AddPart failure, "
		errmsg += "wanted Domain %s with min ==%v and max ==%v, "
		errmsg += "got %s with min ==%v and max ==%v"
		t.Errorf(errmsg,
			dExp, dExp.GetMin(), dExp.GetMax(), d, d.GetMin(), d.GetMax())
	}
}

func Test_IvDomPartAddPart(t *testing.T) {
	setup()
	defer teardown()
	log("IvDomain-IvDomPartAddPart")

	ivdompartAddPartTest(t, [][]int{{10}}, []int{7}, [][]int{{7}, {10}})
	ivdompartAddPartTest(t, [][]int{{}}, []int{10}, [][]int{{10}})
	ivdompartAddPartTest(t, [][]int{{7}, {10}}, []int{8},
		[][]int{{7, 8}, {10}})
	ivdompartAddPartTest(t, [][]int{{10}}, []int{9}, [][]int{{9, 10}})
	ivdompartAddPartTest(t, [][]int{{9, 10}}, []int{11}, [][]int{{9, 11}})
	ivdompartAddPartTest(t, [][]int{{9, 11}}, []int{11}, [][]int{{9, 11}})
	ivdompartAddPartTest(t, [][]int{{9, 11}}, []int{5}, [][]int{{5}, {9, 11}})
	ivdompartAddPartTest(t, [][]int{{5}, {9, 11}}, []int{7},
		[][]int{{5}, {7}, {9, 11}})
	ivdompartAddPartTest(t, [][]int{{5}, {7}, {9, 11}}, []int{6},
		[][]int{{5, 7}, {9, 11}})
	ivdompartAddPartTest(t, [][]int{{5, 7}, {9, 11}}, []int{20},
		[][]int{{5, 7}, {9, 11}, {20}})
	ivdompartAddPartTest(t, [][]int{{5, 7}, {9, 11}, {20}}, []int{15},
		[][]int{{5, 7}, {9, 11}, {15}, {20}})
	ivdompartAddPartTest(t, [][]int{{1, 4}, {10, 15}}, []int{0, 20},
		[][]int{{0, 20}})
	ivdompartAddPartTest(t, [][]int{{1, 4}, {10, 15}, {20, 40}},
		[]int{0, 100}, [][]int{{0, 100}})
	ivdompartAddPartTest(t, [][]int{{1, 4}, {10, 15}, {20, 40}},
		[]int{9, 16}, [][]int{{1, 4}, {9, 16}, {20, 40}})
	ivdompartAddPartTest(t, [][]int{{1, 4}, {10, 15}, {20, 40}},
		[]int{3, 41}, [][]int{{1, 41}})
	ivdompartAddPartTest(t, [][]int{{1, 4}, {10, 15}, {20, 40}},
		[]int{0, 50}, [][]int{{0, 50}})
	ivdompartAddPartTest(t, [][]int{{1, 4}, {10, 15}, {20, 40}},
		[]int{0, 18}, [][]int{{0, 18}, {20, 40}})
}

func getDomainOutOfBounds_test(t *testing.T, d Domain, min, max int,
	expD Domain) {
	calcD := d.GetDomainOutOfBounds(min, max)
	if !calcD.Equals(expD) {
		t.Errorf("IGetDomainOutOfBounds: got %s, want %s", calcD, expD)
	}
}

func Test_IvGetDomainOutOfBounds(t *testing.T) {
	setup()
	defer teardown()
	log("IvGetDomainOutOfBounds")

	d := CreateIvDomainFromIntArr([]int{1, 2, 3, 4, 5, 10, 11, 12})
	min := 3
	max := 11
	expD := CreateIvDomainFromIntArr([]int{1, 2, 12})
	getDomainOutOfBounds_test(t, d, min, max, expD)

	d = CreateIvDomainFromIntArr([]int{1, 2, 3, 4, 5, 10, 11, 12})
	min = 0
	max = 13
	expD = CreateIvDomain()
	getDomainOutOfBounds_test(t, d, min, max, expD)

	d = CreateIvDomainFromIntArr([]int{1, 2, 3, 4, 5, 10, 11, 12})
	min = 1
	max = 5
	expD = CreateIvDomainFromIntArr([]int{10, 11, 12})
	getDomainOutOfBounds_test(t, d, min, max, expD)

	d = CreateIvDomainFromIntArr([]int{1, 2, 3, 4, 5, 10, 11, 12})
	min = 0
	max = 5
	expD = CreateIvDomainFromIntArr([]int{10, 11, 12})
	getDomainOutOfBounds_test(t, d, min, max, expD)

	d = CreateIvDomainFromIntArr([]int{1, 2, 3, 4, 5, 10, 11, 12})
	min = 10
	max = 12
	expD = CreateIvDomainFromIntArr([]int{1, 2, 3, 4, 5})
	getDomainOutOfBounds_test(t, d, min, max, expD)

	d = CreateIvDomainFromIntArr([]int{1, 2, 3, 4, 5, 10, 11, 12})
	min = 10
	max = 13
	expD = CreateIvDomainFromIntArr([]int{1, 2, 3, 4, 5})
	getDomainOutOfBounds_test(t, d, min, max, expD)

	d = CreateIvDomainFromIntArr([]int{1, 2, 3, 4, 5, 10, 11, 12})
	min = 5
	max = 5
	expD = CreateIvDomainFromIntArr([]int{1, 2, 3, 4, 10, 11, 12})
	getDomainOutOfBounds_test(t, d, min, max, expD)
}

func CreateExDomainUnion_test(t *testing.T, fromTos [][]int,
	expFromTos [][]int) {
	parts := make([]*IvDomPart, 0)
	for _, fromTo := range fromTos {
		from := fromTo[0]
		to := fromTo[0]
		if len(fromTo) == 2 {
			to = fromTo[1]
		}
		parts = append(parts, CreateIvDomPart(from, to))
	}
	calcDom := CreateIvDomainUnion(parts)
	expDom := CreateIvDomainFromTos(expFromTos)
	if !calcDom.Equals(expDom) {
		t.Errorf("CreateIvDomainUnion: got %s, want %s", calcDom, expDom)
	}
}

func Test_CreateIvDomainUnion(t *testing.T) {
	setup()
	defer teardown()
	log("CreateIvDomainUnion")
	CreateExDomainUnion_test(t, [][]int{{12, 13}, {21}, {27}, {21}, {32}},
		[][]int{{12, 13}, {21}, {27}, {32}})
	CreateExDomainUnion_test(t, [][]int{{11, 17}, {7, 11}}, [][]int{{7, 17}})
	CreateExDomainUnion_test(t, [][]int{{11, 17}, {7, 10}}, [][]int{{7, 17}})
	CreateExDomainUnion_test(t, [][]int{{11, 17}, {7, 16}}, [][]int{{7, 17}})
	CreateExDomainUnion_test(t, [][]int{{11, 17}, {7, 17}}, [][]int{{7, 17}})
	CreateExDomainUnion_test(t, [][]int{{11, 17}, {12, 20}}, [][]int{{11, 20}})
	CreateExDomainUnion_test(t, [][]int{{11, 17}, {16, 20}}, [][]int{{11, 20}})
	CreateExDomainUnion_test(t, [][]int{{11, 17}, {17, 20}}, [][]int{{11, 20}})
	CreateExDomainUnion_test(t, [][]int{{11, 17}, {18, 20}}, [][]int{{11, 20}})
	CreateExDomainUnion_test(t, [][]int{{11, 17}, {18, 20}}, [][]int{{11, 20}})
	CreateExDomainUnion_test(t, [][]int{{11, 17}, {12, 16}}, [][]int{{11, 17}})
	CreateExDomainUnion_test(t, [][]int{{11, 17}, {10, 18}}, [][]int{{10, 18}})
	CreateExDomainUnion_test(t, [][]int{{11, 17}, {11, 16}}, [][]int{{11, 17}})
	CreateExDomainUnion_test(t, [][]int{{11, 17}, {12, 17}}, [][]int{{11, 17}})
}

func addDomain_test(t *testing.T, fromTosD1 [][]int, fromTosD2 [][]int,
	expFromTos [][]int) {
	d1 := CreateIvDomainFromTos(fromTosD1)
	d2 := CreateIvDomainFromTos(fromTosD2)
	expD := CreateIvDomainFromTos(expFromTos)
	calcD := d1.ADD(d2)
	if !calcD.Equals(expD) {
		t.Errorf("Domain-ADD failed: %s + %s results in %s, want %s",
			d1, d2, calcD, expD)
	}
}

func Test_DomainADD(t *testing.T) {
	setup()
	defer teardown()
	log("DomainADD")
	addDomain_test(t, [][]int{{1, 2}, {4, 6}}, [][]int{{12, 13}},
		[][]int{{13, 19}}) // {13, 15}, {16, 19}
	addDomain_test(t, [][]int{{1, 2}}, [][]int{{3, 3}},
		[][]int{{4, 5}}) // {4, 5}
	addDomain_test(t, [][]int{{1, 2}, {4, 5}}, [][]int{{3, 3}, {5, 6}},
		[][]int{{4, 11}}) // {4, 5}, {6, 8}, {7, 8}, {9, 11}
	addDomain_test(t, [][]int{{1, 2}, {8, 9}}, [][]int{{1, 1}, {3, 4}},
		[][]int{{2, 6}, {9, 13}}) // {2, 3}, {4, 6}, {9, 10}, {11, 13}
	addDomain_test(t, [][]int{{1, 2}, {9, 10}}, [][]int{{1, 1}, {4, 5}},
		[][]int{{2, 3}, {5, 7}, {10, 11}, {13, 15}})
	// {2, 3}, {5, 6}, {10, 11}, {13, 15}
	// ToDo 5, 6 versus 5,7 ok?
}

func subtractDomainCheck(t *testing.T, d1, d2, expD *IvDomain) {
	calcD := d1.SUBTRACT(d2)
	if !calcD.Equals(expD) {
		t.Errorf("Domain-SUBTRACT failed: %s - %s results in %s, want %s",
			d1, d2, calcD, expD)
	}
}

func subtractDomain_test(t *testing.T, fromTosD1 [][]int, fromTosD2 [][]int,
	expFromTos [][]int, expFromTos2 [][]int) {
	d1 := CreateIvDomainFromTos(fromTosD1)
	d2 := CreateIvDomainFromTos(fromTosD2)
	expD := CreateIvDomainFromTos(expFromTos)
	subtractDomainCheck(t, d1, d2, expD)
	expD = CreateIvDomainFromTos(expFromTos2)
	subtractDomainCheck(t, d2, d1, expD)
}

func Test_DomainSUBTRACT(t *testing.T) {
	setup()
	defer teardown()
	log("SUBTRACT")
	subtractDomain_test(t, [][]int{{1, 2}, {4, 6}}, [][]int{{12, 13}},
		[][]int{{-12, -6}}, [][]int{{6, 12}})
	// {-12,-10}, {-9,-6} ->
	// {10, 12}, {6, 9}	<-

	subtractDomain_test(t, [][]int{{1, 2}}, [][]int{{3, 3}},
		[][]int{{-2, -1}}, [][]int{{1, 2}})

	subtractDomain_test(t, [][]int{{1, 2}, {4, 5}}, [][]int{{3, 3}, {5, 6}},
		[][]int{{-5, 2}}, [][]int{{-2, 5}})
	// {-2,-1}, {-5,-3}, {1, 2}, {-2, 0} ->
	// {1, 2}, {-2,-1}, {3, 5}, {0, 2} <-

	subtractDomain_test(t, [][]int{{1, 2}, {8, 9}}, [][]int{{1, 1}, {3, 4}},
		[][]int{{-3, 1}, {4, 8}}, [][]int{{-8, -4}, {-1, 3}})
	// {0, 1}{-3,-1}{7, 8}{4, 6} ->
	// {-1, 0}{-8,-7}{1, 3}{-6,-4} <-

	subtractDomain_test(t, [][]int{{1, 2}, {9, 10}}, [][]int{{1, 1}, {4, 5}},
		[][]int{{-4, -2}, {0, 1}, {4, 6}, {8, 9}},
		[][]int{{-9, -8}, {-6, -4}, {-1, 0}, {2, 4}})
	// {0, 1}{-4, -2}{8, 9}{4, 6} ->
	// {-1, 0}{-9, -8}{2, 4}{-6,-4} <-
}

func negateD_test(t *testing.T, fromTos [][]int, expFromTos [][]int) {
	d := CreateIvDomainFromTos(fromTos)
	expD := CreateIvDomainFromTos(expFromTos)
	calcD := d.NEGATE()
	if !calcD.Equals(expD) {
		t.Errorf("Domain-NEGATE failed: %s results in %s, want %s",
			d, calcD, expD)
	}
}

func Test_DomainNEGATE(t *testing.T) {
	setup()
	defer teardown()
	log("NEGATE")

	negateD_test(t, [][]int{{0, 3}, {5, 7}, {10, 12}},
		[][]int{{-12, -10}, {-7, -5}, {-3, 0}})
	negateD_test(t, [][]int{{0, 0}, {2, 2}, {4, 4}},
		[][]int{{-4, -4}, {-2, -2}, {0, 0}})
	negateD_test(t, [][]int{{0, 0}, {2, 2}, {4, 4}},
		[][]int{{-4, -4}, {-2, -2}, {0, 0}})
	negateD_test(t, [][]int{{0, 20}}, [][]int{{-20, 0}})
	negateD_test(t, [][]int{{1, 20}}, [][]int{{-20, -1}})
	negateD_test(t, [][]int{{-20, -10}, {-5, 2}, {5, 10}},
		[][]int{{-10, -5}, {-2, 5}, {10, 20}})
	negateD_test(t, [][]int{{-5, 0}}, [][]int{{0, 5}})
	negateD_test(t, [][]int{{-5, 2}}, [][]int{{-2, 5}})
}

func absD_test(t *testing.T, fromTos [][]int, expFromTos [][]int) {
	d := CreateIvDomainFromTos(fromTos)
	expD := CreateIvDomainFromTos(expFromTos)
	calcD := d.ABS()
	if !calcD.Equals(expD) {
		t.Errorf("Domain-ABS failed: %s results in %s, want %s",
			d, calcD, expD)
	}
}

func Test_DomainABS(t *testing.T) {
	setup()
	defer teardown()
	log("ABS")
	absD_test(t, [][]int{{0, 3}, {5, 7}, {10, 12}},
		[][]int{{0, 3}, {5, 7}, {10, 12}})
	absD_test(t, [][]int{{0, 0}, {2, 2}, {4, 4}},
		[][]int{{0, 0}, {2, 2}, {4, 4}})
	absD_test(t, [][]int{{-20, -10}, {-5, 2}, {4, 8}},
		[][]int{{0, 8}, {10, 20}})
	absD_test(t, [][]int{{-20, -10}, {-5, 0}, {4, 8}},
		[][]int{{0, 8}, {10, 20}})
	absD_test(t, [][]int{{-20, -10}, {-5, 0}}, [][]int{{0, 5}, {10, 20}})
	absD_test(t, [][]int{{0, 20}}, [][]int{{0, 20}})
	absD_test(t, [][]int{{1, 20}}, [][]int{{1, 20}})
}

func intersectionTest(t *testing.T, fromTos1 [][]int, fromTos2 [][]int,
	expFromTos [][]int) {
	d1 := CreateIvDomainFromTos(fromTos1)
	d2 := CreateIvDomainFromTos(fromTos2)
	expD := CreateIvDomainFromTos(expFromTos)
	calcD := d1.IntersectionIvDomain(d2)
	if !calcD.Equals(expD) {
		msg := "IntersectionIvDomain: %s intersec %s results in %s, want %s"
		t.Errorf(msg, d1, d2, calcD, expD)
	}
}

func Test_IntersectionIvDomain(t *testing.T) {
	setup()
	defer teardown()
	log("IntersectionIvDomain")
	intersectionTest(t, [][]int{{1, 3}}, [][]int{{2, 4}}, [][]int{{2, 3}})
	intersectionTest(t, [][]int{{1, 3}, {5, 10}, {13, 20}},
		[][]int{{2, 5}, {9, 15}},
		[][]int{{2, 3}, {5, 5}, {9, 10}, {13, 15}})
}

func notTest(t *testing.T, fromTos1 [][]int, expFromTos [][]int) {
	d1 := CreateIvDomainFromTos(fromTos1)
	expD := CreateIvDomainFromTos(expFromTos)
	notD := d1.NOT()
	if !notD.Equals(expD) {
		msg := "NOTIvDomain: not(%s) results in %s, want %s"
		t.Errorf(msg, d1, notD, expD)
	}
}

func Test_NOTIvDomain(t *testing.T) {
	setup()
	defer teardown()
	log("NOTIvDomain")
	notTest(t, [][]int{{1, 3}}, [][]int{{NEG_INFINITY, 0}, {4, INFINITY}})
	notTest(t, [][]int{{1, 3}, {5, 10}, {13, 20}},
		[][]int{{NEG_INFINITY, 0}, {4, 4}, {11, 12}, {21, INFINITY}})
	notTest(t, [][]int{{NEG_INFINITY, 3}, {5, 10}, {13, 20}},
		[][]int{{4, 4}, {11, 12}, {21, INFINITY}})
	notTest(t, [][]int{{1, 3}, {5, 10}, {13, INFINITY}},
		[][]int{{NEG_INFINITY, 0}, {4, 4}, {11, 12}})
	notTest(t, [][]int{{NEG_INFINITY, 3}, {5, 10}, {13, INFINITY}},
		[][]int{{4, 4}, {11, 12}})
}

func Test_AppendIvDomain(t *testing.T) {
	setup()
	defer teardown()
	log("AppendIvDomain")

	appendTest(t, [][]int{{0, 5}, {10, 20}}, []int{25, 50},
		[][]int{{0, 5}, {10, 20}, {25, 50}})
	appendTest(t, [][]int{{0, 5}}, []int{25, 50},
		[][]int{{0, 5}, {25, 50}})
	appendTest(t, [][]int{{0, 1}}, []int{3, 4},
		[][]int{{0, 1}, {3, 4}})
	//fail-cases, ToDo: panic-testing not included yet
	//appendTest(t, [][]int{{0,1}}, []int{2,3}, nil)
	//appendTest(t, [][]int{{2,3}}, []int{0,1}, nil)

}

func appendTest(t *testing.T, fromTos1 [][]int,
	appendIv []int, expFromTos [][]int) {
	d1 := CreateIvDomainFromTos(fromTos1)
	iv := CreateIvDomPart(appendIv[0], appendIv[1])
	expD := CreateIvDomainFromTos(expFromTos)

	dNew := d1.Copy().(*IvDomain)
	dNew.Append(iv)

	if !dNew.Equals(expD) {
		msg := "AppendIvDomain: %s.Append(%s) results in %s, want %s"
		t.Errorf(msg, d1, iv, dNew, expD)
	}
}

func Test_AppendsIvDomain(t *testing.T) {
	setup()
	defer teardown()
	log("AppendIvDomain")

	appendsTest(t, [][]int{{0, 5}, {10, 20}}, [][]int{{25, 50}},
		[][]int{{0, 5}, {10, 20}, {25, 50}})
	appendsTest(t, [][]int{{0, 5}}, [][]int{{25, 50}},
		[][]int{{0, 5}, {25, 50}})
	appendsTest(t, [][]int{{0, 1}}, [][]int{{3, 4}, {7, 9}},
		[][]int{{0, 1}, {3, 4}, {7, 9}})
	appendsTest(t, [][]int{{0, 1}, {3, 4}}, [][]int{{7, 9}, {20, 30}},
		[][]int{{0, 1}, {3, 4}, {7, 9}, {20, 30}})
	//fail-cases, ToDo: panic-testing not included yet
	//appendsTest(t, [][]int{{0,1}}, [][]int{{2,3}}, nil)
	//appendsTest(t, [][]int{{2,3}}, [][]int{{0,1}}, nil)

}

func appendsTest(t *testing.T, fromTos1 [][]int, appendIvs [][]int,
	expFromTos [][]int) {
	d1 := CreateIvDomainFromTos(fromTos1)
	ivs := CreateIvDomainFromTos(appendIvs)
	expD := CreateIvDomainFromTos(expFromTos)

	dNew := d1.Copy().(*IvDomain)
	dNew.Appends(ivs.GetParts())

	if !dNew.Equals(expD) {
		msg := "AppendIvDomain: %s.Append(%s) results in %s, want %s"
		t.Errorf(msg, d1, ivs, dNew, expD)
	}
}
