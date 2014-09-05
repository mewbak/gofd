package core

import (
	"fmt"
	"time"
)

// StoreStatistics contains all statistical information about the
// contents and execution of a store.
type StoreStatistics struct {
	variables          int   // cumulated number of variables
	act_variables      int   // computed on retrieve event
	propagators        int   // cumulated number of propagators
	act_propagators    int   // computed on retrieve event
	sizeChannels       int   // cumulated size of channels
	changeEvents       int   // cumulated number of ChangeEvents
	emptyChangeEvents  int   // cumulated number of empty ChangeEvents
	changeEntries      int   // cumulated number of ChangeEntries
	domainReductions   int   // cumulated number of requested domain reductions
	domainValsRemoved  int   // cumulated number of performed domain reductions
	controlEvents      int   // cumulated number of ControlEvents
	working_time       int64 // working time in msecs
	idle_time          int64 // idle time in msecs
	start_working_time int64 // current start working time in msecs for calcs
	end_working_time   int64 // current end working time in msecs for calcs
}

// CreateStoreStatistics creates a new empty StoreStatistics instance
func CreateStoreStatistics() *StoreStatistics {
	stat := new(StoreStatistics)
	stat.resetStat()
	return stat
}

// resetStat clears a statistical values (sets to zero)
func (this *StoreStatistics) resetStat() {
	this.variables = 0
	this.propagators = 0
	// act might also be ignored
	this.act_variables = 0
	this.act_propagators = 0
	this.sizeChannels = 0
	this.changeEvents = 0
	this.emptyChangeEvents = 0
	this.changeEntries = 0
	this.domainReductions = 0
	this.domainValsRemoved = 0
	this.controlEvents = 0
	this.working_time = 0
	this.idle_time = 0
}

// Clone creates a new copy of the StoreStatistics,
// act* variables that count current contents are freshly updated.
func (this *StoreStatistics) Clone(store *Store) *StoreStatistics {
	stat := new(StoreStatistics)
	other := store.stat
	stat.variables = other.variables
	stat.propagators = other.propagators
	stat.act_variables = len(store.iDToIntVar)    // fresh on Clone
	stat.act_propagators = len(store.propagators) // fresh on Clone
	stat.sizeChannels = other.sizeChannels
	stat.changeEvents = other.changeEvents
	stat.emptyChangeEvents = other.emptyChangeEvents
	stat.changeEntries = other.changeEntries
	stat.domainReductions = other.domainReductions
	stat.domainValsRemoved = other.domainValsRemoved
	stat.controlEvents = other.controlEvents
	stat.working_time = other.working_time
	stat.idle_time = other.idle_time
	return stat
}

// UpdateStoreStatistics cumulates statics obtained from other
// single independant stores in one entry.
func (this *StoreStatistics) UpdateStoreStatistics(other *StoreStatistics) {
	this.variables += other.variables
	this.propagators += other.propagators
	// ignore act
	this.sizeChannels += other.sizeChannels
	this.changeEvents += other.changeEvents
	this.emptyChangeEvents += other.emptyChangeEvents
	this.changeEntries += other.changeEntries
	this.domainReductions += other.domainReductions
	this.domainValsRemoved += other.domainValsRemoved
	this.controlEvents += other.controlEvents
	this.working_time += other.working_time
	this.idle_time += other.idle_time
}

// String returns compact readable representation of the statistics
// of a store.
func (this *StoreStatistics) String() string {
	// There are V variables and P propagators with ChanSize Channelbuffer
	// There have been CEvt ChangeEvents of which empty are empty
	// There have been CEnt ChangeEntries with red reductions of
	//    which rem actually removed a value
	// There have been CtrlEvt ControlEvents
	// All for one Store, resetted while cloning
	f := "%d variables with %d propagators buffer %d elements on channels\n"
	// f += "      currently %d variables with %d propagators\n"
	f += "      %3d/%3d ChEvt/empty, %3d/%3d/%3d ChEnt/red/rem;"
	f += " %3d CtrlEvt\n"
	f += " %10dms Store-Idle-Time, %10dms Store-Working-Time"
	s := fmt.Sprintf(f,
		this.variables, this.propagators, this.sizeChannels,
		// this.act_variables, this.act_propagators,
		this.changeEvents, this.emptyChangeEvents,
		this.changeEntries, this.domainReductions, this.domainValsRemoved,
		this.controlEvents, this.idle_time/1000000, this.working_time/1000000)
	return s
}

// InitStatTime sets the start time.
func (this *StoreStatistics) InitStatTime() {
	this.end_working_time = time.Now().UnixNano()
}

// LogIdleTime logs the idle time in the store.
func (this *StoreStatistics) LogIdleTime() {
	this.start_working_time = time.Now().UnixNano()
	this.idle_time += this.start_working_time - this.end_working_time
}

// LogWorkingTime logs the working time in the store.
func (this *StoreStatistics) LogWorkingTime() {
	this.end_working_time = time.Now().UnixNano()
	this.working_time += this.end_working_time - this.start_working_time
}

// AddIdleTime adds a delta to the idle time.
func (this *StoreStatistics) AddIdleTime(v int64) {
	this.idle_time += v
}

// AddWorkingTime adds a delta to the working time.
func (this *StoreStatistics) AddWorkingTime(v int64) {
	this.working_time += v
}

// GetVariables gets a reference to the variables
func (this *StoreStatistics) GetVariables() int {
	return this.variables
}

// GetActVariables gets a reference to the current variables
func (this *StoreStatistics) GetActVariables() int {
	return this.act_variables
}

// GetPropagators gets the number of propagators.
func (this *StoreStatistics) GetPropagators() int {
	return this.propagators
}

// GetActPropagators gets the number of actual propagators.
func (this *StoreStatistics) GetActPropagators() int {
	return this.act_propagators
}
