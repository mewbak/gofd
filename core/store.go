package core

import (
	"fmt"
)

// IntChannelSet type that maps each propagator ID to a
// channel of ChangeEntry's.
type IntChannelSet map[PropId]chan *ChangeEntry

const ReadChannelBuffer int = 17 // some buffering on store input

// Store represents a constraint store
type Store struct {
	iDToIntVar       map[VarId]*IntVar
	iDToWriteChannel map[VarId]IntChannelSet // multiplexer
	propagators      map[PropId]Constraint
	propCounter      PropId
	iDCounter        VarId
	registryStore    *RegistryStore
	// propVarIds, varids of unfixed variables per propagator
	propVarIds     map[PropId]VarIdSet
	readChannel    chan *ChangeEvent // incoming domain reductions
	controlChannel chan ControlEvent // for commands such as register
	// eventCounter increment per outgoing even, expect one incoming
	eventCounter int
	readyChannel chan bool // no further prop., eventCounter == 0
	// true if events in process but result on readyChannel expected
	communicating bool
	closed        bool // true iff store has been closed
	stat          *StoreStatistics
	loggingStats  bool
}

// GetLoggingStat gets some statistics from logging.
func (this *Store) GetLoggingStat() bool {
	return this.loggingStats
}

// CreateStoreWithoutLogging creates a store, but does not log.
func CreateStoreWithoutLogging() *Store {
	if logger.DoDebug() {
		logger.Dln("STORE_CreateStore")
	}
	store := new(Store)
	store.loggingStats = false
	store.iDToIntVar = make(map[VarId]*IntVar)
	store.iDToWriteChannel = make(map[VarId]IntChannelSet)
	store.propagators = make(map[PropId]Constraint)
	store.readyChannel = make(chan bool)
	store.readChannel = make(chan *ChangeEvent, ReadChannelBuffer)
	store.controlChannel = make(chan ControlEvent)
	store.propVarIds = make(map[PropId]VarIdSet)
	store.iDCounter = -1  // will begin with zero for the first var
	store.propCounter = 1 // >0 to distinguish not yet initialized ones
	store.eventCounter = 0
	store.communicating = false
	store.closed = false
	store.stat = CreateStoreStatistics()
	store.registryStore = CreateRegistryStore()
	go store.propagate() // always start the store main loop
	return store
}

// CreateStore creates a new empty store
func CreateStore() *Store {
	if logger.DoDebug() {
		logger.Dln("STORE_CreateStore")
	}
	store := new(Store)
	store.loggingStats = true
	store.iDToIntVar = make(map[VarId]*IntVar)
	store.iDToWriteChannel = make(map[VarId]IntChannelSet)
	store.propagators = make(map[PropId]Constraint)
	store.readyChannel = make(chan bool)
	store.readChannel = make(chan *ChangeEvent, ReadChannelBuffer)
	store.controlChannel = make(chan ControlEvent)
	store.propVarIds = make(map[PropId]VarIdSet)
	store.iDCounter = -1  // will begin with zero for the first var
	store.propCounter = 1 // >0 to distinguish not yet initialized ones
	store.eventCounter = 0
	store.communicating = false
	store.closed = false
	store.stat = CreateStoreStatistics()
	store.registryStore = CreateRegistryStore()
	go store.propagate() // always start the store main loop
	return store
}

// close closes the store and prepares it for gargabe collection.
// No new propagators shall be added, status calls shall still be possible.
// Thus, the control channel remains open, which introduces a leak.
func (this *Store) close() {
	if this.closed { // already closed
		return
	}
	for propId, varIdSet := range this.propVarIds {
		for varId := range varIdSet {
			// propagator no longer listens to variable
			delete(this.propVarIds[propId], varId)
			// this was the last variable the propagator listend to
			// thus, the last accessible channel reference
			if len(this.propVarIds[propId]) == 0 {
				// no longer send to that propagator, propagator closes down
				close(this.iDToWriteChannel[varId][propId])
				// remove any dangling reference to propagator,
				// ready to be garbage collected
				delete(this.propagators, propId)
			}
			// removes channel reference from multiplexer for this variable
			delete(this.iDToWriteChannel[varId], propId)
		}
	}
	// ToDo: still fails, why?
	// by preventing pending change events to have effect there is no
	// longer a nil or deadlock, but for e.g. nqueens I get more solutions...
	// this.iDToIntVar = make(map[VarId]*IntVar)
	this.closed = true
}

// IsConsistent checks whether the store is consistent by asking for
// and retrieving a GetReadyEvent
func (this *Store) IsConsistent() bool {
	this.controlChannel <- createGetReadyEvent()
	return <-this.getReadyChannel()
}

