// package interval provides propagators that are geared towards
// using finite domain variables represented with an intervals domain
package interval

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"strings"
)

// Alldifferent is a proper constraint for pairwise difference,
// i.e. X≠Y, X≠Z, Y≠Z for X,Y,Z. It behaves like quadratically many
// "not equal" constraints, but just holds one copy of the involved
// variables and removes all values of ground variables in all other
// variables.
// No stronger propagation techniques of a global constraint are used.
type Alldifferent struct {
	vars             []core.VarId
	outCh            chan<- *core.ChangeEvent
	inCh             <-chan *core.ChangeEntry
	varidToDomainMap map[core.VarId]*core.IvDomain
	id               core.PropId
}

func (this *Alldifferent) Start(store *core.Store) {
	core.LogInitConsistency(this)
	// initial check
	evt := core.CreateChangeEvent()
	this.initialCheck(evt)
	core.SendChangesToStore(evt, this)

	for changeEntry := range this.inCh {
		core.LogIncomingChange(this, store, changeEntry)

		varidChanged := changeEntry.GetID()
		changedDomain := this.varidToDomainMap[varidChanged]
		changedDomain.Removes(changeEntry.GetValues())
		evt = core.CreateChangeEvent()
		this.inOutIvAll(varidChanged, evt)
		core.SendChangesToStore(evt, this)
	}
}

// check for each variable whether it is ground and thus might
// propagate deletions
func (this *Alldifferent) initialCheck(evt *core.ChangeEvent) {
	for _, varId := range this.vars {
		this.inOutIvAll(varId, evt)
	}
}

// Register registers the propagator at the store. Here, the propagator gets
// his needed channels and domains and stores them in his struct
func (this *Alldifferent) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap(this.vars, this.id)
	this.varidToDomainMap = core.GetVaridToIntervalDomains(domains)
}

// SetID is used by the store to set the propagator's ID, don't use it
// yourself or bad things will happen
func (this *Alldifferent) SetID(propID core.PropId) {
	this.id = propID
}

func (this *Alldifferent) GetID() core.PropId {
	return this.id
}

func (this *Alldifferent) String() string {
	vars_str := make([]string, len(this.vars))
	for i, var_id := range this.vars {
		vars_str[i] = core.GetNameRegistry().GetName(var_id)
	}
	return fmt.Sprintf("PROP_%d %s",
		this.id, strings.Join(vars_str, "!="))
}

// CreateAlldifferent2 creates one single propagator, that for each variable
// that becomes ground removes that value from all other variables.
// Note: Alldifferent is not using stronger propagation techniques of
// a global constraint.
func CreateAlldifferent(vars ...core.VarId) *Alldifferent {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateAlldifferent_interval-propagator")
	}
	prop := new(Alldifferent)
	prop.vars = vars
	return prop
}

func (this *Alldifferent) Clone() core.Constraint {
	prop := new(Alldifferent)
	prop.vars = make([]core.VarId, len(this.vars))
	for i, single_var := range this.vars {
		prop.vars[i] = single_var
	}
	return prop
}

// inOutAll if inDomain is Fixed remove that value from all other domains
// Example: X:{6}, Y:{6,7}, Z:{5,6,7} --> X:{6}, Y:{7}, Z:{5,7}
func (this *Alldifferent) inOutIvAll(inVarId core.VarId, evt *core.ChangeEvent) {
	inDomain := this.varidToDomainMap[inVarId]
	if inDomain.IsGround() {
		fixed_value := inDomain.GetMin()
		for _, outVarId := range this.vars {
			if inVarId == outVarId {
				continue
			}
			outDomain := this.varidToDomainMap[outVarId]
			if outDomain.Contains(fixed_value) {
				chEntry := core.CreateChangeEntryWithValues(outVarId, core.CreateIvDomainFromTo(fixed_value, fixed_value))
				evt.AddChangeEntry(chEntry)
			}
		}
	}
}

func (this *Alldifferent) GetVarIds() []core.VarId {
	return this.vars
}

func (this *Alldifferent) GetDomains() []core.Domain {
	return core.ValuesOfMapVarIdToIvDomain(this.vars, this.varidToDomainMap)
}

func (this *Alldifferent) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *Alldifferent) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
