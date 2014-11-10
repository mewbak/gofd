package core

import (
	"testing"
)

func checkVarNamesRegistryStore(t *testing.T, msg, got, want string) {
	if want != got {
		t.Error("RegistryStore %s: got %s, want %s",
			msg, want, got)
	}
}

func Test_GetVarIdToNameMap(t *testing.T) {
	setup()
	defer teardown()
	log("RegistryStoreGetVarIdToNameMap")
	rs := CreateRegistryStore()
	rs.SetVarName(1, "X")
	rs.SetVarName(2, "Y")
	rs.SetVarName(3, "Z")
	rsIdToName := rs.GetVarIdToNameMap()
	want := 3
	if len(rsIdToName) != want {
		t.Errorf("RegistryStore count varnames got %d, want %d",
			len(rsIdToName), want)
	}
	msg := "inserted"
	checkVarNamesRegistryStore(t, msg, rsIdToName[1], "X")
	checkVarNamesRegistryStore(t, msg, rsIdToName[2], "Y")
	checkVarNamesRegistryStore(t, msg, rsIdToName[3], "Z")
}

func Test_GetSetHasName(t *testing.T) {
	setup()
	defer teardown()
	log("RegistryStoreGetSetHasName")
	rs := CreateRegistryStore()
	rs.SetVarName(1, "X")
	rs.SetVarName(2, "Y")
	rs.SetVarName(3, "Z")
	msg := "SetGetVarName"
	checkVarNamesRegistryStore(t, msg, rs.GetVarName(1), "X")
	checkVarNamesRegistryStore(t, msg, rs.GetVarName(2), "Y")
	checkVarNamesRegistryStore(t, msg, rs.GetVarName(3), "Z")
	v, k := rs.HasVarName(1)
	if !k || v != "X" {
		t.Errorf("RegistryStore.HasName(%s) not found", "X")
	}
}

func Test_RemoveFixedRelations0(t *testing.T) {
	setup()
	defer teardown()
	log("RegistryStoreRemoveFixedRelations0")

	rs := CreateRegistryStore()
	// setup... X=0,Y=1,Z=2 and constraint X+Y=Z
	varids := []VarId{0, 1, 2}
	c1, _ := createPropagatorDummy(varids, t)
	rs.constraints[c1.GetID()] = c1
	interestedVarids := make([]VarId, len(varids))
	copy(interestedVarids, varids)
	writeChannel := make(chan *ChangeEntry, 10)

	rs.RegisterVarIdWithConstraint(c1.GetID(), writeChannel,
		varids, interestedVarids)

	constraintData := rs.varIdsToConstraints[0][0]

	// now, X will be fixed (then, no constraint should listen on X and
	// c1 should only listen to Y, Z)
	rs.RemoveFixedRelations(0) //supposed, X is fixed
	m := make(map[*ConstraintData]int)
	m[constraintData] = 2             // c1 interested in Y,Z anymore (len=2)
	testRegistrations(t, rs, 1, m, 0) // numberConstraints:1, fixedVar:0

	rs.RemoveFixedRelations(2)
	m = make(map[*ConstraintData]int)
	m[constraintData] = 1
	testRegistrations(t, rs, 1, m, 2)

	rs.RemoveFixedRelations(1)
	m = make(map[*ConstraintData]int)
	m[constraintData] = 0
	testRegistrations(t, rs, 0, m, 1)
}

func Test_RemoveFixedRelations1(t *testing.T) {
	setup()
	defer teardown()
	log("RegistryStoreRemoveFixedRelations1")

	rs := CreateRegistryStore()
	//setup... with e.g. X=0, Y=1, Z=2, Q=3, P=4
	//e.g. constraint X+Y=Z
	varids := []VarId{0, 1, 2}
	c1, _ := createPropagatorDummy(varids, t)
	rs.constraints[c1.GetID()] = c1
	interestedVarids := make([]VarId, len(varids))
	copy(interestedVarids, varids)
	writeChannel := make(chan *ChangeEntry, 10)
	rs.RegisterVarIdWithConstraint(c1.GetID(), writeChannel,
		varids, interestedVarids)

	//e.g. constraint Y+Z+Q=P
	varids = []VarId{1, 2, 3, 4}
	c2, _ := createPropagatorDummy(varids, t)
	rs.constraints[c2.GetID()] = c2
	interestedVarids = make([]VarId, len(varids))
	copy(interestedVarids, varids)
	writeChannel = make(chan *ChangeEntry, 10)
	rs.RegisterVarIdWithConstraint(c2.GetID(), writeChannel,
		varids, interestedVarids)

	constraintData1 := rs.varIdsToConstraints[0][0]
	constraintData2 := rs.varIdsToConstraints[4][0]

	// now, X will be fixed (then, no constraint should listen on X and
	// c1 should only listen to Y, Z)
	rs.RemoveFixedRelations(0)
	m := make(map[*ConstraintData]int)
	m[constraintData1] = 2
	m[constraintData2] = 4
	testRegistrations(t, rs, 2, m, 0)

	rs.RemoveFixedRelations(2)
	m = make(map[*ConstraintData]int)
	m[constraintData1] = 1
	m[constraintData2] = 3
	testRegistrations(t, rs, 2, m, 2)

	rs.RemoveFixedRelations(1)
	m = make(map[*ConstraintData]int)
	m[constraintData1] = 0
	m[constraintData2] = 2
	testRegistrations(t, rs, 1, m, 1)

	rs.RemoveFixedRelations(3)
	m = make(map[*ConstraintData]int)
	m[constraintData2] = 1
	testRegistrations(t, rs, 1, m, 3)

	rs.RemoveFixedRelations(4)
	m = make(map[*ConstraintData]int)
	m[constraintData2] = 0
	testRegistrations(t, rs, 0, m, 4)
}

func testRegistrations(t *testing.T, rs *RegistryStore,
	expNumberConstraints int, expNumberVaridsPerPropId map[*ConstraintData]int,
	removedVarId VarId) {
	if len(rs.constraintsToVarIds) != expNumberConstraints {
		t.Errorf("RegistryStore.constraintsToVarIds has wrong number " +
			"constraints")
	}
	for cd, numberVarids := range expNumberVaridsPerPropId {
		if len(rs.constraintsToVarIds[cd]) != numberVarids {
			t.Errorf("RegistryStore.constraintsToVarIds[cd] has wrong number " +
				"of varids")
		}
		if numberVarids == 0 {
			// expected result: closed writeChannel, constraintsToVarIds
			// reduced (cd not contained anymore) and constraints reduced
			// (propId not contained anymore)
			if _, k := rs.constraintsToVarIds[cd]; k {
				t.Errorf("RegistryStore.constraintsToVarIds: contains wrong " +
					"number of elements")
			}
			if _, k := rs.constraints[cd.constraint.GetID()]; k {
				t.Errorf("RegistryStore.constraints: contains wrong number " +
					"of elements")
			}
		}
	}

	if rs.varIdsToConstraints[removedVarId] != nil {
		t.Errorf("RegistryStore.constraints: varid still contained")
	}
}