// AddPropagator registers a propagator to the constraint store.
func (this *Store) AddPropagator(prop Constraint) {
	if (*this).closed {
		panic("cannot add propagator to closed store")
	}
	evt := createRegisterPropagatorEvent(prop)
	this.controlChannel <- evt
	propId := <-evt.channel
	if logger.DoDebug() {
		logger.Df("propagator registered with id %d", propId)
	}
}

// AddPropagators registers Propagators to the constraint store.
func (this *Store) AddPropagators(props ...Constraint) {
	evt := createRegisterPropagatorsEvent(props)
	this.controlChannel <- evt
	propIds := <-evt.channel
	if logger.DoDebug() {
		logger.Df("%d propagators registered", len(propIds))
	}
}

// RegisterPropagator to be called from an Constraint instance
// in its Register() function  before propagation, to obtain the
// local copy of the Domains and the event channels: its input for
// change entries removing values from domains and its output
// where the propagator puts its computed change events.
// Thus, the first channel returned is to be used read-only, while
// the last channel returned is to be used write-only
func (this *Store) RegisterPropagator(varIds []VarId,
	propId PropId) (<-chan *ChangeEntry, []Domain, chan<- *ChangeEvent) {
	this.stat.propagators++
	lenvarIds := len(varIds)
	domains := make([]Domain, lenvarIds)
	varIdSet := make(VarIdSet, lenvarIds)
	channelSize := 0
	loggerDebug := logger.DoDebug()
	for i, varId := range varIds {
		if loggerDebug {
			logger.Df("STORE_registerPropagator: var %s for propId %d",
				this.registryStore.GetVarName(varId), propId)
		}
		domains[i] = this.iDToIntVar[varId].Domain.Copy()
		channelSize += domains[i].Size()
		varIdSet[varId] = true
	}
	// writeChannel, for the store to write and the propagator to read
	writeChannel := make(chan *ChangeEntry, channelSize)
	this.stat.sizeChannels += channelSize
	for _, varId := range varIds {
		if _, exists := this.iDToWriteChannel[varId]; !exists {
			this.iDToWriteChannel[varId] = make(IntChannelSet)
		}
		this.iDToWriteChannel[varId][propId] = writeChannel
	}
	this.propVarIds[propId] = varIdSet
	if loggerDebug {
		logger.Df("STORE_registerPropagator writeChannel: %v", writeChannel)
	}
	return writeChannel, domains, this.readChannel
}

// RegisterPropagatorMap similar to RegisterPropagator but returns a map
// of VarId to Domain in case the propagator prefers that
func (this *Store) RegisterPropagatorMap(varIds []VarId,
	propId PropId) (<-chan *ChangeEntry,
	map[VarId]Domain, chan<- *ChangeEvent) {
	inCh, domains, outCh := this.RegisterPropagator(varIds, propId)
	domainsMap := make(map[VarId]Domain, len(varIds))
	for i, varId := range varIds {
		domainsMap[varId] = domains[i]
	}
	return inCh, domainsMap, outCh
}

// registerIntVarAtStore adds a new IntVar at its Creation to the Store and
// associates it with the given name. In addition generates a unique ID in
// this store (increases internal counter) for the added variable
func (this *Store) registerIntVarAtStore(name string, intVar *IntVar) VarId {
	evt := createRegisterIntVarEvent(name, intVar)
	this.controlChannel <- evt
	return <-evt.channel
}

// registerAuxIntVarAtStore adds a new auxiliary IntVar at its Creation to the
// Store and associates it with the given name. In addition generates a unique
// ID in this store (increases internal counter) for the added variable.
func (this *Store) registerAuxIntVarAtStore(intVar *IntVar) VarId {
	evt := createRegisterAuxIntVarEvent(intVar)
	this.controlChannel <- evt
	return <-evt.channel
}

// makeAuxVariableName generates a name for a variable (used for generating
// names for auxiliary variables). These variable should not be communicated
// to users, only for internal use.
func (this *Store) generateAuxVariableName(varId VarId) string {
	return fmt.Sprintf("_%d", varId)
}

// generateNewVariableName generates a name for a variable (used, if
// the specified variable name is not valid).
func (this *Store) generateNewVariableName() string {
	evt := createGetNewIdEvent()
	this.controlChannel <- evt
	return fmt.Sprintf("V%d", <-evt.channel)
}

