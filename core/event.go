package core

import (
	"fmt"
)

// ControlEvent represents all events that do not remove values from
// a domain but still access or modifiy a store and thus must be
// executed in the central routine of the store to avoid concurrency
// issues.
type ControlEvent interface {
	run(store *Store) // execute whatever needs to be run in the store context
}

// --- Info-Events ---

// GetNameEvent retrieves information about a variable name.
type GetNameEvent struct {
	varId   VarId
	channel chan string
}

func createGetNameEvent(varId VarId) *GetNameEvent {
	infoEvent := new(GetNameEvent)
	infoEvent.varId = varId
	infoEvent.channel = make(chan string)
	return infoEvent
}

func (this *GetNameEvent) run(store *Store) {
	this.channel <- store.iDToName[this.varId]
}

func (this *GetNameEvent) String() string {
	return fmt.Sprintf("GetNameEvent: varid %d", this.varId)
}

// GetNewIdEvent generates a new name, if a variable is named incorrectly.
type GetNewIdEvent struct {
	channel chan VarId
}

func createGetNewIdEvent() *GetNewIdEvent {
	infoEvent := new(GetNewIdEvent)
	infoEvent.channel = make(chan VarId)
	return infoEvent
}

func (this *GetNewIdEvent) run(store *Store) {
	store.iDCounter += 1
	this.channel <- store.iDCounter // unique id, increasing
}

func (this *GetNewIdEvent) String() string {
	return fmt.Sprintf("GetNewIdEvent")
}

// GetDomainEvent to retrieve a copy of a domain from the store
type GetDomainEvent struct {
	varId   VarId
	channel chan Domain
}

func createGetDomainEvent(varId VarId) *GetDomainEvent {
	infoEvent := new(GetDomainEvent)
	infoEvent.varId = varId
	infoEvent.channel = make(chan Domain)
	return infoEvent
}

func (this *GetDomainEvent) run(store *Store) {
	this.channel <- store.iDToIntVar[this.varId].Domain.Copy()
}

func (this *GetDomainEvent) String() string {
	return fmt.Sprintf("GetDomainEvent: varid %d", this.varId)
}

// GetDomainsEvent to retrieve a slice of domain copies from the store
type GetDomainsEvent struct {
	varIds  []VarId
	channel chan []Domain
}

func createGetDomainsEvent(varIds []VarId) *GetDomainsEvent {
	infoEvent := new(GetDomainsEvent)
	infoEvent.varIds = varIds
	infoEvent.channel = make(chan []Domain)
	return infoEvent
}

func (this *GetDomainsEvent) run(store *Store) {
	domains := make([]Domain, len(this.varIds))
	for i, varId := range this.varIds {
		domains[i] = store.iDToIntVar[varId].Domain.Copy()
	}
	this.channel <- domains
}

func (this *GetDomainsEvent) String() string {
	return fmt.Sprintf("GetDomainsEvent: for %d domains", len(this.varIds))
}

// GetMinMaxDomainEvent to retrieve the extreme values of a
// domain from the store
type GetMinMaxDomainEvent struct {
	varId   VarId
	channel chan int
}

func createGetMinMaxDomainEvent(varId VarId) *GetMinMaxDomainEvent {
	infoEvent := new(GetMinMaxDomainEvent)
	infoEvent.varId = varId
	infoEvent.channel = make(chan int)
	return infoEvent
}

func (this *GetMinMaxDomainEvent) run(store *Store) {
	domain := store.iDToIntVar[this.varId].Domain
	this.channel <- domain.GetMin()
	this.channel <- domain.GetMax()
}

func (this *GetMinMaxDomainEvent) String() string {
	return fmt.Sprintf("GetMinMaxDomainEvent: %d", this.varId)
}

// SelectVarIdUnfixedDomainEvent retrieves the id of a variable
// with a domain that is not yet fixed, which means it contains
// at least two elements. It selects either the if of a variable with
// a domain that contains the most or the least number of values.
// -1 if there is no unfixed variable left.
type SelectVarIdUnfixedDomainEvent struct {
	min     bool
	channel chan VarId
}

