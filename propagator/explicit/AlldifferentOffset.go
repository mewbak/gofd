package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

type AlldifferentOffset struct {
	vars             []core.VarId
	offsets          []int
	outCh            chan<- *core.ChangeEvent
	inCh             <-chan *core.ChangeEntry
	varidToOffsetMap map[core.VarId]int
	varidToDomainMap map[core.VarId]*core.ExDomain
	id               core.PropId
	store            *core.Store
}

func (this *AlldifferentOffset) Start() {
	core.LogInitConsistency(this)
	// initial check
	evt := core.CreateChangeEvent()
	this.initialCheck(evt)
	core.SendChangesToStore(evt, this)
	for changeEntry := range this.inCh {
		core.LogIncomingChange(this, this.store, changeEntry)
		varidChanged := changeEntry.GetID()
		changedDomain := this.varidToDomainMap[varidChanged]
		changedDomain.Removes(changeEntry.GetValues())
		evt = core.CreateChangeEvent()
		this.inOutAll(varidChanged, evt)
		core.SendChangesToStore(evt, this)
	}
}

func (this *AlldifferentOffset) initialCheck(evt *core.ChangeEvent) {
	for _, varId := range this.vars {
		this.inOutAll(varId, evt)
	}
}

// inOutAll if inDomain is Fixed remove that value (plus/minus offset from
// all other domains.
func (this *AlldifferentOffset) inOutAll(inVarId core.VarId,
	evt *core.ChangeEvent) {
	inDomain := this.varidToDomainMap[inVarId]
	if inDomain.IsGround() {
		fixed_value := inDomain.GetAnyElement()
		fixed_value += this.varidToOffsetMap[inVarId] // value of left hand side
		for _, outVarId := range this.vars {
			if inVarId == outVarId {
				continue
			}
			valToRemove := fixed_value - this.varidToOffsetMap[outVarId]
			outDomain := this.varidToDomainMap[outVarId]
			if outDomain.Contains(valToRemove) {
				chEntry := core.CreateChangeEntryWithIntValue(outVarId, valToRemove)
				evt.AddChangeEntry(chEntry)
			}
		}
	}
}

// CreateAlldifferentOffset creates one propagator that ensure that each
// variable subject to an individual offset if not equal to any other
// variable and its offset.
func CreateAlldifferentOffset(vars []core.VarId,
	offsets []int) *AlldifferentOffset {
	if len(vars) != len(offsets) {
		panic("AlldifferentOffset-Creation: len(vars) != len(offsets)")
	}
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateAlldifferentOffset-propagator")
	}
	prop := new(AlldifferentOffset)
	prop.vars = vars
	prop.offsets = offsets
	prop.varidToOffsetMap = make(map[core.VarId]int, len(prop.vars))
	for i, varid := range vars {
		prop.varidToOffsetMap[varid] = offsets[i]
	}
	return prop
}

func (this *AlldifferentOffset) Clone() core.Constraint {
	prop := new(AlldifferentOffset)
	prop.vars = make([]core.VarId, len(this.vars))
	for i, varid := range this.vars {
		prop.vars[i] = varid
	}
	prop.offsets = make([]int, len(this.offsets))
	for i, off := range this.offsets {
		prop.offsets[i] = off
	}
	prop.varidToOffsetMap = make(map[core.VarId]int, len(this.vars))
	for i, varid := range this.vars {
		prop.varidToOffsetMap[varid] = this.offsets[i]
	}
	return prop
}

// nothing to change below

// Register registers the propagator at the store.
func (this *AlldifferentOffset) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap(this.vars, this.id)
	this.varidToDomainMap = core.GetVaridToExplicitDomainsMap(domains)
	this.store = store
}

// SetID is used by the store to set the propagator's ID, don't use it
// yourself or bad things will happen.
func (this *AlldifferentOffset) SetID(propID core.PropId) {
	this.id = propID
}

func (this *AlldifferentOffset) GetID() core.PropId {
	return this.id
}

func (this *AlldifferentOffset) String() string {
	return fmt.Sprintf("PROP_AlldifferentOffset-ID: %d ", this.id)
}

func (this *AlldifferentOffset) GetVarIds() []core.VarId {
	return this.vars
}

func (this *AlldifferentOffset) GetDomains() []core.Domain {
	return core.ValuesOfMapVarIdToExDomain(this.vars, this.varidToDomainMap)
}

func (this *AlldifferentOffset) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *AlldifferentOffset) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
