package ixrange

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"testing"
)

// ToDo
// ProcessingTest integrated, missing: GetValue Test
// Evaluable-Test

func Test_DomRange(t *testing.T) {
	setup()
	defer teardown()
	log("DomRange and AddRange")

	//1..10 in 1..10+5
	domRange_test(t, []int{1, 10}, []int{1, 10}, 5, []int{6, 10})
	//1..10 in 1..10+0
	domRange_test(t, []int{1, 10}, []int{1, 10}, 0, []int{1, 10})
	//1..10 in 1..10-5
	domRange_test(t, []int{1, 10}, []int{1, 10}, -5, []int{1, 5})
	//1..10 in 20..30+2
	domRange_test(t, []int{1, 10}, []int{20, 30}, 2, []int{})
}

func domRange_test(t *testing.T, outVar []int, inVar []int, c int, expOut []int) {
	outVarD := core.CreateIvDomainFromTo(outVar[0], outVar[1])
	inVarD := core.CreateIvDomainFromTo(inVar[0], inVar[1])
	var expOutD *core.IvDomain
	if len(expOut) == 2 {
		expOutD = core.CreateIvDomainFromTo(expOut[0], expOut[1])
	} else {
		expOutD = core.CreateIvDomain()
	}

	domR := CreateDomRange(1, inVarD)
	valueT := ixterm.CreateValueTerm(c)
	valueR := CreateSingleValueRange(valueT)

	addR := CreateAddRange(domR, valueR)

	processing_test(t, addR, outVarD, expOutD)
}

func processing_test(t *testing.T, r IRange, outVarD *core.IvDomain, expOutD *core.IvDomain) {
	outVarBefore := outVarD.Copy()
	parts := r.Process(outVarD)
	removingD := core.CreateIvDomainDomParts(parts)
	outVarD.Removes(removingD)

	if !outVarD.Equals(expOutD) {
		t.Errorf("Range-Processing failed: %s in %s, "+
			"resulted in %s, wanted %s", outVarBefore, r, outVarD, expOutD)
	}
}

func simpleFromToRange_test(t *testing.T, fromToOutVar []int, fromToRange []int, expFromTo []int) {
	outVarD := core.CreateIvDomainFromTo(fromToOutVar[0], fromToOutVar[1])
	fromTerm := ixterm.CreateValueTerm(fromToRange[0])
	toTerm := ixterm.CreateValueTerm(fromToRange[1])
	r := CreateFromToRange(fromTerm, toTerm)
	expOutD := core.CreateIvDomainFromTo(expFromTo[0], expFromTo[1])

	processing_test(t, r, outVarD, expOutD)
}

func extendedFromToRange_test(t *testing.T, fromToOutVar []int, from int, toTerm ixterm.ITerm, addV int, expOut []int) {
	outVarD := core.CreateIvDomainFromTo(fromToOutVar[0], fromToOutVar[1])
	fromTerm := ixterm.CreateValueTerm(from)

	var term ixterm.ITerm
	if addV < 0 {
		addV = -addV //--5
		valueT := ixterm.CreateValueTerm(addV)
		term = ixterm.CreateSubtractionTerm(toTerm, valueT)
	} else {
		valueT := ixterm.CreateValueTerm(addV)
		term = ixterm.CreateAdditionTerm(toTerm, valueT)
	}

	expOutD := core.CreateIvDomain()
	if expOut != nil {
		expOutD = core.CreateIvDomainFromTo(expOut[0], expOut[1])
	}
	r := CreateFromToRange(fromTerm, term)
	processing_test(t, r, outVarD, expOutD)
}

