package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical"
	"bitbucket.org/gofd/gofd/indexical/ixrange"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
)

// XplusCeqY_IC represents the propagator for the constraint X + C = Y
// example: X + 3 = Y
type XplusCeqY struct {
	x, y               core.VarId
	c                  int
	x_Domain, y_Domain *core.IvDomain
	iColl              *indexical.IndexicalCollection
	outCh              chan<- *core.ChangeEvent
	inCh               <-chan *core.ChangeEntry
	id                 core.PropId
	store              *core.Store
}

func (this *XplusCeqY) GetIndexicalCollection() *indexical.IndexicalCollection {
	return this.iColl
}

func (this *XplusCeqY) Clone() core.Constraint {
	prop := new(XplusCeqY)
	prop.x = this.x
	prop.y = this.y
	prop.c = this.c
	return prop
}

func (this *XplusCeqY) Start(store *core.Store) {
	indexical.InitProcessConstraint(this, false)
	indexical.ProcessConstraint(this, false)
}

// Register registers the propagator at the store.
func (this *XplusCeqY) Register(store *core.Store) {
	var domains map[core.VarId]core.Domain
	this.inCh, domains, this.outCh =
		store.RegisterPropagatorMap([]core.VarId{this.x, this.y}, this.id)

	varidToDomainMap := core.GetVaridToIntervalDomains(domains)

	this.x_Domain = varidToDomainMap[this.x]
	this.y_Domain = varidToDomainMap[this.y]
	this.store = store

	this.iColl = indexical.CreateIndexicalCollection()

	this.iColl.AddIndexicalsAtPrio(this.MakeXplusCeqYBoundsIndexicals(), indexical.HIGHEST)
	this.iColl.AddIndexicalsAtPrio(this.MakeXplusCeqYArcIndexicals(), indexical.HIGH)
}

func (this *XplusCeqY) MakeXplusCeqYBoundsIndexicals() []*indexical.Indexical {

	indexicals := make([]*indexical.Indexical, 2)

	//X in (min(Y)-1)..(max(Y)-1)
	minYT := ixterm.CreateMinTerm(this.y, this.y_Domain)
	valueT := ixterm.CreateValueTerm(this.c)
	subT1 := ixterm.CreateSubtractionTerm(minYT, valueT)

	maxYT := ixterm.CreateMaxTerm(this.y, this.y_Domain)
	subT2 := ixterm.CreateSubtractionTerm(maxYT, valueT)

	ftR := ixrange.CreateFromToRange(subT1, subT2)
	indexicals[0] = indexical.CreateIndexical(this.x, this.x_Domain, ftR)

	//Y in (min(X)+1)..(max(X)+1)
	minXT := ixterm.CreateMinTerm(this.x, this.x_Domain)
	addT1 := ixterm.CreateAdditionTerm(minXT, valueT)

	maxXT := ixterm.CreateMaxTerm(this.x, this.x_Domain)
	addT2 := ixterm.CreateAdditionTerm(maxXT, valueT)

	ftR = ixrange.CreateFromToRange(addT1, addT2)
	indexicals[1] = indexical.CreateIndexical(this.y, this.y_Domain, ftR)

	return indexicals
}

func (this *XplusCeqY) MakeXplusCeqYArcIndexicals() []*indexical.Indexical {

	indexicals := make([]*indexical.Indexical, 2)

	//X in dom(Y)-C
	yDomR := ixrange.CreateDomRange(this.y, this.y_Domain)
	cTerm := ixterm.CreateValueTerm(this.c)
	cSingleR := ixrange.CreateSingleValueRange(cTerm)
	subR := ixrange.CreateSubRange(yDomR, cSingleR)
	indexicals[0] = indexical.CreateIndexical(this.x, this.x_Domain, subR)

	//Y in dom(X)+C
	xDomR := ixrange.CreateDomRange(this.x, this.x_Domain)
	cTerm = ixterm.CreateValueTerm(this.c)
	cSingleR = ixrange.CreateSingleValueRange(cTerm)
	addR := ixrange.CreateAddRange(xDomR, cSingleR)
	indexicals[1] = indexical.CreateIndexical(this.y, this.y_Domain, addR)

	return indexicals
}

func (this *XplusCeqY) SetID(propID core.PropId) {
	this.id = propID
}

func (this *XplusCeqY) GetID() core.PropId {
	return this.id
}

func (this *XplusCeqY) GetVarIds() []core.VarId {
	return []core.VarId{this.x, this.y}
}

func (this *XplusCeqY) GetDomains() []core.Domain {
	return []core.Domain{this.x_Domain, this.y_Domain}
}

func (this *XplusCeqY) GetInCh() <-chan *core.ChangeEntry {
	return this.inCh
}

func (this *XplusCeqY) GetOutCh() chan<- *core.ChangeEvent {
	return this.outCh
}

func (this *XplusCeqY) String() string {
	return fmt.Sprintf("PROP_%d %s+%d = %s",
		this.id, this.store.GetName(this.x), this.c,
		this.store.GetName(this.y))
}

//X+C=Z
func CreateXplusCeqY(x core.VarId, c int, y core.VarId) *XplusCeqY {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXplusCeqY_propagator")
	}
	prop := new(XplusCeqY)
	prop.x = x
	prop.y = y
	prop.c = c
	return prop
}

//X=Z
func CreateXeqY_IC(x core.VarId, y core.VarId) *XplusCeqY {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateXplusCeqY_propagator")
	}
	prop := new(XplusCeqY)
	prop.x = x
	prop.y = y
	prop.c = 0
	return prop
}
