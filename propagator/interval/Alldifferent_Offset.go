package interval

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"strings"
)

// Alldifferent_Offset is similar to Alldifferent, but allows to specify
// an offset for each variable. Thus, e.g. X+dX≠Y+dY, X+dX≠Z+dZ, Y+dY≠Z+dZ
// must hold for three variales {X, Y, Z} and offsets {dX, dY, dZ}.
// No stronger propagation techniques of a global constraint are used.
type Alldifferent_Offset struct {
	vars             []core.VarId
	offsets          []int
	outCh            chan<- *core.ChangeEvent
	inCh             <-chan *core.ChangeEntry
	varidToOffsetMap map[core.VarId]int
	varidToDomainMap map[core.VarId]*core.IvDomain
	id               core.PropId
	store            *core.Store
}

func (this *Alldifferent_Offset) Start() {
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

func (this *Alldifferent_Offset) initialCheck(evt *core.ChangeEvent) {
	for _, varId := range this.vars {
		this.inOutAll(varId, evt)
	}
}

// inOutAll if inDomain is Fixed remove that value (plus/minus offset) from
// all other domains.
func (this *Alldifferent_Offset) inOutAll(inVarId core.VarId,
	evt *core.ChangeEvent) {
	inDomain := this.varidToDomainMap[inVarId]
	if inDomain.IsGround() {
		fixed_value := inDomain.GetMin()
		fixed_value += this.varidToOffsetMap[inVarId] // value of left hand side
		// println("new ground var id", inVarId, "is", inDomain.GetMin())
		for _, outVarId := range this.vars {
			if inVarId == outVarId {
				continue
			}
			valToRemove := fixed_value - this.varidToOffsetMap[outVarId]
			outDomain := this.varidToDomainMap[outVarId]
			if outDomain.Contains(valToRemove) {
				chEntry := core.CreateChangeEntryWithIntValue(outVarId, valToRemove)
				// println("      changeEntry", chEntry.String())
				evt.AddChangeEntry(chEntry)
			}
		}
	}
}

// CreateAlldifferentOffset creates one propagator that ensure that each
// variable subject to an individual offset if not equal to any other
// variable and its offset.
func CreateAlldifferent_Offset(vars []core.VarId,
	offsets []int) *Alldifferent_Offset {
	if len(vars) != len(offsets) {
		panic("AlldifferentOffset-Creation: len(vars) != len(offsets)")
	}
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateAlldifferent_Offset-propagator")
	}
	prop := new(Alldifferent_Offset)
	prop.vars = vars
	prop.offsets = offsets
	prop.varidToOffsetMap = make(map[core.VarId]int, len(prop.vars))
	for i, varid := range vars {
		prop.varidToOffsetMap[varid] = offsets[i]
	}
	return prop
}

func (this *Alldifferent_Offset) Clone() core.Constraint {
	prop := new(Alldifferent_Offset)
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

// Register registers the propagator at the store.
func (this *Alldifferent_Offset) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap(this.vars, this.id)
	this.varidToDomainMap = core.GetVaridToIntervalDomains(domains)
	this.store = store
}

// SetID is used by the store to set the propagator's ID, don't use it
// yourself or bad things will happen.
func (this *Alldifferent_Offset) SetID(propID core.PropId) {
	this.id = propID
}

func (this *Alldifferent_Offset) GetID() core.PropId {
	return this.id
}

func (this *Alldifferent_Offset) String() string {
	vars_str := make([]string, len(this.vars))
	for i, var_id := range this.vars {
		vars_str[i] = fmt.Sprintf("%d*%s",
			this.varidToOffsetMap[var_id], this.store.GetName(var_id))
	}
	return fmt.Sprintf("PROP_Alldifferent_Offset %d %s",
		this.id, strings.Join(vars_str, ", "))
}

func (this *Alldifferent_Offset) GetVarIds() []core.VarId {
	return this.vars
}

func (this *Alldifferent_Offset) GetDomains() []core.Domain {
	return core.ValuesOfMapVarIdToIvDomain(this.vars, this.varidToDomainMap)
}

func (this *Alldifferent_Offset) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *Alldifferent_Offset) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
