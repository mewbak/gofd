package ixrange

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
)

// Examples: X in dom(Y)-dom(Z)

type SubRange struct {
	r1 IRange //input-domain (right-side of indexical)
	r2 IRange //input-domain (right-side of indexical)
}

// CreateFromToRange creates a FromToRange with the given from- and toTerm
func CreateSubRange(r1, r2 IRange) *SubRange {
	newr := new(SubRange)
	newr.r1 = r1
	newr.r2 = r2

	return newr
}

func (this *SubRange) CheckEntail(outDom *core.IvDomain) bool {
	f, t := GetMinMaxTermsWithoutVarID(this.GetValue())
	ftR := CreateFromToRange(f, t)
	return ftR.CheckEntail(outDom)
}

func (this *SubRange) CheckDisentail(outDom *core.IvDomain) bool {
	f, t := GetMinMaxTermsWithoutVarID(this.GetValue())
	ftR := CreateFromToRange(f, t)
	return ftR.CheckDisentail(outDom)
}

// Process collectChanges, if it has the changed variable/domain (where some
// values has been removed) as input-variable. Have to return slice, because
// of usage of append
// dom: output-variable/domain
func (this *SubRange) Process(dom *core.IvDomain) []*core.IvDomPart {
	rparts := dom.GetParts()
	beginRparts := rparts[0]
	endRparts := rparts[len(rparts)-1]
	translatedParts := make([]*core.IvDomPart, 0)
	for _, resultInP := range this.r1.GetValue().GetParts() {
		for _, firstInP := range this.r2.GetValue().GetParts() {
			tPart := resultInP.SUBTRACT(firstInP)
			if tPart.LT_DP(beginRparts) || tPart.GT_DP(endRparts) {
				continue
			}
			translatedParts = append(translatedParts, tPart)
		}
	}
	var removingD *core.IvDomain
	if len(translatedParts) != 0 {
		unionD := core.CreateIvDomainUnion(translatedParts)
		removingD = dom.DifferenceWithIvDomain(unionD)
		if removingD.IsEmpty() {
			return nil
		}
	} else {
		removingD = dom.Copy().(*core.IvDomain)
	}
	return removingD.GetParts()
}

// HasVarAsInput returns, if the specific Range has a specific variable as
// input-variable (right side of expression)
func (this *SubRange) HasVarAsInput(varid core.VarId) bool {
	return this.r1.HasVarAsInput(varid) || this.r2.HasVarAsInput(varid)
}

func (this *SubRange) GetValue() *core.IvDomain {
	return this.r1.GetValue().SUBTRACT(this.r2.GetValue())
}

func (this *SubRange) Evaluable() ixterm.EvalState {
	if this.r1.Evaluable() == ixterm.EMPTY ||
		this.r2.Evaluable() == ixterm.EMPTY {
		return ixterm.EMPTY
	}
	if this.r1.Evaluable() == ixterm.NOT_EVALUABLE_YET ||
		this.r2.Evaluable() == ixterm.NOT_EVALUABLE_YET {
		return ixterm.NOT_EVALUABLE_YET
	}
	return ixterm.EVALUABLE
}

func (this *SubRange) String() string {
	return fmt.Sprintf("%s -R %s)", this.r1, this.r2)
}
