package ixrange

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
)

//examples:
// X in dom(Y)

type DomRange struct {
	varId core.VarId
	dom   *core.IvDomain
}

// CreateFromToRange creates a FromToRange with the given from- and toTerm
func CreateDomRange(varid core.VarId, dom *core.IvDomain) *DomRange {
	r := new(DomRange)
	r.dom = dom
	r.varId = varid

	return r
}

func (this *DomRange) CheckEntail(outDom *core.IvDomain) bool {
	f, t := GetMinMaxTermsWithoutVarID(this.GetValue())
	ftR := CreateFromToRange(f, t)
	return ftR.CheckEntail(outDom)
}

func (this *DomRange) CheckDisentail(outDom *core.IvDomain) bool {
	f, t := GetMinMaxTermsWithoutVarID(this.GetValue())
	ftR := CreateFromToRange(f, t)
	return ftR.CheckDisentail(outDom)
}

// Process collectChanges, if it has the changed variable/domain (where some
// values has been removed) as input-variable. Have to return slice, because
// of usage of append
// dom: output-variable/domain
func (this *DomRange) Process(dom *core.IvDomain) []*core.IvDomPart {

	removingDom := dom.DifferenceWithIvDomain(this.dom)
	return removingDom.GetParts()
}

// HasVarAsInput returns, if the specific Range has a specific variable as
// input-variable (right side of expression)
func (this *DomRange) HasVarAsInput(varid core.VarId) bool {
	return this.varId == varid
}

func (this *DomRange) String() string {
	return fmt.Sprintf("DomRange(%s)", this.dom)
}

func (this *DomRange) Evaluable() ixterm.EvalState {
	if this.dom.IsEmpty() {
		return ixterm.EMPTY
	}

	return ixterm.EVALUABLE
}

func (this *DomRange) GetValue() *core.IvDomain {
	return this.dom
}
