package explicit

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// simple implementation for the Among constraint
// signature: Among({X1,...,Xi}, K, N)
// a value from K can be assigned to to N variables in {X1,...,Xi}
type AmongSimple struct {
	xi         []core.VarId
	n          core.VarId
	k          core.Domain
	outCh      chan<- *core.ChangeEvent
	inCh       <-chan *core.ChangeEntry
	xi_Domains map[core.VarId]core.Domain
	n_Domain   core.Domain
	id         core.PropId
	store      *core.Store
}

func (this *AmongSimple) Clone() core.Constraint {
	prop := new(AmongSimple)
	prop.xi, prop.k, prop.n = this.xi, this.k, this.n
	return prop
}

// Start checks for consistency and listens to channels for incoming domain
// updates.
func (this *AmongSimple) Start(store *core.Store) {
	loggerDebug := core.GetLogger().DoDebug()
	if loggerDebug {
		core.GetLogger().Df("%s_'initial consistency check'", this)
	}
	evt := core.CreateChangeEvent()
	this.makeConsistentIfGround(evt)
	this.sendChangesToStore(evt)
	for changeEntry := range this.inCh {
		if loggerDebug {
			msg := "%s_'Incoming Change for %s'"
			core.GetLogger().Df(msg, this, core.GetNameRegistry().GetName(changeEntry.GetID()))
		}
		evt = core.CreateChangeEvent()
		switch var_id := changeEntry.GetID(); var_id {
		case this.n:
			this.n_Domain.Removes(changeEntry.GetValues())
			this.makeConsistentIfGround(evt)
			break
		default:
			this.xi_Domains[var_id].Removes(changeEntry.GetValues())
			this.makeConsistentIfGround(evt)
			break
		}
		this.sendChangesToStore(evt)
	}
}

// sendChangesToStore send the collected changes (stored in the event)
// to the store
func (this *AmongSimple) sendChangesToStore(evt *core.ChangeEvent) {
	if core.GetLogger().DoDebug() {
		msg := "%s_propagate_'communicate change, evt-value: %s'"
		core.GetLogger().Df(msg, this, evt)
	}
	this.outCh <- evt
}

// makeConsistentIfGround checks if X_i is consistent
// and makes it consistent only if N is ground
func (this *AmongSimple) makeConsistentIfGround(evt *core.ChangeEvent) {
	if !this.n_Domain.IsGround() {
		return
	}

	for _, domain := range this.xi_Domains {
		if !domain.IsGround() {
			return
		}
	}

	countSubsets := 0
	for _, domain := range this.xi_Domains {
		if domain.IsSubset(this.k) {
			countSubsets += 1
		}
	}

	var chEntry *core.ChangeEntry = nil
	//in case that the amount of subsets does not equal the only value in N
	nValue := this.n_Domain.GetAnyElement()
	if nValue != countSubsets {
		chEntry = core.CreateChangeEntry(this.n)
		chEntry.Add(nValue)
		evt.AddChangeEntry(chEntry)
		return
	}
}

// Register registers the propagator at the store.
func (this *AmongSimple) Register(store *core.Store) {
	var domains []core.Domain
	idSlice := make([]core.VarId, len(this.xi)+1)
	i := 0
	for i = 0; i < len(this.xi); i++ {
		idSlice[i] = this.xi[i]
	}
	idSlice[i] = this.n
	this.inCh, domains, this.outCh =
		store.RegisterPropagator(idSlice, this.id)

	this.xi_Domains = make(map[core.VarId]core.Domain)
	for i = 0; i < (len(idSlice) - 1); i++ {
		this.xi_Domains[idSlice[i]] = domains[i]
	}
	this.n_Domain = domains[len(idSlice)-1]

}

func (this *AmongSimple) SetID(propID core.PropId) {
	this.id = propID
}

func (this *AmongSimple) GetID() core.PropId {
	return this.id
}

func CreateAmongSimple(xi []core.VarId, k []int, n core.VarId) *AmongSimple {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateCmultXeqCmultY_propagator")
	}
	prop := new(AmongSimple)
	prop.xi = xi
	prop.k = core.CreateExDomainAdds(k)
	prop.n = n
	return prop
}

func (this *AmongSimple) String() string {
	var s string
	for i := 0; i < len(this.xi); i++ {
		s += core.GetNameRegistry().GetName(this.xi[i])
	}
	var kstring string
	for val, _ := range this.k.Values_asMap() {
		kstring += fmt.Sprintf("%v,", val)
	}
	return fmt.Sprintf("PROP_AMONG({%s}, {%s}, %s)",
		s,
		kstring,
		core.GetNameRegistry().GetName(this.n))
}

func (this *AmongSimple) GetVarIds() []core.VarId {
	varIds := make([]core.VarId, len(this.xi))
	for _, var_id := range this.xi {
		varIds = append(varIds, var_id)
	}
	varIds = append(varIds, this.n)
	return varIds
}

func (this *AmongSimple) GetDomains() []core.Domain {
	domains := make([]core.Domain, len(this.xi_Domains))
	for _, var_id := range this.xi {
		domains = append(domains, this.xi_Domains[var_id])
	}
	domains = append(domains, this.n_Domain)
	return domains
}

func (this *AmongSimple) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *AmongSimple) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}
