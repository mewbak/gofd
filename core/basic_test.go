package core

import (
	"testing"
)

// tests basic functionality

func Test_BasicVarNameTest(t *testing.T) {
	setup()
	defer teardown()
	log("BasicVarNameTest")
	X := CreateIntVarValues("_X", store, []int{1, 2, 3})
	if "_X" == store.GetName(X) {
		t.Errorf("CreateIntVarValues: store.GetName() == '%v'"+
			", want != '%v'",
			"_X", "_X")
	}
	Y := CreateIntVarValues("", store, []int{1, 2, 3})
	if "" == store.GetName(Y) {
		t.Errorf("CreateIntVarValues: store.GetName() == '%v'"+
			", want !='%v'",
			"", "")
	}
}

func Test_BasicCreateIntVarValues(t *testing.T) {
	setup()
	defer teardown()
	log("CreateIntVarValues")
	X := CreateIntVarValues("X", store, []int{1, 2, 3})
	ivar, exists := store.GetIntVar(X)
	if exists && ivar.ID != X {
		t.Errorf("CreateIntVarValues: store.GetIntVar() = %v"+
			", want %v",
			ivar.ID, X)
	}
}

func Test_BasicCloneIntVar(t *testing.T) {
	setup()
	defer teardown()
	log("BasicCloneIntVar")
	Xid := CreateIntVarValues("X", store, []int{1, 2, 3})
	X, _ := store.GetIntVar(Xid)
	Xcloned := X.Clone()
	checkVarId(t, X, Xcloned)
	checkDomainEquals(t, X, Xcloned, true)
	X.Domain = CreateExDomainFromTo(0, 9)
	checkDomainEquals(t, X, Xcloned, false)
}

func checkVarId(t *testing.T, X *IntVar, Xcloned *IntVar) {
	if X == Xcloned {
		t.Errorf("checkVarId-object_check: result %v, want %v",
			X == Xcloned, X != Xcloned)
	}
	if X.ID != Xcloned.ID {
		t.Errorf("checkVarId-ID_check: result %v, want %v",
			X.ID != Xcloned.ID, X.ID == Xcloned.ID)
	}
}

func checkDomainEquals(t *testing.T, X *IntVar, Xc *IntVar, exp bool) {
	if X.Domain.Equals(Xc.Domain) != exp {
		t.Errorf("%v.Equals(%v) = %v, want %v",
			X.Domain, Xc.Domain,
			!exp, exp)
	}
}

func Test_BasicCreateIntVarsIvFromTo(t *testing.T) {
	setup()
	defer teardown()
	log("CreateIntVarsIvFromTo")
	var X, Y, Z, Q VarId
	CreateIntVarsIvFromTo([]*VarId{&X, &Y, &Z},
		[]string{"X", "Y", "Z"}, store, 0, 6)
	checkIntVar(t, X)
	checkIntVar(t, Y)
	checkIntVar(t, Z)
	CreateIntVarsIvFromTo([]*VarId{&Q},
		[]string{"Q"}, store, 0, 6)
	checkIntVar(t, Q)
	checkIntVar(t, Z)
}

func checkIntVar(t *testing.T, VAR VarId) {
	ivar, exists := store.GetIntVar(VAR)
	if exists && ivar.ID != VAR {
		t.Errorf("CreateIntVarsIvFromTo: store.GetIntVar() = %v"+
			", want %v",
			ivar.ID, VAR)
	}
}