// propagate is the main execution thread in the store, which initially starts
// all added propagators, handles and dispatches ChangeEvents and handles
// new registrations of further propagators.
func (this *Store) propagate() {
	loggerDebug := logger.DoDebug()
	if loggerDebug {
		logger.Dln("STORE_Propagate'")
	}
	if this.loggingStats {
		this.stat.InitStatTime()
	}
	i := 0
	for { // ToDo: store will *never* end (because of controlChannel); leak
		i += 1
		select {
		case event := <-this.controlChannel: // commands
			if this.loggingStats {
				this.stat.LogIdleTime()
			}
			this.stat.controlEvents++
			if loggerDebug {
				logger.Df("STORE_controlEvent caught: %s", event)
			}
			event.run(this) // execute task in central thread
			if this.loggingStats {
				this.stat.LogWorkingTime()
			}
		case event := <-this.readChannel: // domain reductions
			if this.loggingStats {
				this.stat.LogIdleTime()
			}
			this.stat.changeEvents++
			if loggerDebug {
				logger.Df("STORE_changeEvent caught: %v", event.changes)
			}
			if !this.closed {
				if len(event.changes) > 0 {
					this.processChanges(event.changes, loggerDebug)
				} else {
					this.stat.emptyChangeEvents++
				}
			}
			this.eventCounter -= 1
			if loggerDebug {
				logger.Df("STORE_eventcounter reached %v", this.eventCounter)
			}
			if this.eventCounter == 0 && this.communicating {
				if loggerDebug {
					logger.Iln("STORE_eventcounter reached 0, send ready event")
				}
				this.readyChannel <- !this.isInconsistent()
				this.communicating = false
			}
			if this.loggingStats {
				this.stat.LogWorkingTime()
			}
		}
	}
}

// processChanges processes domain reductions from any propagator
func (this *Store) processChanges(changes []*ChangeEntry, loggerDebug bool) {
	for _, changeEntry := range changes {
		this.stat.changeEntries++
		varId := changeEntry.varId
		domain := this.iDToIntVar[varId].Domain
		// removing values from changeEntry
		toremovevalues := changeEntry.GetValues()
		// too high, wraps around, useless for unbound domains
		this.stat.domainReductions += toremovevalues.Size()
		beforeSize := domain.Size()
		domain.RemovesWithOther(toremovevalues)
		numberOfRemoved := beforeSize - domain.Size()
		this.stat.domainValsRemoved += numberOfRemoved
		// ToDo: with this.stat.domainReductions += 1 this number should
		// coincide with stat.changeEntries, it doesn't when enumerating.
		// Thus, there must be a concurrency bug somewhere
		// (prime suspect is labeling).
		if domain.IsEmpty() { // store is failed
			this.close()
			return // stop processing
		}
		if !toremovevalues.IsEmpty() { // reductions to communicate
			if numberOfRemoved == 0 {
				logger.Ef("store.processChanges: no change %v\n", changeEntry)
			}
			this.multiplex(changeEntry) // distribute to propagators
		}
		if domain.IsGround() {
			if loggerDebug {
				logger.Df("STORE_fixed domain, var: %s",
					this.registryStore.GetVarName(varId))
			}
			for propId := range this.iDToWriteChannel[varId] {
				// propagator no longer listens to fixed variable
				delete(this.propVarIds[propId], varId)
				// if it was the last variable the propagator listend to
				if len(this.propVarIds[propId]) == 0 {
					// no longer send to that propagator, close it down
					close(this.iDToWriteChannel[varId][propId])
					this.stat.propagators -= 1
					// remove dangling reference to propagator, allow gc
					delete(this.propagators, propId)
				}
				// remove channel (only) from that multiplexer
				delete(this.iDToWriteChannel[varId], propId)
			}
		}
	}
}

// isReady returns true iff there are no outstanding propagations to be done
// Checks the current state, nonblocking. Use IsConsistent() for a blocking
// information with status as soon as propagation has come to an end.
func (this *Store) isReady() bool {
	return this.eventCounter == 0
}

// isInconsistent returns true iff one IntVar in the store has an empty domain
func (this *Store) isInconsistent() bool {
	// Note that precomputing the result (during reductions) is not worth
	// the effort (benched it) and complicates cloning
	for _, intVar := range this.iDToIntVar {
		if intVar.Domain.IsEmpty() {
			return true
		}
	}
	return false
}

// multiplex multiplexes and dispatches copied value removal messages to all
// registered propagators containing copies of those domains
func (this *Store) multiplex(changeEntry *ChangeEntry) {
	if logger.DoDebug() {
		logger.Df("STORE_multiplex change on %s, value %s to %v channels",
			this.registryStore.GetVarName(changeEntry.varId), changeEntry,
			len(this.iDToWriteChannel[changeEntry.varId]))
	}
	this.eventCounter += len(this.iDToWriteChannel[changeEntry.varId])
	for _, readChannel := range this.iDToWriteChannel[changeEntry.varId] {
		readChannel <- changeEntry
	}
}