func Test_FromToRange(t *testing.T) {
	setup()          // common first two
	defer teardown() // lines for every test
	log("FromToRange")

	//0..100 in [50,60]
	simpleFromToRange_test(t, []int{0, 100}, []int{50, 60}, []int{50, 60})
	//0..100 in [50,120]
	//0..100 in [0,60]
	simpleFromToRange_test(t, []int{0, 100}, []int{50, 120}, []int{50, 100})
	simpleFromToRange_test(t, []int{0, 100}, []int{0, 60}, []int{0, 60})

	//0..100 in 50..max([80..90])+5
	inVarD := core.CreateIvDomainFromTo(80, 90)
	toTermMax := ixterm.CreateMaxTerm(1, inVarD)
	extendedFromToRange_test(t, []int{0, 100}, 50, toTermMax, 5, []int{50, 95})

	//0..100 in 50..min([80..90])-5
	inVarD = core.CreateIvDomainFromTo(80, 90)
	toTermMin := ixterm.CreateMinTerm(1, inVarD)
	extendedFromToRange_test(t, []int{0, 100}, 50, toTermMin, -5, []int{50, 75})

	//0..100 in 110..min([120..190])-5
	inVarD = core.CreateIvDomainFromTo(120, 190)
	toTermMin = ixterm.CreateMinTerm(1, inVarD)
	extendedFromToRange_test(t, []int{0, 100}, 110, toTermMin, -5, nil)
}

func singleValueRange_test(t *testing.T, fromTo []int, value int, expFromTo []int) {

	outVarD := core.CreateIvDomainFromTo(fromTo[0], fromTo[1])
	valTerm := ixterm.CreateValueTerm(value)
	expOutD := core.CreateIvDomain()
	if expFromTo != nil {
		expOutD = core.CreateIvDomainFromTo(expFromTo[0], expFromTo[1])
	}

	r := CreateSingleValueRange(valTerm)
	processing_test(t, r, outVarD, expOutD)
}

func Test_SingleValueRange(t *testing.T) {
	setup()          // common first two
	defer teardown() // lines for every test
	log("SingleValueRange")

	singleValueRange_test(t, []int{0, 100}, 50, []int{50, 50})
	singleValueRange_test(t, []int{0, 100}, 0, []int{0, 0})
	singleValueRange_test(t, []int{0, 100}, 100, []int{100, 100})
	singleValueRange_test(t, []int{0, 100}, 105, nil)
}

func addTest(t *testing.T, outVarFromTo, fromToT1, fromToT2, expOutVarFromTo [][]int) {
	d1 := core.CreateIvDomainFromTos(fromToT1)
	d2 := core.CreateIvDomainFromTos(fromToT2)

	dR1 := CreateDomRange(0, d1)
	dR2 := CreateDomRange(1, d2)

	addR := CreateAddRange(dR1, dR2)

	outVarD := core.CreateIvDomainFromTos(outVarFromTo)
	expOutVarD := core.CreateIvDomainFromTos(expOutVarFromTo)

	processing_test(t, addR, outVarD, expOutVarD)
}

func addTest2(t *testing.T, outVarFromTo []int, fromToT1 []int, toTerm ixterm.ITerm, expOutVarFromTo []int) {
	d1 := core.CreateIvDomainFromTo(fromToT1[0], fromToT1[1])

	dR1 := CreateDomRange(0, d1)
	singleValueRange := CreateSingleValueRange(toTerm)

	addR := CreateAddRange(dR1, singleValueRange)

	outVarD := core.CreateIvDomainFromTo(outVarFromTo[0], outVarFromTo[1])
	expOutVarD := core.CreateIvDomainFromTo(expOutVarFromTo[0], expOutVarFromTo[1])

	processing_test(t, addR, outVarD, expOutVarD)
}

