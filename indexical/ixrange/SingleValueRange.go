package ixrange

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
)

//examples:
// X in val(Y)
// X in 5
// X in min(Y)

//neg
// X in \val(Y)	: if Y ground, remove Y-value from X
// X in \5		: remove 5 from X
// X in \min(Y)	: remove  min(Y) from X

type SingleValueRange struct {
	valTerm ixterm.ITerm //1, min(X), val(X)
}

// CreateFromToRange creates a FromToRange with the given from- and toTerm
func CreateSingleValueRange(valTerm ixterm.ITerm) *SingleValueRange {
	r := new(SingleValueRange)
	r.valTerm = valTerm
	return r
}

func (this *SingleValueRange) CheckEntail(outDom *core.IvDomain) bool {
	f, t := GetMinMaxTermsWithoutVarID(this.GetValue())
	ftR := CreateFromToRange(f, t)
	return ftR.CheckEntail(outDom)
}

func (this *SingleValueRange) CheckDisentail(outDom *core.IvDomain) bool {
	f, t := GetMinMaxTermsWithoutVarID(this.GetValue())
	ftR := CreateFromToRange(f, t)
	return ftR.CheckDisentail(outDom)
}

func (this *SingleValueRange) process(dom *core.IvDomain, value int, removingParts []*core.IvDomPart) []*core.IvDomPart {
	for _, part := range dom.GetParts() {
		if part.ContainsInt(value) {
			if part.ContainsInt(value - 1) {
				p1 := core.CreateIvDomPart(part.From, value-1)
				removingParts = append(removingParts, p1)
			}
			if part.ContainsInt(value + 1) {
				p2 := core.CreateIvDomPart(value+1, part.To)
				removingParts = append(removingParts, p2)
			}
		} else {
			removingParts = append(removingParts, part.Copy())
		}
	}

	return removingParts
}

// Process collectChanges
func (this *SingleValueRange) Process(dom *core.IvDomain) []*core.IvDomPart {
	removingParts := make([]*core.IvDomPart, 0)

	value := this.valTerm.GetValue().GetMax()

	return this.process(dom, value, removingParts)
}

// HasVarAsInput returns, if the specific Range has a specific variable as
// input-variable (right side of expression)
func (this *SingleValueRange) HasVarAsInput(varid core.VarId) bool {
	return this.valTerm.HasVarId(varid)
}

func (this *SingleValueRange) String() string {
	return fmt.Sprintf("%s", this.valTerm)
}

func (this *SingleValueRange) GetValue() *core.IvDomain {
	return this.valTerm.GetValue()
}

func (this *SingleValueRange) Evaluable() ixterm.EvalState {
	return this.valTerm.Evaluable()
}
