package core

import (
	"fmt"
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
	if !event.IsEmpty() {
		t.Errorf("ChangeEvent %s should initially be empty", event.String())
	}
	entry0 := CreateChangeEntry(varIds[0])
	if !entry0.IsEmpty() {
		t.Errorf("ChangeEntry %s should initially be empty", entry0.String())
	}
	if entry0.GetID() != varIds[0] {
		t.Errorf("ChangeEntry Varid is %s, want %s",
			entry0.GetID(), varIds[0])
	}
	v0min := vars[0].Domain.GetMin()
	entry0 = CreateChangeEntryWithIntValue(varIds[0], v0min)
	if entry0.IsEmpty() {
		t.Errorf("ChangeEntry %s should contain a value", entry0.String())
	}
	event.AddChangeEntry(entry0)
	if event.IsEmpty() {
		t.Errorf("ChangeEvent %s should no longer be empty", event.String())
	}
	for i := 1; i < noVars; i++ {
		entry := CreateChangeEntry(varIds[i])
		dom := CreateIvDomainFromTo(i, noVars-1)
		entry.SetValues(dom)
		if !entry.GetValues().Equals(dom) {
			t.Errorf("ChangeEntry domain is %s, want %s",
				entry.GetValues().String(), dom.String())
		}
		event.AddChangeEntry(entry)
	}

	// Equals-Test
	entry1 := CreateChangeEntryWithIntValues(varIds[0], []int{v0min})
	if !entry0.Equals(entry1) {
		t.Errorf("ChangeEntry %s should be equal to ChangeEntry %s",
			entry0, entry1)
	}
	event2 := CreateChangeEvent()
	if event.Equals(event2) {
		t.Errorf("ChangeEvent %s should not be equal to ChangeEvent %s",
			event, event2)
	}
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
	// change one entry; Note, do not do it on 0 as that one is shared
	event2.changes[1].varId += 1
	if event.Equals(event2) {
		t.Errorf("ChangeEvent %s should no be equal to ChangeEvent %s",
			event, event2)
	}
	event2.changes[1].varId -= 1
	event2.changes[1].SetValuesByIntArr([]int{17, 1717, 171717})
	if event.Equals(event2) {
		t.Errorf("ChangeEvent %s should no be equal to ChangeEvent %s",
			event, event2)
	}

	// Clone-Test
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

	// String-Test
	showEvent := CreateChangeEvent()
	showEventString := showEvent.String()
	want := "ChangeEvent[]"
	if want != showEventString {
		t.Errorf("ChangeEvent:String got %s, want %s",
			showEventString, want)
	}
	dom := CreateIvDomainFromTo(17, 34)
	simpleEntry := CreateChangeEntryWithValues(6, dom)
	simpleEntry.Add(35)
	simpleEntry.AddValues(CreateIvDomainFromTo(36, 42))
	showEvent.AddChangeEntry(simpleEntry)
	showEventString = showEvent.String()
	want = "ChangeEvent[ChangeEntry{6, [17..42]}]"
	if want != showEventString {
		t.Errorf("ChangeEvent:String got %s, want %s",
			showEventString, want)
	}
}

func Test_GetNameEvent(t *testing.T) {
	setup()
	defer teardown()
	log("GetNameEventTest")
	Xid := CreateIntVarValues("X", store, []int{1, 2, 3})
	Yid := CreateIntVarValues("Y", store, []int{4, 5, 6})
	getNameEvent := createGetNameEvent(Xid)
	expected := fmt.Sprintf("GetNameEvent: varid %d", Xid)
	if expected != getNameEvent.String() {
		t.Errorf("GetNameEvent:String got %s, want %s",
			getNameEvent.String(), expected)
	}
	// execute on store
	store.controlChannel <- getNameEvent
	name := <-getNameEvent.channel
	if name != "X" {
		t.Errorf("GetNameEvent:Name got %s, want %s", name, "X")
	}
	// through public API, which is implemented by direct access
	name = store.GetName(Yid)
	if name != "Y" {
		t.Errorf("GetNameEvent:Name got %s, want %s", name, "Y")
	}
}

func Test_GetNewIdEvent(t *testing.T) {
	setup()
	defer teardown()
	log("GetNewIdEvent")
	getNewIdEvent := createGetNewIdEvent()
	expected := "GetNewIdEvent"
	if getNewIdEvent.String() != expected {
		t.Errorf("GetNewIdEvent:Name got %s, want %s",
			getNewIdEvent.String(), expected)
	}
	// execute on store
	store.controlChannel <- getNewIdEvent
	newId := <-getNewIdEvent.channel
	ids := make(map[VarId]bool)
	ids[newId] = true
	names := make(map[string]bool)
	names["X"] = true
	for i := 0; i < 1000; i += 1 {
		getNewIdEvent = createGetNewIdEvent()
		store.controlChannel <- getNewIdEvent
		newId = <-getNewIdEvent.channel
		ids[newId] = true
		names[store.generateNewVariableName()] = true
	}
	if len(ids) != 1001 || len(names) != 1001 || newId+1 != 2000 {
		t.Errorf("GetNewIdEvent: len(ids) %d, len(names) %d, lastNewId %d",
			len(ids), len(names), newId)
	}
}

func Test_GetDomainEvent(t *testing.T) {
	setup()
	defer teardown()
	log("GetDomainEventTest")
	Xid := CreateIntVarValues("X", store, []int{1, 2, 3})
	Yid := CreateIntVarValues("X", store, []int{4, 5, 6})
	getDomainEvent := createGetDomainEvent(Xid)
	expected := fmt.Sprintf("GetDomainEvent: varid %d", Xid)
	if expected != getDomainEvent.String() {
		t.Errorf("GetDomainEvent:String got %s, want %s",
			getDomainEvent.String(), expected)
	}
	// execute on store
	store.controlChannel <- getDomainEvent
	domainX := <-getDomainEvent.channel
	expectedDomainX := CreateIvDomainFromTo(1, 3)
	if !domainX.Equals(expectedDomainX) {
		t.Errorf("GetDomainEvent: got %s, want %s",
			domainX.String(), expectedDomainX.String())
	}
	// directly
	domainX = store.GetDomain(Xid)
	if !domainX.Equals(expectedDomainX) {
		t.Errorf("GetDomain: got %s, want %s",
			domainX.String(), expectedDomainX.String())
	}
	domains := store.GetDomains([]VarId{Xid, Yid})
	expectedDomainY := CreateIvDomainFromTo(4, 6)
	if !domains[1].Equals(expectedDomainY) {
		t.Errorf("GetDomainsY: got %s, want %s",
			domains[1].String(), expectedDomainY.String())
	}
	if !domains[0].Equals(expectedDomainX) {
		t.Errorf("GetDomainsX: got %s, want %s",
			domains[0].String(), expectedDomainX.String())
	}
	domainsEvent := createGetDomainsEvent([]VarId{Xid, Yid})
	expected = "GetDomainsEvent: for 2 domains"
	if expected != domainsEvent.String() {
		t.Errorf("GetDomainsEvent.String: got %s, want %s",
			domainsEvent.String(), expected)
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