func Test_AddRange(t *testing.T) {
	setup()          // common first two
	defer teardown() // lines for every test
	log("AddRange")

	//1..100 in dom([1..10])+dom([2..5])
	addTest(t, [][]int{{1, 100}}, [][]int{{1, 10}}, [][]int{{2, 5}}, [][]int{{3, 15}})
	//1..100 in dom([1..10])+dom([3,3])
	addTest(t, [][]int{{1, 100}}, [][]int{{1, 10}}, [][]int{{3, 3}}, [][]int{{4, 13}})
	//1..100 in dom([10,10])+dom([2..5])
	addTest(t, [][]int{{1, 100}}, [][]int{{10, 10}}, [][]int{{2, 5}}, [][]int{{12, 15}})
	//1..100 in dom([10,10])+dom([3,3])
	addTest(t, [][]int{{1, 100}}, [][]int{{10, 10}}, [][]int{{3, 3}}, [][]int{{13, 13}})

	//d1,d2,exp: d1+d2=exp
	addTest(t, [][]int{{1, 100}}, [][]int{{1, 2}, {4, 6}}, [][]int{{12, 13}},
		[][]int{{13, 19}}) //{13,15},{16,19}
	addTest(t, [][]int{{1, 100}}, [][]int{{1, 2}}, [][]int{{3, 3}},
		[][]int{{4, 5}}) //{4,5}
	addTest(t, [][]int{{1, 100}}, [][]int{{1, 2}, {4, 5}}, [][]int{{3, 3}, {5, 6}},
		[][]int{{4, 11}}) //{4,5},{6,8},{7,8},{9,11}
	addTest(t, [][]int{{1, 100}}, [][]int{{1, 2}, {8, 9}}, [][]int{{1, 1}, {3, 4}},
		[][]int{{2, 6}, {9, 13}}) //{2,3},{4,6},{9,10},{11,13}
	addTest(t, [][]int{{1, 100}}, [][]int{{1, 2}, {9, 10}}, [][]int{{1, 1}, {4, 5}},
		[][]int{{2, 3}, {5, 7}, {10, 11}, {13, 15}}) //{2,3},{5,6}, {10,11},{13,15}

	//1..100 in dom([10,10])+singlevalueRange(3)
	toTerm := ixterm.CreateValueTerm(3)
	addTest2(t, []int{1, 100}, []int{10, 10}, toTerm, []int{13, 13})

	//1..100 in dom([10,10])+singlevalueRange(val([3,3]))
	d := core.CreateIvDomainFromTo(3, 3)
	toValTerm := ixterm.CreateValTerm(1, d)
	addTest2(t, []int{1, 100}, []int{10, 10}, toValTerm, []int{13, 13})
}

func subTest(t *testing.T, outVarFromTo, fromToT1, fromToT2, expOutVarFromTo, expOutVarFromToReverse [][]int) {
	outVarD := core.CreateIvDomainFromTos(outVarFromTo)
	d1 := core.CreateIvDomainFromTos(fromToT1)
	d2 := core.CreateIvDomainFromTos(fromToT2)

	outVarDc := outVarD.Copy().(*core.IvDomain)

	dR1 := CreateDomRange(0, d1)
	dR2 := CreateDomRange(1, d2)

	subR := CreateSubRange(dR1, dR2)
	expOutVarD := core.CreateIvDomainFromTos(expOutVarFromTo)
	processing_test(t, subR, outVarDc, expOutVarD)

	dR1 = CreateDomRange(0, d1)
	dR2 = CreateDomRange(1, d2)

	subR = CreateSubRange(dR2, dR1)
	expOutVarD = core.CreateIvDomainFromTos(expOutVarFromToReverse)
	processing_test(t, subR, outVarD, expOutVarD)
}

func subTest2(t *testing.T, outVarFromTo [][]int, fromToT1 [][]int, toTerm ixterm.ITerm, expOutVarFromTo [][]int, expOutVarFromToReverse [][]int) {
	outVarD := core.CreateIvDomainFromTos(outVarFromTo)
	d1 := core.CreateIvDomainFromTos(fromToT1)

	outVarDc := outVarD.Copy().(*core.IvDomain)

	dR1 := CreateDomRange(0, d1)
	singleValueRange := CreateSingleValueRange(toTerm)

	subR := CreateSubRange(dR1, singleValueRange)
	expOutVarD := core.CreateIvDomainFromTos(expOutVarFromTo)
	processing_test(t, subR, outVarDc, expOutVarD)

	subR = CreateSubRange(singleValueRange, dR1)
	expOutVarD = core.CreateIvDomainFromTos(expOutVarFromToReverse)
	processing_test(t, subR, outVarD, expOutVarD)
}

