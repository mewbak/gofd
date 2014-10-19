package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"strings"
)

// AllDistinct is a stronger propagator for the semantically equivalent
// AllDifferent. In addition to removing ground values from other variables
// it also checks whether the union of non ground values is at least as
// large as the number of non ground variables. Furthermore, if the number
// of non ground values equals the number of variables and a value just occurs
// once it forces that variable to be ground with that value.

type AllDistinct struct {
	vars             []core.VarId
	outCh            chan<- *core.ChangeEvent
	inCh             <-chan *core.ChangeEntry
	varidToDomainMap map[core.VarId]*core.ExDomain
	nonGroundVars    int // number of nonground variables
	// union of all values in nonground variables
	nonGroundValues *core.ExDomain
	// for each value in nonGroundValues in which variables it occurs
	witness map[int]map[core.VarId]bool
	id      core.PropId
	store   *core.Store
}

func (this *AllDistinct) Start() {
	loggerDebug := core.GetLogger().DoDebug()
	if loggerDebug {
		core.GetLogger().Dln("AllDistinct_'initial consistency check'")
	}
	// setup and initial check
	this.setup()
	evt := core.CreateChangeEvent()
	this.initialCheck(evt)
	this.outCh <- evt // send changes to store
	// continuous propagation
	for changeEntry := range this.inCh {
		if loggerDebug {
			core.GetLogger().Df("%s_'Incoming Change for %s'",
				this, this.store.GetName(changeEntry.GetID()))
		}
		varidChanged := changeEntry.GetID()
		evt = core.CreateChangeEvent()
		vals := changeEntry.GetValues().Values_asMap()
		cons := this.removeValues(varidChanged, vals, evt)
		if cons { // only if it has not yet failed
			this.inOutAll(varidChanged, evt)
		}
		if loggerDebug {
			msg := "%s_propagate_'communicate change, evt-value: %s'"
			core.GetLogger().Df(msg, this, evt)
		}
		this.outCh <- evt // send changes to store
	}
}

// Register registers the propagator at the store.
func (this *AllDistinct) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap(this.vars, this.id)
	this.varidToDomainMap = core.GetVaridToExplicitDomainsMap(domains)
	this.store = store
}

// SetID is only used by the store to set the propagator's ID
func (this *AllDistinct) SetID(propID core.PropId) {
	this.id = propID
}

func (this *AllDistinct) GetID() core.PropId {
	return this.id
}

func (this *AllDistinct) String() string {
	vars_str := make([]string, len(this.vars))
	for i, var_id := range this.vars {
		vars_str[i] = this.store.GetName(var_id)
	}
	return fmt.Sprintf("PROP_%d %s",
		this.id, strings.Join(vars_str, "!!="))
}

func (this *AllDistinct) StringDomains() string {
	vars_domains_str := make([]string, len(this.vars))
	for i, var_id := range this.vars {
		vars_domains_str[i] = fmt.Sprintf("%s{%s} ",
			this.store.GetName(var_id),
			this.varidToDomainMap[var_id].String())
	}
	return fmt.Sprintf("PROP_%d %s\n  nonGVars %d, nonGValues %s",
		this.id, strings.Join(vars_domains_str, "!!="),
		this.nonGroundVars, this.nonGroundValues.String())
}

// CreateAlldistinct creates single propagator with stronger propagation
// compared to AllDifferent.
func CreateAlldistinct(vars ...core.VarId) *AllDistinct {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateAllDistinct-propagator")
	}
	prop := new(AllDistinct)
	prop.vars = vars
	return prop
}

func (this *AllDistinct) Clone() core.Constraint {
	prop := new(AllDistinct)
	prop.vars = make([]core.VarId, len(this.vars))
	for i, varId := range this.vars {
		prop.vars[i] = varId
	}
	return prop
}

