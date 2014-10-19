package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// implementation for the Among constraint with local consistency
// signature: Among({X1,...,Xi}, K, N)
// values from K can be assigned to to N variables in {X1,...,Xi}
type Among struct {
	xi         []core.VarId
	n          core.VarId
	k          *core.ExDomain
	outCh      chan<- *core.ChangeEvent
	inCh       <-chan *core.ChangeEntry
	xi_Domains map[core.VarId]*core.ExDomain
	n_Domain   *core.ExDomain
	id         core.PropId
	store      *core.Store
	lb, ub     int
	minN, maxN int
}

func (this *Among) Clone() core.Constraint {
	prop := new(Among)
	prop.xi, prop.k, prop.n = this.xi, this.k, this.n
	return prop
}

// Start propagates for all variables initially and listens to channels
// for incoming domain updates in order to propagate accordingly.
func (this *Among) Start() {
	loggerDebug := core.GetLogger().DoDebug()
	if loggerDebug {
		core.GetLogger().Df("%s_'initial consistency check'", this)
	}
	// initial propagation
	evt := core.CreateChangeEvent()
	this.initialConsistencyCheck(evt)
	this.sendChangesToStore(evt)
	for changeEntry := range this.inCh {
		if loggerDebug {
			msg := "%s_'Incoming Change for %s'"
			core.GetLogger().Df(msg, this, 
				this.store.GetName(changeEntry.GetID()))
		}
		// handle incoming events and propagate if necessary
		evt = core.CreateChangeEvent()
		switch var_id := changeEntry.GetID(); var_id {
		case this.n:
			this.n_Domain.Removes(changeEntry.GetValues())
			// in case that N has changed, propagate the Xi
			changeValues := changeEntry.GetValues().Values_asMap()
			if changeValues[this.minN] || changeValues[this.maxN] {
				this.minN, this.maxN = this.n_Domain.GetMinAndMax()
				this.NinXiout(evt)
			}
			break
		default:
			// check is necessary since irrelevant domains were removed
			// from this map
			if this.xi_Domains[var_id] != nil {
				xDom := this.xi_Domains[var_id]
				xDom.Removes(changeEntry.GetValues())
				// in case that an xi has changed, propagate N and the Xi
				if xDom.IsSubset(this.k) {
					// if xDom is now a subset of K, then increase lb and
					// remove the old value of lb from the domain of N
					this.XiinNout(evt, true, false)
					this.lb += 1
					delete(this.xi_Domains, var_id)
					this.NinXiout(evt)
				} else if xDom.Intersection(this.k).IsEmpty() {
					// if xDom has no value with K in common anymore, then
					// decrease ub and remove the old value of ub from the
					// domain of N
					this.XiinNout(evt, false, true)
					this.ub -= 1
					delete(this.xi_Domains, var_id)
					this.NinXiout(evt)
				}
			}
			break
		}
		this.sendChangesToStore(evt)
	}
}

// sendChangesToStore send the collected changes (stored in the event)
// to the store
func (this *Among) sendChangesToStore(evt *core.ChangeEvent) {
	if core.GetLogger().DoDebug() {
		msg := "%s_propagate_'communicate change, evt-value: %s'"
		core.GetLogger().Df(msg, this, evt)
	}
	this.outCh <- evt
}

// removeIrrelevantDomains removes the domains from the propagator's local copy
// that have no value in common with K or are subsets of K
// the method also propagates N initially and removes old lbs and ubs from
// its domain.
func (this *Among) removeIrrelevantDomains(evt *core.ChangeEvent) {
	valuesToRemove := make([]core.VarId, 0)
	// store the number of subsets and intersecting domains in order
	// to compute lb and ub correctly in the future
	this.lb = 0
	this.ub = len(this.xi_Domains)
	var chEntry *core.ChangeEntry = nil
	for varid, domain := range this.xi_Domains {
		tempDomain := domain.Intersection(this.k)
		if domain.IsSubset(this.k) {
			// if D(N) contains old lb, remove it
			if this.n_Domain.Contains(this.lb) {
				if chEntry == nil {
					chEntry = core.CreateChangeEntry(this.n)
				}
				chEntry.Add(this.lb)
			}
			this.lb += 1
			valuesToRemove = append(valuesToRemove, varid)
		} else if tempDomain.IsEmpty() {
			// if D(N) contains old ub, remove it
			if this.n_Domain.Contains(this.ub) {
				if chEntry == nil {
					chEntry = core.CreateChangeEntry(this.n)
				}
				chEntry.Add(this.ub)
			}
			this.ub -= 1
			valuesToRemove = append(valuesToRemove, varid)
		}
	}

	if chEntry != nil {
		evt.AddChangeEntry(chEntry)
	}

	for _, varid := range valuesToRemove {
		delete(this.xi_Domains, varid)
	}
}

