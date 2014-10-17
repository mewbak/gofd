// package indexical provides propagators implemented using a
// higher level abstraction as provided in the indexical package.
package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical"
	"bitbucket.org/gofd/gofd/indexical/ixrange"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
	"strings"
)

// Alldifferent is a proper constraint that is semantically equivalent to
// CreateAlldifferent (quadratically many "not equal" constraints), but
// just holds one copy of the involved variables and removes all values
// of ground variables in all other variables.
// No stronger propagation techniques of a global constraint are used.
type Alldifferent struct {
	vars             []core.VarId
	outCh            chan<- *core.ChangeEvent
	inCh             <-chan *core.ChangeEntry
	varidToDomainMap map[core.VarId]*core.IvDomain
	id               core.PropId
	iColl            *indexical.IndexicalCollection
	store            *core.Store
}

func (this *Alldifferent) GetIndexicalCollection() *indexical.IndexicalCollection {
	return this.iColl
}

func (this *Alldifferent) Start(store *core.Store) {

	indexical.InitProcessConstraint(this, false)
	indexical.ProcessConstraint(this, false)
}

// Register registers the propagator at the store. Here, the propagator gets
// his needed channels and domains and stores them in his struct
func (this *Alldifferent) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain

	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap(this.vars, this.id)

	this.varidToDomainMap = core.GetVaridToIntervalDomains(domains)

	this.store = store

	this.iColl = indexical.CreateIndexicalCollection()

	for index, i := range this.vars {
		//make X!=Y
		for _, j := range this.vars[index+1:] {
			xId := i
			yId := j
			xDom := this.varidToDomainMap[i]
			yDom := this.varidToDomainMap[j]

			//make X->Y
			valY := ixterm.CreateValTerm(yId, yDom)
			rY := ixrange.CreateSingleValueRange(valY)
			nrY := ixrange.CreateNotRange(rY)
			this.iColl.CreateAndAddIndexical(xId, xDom, nrY, indexical.HIGHEST)

			//make X<-Y
			valX := ixterm.CreateValTerm(xId, xDom)
			rX := ixrange.CreateSingleValueRange(valX)
			nrX := ixrange.CreateNotRange(rX)
			this.iColl.CreateAndAddIndexical(yId, yDom, nrX, indexical.HIGHEST)
		}
	}
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
		vars_str[i] = this.store.GetName(var_id)
	}
	return fmt.Sprintf("PROP_%d %s",
		this.id, strings.Join(vars_str, "!="))
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

// CreateAlldifferent2 creates one single propagator, that for each variable
// that becomes ground removes that value from all other variables.
// Note: Alldifferent is not using stronger propagation techniques of
// a global constraint.
func CreateAlldifferent(vars ...core.VarId) *Alldifferent {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateAlldifferent_IC-propagator")
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