func createSelectVarIdUnfixedDomainEvent(min bool) *SelectVarIdUnfixedDomainEvent {
	infoEvent := new(SelectVarIdUnfixedDomainEvent)
	infoEvent.min = min
	infoEvent.channel = make(chan VarId)
	return infoEvent
}

func (this *SelectVarIdUnfixedDomainEvent) run(store *Store) {
	varId := VarId(-1)
	if this.min {
		varId = store.getMinUnfixedDomain()
	} else {
		varId = store.getMaxUnfixedDomain()
	}
	this.channel <- varId
}

func (this *SelectVarIdUnfixedDomainEvent) String() string {
	return fmt.Sprintf("SelectVarIdUnfixedDomainEvent: %v", this.min)
}

func (this *Store) getMinUnfixedDomain() VarId {
	id, size := VarId(-1), MaxInt
	for varId, intVar := range this.iDToIntVar {
		dSize := intVar.Domain.Size()
		if 1 < dSize && dSize < size {
			id, size = varId, dSize
		}
	}
	return id // -1 iff none is found
}

func (this *Store) getMaxUnfixedDomain() VarId {
	id, size := VarId(-1), 1
	for varId, intVar := range this.iDToIntVar {
		dSize := intVar.Domain.Size()
		if dSize > size {
			id, size = varId, dSize
		}
	}
	return id // -1 iff none is found
}

// GetVariableIDsEvent to retrieve a slice of all variable ids of a store.
type GetVariableIDsEvent struct {
	channel chan []VarId
}

func createGetVariableIDsEvent() *GetVariableIDsEvent {
	infoEvent := new(GetVariableIDsEvent)
	infoEvent.channel = make(chan []VarId, 1)
	return infoEvent
}

func (this *GetVariableIDsEvent) run(store *Store) {
	this.channel <- store.getVariableIDs()
}

func (this *GetVariableIDsEvent) String() string {
	return fmt.Sprintf("GetVariableIDsEvent")
}

// GetIntVarEvent to retrieve an IntVar given a variable id
type GetIntVarEvent struct {
	varId   VarId
	channel chan *IntVar
}

func createGetIntVarEvent(varId VarId) *GetIntVarEvent {
	infoEvent := new(GetIntVarEvent)
	infoEvent.varId = varId
	infoEvent.channel = make(chan *IntVar, 1)
	return infoEvent
}

func (this *GetIntVarEvent) run(store *Store) {
	this.channel <- store.iDToIntVar[this.varId].Clone()
}

func (this *GetIntVarEvent) String() string {
	return fmt.Sprintf("GetIntVarEvent: %d", this.varId)
}

// GetNumPropagatorsEvent to retrieve the number of propagators
type GetNumPropagatorsEvent struct {
	channel chan int
}

func createGetNumPropagatorsEvent() *GetNumPropagatorsEvent {
	infoEvent := new(GetNumPropagatorsEvent)
	infoEvent.channel = make(chan int, 1)
	return infoEvent
}

func (this *GetNumPropagatorsEvent) run(store *Store) {
	this.channel <- len(store.propagators)
}

func (this *GetNumPropagatorsEvent) String() string {
	return fmt.Sprintf("GetNumPropagatorsEvent")
}

// GetStatEvent to retrieve some statistics
type GetStatEvent struct {
	channel chan *StoreStatistics
}

func createGetStatEvent() *GetStatEvent {
	statEvent := new(GetStatEvent)
	statEvent.channel = make(chan *StoreStatistics, 1)
	return statEvent
}

func (this *GetStatEvent) run(store *Store) {
	this.channel <- store.stat.Clone(store) // a fresh copy, current values
}

func (this *GetStatEvent) String() string {
	return fmt.Sprintf("GetStatEvent")
}

// GetReadyEvent to retrieve a message as soon as propagation has ended
type GetReadyEvent struct{}