func Test_SubRange(t *testing.T) {
	setup()          // common first two
	defer teardown() // lines for every test
	log("SubRange")

	//-100..100 in dom([10..20])-dom([2..5]), -100..100 in dom([2..5])-dom([10..20])
	subTest(t, [][]int{{-100, 100}}, [][]int{{10, 20}}, [][]int{{2, 5}},
		[][]int{{5, 18}}, [][]int{{-18, -5}})
	//-100..100 in dom([10..20])-dom([3,3]), -100..100 in dom([3,3])-dom([10..20])
	subTest(t, [][]int{{-100, 100}}, [][]int{{10, 20}}, [][]int{{3, 3}},
		[][]int{{7, 17}}, [][]int{{-17, -7}})
	//-100..100 in dom([10,10])-dom([2..5]), -100..100 in dom([2..5])-dom([10,10])
	subTest(t, [][]int{{-100, 100}}, [][]int{{10, 10}}, [][]int{{2, 5}},
		[][]int{{5, 8}}, [][]int{{-8, -5}})
	//-100..100 in dom([10,10])-dom([3,3]), -100..100 in dom([3,3])-dom([10,10])
	subTest(t, [][]int{{-100, 100}}, [][]int{{10, 10}}, [][]int{{3, 3}},
		[][]int{{7, 7}}, [][]int{{-7, -7}})

	//...
	subTest(t, [][]int{{-100, 100}}, [][]int{{1, 2}, {4, 6}}, [][]int{{12, 13}},
		[][]int{{-12, -6}}, [][]int{{6, 12}})
	//{-12,-10},{-9,-6} ->
	//{10,12},{6,9}	<-

	subTest(t, [][]int{{-100, 100}}, [][]int{{1, 2}}, [][]int{{3, 3}},
		[][]int{{-2, -1}}, [][]int{{1, 2}})

	subTest(t, [][]int{{1 - 100, 100}}, [][]int{{1, 2}, {4, 5}}, [][]int{{3, 3}, {5, 6}},
		[][]int{{-5, 2}}, [][]int{{-2, 5}})
	//{-2,-1},{-5,-3},{1,2},{-2,0} ->
	//{1,2},{-2,-1},{3,5},{0,2} <-

	subTest(t, [][]int{{-100, 100}}, [][]int{{1, 2}, {8, 9}}, [][]int{{1, 1}, {3, 4}},
		[][]int{{-3, 1}, {4, 8}}, [][]int{{-8, -4}, {-1, 3}})
	//{0,1}{-3,-1}{7,8}{4,6} ->
	//{-1,0}{-8,-7}{1,3}{-6,-4} <-

	subTest(t, [][]int{{-100, 100}}, [][]int{{1, 2}, {9, 10}}, [][]int{{1, 1}, {4, 5}},
		[][]int{{-4, -2}, {0, 1}, {4, 6}, {8, 9}}, [][]int{{-9, -8}, {-6, -4}, {-1, 0}, {2, 4}})
	//{0,1}{-4,-2}{8,9}{4,6} ->
	//{-1,0}{-9,-8}{2,4}{-6,-4} <-*/

	//1..100 in dom([10,10])-singlevalueRange(3), 1..100 in singlevalueRange(3)-dom([10,10])
	toTerm := ixterm.CreateValueTerm(3)
	subTest2(t, [][]int{{-100, 100}}, [][]int{{10, 10}}, toTerm, [][]int{{7, 7}}, [][]int{{-7, -7}})
	//1..100 in dom([10,10])-singlevalueRange(val([3,3])), 1..100 in singlevalueRange(val([3,3]))-dom([10,10])
	d := core.CreateIvDomainFromTo(3, 3)
	toValTerm := ixterm.CreateValTerm(1, d)
	subTest2(t, [][]int{{-100, 100}}, [][]int{{10, 10}}, toValTerm, [][]int{{7, 7}}, [][]int{{-7, -7}})
}

func multTest(t *testing.T, outVarFromTo []int, fromToT1 []int, fromToT2 []int, expOutVarFromTo [][]int) {
	d1 := core.CreateIvDomainFromTo(fromToT1[0], fromToT1[1])
	d2 := core.CreateIvDomainFromTo(fromToT2[0], fromToT2[1])

	dR1 := CreateDomRange(0, d1)
	dR2 := CreateDomRange(1, d2)

	multR := CreateMultRange(dR1, dR2)

	outVarD := core.CreateIvDomainFromTo(outVarFromTo[0], outVarFromTo[1])
	expOutVarD := core.CreateIvDomainFromTos(expOutVarFromTo)

	processing_test(t, multR, outVarD, expOutVarD)
}