// setup initializes a registered propagator to setup additional data
// structures preparing for the initial check and further propagation.
func (this *AllDistinct) setup() {
	this.nonGroundVars = len(this.vars)
	this.nonGroundValues = core.CreateExDomain()
	this.witness = make(map[int]map[core.VarId]bool)
	for varId, domain := range this.varidToDomainMap {
		for value, _ := range domain.Values {
			this.nonGroundValues.Add(value)
			if _, exists := this.witness[value]; !exists {
				this.witness[value] = make(map[core.VarId]bool)
			}
			this.witness[value][varId] = true
		}
	}
}

func (this *AllDistinct) removeValues(varId core.VarId, values map[int]bool,
	evt *core.ChangeEvent) bool {
	changedDomain := this.varidToDomainMap[varId]
	for value := range values {
		changedDomain.Remove(value)
		delete(this.witness[value], varId)
		if len(this.witness[value]) == 0 {
			// must have been nonGround, otherwise already inconsistent
			this.nonGroundValues.Remove(value)
			delete(this.witness, value)
			if this.checkNonGroundSize(evt) == false {
				return false // fail
			}
		}
	}
	return true
}

// initialCheck assumes that any variable could be ground and propagates
func (this *AllDistinct) initialCheck(evt *core.ChangeEvent) {
	for _, varId := range this.vars {
		this.inOutAll(varId, evt)
	}
	// may have initial hit
	this.checkNonGroundSize(evt)
	this.checkNonGroundFix(evt)
}

// checkNonGroundSize checks if the set of values of all unfixed values is
// smaller than the number of unfixed variables, and then fails early.
func (this *AllDistinct) checkNonGroundSize(evt *core.ChangeEvent) bool {
	if this.nonGroundValues.Size() < this.nonGroundVars { // fail
		// fail by removing values of any, here the first, variable
		outVarId := this.vars[0]
		chEntry := core.CreateChangeEntry(outVarId)
		for value := range this.varidToDomainMap[outVarId].Values {
			chEntry.Add(value)
		}
		evt.AddChangeEntry(chEntry)
		return false
	}
	return true
}

// checkNonGroundFix checks if the set of values of all unfixed values
// equals the number of unfixed variables and a value occurs just once,
// and then fixes that variable to that value.
func (this *AllDistinct) checkNonGroundFix(evt *core.ChangeEvent) {
	if this.nonGroundValues.Size() == this.nonGroundVars { // may force
		for value, vars := range this.witness {
			if len(vars) == 1 { // value occurs once and we are tight
				// force, fix that single occurence
				var varId = core.VarId(-1)
				for vid := range vars { // ugly (?), get single entry
					varId = vid // will find it; just once!
				}
				chEntry := core.CreateChangeEntry(varId)
				for toremove := range this.varidToDomainMap[varId].Values {
					if toremove != value {
						chEntry.Add(toremove)
					}
				}
				evt.AddChangeEntry(chEntry)
			}
		}
	}
}

// inOutAll, as AllDistinct plus more powerful propagators
func (this *AllDistinct) inOutAll(inVarId core.VarId,
	evt *core.ChangeEvent) {
	inDomain := this.varidToDomainMap[inVarId]
	if inDomain.IsGround() { // newly identified as ground
		this.nonGroundVars -= 1     // one less non ground
		fixed_value := inDomain.Min // the ground value
		// remove from non ground vars, alldifferent propagation
		for outVarId, outDomain := range this.varidToDomainMap {
			if inVarId == outVarId { // must be a different one
				continue
			}
			if outDomain.Contains(fixed_value) {
				chEntry := core.CreateChangeEntry(outVarId)
				chEntry.Add(fixed_value)
				evt.AddChangeEntry(chEntry)
			}
		}
		this.nonGroundValues.Remove(fixed_value) // no longer non ground
		this.checkNonGroundSize(evt)
		this.checkNonGroundFix(evt)
	}
}

func (this *AllDistinct) GetVarIds() []core.VarId {
	return this.vars
}

func (this *AllDistinct) GetDomains() []core.Domain {
	return core.ValuesOfMapVarIdToExDomain(this.vars, this.varidToDomainMap)
}

func (this *AllDistinct) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *AllDistinct) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