// XiinNout propagates N
// lbUp indicates whether lb was increased or not,
// ubDown indicates whether ub was decreased or not.
func (this *Among) XiinNout(evt *core.ChangeEvent, lbUp, ubDown bool) {
	// remove all values from D(N) that are smaller than lb or greater than ub
	// and therefore can not be part of any satisfying assignment
	var chEntry *core.ChangeEntry = nil
	if lbUp {
		//if D(N) contains the old lb, remove it from D(N)
		if this.n_Domain.Contains(this.lb) {
			if chEntry == nil {
				chEntry = core.CreateChangeEntry(this.n)
			}
			chEntry.Add(this.lb)
		}
	}

	if ubDown {
		// if D(N) contains the old ub, remove it from D(N)
		if this.n_Domain.Contains(this.ub) {
			if chEntry == nil {
				chEntry = core.CreateChangeEntry(this.n)
			}
			chEntry.Add(this.ub)
		}
	}

	if chEntry != nil {
		evt.AddChangeEntry(chEntry)
	}

}

// NinXiout propagates the Xi
func (this *Among) NinXiout(evt *core.ChangeEvent) {
	valuesToRemove := make([]core.VarId, 0)
	subtractUb := 0
	addLb := 0
	var chEntry *core.ChangeEntry = nil
	// if only one element exists in D(N) which has the same value
	// as the lower bound
	if this.lb == this.maxN {
		for varid, domain := range this.xi_Domains {
			chEntry = nil
			// the domains of the xi which are no subsets of D(K) are checked
			// and all values intersecting with D(K) are removed
			if !domain.IsSubset(this.k) &&
				(!domain.Intersection(this.k).IsEmpty()) {
				if chEntry == nil {
					chEntry = core.CreateChangeEntry(varid)
				}
				newDomain := domain.Intersection(this.k)
				chEntry.AddValues(newDomain)
				evt.AddChangeEntry(chEntry)
				// decrease ub because the current domain is now in state f and
				// has become irrelevant
				subtractUb += 1
				valuesToRemove = append(valuesToRemove, varid)
			}
		}
	}

	// if only one element exists in D(N) which has the same value
	// as the upper bound
	if this.ub == this.minN {
		for varid, domain := range this.xi_Domains {
			chEntry = nil
			newDomain := domain.Intersection(this.k)
			// the domains of the xi which are intersecting with D(K)
			// are checked and all values outside of D(K) are removed
			if !newDomain.IsEmpty() {
				if chEntry == nil {
					chEntry = core.CreateChangeEntry(varid)
				}
				newDomain = domain.Difference(this.k)
				chEntry.AddValues(newDomain)
				valuesToRemove = append(valuesToRemove, varid)
				// increase lb because the current domain is now in state t
				// and has become irrelevant
				addLb += 1
			}
			if chEntry != nil {
				evt.AddChangeEntry(chEntry)
			}
		}
	}

	this.lb += addLb
	this.ub -= subtractUb

	for _, varid := range valuesToRemove {
		delete(this.xi_Domains, varid)
	}
}

// initalConsistencyCheck checks if the constraint is arc consistent
// and makes it arc consistent if possible (only used initially)
func (this *Among) initialConsistencyCheck(evt *core.ChangeEvent) {
	this.minN, this.maxN = this.n_Domain.GetMinAndMax()
	// remove domains that don't contain values from D(K) or are subsets
	// and remove values from N that are smaller than lb and greater than ub
	this.removeIrrelevantDomains(evt)
	// propagate Xi
	this.NinXiout(evt)
}

// Register registers the propagator at the store.
func (this *Among) Register(store *core.Store) {
	var domains []core.Domain
	// create one big slice containing all needed
	// VarIds and send it to the store
	idSlice := make([]core.VarId, len(this.xi)+1)
	i := 0
	//create the map for the switch-case in Start
	for i = 0; i < len(this.xi); i++ {
		idSlice[i] = this.xi[i]
	}
	idSlice[i] = this.n
	this.inCh, domains, this.outCh =
		store.RegisterPropagator(idSlice, this.id)

	ds := core.GetExDomainSlice(domains)

	this.xi_Domains = make(map[core.VarId]*core.ExDomain)
	for i = 0; i < (len(idSlice) - 1); i++ {
		this.xi_Domains[idSlice[i]] = ds[i]
	}
	this.n_Domain = domains[i].(*core.ExDomain)
	this.store = store
}

func (this *Among) SetID(propID core.PropId) {
	this.id = propID
}

func (this *Among) GetID() core.PropId {
	return this.id
}

func CreateAmong(xi []core.VarId, k []int, n core.VarId) *Among {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateAmong_propagator")
	}
	prop := new(Among)
	prop.xi = xi
	prop.k = core.CreateExDomainAdds(k)
	prop.n = n
	return prop
}

func (this *Among) String() string {
	var s string
	for i := 0; i < len(this.xi); i++ {
		s += this.store.GetName(this.xi[i])
	}
	var kstring string
	for val, _ := range this.k.Values {
		kstring += fmt.Sprintf("%v,", val)
	}
	return fmt.Sprintf("PROP_AMONG({%s}, {%s}, %s)",
		s,
		kstring,
		this.store.GetName(this.n))
}

func (this *Among) GetVarIds() []core.VarId {
	varIds := make([]core.VarId, len(this.xi))
	for _, var_id := range this.xi {
		varIds = append(varIds, var_id)
	}
	varIds = append(varIds, this.n)
	return varIds
}

func (this *Among) GetDomains() []core.Domain {
	domains := make([]core.Domain, len(this.xi_Domains))
	for _, var_id := range this.xi {
		domains = append(domains, this.xi_Domains[var_id])
	}
	domains = append(domains, this.n_Domain)
	return domains
}

func (this *Among) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *Among) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