func multTest2(t *testing.T, outVarFromTo []int, fromToT1 []int, toTerm ixterm.ITerm, expOutVarFromTo [][]int) {
	d1 := core.CreateIvDomainFromTo(fromToT1[0], fromToT1[1])

	dR1 := CreateDomRange(0, d1)
	singleValueRange := CreateSingleValueRange(toTerm)

	multR := CreateMultRange(dR1, singleValueRange)

	outVarD := core.CreateIvDomainFromTo(outVarFromTo[0], outVarFromTo[1])
	expOutVarD := core.CreateIvDomainFromTos(expOutVarFromTo)

	processing_test(t, multR, outVarD, expOutVarD)
}

func divTest(t *testing.T, outVarFromTo []int, fromToT1 []int, fromToT2 []int, expOutVarFromTo []int) {
	d1 := core.CreateIvDomainFromTo(fromToT1[0], fromToT1[1])
	d2 := core.CreateIvDomainFromTo(fromToT2[0], fromToT2[1])

	dR1 := CreateDomRange(0, d1)
	dR2 := CreateDomRange(1, d2)

	divR := CreateDivRange(dR1, dR2)

	outVarD := core.CreateIvDomainFromTo(outVarFromTo[0], outVarFromTo[1])
	expOutVarD := core.CreateIvDomain()
	if expOutVarFromTo != nil {
		expOutVarD = core.CreateIvDomainFromTo(expOutVarFromTo[0], expOutVarFromTo[1])
	}

	processing_test(t, divR, outVarD, expOutVarD)
}

func divTest2(t *testing.T, outVarFromTo []int, fromToT1 []int, toTerm ixterm.ITerm, expOutVarFromTo []int) {
	d1 := core.CreateIvDomainFromTo(fromToT1[0], fromToT1[1])

	dR1 := CreateDomRange(0, d1)
	singleValueRange := CreateSingleValueRange(toTerm)

	divR := CreateDivRange(dR1, singleValueRange)

	outVarD := core.CreateIvDomainFromTo(outVarFromTo[0], outVarFromTo[1])
	expOutVarD := core.CreateIvDomainFromTo(expOutVarFromTo[0], expOutVarFromTo[1])

	processing_test(t, divR, outVarD, expOutVarD)
}

func Test_MultRange(t *testing.T) {
	setup()          // common first two
	defer teardown() // lines for every test
	log("MultRange")

	//1..100 in dom([1..5])*dom([3,3])
	multTest(t, []int{1, 100}, []int{1, 5}, []int{3, 3}, [][]int{{3, 3}, {6, 6}, {9, 9}, {12, 12}, {15, 15}})
	//1..100 in dom([5,5])*dom([3,3])
	multTest(t, []int{1, 100}, []int{5, 5}, []int{3, 3}, [][]int{{15, 15}})
	//1..100 in dom([1..5])*singlevalueRange(val([3,3]))
	d := core.CreateIvDomainFromTo(3, 3)
	toValTerm := ixterm.CreateValTerm(1, d)
	multTest2(t, []int{1, 100}, []int{1, 5}, toValTerm, [][]int{{3, 3}, {6, 6}, {9, 9}, {12, 12}, {15, 15}})

}

func Test_DivRange(t *testing.T) {
	setup()          // common first two
	defer teardown() // lines for every test
	log("DivRange")

	//1..100 in dom([10..20])/dom([3,3])
	divTest(t, []int{1, 100}, []int{10, 20}, []int{3, 3}, []int{4, 6})
	//1..100 in dom([9,9])/dom([3,3])
	divTest(t, []int{1, 100}, []int{9, 9}, []int{3, 3}, []int{3, 3})
	//1..100 in dom([5,5])/dom([3,3])
	divTest(t, []int{1, 100}, []int{5, 5}, []int{3, 3}, nil)

	//1..100 in dom([1..5])/singlevalueRange(val([3,3]))
	d := core.CreateIvDomainFromTo(3, 3)
	toValTerm := ixterm.CreateValTerm(1, d)
	divTest2(t, []int{1, 100}, []int{1, 5}, toValTerm, []int{1, 1})
}

