package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func clone_test(t *testing.T, prop core.Constraint) {
	prop1Cloned := prop.Clone()
	prop1ClonedCloned := prop1Cloned.Clone()

	test(t, core.Constraint(prop), prop1Cloned, "C1XplusC2YeqC3-cloned")
	test(t, prop1ClonedCloned, prop1Cloned, "C1XplusC2YeqC3-clone-cloned")

	//	if prop2!=nil{
	//		prop2Cloned:=prop2.Clone()
	//		prop2Cloned2:=prop2.Clone()
	//
	//		test(t,prop2Cloned,prop2,"C1XplusC2YeqC3-right")
	//		test(t,prop2Cloned2,prop2,"C1XplusC2YeqC3-right")
	//	}
}

//func Test_Clone(t *testing.T) {
//	setup()
//	defer teardown()
//
//	log("Clone-Test for Propagators (ToDo: Equals on Propagators):")
//
//	c1XplusC2YeqC3_Clone(t)
//	xEqC_Clone(t)
//	xGtC_Clone(t)
//	xGtY_Clone(t)
//	xPlusCEqY_Clone(t)
//	xPlusCNeqY_Clone(t)
//}

func c1XplusC2YeqC3_Clone(t *testing.T) {
	log("C1XplusC2YeqC3")
	//
	//	X_ForCloneTest = core.CreateIntVarExFromTo("X", store, 0, 9)
	//	Y_ForCloneTest = core.CreateIntVarExFromTo("Y", store, 0, 9)
	//
	//	prop := CreateC1XplusC2YeqC3(1, X_ForCloneTest, 2, Y_ForCloneTest, 12)

	//	if clonedProp, ok := prop.Clone().(C1XplusC2YeqC3); ok {
	//		equals := true
	//		if prop.c1 != clonedProp.c1 || prop.c2 != clonedProp.c2 || prop.c3 != clonedProp.c3 {
	//			equals = false
	//		}
	//		if prop.x != clonedProp.x || prop.y != clonedProp.y {
	//			equals = false
	//		}
	//
	//		if !equals {
	//			t.Errorf("Clone_test: C1XplusC2YeqC3 equals cloned C1XplusC2YeqC3 == %v, want %v",
	//				equals, true)
	//		}
	//	}

	//	clone_test(t, prop)
}

func xEqC_Clone(t *testing.T) {
	log("XeqC")
	prop := CreateXeqC(X_ForCloneTest, 2)
	clone_test(t, prop)
}

func xGtC_Clone(t *testing.T) {
	log("XgtC")
	prop := CreateXgtC(X_ForCloneTest, 2)
	clone_test(t, prop)
}

func xGtY_Clone(t *testing.T) {
	log("XgtY")
	prop := CreateXgtY(X_ForCloneTest, Y_ForCloneTest)
	clone_test(t, prop)
}

func xPlusCEqY_Clone(t *testing.T) {
	log("XplusCeqY")
	prop := CreateXplusCeqY(X_ForCloneTest, 5, Y_ForCloneTest)
	clone_test(t, prop)
}

func xPlusCNeqY_Clone(t *testing.T) {
	log("XplusCneqY")
	prop := CreateXplusCneqY(X_ForCloneTest, 5, Y_ForCloneTest)
	clone_test(t, prop)
}

func test(t *testing.T, prop core.Constraint, propCloned core.Constraint, which_propagator string) {
	if prop == propCloned {
		t.Errorf("Clone_test: %v prop==propCloned is %v, want %v",
			which_propagator, prop == propCloned, prop != propCloned)
	}
}