// createGetReadyEvent creates an empty event forcing the store
// to generate an event on the ready channel.
func createGetReadyEvent() *GetReadyEvent {
	infoEvent := new(GetReadyEvent)
	return infoEvent
}

func (this *GetReadyEvent) run(store *Store) {
	if store.isReady() { // no more events pending
		store.readyChannel <- !store.isInconsistent()
	} else {
		// events pending, mark that answer on ready channel is expected
		store.communicating = true
	}
}

func (this *GetReadyEvent) String() string {
	return fmt.Sprintf("GetReadyEvent")
}

// CloneEvent clones a complete store
type CloneEvent struct {
	store   *Store
	chEvt   *ChangeEvent
	channel chan *Store
}

func createCloneEvent(store *Store, chEvt *ChangeEvent) *CloneEvent {
	evt := new(CloneEvent)
	evt.chEvt = chEvt
	evt.store = store
	evt.channel = make(chan *Store)
	return evt
}

func (this *CloneEvent) run(store *Store) {
	/* ToDo: Closing store is still a problem. Also
	   working on closed store. Should be fixed soon.
	   But with this solution, panic-ing while cloning
	   because of closed store, brings up problems.
	   Problem is in Labeling. Cloning a closed store
	   is not checked and wouldn't even work with a
	   check, if store is closed. Should be fixed!

	if (store.closed) {
		panic(fmt.Sprintf("cannot clone closed store %p", store))
	} */
	newStore := new(Store)
	newStoreStat := CreateStoreStatistics() // empty statistics
	newStore.stat = newStoreStat
	newStore.readChannel = make(chan *ChangeEvent, ReadChannelBuffer)
	newStore.controlChannel = make(chan ControlEvent)
	newStore.iDToIntVar = make(map[VarId]*IntVar)
	// copy variables on dead store
	for k, v := range this.store.iDToIntVar {
		newStore.iDToIntVar[k] = v.Clone()
		newStoreStat.variables += 1
	}
	// do the changes on the domains if there are any
	if this.chEvt != nil {
		newStoreStat.changeEvents += 1
		for _, changeEntry := range this.chEvt.changes {
			domain := newStore.iDToIntVar[changeEntry.GetID()].Domain
			domain.Removes(changeEntry.GetValues())
			newStoreStat.changeEntries += 1
		}
	}
	newStore.propVarIds = make(map[PropId]VarIdSet, len(store.propVarIds))
	newStore.iDToWriteChannel =
		make(map[VarId]IntChannelSet, len(store.iDToWriteChannel))
	newStore.propagators =
		make(map[PropId]Constraint, len(store.propagators))
	// use old variable IDs and name mapping
	newStore.nameToID = this.store.nameToID
	newStore.iDToName = this.store.iDToName
	newStore.iDCounter = this.store.iDCounter
	newStore.loggingStats = this.store.GetLoggingStat()
	// use new propIds
	newStore.propCounter = this.store.propCounter
	newStore.readyChannel = make(chan bool)
	clonedProps := make([]Constraint, len(this.store.propagators))
	i := 0
	for _, prop := range this.store.propagators {
		clonedProp := prop.Clone()
		clonedProp.SetID(0) // ignore existing id, treat as new
		clonedProps[i] = clonedProp
		i++
	}
	go newStore.propagate()
	newStore.AddPropagators(clonedProps...)
	// println(len(cloned_props), newStore.stat.propagators) // why diff?
	this.channel <- newStore
}

func (this *CloneEvent) String() string {
	return fmt.Sprintf("CloneEvent: cloning %p", this.store)
}

type CloseEvent struct {
	store   *Store
	channel chan bool
}

func createCloseEvent(store *Store) *CloseEvent {
	evt := new(CloseEvent)
	evt.store = store
	evt.channel = make(chan bool)
	return evt
}

func (this *CloseEvent) run(store *Store) {
	store.close()
	this.channel <- true
}

func (this *CloseEvent) String() string {
	return fmt.Sprintf("CloseEvent: closing %p", this.store)
}