func negTest(t *testing.T, outFromTos, inFromTos, expOutFromTos [][]int) {
	outD := core.CreateIvDomainFromTos(outFromTos)
	inD := core.CreateIvDomainFromTos(inFromTos)
	expOutD := core.CreateIvDomainFromTos(expOutFromTos)

	domD := CreateDomRange(0, inD)
	negD := CreateInverseRange(domD)

	processing_test(t, negD, outD, expOutD)
}

func Test_InverseRange(t *testing.T) {
	setup()          // common first two
	defer teardown() // lines for every test
	log("InverseRange")

	//1..100 in neg(-20..5)
	negTest(t, [][]int{{1, 100}}, [][]int{{-20, 5}}, [][]int{{1, 20}})

	//1..100 in neg(-20..20)
	negTest(t, [][]int{{1, 100}}, [][]int{{-20, 20}}, [][]int{{1, 20}})

	//5..30 in neg(-5..20)
	negTest(t, [][]int{{5, 30}}, [][]int{{-5, 20}}, [][]int{{5, 5}})

	//5..30 in neg(0..20)
	negTest(t, [][]int{{5, 30}}, [][]int{{0, 20}}, [][]int{{}})
}

func absTest(t *testing.T, outFromTos, inFromTos, expOutFromTos [][]int) {
	outD := core.CreateIvDomainFromTos(outFromTos)
	inD := core.CreateIvDomainFromTos(inFromTos)
	expOutD := core.CreateIvDomainFromTos(expOutFromTos)

	domD := CreateDomRange(0, inD)
	negD := CreateAbsRange(domD)

	processing_test(t, negD, outD, expOutD)
}

func Test_AbsRange(t *testing.T) {
	setup()          // common first two
	defer teardown() // lines for every test
	log("AbsRange")

	//1..100 in abs(-20..5)
	absTest(t, [][]int{{1, 100}}, [][]int{{-20, 5}}, [][]int{{1, 20}})

	//1..100 in abs(-20..20)
	absTest(t, [][]int{{1, 100}}, [][]int{{-20, 20}}, [][]int{{1, 20}})

	//5..30 in abs(-5..20)
	absTest(t, [][]int{{5, 30}}, [][]int{{-5, 20}}, [][]int{{5, 20}})

	//5..30 in abs(0..20)
	absTest(t, [][]int{{5, 30}}, [][]int{{0, 20}}, [][]int{{5, 20}})

	//5..30 in abs(-4..5)
	absTest(t, [][]int{{5, 30}}, [][]int{{-4, 5}}, [][]int{{5, 5}})

	//5..30 in abs(-4..4)
	absTest(t, [][]int{{5, 30}}, [][]int{{-4, 4}}, [][]int{{}})
}

func unionTest(t *testing.T, outFromTos [][]int, inRanges []IRange, expOutFromTos [][]int) {
	outD := core.CreateIvDomainFromTos(outFromTos)
	expOutD := core.CreateIvDomainFromTos(expOutFromTos)

	unionRange := CreateUnionRange(inRanges...)

	processing_test(t, unionRange, outD, expOutD)
}

// ToDo: more tests for union
func Test_UnionRange(t *testing.T) {
	setup()
	defer teardown()
	log("UnionRange")
	//1..100 in -20..5, 10..100
	r1 := CreateFromToRangeInts(-20, 5)
	r2 := CreateFromToRangeInts(10, 100)
	unionTest(t, [][]int{{1, 100}}, []IRange{r1, r2}, [][]int{{1, 5}, {10, 100}})
}

func notTest(t *testing.T, outFromTos [][]int, inRange IRange, expOutFromTos [][]int) {
	outD := core.CreateIvDomainFromTos(outFromTos)
	expOutD := core.CreateIvDomainFromTos(expOutFromTos)
	NotRange := CreateNotRange(inRange)
	processing_test(t, NotRange, outD, expOutD)
}

func Test_NotRange(t *testing.T) {
	setup()
	defer teardown()
	log("NotRange")
	d := core.CreateIvDomainFromTo(1, 5)
	r := CreateDomRange(1, d)
	notTest(t, [][]int{{0, 10}}, r, [][]int{{0, 0}, {6, 10}})
	rft := CreateFromToRangeInts(1, 5)
	notTest(t, [][]int{{0, 10}}, rft, [][]int{{0, 0}, {6, 10}})
}
