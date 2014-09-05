package ixterm

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

//weakly tested... more tests are needed

func Test_BasicTermCreation(t *testing.T) {
	setup()          // common first two
	defer teardown() // lines for every test
	log("BasicTermCreation: Max, Min, Val, Dom, Value")

	d := core.CreateIvDomainFromTo(1, 10)
	var varid core.VarId
	varid = 1

	maxterm := CreateMaxTerm(varid, d)
	if !maxterm.HasVarId(varid) || maxterm.GetValue().GetAnyElement() != d.GetMax() {
		t.Errorf("MaxTerm: Creation failed")
	}

	minterm := CreateMinTerm(varid, d)
	if !minterm.HasVarId(varid) || minterm.GetValue().GetAnyElement() != d.GetMin() {
		t.Errorf("MinTerm: Creation failed")
	}

	val := 10
	valueterm := CreateValueTerm(val)
	if valueterm.HasVarId(123) || valueterm.GetValue().GetAnyElement() != val {
		t.Errorf("ValueTerm: Creation failed")
	}

	d = core.CreateIvDomainFromTo(5, 5)
	valTerm := CreateValTerm(varid, d)
	if valTerm.GetValue().GetAnyElement() != 5 {
		t.Errorf("ValTerm: Creation failed")
	}
}

func Test_SumTerm(t *testing.T) {
	setup()          // common first two
	defer teardown() // lines for every test
	log("SumTerm")
	var varid, varid2, varid3 core.VarId

	d := core.CreateIvDomainFromTo(2, 10)
	varid = 1
	d2 := core.CreateIvDomainFromTo(1, 20)
	varid2 = 2
	d3 := core.CreateIvDomainFromTo(0, 5)
	varid3 = 3

	maxT := CreateMaxTerm(varid, d)
	maxT2 := CreateMaxTerm(varid2, d2)
	minT := CreateMinTerm(varid3, d3)

	sumTest(t, []ITerm{maxT, maxT2, minT}, 30)

	d = core.CreateIvDomainFromTo(2, 10)
	varid = 1
	d2 = core.CreateIvDomainFromTo(1, 20)
	varid2 = 2
	d3 = core.CreateIvDomainFromTo(30, 50)
	varid3 = 3

	maxT = CreateMaxTerm(varid, d)
	maxT2 = CreateMaxTerm(varid2, d2)

	sumT := CreateSumTerm(maxT, maxT2)
	minT = CreateMinTerm(varid3, d3)

	subTest(t, minT, sumT, 0)
}

func Test_CreateOperationalTerms(t *testing.T) {
	setup()          // common first two
	defer teardown() // lines for every test
	log("CreateOperationalTerms")

	d := core.CreateIvDomainFromTo(2, 10)
	var varid core.VarId
	varid = 1

	//max(d) op min(d)
	maxT := CreateMaxTerm(varid, d)
	minT := CreateMinTerm(varid, d)
	test(t, maxT, minT)

	//20 op 5
	t1 := CreateValueTerm(20)
	t2 := CreateValueTerm(5)
	test(t, t1, t2)

	//val(d) op 12
	d = core.CreateIvDomainFromTo(5, 5)
	varid = 2
	valt1 := CreateValTerm(varid, d)
	t2 = CreateValueTerm(12)
	test(t, valt1, t2)
}

func test(t *testing.T, t1, t2 ITerm) {
	addTest(t, t1, t2, t1.GetValue().GetAnyElement()+t2.GetValue().GetAnyElement())
	subTest(t, t1, t2, t1.GetValue().GetAnyElement()-t2.GetValue().GetAnyElement())
	multTest(t, t1, t2, t1.GetValue().GetAnyElement()*t2.GetValue().GetAnyElement())
	divTest(t, t1, t2, t1.GetValue().GetAnyElement()/t2.GetValue().GetAnyElement())
}

func addTest(t *testing.T, t1, t2 ITerm, expVal int) {
	addT := CreateAdditionTerm(t1, t2)
	if addT.GetValue().GetAnyElement() != expVal {
		t.Errorf("AdditionTerm '%s+%s': Creation failed", t1.GetValue(), t2.GetValue())
	}
}

