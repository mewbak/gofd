package ixrange

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
)

// Examples: X in dom(Y)*dom(Z)

type MultRange struct {
	r1 IRange //input-domain (right-side of indexical)
	r2 IRange //input-domain (right-side of indexical)
}

// CreateFromToRange creates a FromToRange with the given from- and toTerm
func CreateMultRange(r1, r2 IRange) *MultRange {
	newr := new(MultRange)
	newr.r1 = r1
	newr.r2 = r2

	return newr
}

func (this *MultRange) CheckEntail(outDom *core.IvDomain) bool {
	f, t := GetMinMaxTermsWithoutVarID(this.GetValue())
	ftR := CreateFromToRange(f, t)
	return ftR.CheckEntail(outDom)
}

func (this *MultRange) CheckDisentail(outDom *core.IvDomain) bool {
	f, t := GetMinMaxTermsWithoutVarID(this.GetValue())
	ftR := CreateFromToRange(f, t)
	return ftR.CheckDisentail(outDom)
}

// Process collectChanges, if it has the changed variable/domain (where some
// values has been removed) as input-variable. Have to return slice, because
// of usage of append
// dom: output-variable/domain
// X*C=Z or C*Y=Z or X*Y=Z or X*Y=C
func (this *MultRange) Process(dom *core.IvDomain) []*core.IvDomPart {
	xDom := this.r1.GetValue()
	yDom := this.r2.GetValue()
	zDom := dom
	// X*C=Z
	if yDom.Size() == 1 {
		c := yDom.GetAnyElement()
		vals := make([]int, 0)
		if c == 0 {
			panic("Division by zero! Not allowed")
		}
		for _, z_val := range zDom.GetValues() {
			if !(z_val%c == 0) {
				vals = append(vals, z_val)
			} else {
				if !xDom.Contains(z_val / c) {
					vals = append(vals, z_val)
				}
			}
		}
		d := core.CreateIvDomainFromIntArr(vals)
		return d.GetParts()
	} else if xDom.Size() == 1 {
		// C*Y=Z
		panic("DivRange-Panic:NOT IMPLEMENTED YET")
	} else if zDom.Size() == 1 {
		// X*Y=C
		panic("DivRange-Panic:NOT IMPLEMENTED YET")
	} else {
		// X*Y=Z
		panic("DivRange-Panic:NOT IMPLEMENTED YET")
	}
	return nil
}

// HasVarAsInput returns, if the specific Range has a specific variable as
// input-variable (right side of expression)
func (this *MultRange) HasVarAsInput(varid core.VarId) bool {
	return this.r1.HasVarAsInput(varid) || this.r2.HasVarAsInput(varid)
}

func (this *MultRange) GetValue() *core.IvDomain {
	return this.r1.GetValue().MULTIPLY(this.r2.GetValue())
}

func (this *MultRange) Evaluable() ixterm.EvalState {
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

func (this *MultRange) String() string {
	return fmt.Sprintf("%s *R %s)", this.r1, this.r2)
}
