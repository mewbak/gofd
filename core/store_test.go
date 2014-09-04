package core

import (
	"fmt"
	"testing"
	"time"
)

// PropDummy propagator dummy for testing
type PropDummy struct {
	x        VarId
	id       PropId
	x_domain Domain
	t        *testing.T
	in       <-chan *ChangeEntry
	out      chan<- *ChangeEvent
}

func (this *PropDummy) GetVarIds() []VarId {
	return []VarId{this.x}
}

func (this *PropDummy) GetDomains() []Domain {
	return []Domain{this.x_domain}
}

func (this *PropDummy) GetInCh() <-chan *ChangeEntry {
	return this.in
}

func (this *PropDummy) GetOutCh() chan<- *ChangeEvent {
	return this.out
}

func (this *PropDummy) GetID() PropId {
	return this.id
}

func (this *PropDummy) SetID(id PropId) {
	this.id = id
}

func (this *PropDummy) Start(store *Store) {
	this.out <- CreateChangeEvent() // "nothing happened" on first run
	for _ = range this.in {         // answer to every ChangeEvent
		println("PropDummy: Do Something")
		this.out <- CreateChangeEvent() // with "nothing happened"
	}
}

func (this *PropDummy) Register(store *Store) {
	var domains []Domain
	this.in, domains, this.out =
		store.RegisterPropagator([]VarId{this.x}, this.id)
	this.x_domain = domains[0]
	this.x_domain.Remove(1)
	if this.x_domain.Equals(store.iDToIntVar[this.x].Domain) {
		msg := "Store.PropDummyRegister: "
		msg += "local domain IS NOT a copy of store domain"
		this.t.Errorf(msg)
	}

}

func (this *PropDummy) Clone() Constraint {
	return &PropDummy{id: this.id}
}

func (this *PropDummy) String() string {
	return fmt.Sprintf("PropDummy, Dom: %v", this.x_domain)
}

func createPropagatorDummy(x VarId, t *testing.T) (Constraint, Constraint) {
	prop1, prop2 := new(PropDummy), new(PropDummy)
	prop1.x, prop2.x = x, x
	prop1.t, prop2.t = t, t
	return prop1, prop2
}

func createVarPropagatorDummy(t *testing.T) (Constraint, Constraint) {
	x := CreateAuxIntVarValues(store, []int{0, 1, 2, 3, 4})
	return createPropagatorDummy(x, t)
}

func Test_RegisterPropagator(t *testing.T) {
	setup()
	defer teardown()
	log("StoreRegisterPropagator")

	prop, _ := createVarPropagatorDummy(t)
	store.AddPropagator(prop)
}

func Test_StoreAddPropagator(t *testing.T) {
	setup()
	defer teardown()
	log("StoreAddPropagator")

	left, right := createVarPropagatorDummy(t)
	store.AddPropagator(left)
	time.Sleep(500) // adding propagator is async; ToDo: sleep?
	if len(store.propagators) != 1 {
		t.Errorf("Store.AddPropagator: len(store.propagators)=%v, want %v",
			len(store.propagators), 1)
	}
	if store.propagators[store.propCounter-1] != left {
		t.Errorf("Store.AddPropagator: last added = %v, want %v",
			store.propagators[store.propCounter-1], left)
	}
	store.AddPropagator(right)
	time.Sleep(500) // adding propagator is async; ToDo: sleep?
	if len(store.propagators) != 2 {
		t.Errorf("Store.AddPropagator: len(store.propagators)=%v, want %v",
			len(store.propagators), 2)
	}
	if store.propagators[store.propCounter-1] != right {
		t.Errorf("Store.AddPropagator: last added = %v, want %v",
			store.propagators[store.propCounter-1], right)
	}
}

func Test_StoreAddPropagators(t *testing.T) {
	setup()
	defer teardown()
	log("StoreAddPropagators")

	store.AddPropagators(createVarPropagatorDummy(t))
	time.Sleep(500) //adding propagator is async
	if len(store.propagators) != 2 {
		t.Errorf("Store.AddPropagator: len(store.propagators)=%v, want %v",
			len(store.propagators), 2)
	}
	store.AddPropagators(createVarPropagatorDummy(t))
	time.Sleep(500) //adding propagator is async
	if len(store.propagators) != 4 {
		t.Errorf("Store.AddPropagator: len(store.propagators)=%v, want %v",
			len(store.propagators), 4)
	}
}

func Test_StoreCreateIntVar(t *testing.T) {
	setup()
	defer teardown()
	log("StoreCreateIntVar")
	X := CreateIntVarValues("X", store, []int{1, 2, 3})
	XStored, exists := store.GetIntVar(X)
	if !exists {
		t.Errorf("StoreCreateIntVar: IntVar in store=%v, want %v",
			!exists, exists)
	}
	if X != XStored.ID {
		t.Errorf("StoreCreateIntVar: IntVar X not the X in Store")
	}
}

func Test_StoreIsConsistent(t *testing.T) {
	setup()
	defer teardown()
	log("StoreIsConsistent")
	store.AddPropagators(createVarPropagatorDummy(t))
	if !store.IsConsistent() { // shall eventually be ready and consistent
		t.Errorf("Store.IsConsistent: gives false on Dummy Propagators ")
	}
}

func Test_StoreClose(t *testing.T) {
	setup()
	defer teardown()
	log("StoreClose")
	store.AddPropagators(createVarPropagatorDummy(t))
	log("StoreClose: ToDo!!")
	// store.close()
	if !store.IsConsistent() {
		t.Errorf("Store.close: After closing still not ready ")
	}
}

func Test_CounterLock(t *testing.T) {
	setup()
	defer teardown()
	log("CounterLock")
	vars := make(map[VarId]bool)
	varCreator := func() {
		x := CreateAuxIntVarValues(store, []int{1, 2, 3})
		if vars[x] {
			t.Errorf("Store.CounterLock: same varid allocated several times")
		}
		vars[x] = true
	}
	for i := 0; i < 1000; i++ {
		go varCreator()
	}
}

func Test_StoreClone(t *testing.T) {
	setup()
	defer teardown()
	log("StoreClone")
	store.AddPropagators(createVarPropagatorDummy(t))
	newStore := store.Clone(nil)
	if &store == &newStore {
		t.Errorf("Store.Clone: &newStore=%v identical to &store=%v",
			&newStore, &store)
	}
	if store == newStore {
		t.Errorf("Store.StoreClone: newStore=%v the same as store=%v",
			store, newStore)
	}
}
