// package explicit provides propagators that are geared towards
// using finite domain variables represented with an explicit domain
package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"strings"
)

// AlldifferentPrimitives is a proper propagator that is semantically
// equivalent to CreateAlldifferent
// (quadratically many "not equal" constraints),
// but just holds one copy of the involved variables and removes all values
// of ground variables in all other variables.
// No stronger propagation techniques of a global constraint are used.
type AlldifferentPrimitives struct {
	vars             []core.VarId
	outCh            chan<- *core.ChangeEvent
	inCh             <-chan *core.ChangeEntry
	varidToDomainMap map[core.VarId]*core.ExDomain
	id               core.PropId
	store            *core.Store
}

func (this *AlldifferentPrimitives) Start() {
	loggerDebug := core.GetLogger().DoDebug()
	if loggerDebug {
		msg := "AlldifferentPrimitives_'initial consistency check'"
		core.GetLogger().Dln(msg)
	}
	// initial check
	evt := core.CreateChangeEvent()
	this.initialCheck(evt)
	this.sendChangesToStore(evt)
	for changeEntry := range this.inCh {
		if loggerDebug {
			core.GetLogger().Df("%s_'Incoming Change for %s'",
				this, this.store.GetName(changeEntry.GetID()))
		}
		varidChanged := changeEntry.GetID()
		changedDomain := this.varidToDomainMap[varidChanged]
		changedDomain.Removes(changeEntry.GetValues())
		evt = core.CreateChangeEvent()
		this.inOutAll(varidChanged, evt)
		this.sendChangesToStore(evt)
	}
}

// sendChangesToStore send the collected changes (stored in the event)
// to the store.
func (this *AlldifferentPrimitives) sendChangesToStore(evt *core.ChangeEvent) {
	if core.GetLogger().DoDebug() {
		msg := "%s_propagate_'communicate change, evt-value: %s'"
		core.GetLogger().Df(msg, this, evt)
	}
	this.outCh <- evt
}

// initialCheck checks for each variable whether it is ground and thus might
// propagate deletions.
func (this *AlldifferentPrimitives) initialCheck(evt *core.ChangeEvent) {
	for _, varId := range this.vars {
		this.inOutAll(varId, evt)
	}
}

// Register registers the propagator at the store. Here, the propagator gets
// his needed channels and domains and stores them in his struct.
func (this *AlldifferentPrimitives) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap(this.vars, this.id)
	this.varidToDomainMap = core.GetVaridToExplicitDomainsMap(domains)
	this.store = store
}

// SetID is used by the store to set the propagator's ID, don't use it
// yourself or bad things will happen.
func (this *AlldifferentPrimitives) SetID(propID core.PropId) {
	this.id = propID
}

func (this *AlldifferentPrimitives) GetID() core.PropId {
	return this.id
}

func (this *AlldifferentPrimitives) String() string {
	vars_str := make([]string, len(this.vars))
	for i, var_id := range this.vars {
		vars_str[i] = this.store.GetName(var_id)
	}
	return fmt.Sprintf("PROP_%d %s",
		this.id, strings.Join(vars_str, "!="))
}

// CreateAlldifferent2 creates one single propagator, that for each variable
// that becomes ground removes that value from all other variables.
// Note: Alldifferent is not using stronger propagation techniques of
// a global constraint.
func CreateAlldifferentPrimitives(vars ...core.VarId) *AlldifferentPrimitives {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateAlldifferentPrimitives-propagator")
	}
	prop := new(AlldifferentPrimitives)
	prop.vars = vars
	return prop
}

func (this *AlldifferentPrimitives) Clone() core.Constraint {
	prop := new(AlldifferentPrimitives)
	prop.vars = make([]core.VarId, len(this.vars))
	for i, single_var := range this.vars {
		prop.vars[i] = single_var
	}
	return prop
}

// inOutAll if inDomain is Fixed remove that value from all other domains
// Example: X:{6}, Y:{6,7}, Z:{5,6,7} --> X:{6}, Y:{7}, Z:{5,7}
func (this *AlldifferentPrimitives) inOutAll(inVarId core.VarId, evt *core.ChangeEvent) {
	inDomain := this.varidToDomainMap[inVarId]
	if inDomain.IsGround() {
		fixed_value := inDomain.Min
		for _, outVarId := range this.vars {
			if inVarId == outVarId {
				continue
			}
			outDomain := this.varidToDomainMap[outVarId]
			if outDomain.Contains(fixed_value) {
				chEntry := core.CreateChangeEntry(outVarId)
				evt.AddChangeEntry(chEntry)
				chEntry.Add(fixed_value)
			}
		}
	}
}

// inOut if inDomain is Fixed remove that value from the other domain
// Example: X:{6}, Y:{6,7} --> X:{6}, Y:{7}
// ToDo: no longer used internally - only in NQueensProp.go;
//       to be replaced there as well
func inOut(inDomain *core.ExDomain, outDomain *core.ExDomain,
	outVarId core.VarId, evt *core.ChangeEvent) {
	if inDomain.IsGround() {
		fixed_value := inDomain.Min
		for out_val := range outDomain.Values {
			if fixed_value == out_val {
				chEntry := core.CreateChangeEntry(outVarId)
				evt.AddChangeEntry(chEntry)
				chEntry.Add(out_val)
				return
			}
		}
	}
}

func (this *AlldifferentPrimitives) GetVarIds() []core.VarId {
	return this.vars
}

func (this *AlldifferentPrimitives) GetDomains() []core.Domain {
	return core.ValuesOfMapVarIdToExDomain(this.vars, this.varidToDomainMap)
}

func (this *AlldifferentPrimitives) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *AlldifferentPrimitives) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
