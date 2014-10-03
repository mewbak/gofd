package ixrange

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
)

// examples: X in abs(dom(Y))

type AbsRange struct {
	r1 IRange
}

func CreateAbsRange(r1 IRange) *AbsRange {
	newr := new(AbsRange)
	newr.r1 = r1
	return newr
}

func (this *AbsRange) CheckEntail(outDom *core.IvDomain) bool {
	f, t := GetMinMaxTermsWithoutVarID(this.GetValue())
	ftR := CreateFromToRange(f, t)
	return ftR.CheckEntail(outDom)
}

func (this *AbsRange) CheckDisentail(outDom *core.IvDomain) bool {
	f, t := GetMinMaxTermsWithoutVarID(this.GetValue())
	ftR := CreateFromToRange(f, t)
	return ftR.CheckDisentail(outDom)
}

// Process collectChanges, if it has the changed variable/domain (where some
// values has been removed) as input-variable.
// dom: output-variable/domain
func (this *AbsRange) Process(dom *core.IvDomain) []*core.IvDomPart {
	removingDom := dom.DifferenceWithIvDomain(this.GetValue())
	return removingDom.GetParts()
}

// HasVarAsInput returns, if the specific Range has a specific variable as
// input-variable (right side of expression)
func (this *AbsRange) HasVarAsInput(varid core.VarId) bool {
	return this.r1.HasVarAsInput(varid)
}

func (this *AbsRange) GetValue() *core.IvDomain {
	return this.r1.GetValue().ABS()
}

func (this *AbsRange) Evaluable() ixterm.EvalState {
	return this.r1.Evaluable()
}

func (this *AbsRange) String() string {
	return fmt.Sprintf("abs(%s))", this.r1)
}
