package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// implementation of the Range constraint
// signature: Range({X1,...,Xi}, S, T) iff T = union(X(i)) with i ∈ S
// if the index of a variable in xi appears in S,
// it must be assigned a value from T
// T can only contain values that appear in D(Xi) with i ∈ S
// (note: does not work with setvariables as it should)
type Range struct {
	xi                 []core.VarId
	s, t               core.VarId
	outCh              chan<- *core.ChangeEvent
	inCh               <-chan *core.ChangeEntry
	xi_Domains         map[core.VarId]core.Domain
	s_Domain, t_Domain core.Domain
	id                 core.PropId
	store              *core.Store
}

func (this *Range) Clone() core.Constraint {
	prop := new(Range)
	prop.xi, prop.s, prop.t = this.xi, this.s, this.t
	return prop
}

func (this *Range) Start(store *core.Store) {
	loggerDebug := core.GetLogger().DoDebug()
	if loggerDebug {
		core.GetLogger().Df("%s_'initial consistency check'", this)
	}
	evt := core.CreateChangeEvent()
	this.makeConsistent(evt)
	this.sendChangesToStore(evt)
	for changeEntry := range this.inCh {
		if loggerDebug {
			msg := "%s_'Incoming Change for %s'"
			core.GetLogger().Df(msg, this, store.GetName(changeEntry.GetID()))
		}
		evt = core.CreateChangeEvent()
		switch var_id := changeEntry.GetID(); var_id {
		case this.s:
			this.s_Domain.Removes(changeEntry.GetValues())
			this.makeConsistent(evt)
			break
		case this.t:
			this.s_Domain.Removes(changeEntry.GetValues())
			this.makeConsistent(evt)
			break
		default:
			this.xi_Domains[var_id].Removes(changeEntry.GetValues())
			this.makeConsistent(evt)
			break
		}
		this.sendChangesToStore(evt)
	}
}

// sendChangesToStore send the collected changes (stored in the event)
// to the store
func (this *Range) sendChangesToStore(evt *core.ChangeEvent) {
	if core.GetLogger().DoDebug() {
		msg := "%s_propagate_'communicate change, evt-value: %s'"
		core.GetLogger().Df(msg, this, evt)
	}
	this.outCh <- evt
}

// makeConsistent checks if the constraint is arc consistent
// and makes it arc consistent if possible
func (this *Range) makeConsistent(evt *core.ChangeEvent) bool {
	//if any domain is empty, propagation can be stopped
	if this.s_Domain.IsEmpty() || this.t_Domain.IsEmpty() {
		return false
	}
	for _, domain := range this.xi_Domains {
		if domain.IsEmpty() {
			return false
		}
	}

	countChangeEntries := 0

	// remove all values from D(S) which are bigger than the amount of variables
	// or the indexes of the Xi having domains that do not intersect with T
	var chEntry *core.ChangeEntry = nil
	for index := range this.s_Domain.Values_asMap() {
		if index > len(this.xi) {
			if chEntry == nil {
				chEntry = core.CreateChangeEntry(this.s)
			}
			chEntry.Add(index)
		} else {
			x_domain := this.xi_Domains[this.xi[index-1]]
			if x_domain.Intersection(this.t_Domain).IsEmpty() {
				if chEntry == nil {
					chEntry = core.CreateChangeEntry(this.s)
				}
				chEntry.Add(index)
			}
		}
	}
	if chEntry != nil {
		evt.AddChangeEntry(chEntry)
		countChangeEntries += 1
	}
	// if the index of a variable X in Xi appears in D(S) and its domain
	// contains values from and outside of D(T), all values that are
	// not part of D(T) are removed from D(X).
	for index := range this.s_Domain.Values_asMap() {
		if index <= len(this.xi) {
			x_domain := this.xi_Domains[this.xi[index-1]]
			if !x_domain.IsSubset(this.t_Domain) &&
				(!x_domain.Intersection(this.t_Domain).IsEmpty()) {
				chEntry = nil
				if chEntry == nil {
					chEntry = core.CreateChangeEntry(this.xi[index-1])
				}
				newDomain := this.t_Domain.Difference(x_domain)
				chEntry.AddValues(newDomain.Values_asMap())
				if chEntry != nil {
					evt.AddChangeEntry(chEntry)
					countChangeEntries += 1
				}
			}
		}
	}

	// propagate T
	// create the union of all domains of the Xi which index appears in D(S)
	unionOfXiDoms := core.CreateExDomain()
	for index := range this.s_Domain.Values_asMap() {
		if index <= len(this.xi) {
			x_domain := this.xi_Domains[this.xi[index-1]]
			temp := x_domain.Union(unionOfXiDoms)
			unionOfXiDoms = temp.(*core.ExDomain)
		}
	}
	// compute the values that have to be removed from T since
	// they don't appear in any domain of Xi
	valuesToRemoveFromT := unionOfXiDoms.Difference(this.t_Domain)
	chEntry = nil
	chEntry = core.CreateChangeEntry(this.t)
	chEntry.AddValues(valuesToRemoveFromT.Values_asMap())
	if chEntry != nil {
		evt.AddChangeEntry(chEntry)
		countChangeEntries += 1
	}

	if countChangeEntries != 0 {
		return false
	}
	return true
}

// Register registers the propagator at the store.
func (this *Range) Register(store *core.Store) {
	var domains []core.Domain
	idSlice := make([]core.VarId, len(this.xi)+2)
	i := 0
	for i = 0; i < len(this.xi); i++ {
		idSlice[i] = this.xi[i]
	}
	idSlice[i] = this.s
	idSlice[i+1] = this.t
	this.inCh, domains, this.outCh =
		store.RegisterPropagator(idSlice, this.id)

	this.xi_Domains = make(map[core.VarId]core.Domain)
	for i = 0; i < (len(idSlice) - 2); i++ {
		this.xi_Domains[idSlice[i]] = domains[i]
	}
	this.s_Domain = domains[len(idSlice)-2]
	this.t_Domain = domains[len(idSlice)-1]
	this.store = store
}

func (this *Range) SetID(propID core.PropId) {
	this.id = propID
}

func (this *Range) GetID() core.PropId {
	return this.id
}

func CreateRange(xi []core.VarId, s core.VarId, t core.VarId) *Range {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateRange_propagator")
	}
	prop := new(Range)
	prop.xi = xi
	prop.s = s
	prop.t = t
	return prop
}

func (this *Range) String() string {
	var s string
	for i := 0; i < len(this.xi); i++ {
		s += this.store.GetName(this.xi[i])
	}
	return fmt.Sprintf("PROP_RANGE({%s}, {%s}, %s)",
		s,
		this.store.GetName(this.s),
		this.store.GetName(this.t))
}

func (this *Range) GetVarIds() []core.VarId {
	varIds := make([]core.VarId, len(this.xi))
	for _, var_id := range this.xi {
		varIds = append(varIds, var_id)
	}
	varIds = append(varIds, this.s)
	varIds = append(varIds, this.t)
	return varIds
}

func (this *Range) GetDomains() []core.Domain {
	domains := make([]core.Domain, len(this.xi_Domains))
	for _, var_id := range this.xi {
		domains = append(domains, this.xi_Domains[var_id])
	}
	domains = append(domains, this.s_Domain)
	domains = append(domains, this.t_Domain)
	return domains
}

func (this *Range) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *Range) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
