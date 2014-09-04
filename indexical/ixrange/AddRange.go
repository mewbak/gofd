package ixrange

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
)

// Examples: X in dom(Y)+dom(Z)

type AddRange struct {
	r1 IRange //input-domain (right-side of indexical)
	r2 IRange //input-domain (right-side of indexical)
}

// CreateFromToRange creates a FromToRange with the given from- and toTerm
func CreateAddRange(r1, r2 IRange) *AddRange {
	newr := new(AddRange)
	newr.r1 = r1
	newr.r2 = r2

	return newr
}

func (this *AddRange) CheckEntail(outDom *core.IvDomain) bool {
	f, t := GetMinMaxTermsWithoutVarID(this.GetValue())
	ftR := CreateFromToRange(f, t)
	return ftR.CheckEntail(outDom)
}

func (this *AddRange) CheckDisentail(outDom *core.IvDomain) bool {
	f, t := GetMinMaxTermsWithoutVarID(this.GetValue())
	ftR := CreateFromToRange(f, t)
	return ftR.CheckDisentail(outDom)
}

// Process collectChanges, if it has the changed variable/domain (where some
// values has been removed) as input-variable. Have to return slice, because
// of usage of append
// dom: output-variable/domain
func (this *AddRange) Process(dom *core.IvDomain) []*core.IvDomPart {
	rparts := dom.GetParts()
	beginRparts := rparts[0]
	endRparts := rparts[len(rparts)-1]
	translatedParts := make([]*core.IvDomPart, 0)
	for _, firstInP := range this.r1.GetValue().GetParts() {
		for _, secondInP := range this.r2.GetValue().GetParts() {
			tPart := firstInP.ADD(secondInP)
			// only take parts, which are relevant
			// (are not out of bounds of resultOutDomain)
			if tPart.LT_DP(beginRparts) || tPart.GT_DP(endRparts) {
				continue
			}
			translatedParts = append(translatedParts, tPart)
		}
	}
	var removingD *core.IvDomain
	if len(translatedParts) != 0 {
		// A+B = U			(U: Union)
		// removingD: outD diff V   (Difference)
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
func (this *AddRange) HasVarAsInput(varid core.VarId) bool {
	return this.r1.HasVarAsInput(varid) || this.r2.HasVarAsInput(varid)
}

func (this *AddRange) GetValue() *core.IvDomain {
	return this.r1.GetValue().ADD(this.r2.GetValue())
}

func (this *AddRange) Evaluable() ixterm.EvalState {
	if this.r1.Evaluable() == ixterm.EMPTY || this.r2.Evaluable() == ixterm.EMPTY {
		return ixterm.EMPTY
	}

	if this.r1.Evaluable() == ixterm.NOT_EVALUABLE_YET || this.r2.Evaluable() == ixterm.NOT_EVALUABLE_YET {
		return ixterm.NOT_EVALUABLE_YET
	}

	return ixterm.EVALUABLE
}

func (this *AddRange) String() string {
	return fmt.Sprintf("%s +R %s)", this.r1, this.r2)
}
