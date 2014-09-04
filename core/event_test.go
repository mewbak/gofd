package core

import (
	"testing"
)

func Test_ChangeEvent(t *testing.T) {
	setup()
	defer teardown()
	log("ChangeEventTest")
	noVars := 10
	// noVars variables with domain of size noVars where for the ith variable
	// starting with the second all domain values >= i are to be removed
	varIds := make([]VarId, noVars)
	vars := make([]*IntVar, noVars)
	for i := 0; i < noVars; i++ {
		varIds[i] = CreateAuxIntVarFromTo(store, 1, noVars)
		vars[i], _ = store.GetIntVar(varIds[i])
	}
	event := CreateChangeEvent()
	entry0 := CreateChangeEntry(varIds[0])
	if !entry0.IsEmpty() {
		t.Errorf("ChangeEntry %s should be initially empty", entry0.String())
	}
	dom := CreateIvDomainFromTo(vars[0].Domain.GetMin(),
		vars[0].Domain.GetMin())
	entry0.SetValues(dom)
	if entry0.IsEmpty() {
		t.Errorf("ChangeEntry %s should contain a value", entry0.String())
	}
	event.AddChangeEntry(entry0)
	for i := 1; i < noVars; i++ {
		entry := CreateChangeEntry(varIds[i])
		dom := CreateIvDomainFromTo(i, noVars-1)
		entry.SetValues(dom)
		event.AddChangeEntry(entry)
	}

	//Equals-Test
	entry1 := CreateChangeEntry(varIds[0])
	entry1.SetValues(dom)
	if !entry0.Equals(entry1) {
		t.Errorf("ChangeEntry %s should be equal to ChangeEntry %s",
			entry0, entry1)
	}
	event2 := CreateChangeEvent()
	event2.AddChangeEntry(entry1)
	for i := 1; i < noVars; i++ {
		entry := CreateChangeEntry(varIds[i])
		dom := CreateIvDomainFromTo(i, noVars-1)
		entry.SetValues(dom)
		event2.AddChangeEntry(entry)
	}
	if !event.Equals(event2) {
		t.Errorf("ChangeEvent %s should be equal to ChangeEvent %s",
			event, event2)
	}

	//Clone-Test
	entry1c := entry1.Clone()
	if !entry1.Equals(entry1c) {
		t.Errorf("ChangeEntry %s should be equal to ChangeEntry %s",
			entry1, entry1c)
	}
	eventC := event.Clone()
	if !event.Equals(eventC) {
		t.Errorf("ChangeEvent %s should be equal to ChangeEvent %s",
			event, eventC)
	}
}

func Test_IntVarEvent(t *testing.T) {
	setup()
	defer teardown()
	log("IntVarEventTest")
	Xid := CreateIntVarValues("X", store, []int{1, 2, 3})
	X, _ := store.GetIntVar(Xid)
	rive := createRegisterIntVarEvent("X", X)
	expected := "RegisterIntVarEvent: name X, [1..3]"
	if expected != rive.String() {
		t.Errorf("RegisterIntVarEvent.String = %s, want %s",
			rive.String(), expected)
	}
}
