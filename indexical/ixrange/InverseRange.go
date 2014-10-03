package ixrange

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
)

//examples:
// X in inverse(dom(Y))
//dom(Y):1,2,3 --> inverse(dom(Y)) = -3,-2,-1

type InverseRange struct {
	r1 IRange
}

func CreateInverseRange(r1 IRange) *InverseRange {
	newr := new(InverseRange)
	newr.r1 = r1

	return newr
}

func (this *InverseRange) CheckEntail(outDom *core.IvDomain) bool {
	f, t := GetMinMaxTermsWithoutVarID(this.GetValue())
	ftR := CreateFromToRange(f, t)
	return ftR.CheckEntail(outDom)
}

func (this *InverseRange) CheckDisentail(outDom *core.IvDomain) bool {
	f, t := GetMinMaxTermsWithoutVarID(this.GetValue())
	ftR := CreateFromToRange(f, t)
	return ftR.CheckDisentail(outDom)
}

// Process collectChanges, if it has the changed variable/domain (where some
// values has been removed) as input-variable.
// dom: output-variable/domain
func (this *InverseRange) Process(dom *core.IvDomain) []*core.IvDomPart {

	removingDom := dom.DifferenceWithIvDomain(this.GetValue())
	return removingDom.GetParts()
}

// HasVarAsInput returns, if the specific Range has a specific variable as
// input-variable (right side of expression)
func (this *InverseRange) HasVarAsInput(varid core.VarId) bool {
	return this.r1.HasVarAsInput(varid)
}

func (this *InverseRange) GetValue() *core.IvDomain {
	return this.r1.GetValue().NEGATE()
}

func (this *InverseRange) Evaluable() ixterm.EvalState {
	return this.r1.Evaluable()
}

func (this *InverseRange) String() string {
	return fmt.Sprintf("negate(%s))", this.r1)
}