func sumTest(t *testing.T, terms []ITerm, expVal int) {
	sumT := CreateSumTerm(terms...)

	if sumT.GetValue().GetAnyElement() != expVal {
		t.Errorf("SumTerm '%s': failed", sumT)
	}
}

func subTest(t *testing.T, t1, t2 ITerm, expVal int) {
	subT := CreateSubtractionTerm(t1, t2)
	if subT.GetValue().GetAnyElement() != expVal {
		t.Errorf("SubtractionTerm '%s-%s': Creation failed", t1.GetValue(), t2.GetValue())
	}
}

func multTest(t *testing.T, t1, t2 ITerm, expVal int) {
	multT := CreateMultiplicationTerm(t1, t2)
	if multT.GetValue().GetAnyElement() != expVal {
		t.Errorf("MultiplicationTerm '%s*%s': Creation failed", t1.GetValue(), t2.GetValue())
	}
}

func divTest(t *testing.T, t1, t2 ITerm, expVal int) {
	divT := CreateDivisionTerm(t1, t2, false)
	if divT.GetValue().GetAnyElement() != expVal {
		t.Errorf("DivisionTerm '%s/%s': Creation failed", t1.GetValue(), t2.GetValue())
	}
}

func Test_TermInfSupTest(t *testing.T) {
	setup()          // common first two
	defer teardown() // lines for every test
	log("TermInfSup: Max, Min, Val, Value, Add, Sub, Sum")

	d := core.CreateIvDomainFromTo(1, 10)
	var varid core.VarId
	varid = 1

	maxterm := CreateMaxTerm(varid, d)
	minterm := CreateMinTerm(varid, d)

	if _, k := maxterm.GetInf().(*MinTerm); !k {
		t.Errorf("MaxTerm.GetInf() failed")
	}
	if _, k := maxterm.GetSup().(*MaxTerm); !k {
		t.Errorf("MaxTerm.GetSup() failed")
	}

	if _, k := minterm.GetInf().(*MinTerm); !k {
		t.Errorf("MinTerm.GetInf() failed")
	}
	if _, k := minterm.GetSup().(*MaxTerm); !k {
		t.Errorf("MinTerm.GetSup() failed")
	}

	val := 10
	valueterm := CreateValueTerm(val)
	if _, k := valueterm.GetInf().(*ValueTerm); !k {
		t.Errorf("ValueTerm.GetInf() failed")
	}
	if _, k := valueterm.GetSup().(*ValueTerm); !k {
		t.Errorf("ValueTerm.GetSup() failed")
	}

	d = core.CreateIvDomainFromTo(5, 5)
	valTerm := CreateValTerm(varid, d)
	if _, k := valTerm.GetInf().(*ValTerm); !k {
		t.Errorf("ValTerm.GetInf() failed")
	}
	if _, k := valTerm.GetSup().(*ValTerm); !k {
		t.Errorf("ValTerm.GetSup() failed")
	}

	// --- Add, Sub, Sum ---
	addT := CreateAdditionTerm(maxterm, valueterm)
	if _, k := addT.GetInf().(*AdditionTerm); !k {
		t.Errorf("AdditionTerm.GetInf() failed")
	}
	if _, k := addT.GetSup().(*AdditionTerm); !k {
		t.Errorf("AdditionTerm.GetSup() failed")
	}

	subT := CreateSubtractionTerm(maxterm, valueterm)
	if _, k := subT.GetInf().(*SubtractionTerm); !k {
		t.Errorf("SubtractionTerm.GetInf() failed")
	}
	if _, k := subT.GetSup().(*SubtractionTerm); !k {
		t.Errorf("SubtractionTerm.GetSup() failed")
	}

	sumT := CreateSumTerm(maxterm, valueterm)
	if _, k := sumT.GetInf().(*SumTerm); !k {
		t.Errorf("SumTerm.GetInf() failed")
	}
	if _, k := sumT.GetSup().(*SumTerm); !k {
		t.Errorf("SumTerm.GetSup() failed")
	}
}