// propagatorExistsAlready returns true iff the store has added
// this propagator already; avoids duplicate propagators in store
func (this *Store) propagatorExistsAlready(prop Constraint) bool {
	_, exists := this.propagators[prop.GetID()]
	return exists
}

// Clone creates a running copy of this store and returns pointer to this copy
func (this *Store) Clone(chEvt *ChangeEvent) *Store {
	evt := createCloneEvent(this, chEvt)
	this.controlChannel <- evt
	store := <-evt.channel
	return store
}

// Close closes the store, its channel and enables the release of memory
func (this *Store) Close() bool {
	evt := createCloseEvent(this)
	this.controlChannel <- evt
	return <-evt.channel
}

// getReadyChannel returns a control channel where the store signals that
// constraint propagation has reached a fixpoint
func (this *Store) getReadyChannel() chan bool {
	return this.readyChannel
}

// GetVariableIDs returns all variable ids of a store
func (this *Store) GetVariableIDs() []VarId {
	evt := createGetVariableIDsEvent()
	this.controlChannel <- evt
	return <-evt.channel
}

func (this *Store) getVariableIDs() []VarId {
	value_slice := make([]VarId, len(this.iDToIntVar))
	i := 0
	for k := range this.iDToIntVar {
		value_slice[i] = k
		i += 1
	}
	return value_slice
}

// GetName returns the name of the IntVar with the given id
func (this *Store) GetName(id VarId) string {
	// It is safe to do directly as names never change
	return this.registryStore.GetVarName(id)
	// ToDo: really no race condition during construction?
	// Alternative:
	// evt := createGetNameEvent(id)
	// this.controlChannel <- evt
	// return <-evt.channel
}

// GetIntVar returns a pointer to an IntVar represented by given varId
func (this *Store) GetIntVar(id VarId) (*IntVar, bool) {
	evt := createGetIntVarEvent(id)
	this.controlChannel <- evt
	value := <-evt.channel
	return value, value != nil
}

// GetDomain returns a pointer to the domain of the IntVar with the given id
func (this *Store) GetDomain(id VarId) Domain {
	evt := createGetDomainEvent(id)
	this.controlChannel <- evt
	return <-evt.channel
}

// GetDomains returns a slice of pointers to the domain of the given varIds
func (this *Store) GetDomains(varIds []VarId) []Domain {
	evt := createGetDomainsEvent(varIds)
	this.controlChannel <- evt
	return <-evt.channel
}

// GetMinMaxDomain returns the minimum and the maximum of a domain
func (this *Store) GetMinMaxDomain(varId VarId) (int, int) {
	evt := createGetMinMaxDomainEvent(varId)
	this.controlChannel <- evt
	min := <-evt.channel
	max := <-evt.channel
	return min, max
}

// GetVarIdSmallestUnfixedDomain gets the variable id of the variable
// with the smallest (min=true) or largest (min=false) unfixed domain
func (this *Store) SelectVarIdUnfixedDomain(min bool) VarId {
	evt := createSelectVarIdUnfixedDomainEvent(min)
	this.controlChannel <- evt
	return <-evt.channel
}

// GetNumPropagators returns the number of propagators
func (this *Store) GetNumPropagators() int {
	evt := createGetNumPropagatorsEvent()
	this.controlChannel <- evt
	return <-evt.channel
}

// GetStat returns statistics of processed events per store
func (this *Store) GetStat() *StoreStatistics {
	evt := createGetStatEvent()
	this.controlChannel <- evt
	return <-evt.channel
}

// String returns the current state of the store as a string
func (this *Store) String() string {
	s := ""
	for id, name := range this.registryStore.GetVarIdToNameMap() {
		s += fmt.Sprintf("%s   %s\n", name, this.iDToIntVar[id].Domain)
	}
	msg := "---Store-Status---\n"
	msg += "closed: %v\n"
	msg += "communicating (request from main): %v\n"
	msg += "%s"
	return fmt.Sprintf(msg, this.closed, this.communicating, s)
}

// StringWithSpecVarIds returns the current state of the store as a string.
// But in contrast to normal String function, the returned string only includes
// the variable names of the given varids (for example to avoid
// auxiliary variables in the returned string)
func (this *Store) StringWithSpecVarIds(varids []VarId) string {
	s := ""
	for _, id := range varids {
		if name, k := this.registryStore.HasVarName(id); k {
			s += fmt.Sprintf("%s   %s\n", name, this.iDToIntVar[id].Domain)
		}
	}
	msg := "---Store-Status---\n"
	msg += "closed: %v\n"
	msg += "communicating (request from main): %v\n"
	msg += "%s"
	return fmt.Sprintf(msg, this.closed, this.communicating, s)
}
