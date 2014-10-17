package core

import (
	"fmt"
)

// The "Register"-Events allow to register new proper or auxilliary variables
// and one or many propagators from outside the store.

// RegisterIntVarEvent register a new IntVar to be managed by the store.
type RegisterIntVarEvent struct {
	name    string     // The externally visible name
	intvar  *IntVar    // The finite domain variable with possible values
	channel chan VarId // communicate back the new unique variable id
}

// createRegisterIntVarEvent creates a new instance of RegisterIntVarEvent.
func createRegisterIntVarEvent(name string,
	intvar *IntVar) *RegisterIntVarEvent {
	regEvent := new(RegisterIntVarEvent)
	regEvent.name = name
	regEvent.intvar = intvar
	regEvent.channel = make(chan VarId)
	return regEvent
}

// run registers the IntVar, to be run in store.
func (this *RegisterIntVarEvent) run(store *Store) {
	varId := registerOneIntVar(store, this.name, this.intvar)
	this.channel <- varId
}

// registerOneIntVar registers one IntVar at store, to be run in store,
// provides name if name is empty.
func registerOneIntVar(store *Store, name string, intvar *IntVar) VarId {
	store.stat.variables++
	store.iDCounter += 1
	varId := store.iDCounter
	intvar.ID = varId
	store.iDToIntVar[varId] = intvar
	if name == "" {
		name = store.generateAuxVariableName(varId) // provide a name
	}
	
	store.registryStore.SetVarName(intvar.ID, name)
	if logger.DoDebug() {
		logger.Df("STORE_register[Aux]IntVar: %d, %s", varId, name)
	}
	return varId
}

// String returns a readable representation of the RegisterIntVarEvent.
func (this *RegisterIntVarEvent) String() string {
	msg := "RegisterIntVarEvent: name %s, %s"
	return fmt.Sprintf(msg, this.name, this.intvar.Domain.String())
}

// RegisterAuxIntVarEvent a new Helper IntVar with automatically
// assigned external name that is managed by the store.
type RegisterAuxIntVarEvent struct {
	intvar  *IntVar    // The finite domain variable with possible values
	channel chan VarId // communicate back the new unique variable id
}

// createRegisterAuxIntVarEvent creates a new instance of
// RegisterAuxIntVarEvent to register an auxiliary variable at a store.
func createRegisterAuxIntVarEvent(intvar *IntVar) *RegisterAuxIntVarEvent {
	regEvent := new(RegisterAuxIntVarEvent)
	regEvent.intvar = intvar
	regEvent.channel = make(chan VarId)
	return regEvent
}

// run registers the helper IntVar in the central store goroutine.
func (this *RegisterAuxIntVarEvent) run(store *Store) {
	varId := registerOneIntVar(store, "", this.intvar)
	this.channel <- varId
}

// String returns a readable representation of the RegisterAuxIntVarEvent.
func (this *RegisterAuxIntVarEvent) String() string {
	msg := "RegisterAuxIntVarEvent: auxintvar, %s"
	return fmt.Sprintf(msg, this.intvar.Domain)
}

// RegisterPropagatorEvent a new propagator that is managed by the store.
type RegisterPropagatorEvent struct {
	prop    Constraint  // any provided or custom developed propagator
	channel chan PropId // communicate back the new unique propagator id
}

// createRegisterPropagatorEvent creates a new instance of RegisterEvent
// to register a given propagator prop at a store.
func createRegisterPropagatorEvent(prop Constraint) *RegisterPropagatorEvent {
	regEvent := new(RegisterPropagatorEvent)
	regEvent.prop = prop
	regEvent.channel = make(chan PropId)
	return regEvent
}

// run registers the Propagator, to be run in the store.
func (this *RegisterPropagatorEvent) run(store *Store) {
	if !store.isInconsistent() {
		propId := registerOnePropagator(store, this.prop)
		this.channel <- propId // communicate id back
		return
	}
	this.channel <- -1
}

// registerOnePropagator registers one propagator in the store,
// to be run in the store.
func registerOnePropagator(store *Store, prop Constraint) PropId {
	if store.propagatorExistsAlready(prop) {
		msg := "propagator %d exists already\n"
		panic(fmt.Sprintf(msg, prop.GetID()))
	}
	propId := store.propCounter
	store.propCounter += 1
	store.propagators[propId] = prop
	prop.SetID(propId)
	prop.Register(store)    // weave into internal data structure
	store.eventCounter += 1 // expect event from initial check
	go prop.Start(store)    // make the propagator propagate
	if logger.DoDebug() {
		logger.If("Propagator registered (ec=%d)", store.eventCounter)
	}
	return propId
}

// String returns a readable representation of the RegisterPropagatorEvent.
func (this *RegisterPropagatorEvent) String() string {
	msg := "RegisterPropagatorEvent: "
	return fmt.Sprintf(msg)
}

// RegisterPropagatorsEvent new propagators that are managed by the store.
type RegisterPropagatorsEvent struct {
	props   []Constraint
	channel chan []PropId
}

// createRegisterPropagatorsEvent creates a new instance of
// RegisterPropagatorsEvent to register given propagators props at a store.
func createRegisterPropagatorsEvent(props []Constraint) *RegisterPropagatorsEvent {
	regEvent := new(RegisterPropagatorsEvent)
	regEvent.props = props
	regEvent.channel = make(chan []PropId)
	return regEvent
}

// run registers the propagators in the central store goroutine.
func (this *RegisterPropagatorsEvent) run(store *Store) {
	if !store.isInconsistent() {
		propIds := make([]PropId, len(this.props))
		for i, prop := range this.props {
			propIds[i] = registerOnePropagator(store, prop)
		}
		this.channel <- propIds // communicate ids back
		return
	}
	this.channel <- []PropId{}
}

// String returns a readable representation of the RegisterPropagatorsEvent
func (this *RegisterPropagatorsEvent) String() string {
	msg := "RegisterPropagatorsEvent for %d propagators"
	return fmt.Sprintf(msg, len(this.props))
}
