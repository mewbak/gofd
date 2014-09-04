package ixrange

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
)

//examples:
// X in dom(Y)/dom(Z)

type DivRange struct {
	r1 IRange //input-domain (right-side of indexical)
	r2 IRange //input-domain (right-side of indexical)
}

// CreateFromToRange creates a FromToRange with the given from- and toTerm
func CreateDivRange(r1, r2 IRange) *DivRange {
	newr := new(DivRange)
	newr.r1 = r1
	newr.r2 = r2

	return newr
}

func (this *DivRange) CheckEntail(outDom *core.IvDomain) bool {
	f, t := GetMinMaxTermsWithoutVarID(this.GetValue())
	ftR := CreateFromToRange(f, t)
	return ftR.CheckEntail(outDom)
}

func (this *DivRange) CheckDisentail(outDom *core.IvDomain) bool {
	f, t := GetMinMaxTermsWithoutVarID(this.GetValue())
	ftR := CreateFromToRange(f, t)
	return ftR.CheckDisentail(outDom)
}

// Process collectChanges, if it has the changed variable/domain (where some
// values has been removed) as input-variable. Have to return slice, because
// of usage of append
// dom: output-variable/domain
// Z=X/C or Z=C/X or Z=X/Y or or C=X/Y
func (this *DivRange) Process(dom *core.IvDomain) []*core.IvDomPart {

	xDom := this.r1.GetValue()
	yDom := this.r2.GetValue()
	zDom := dom

	//Z=X/C
	if yDom.Size() == 1 {

		c := yDom.GetAnyElement()
		if c == 0 {
			panic("Division by zero! Not allowed")
		}
		vals := make([]int, 0)
		for _, val := range zDom.GetValues() {
			//println("val:",val,"c:",c)
			if !xDom.Contains(val * c) {
				vals = append(vals, val)
			}
		}
		d := core.CreateIvDomainFromIntArr(vals)
		return d.GetParts()
	} else if xDom.Size() == 1 {
		//C/Y=Z
		panic("DivRange-Panic:NOT IMPLEMENTED YET")
	} else if zDom.Size() == 1 {
		//X/Y=C
		panic("DivRange-Panic:NOT IMPLEMENTED YET")
	} else {
		//X/Y=Z
		panic("DivRange-Panic:NOT IMPLEMENTED YET")
	}
	return nil
}

// HasVarAsInput returns, if the specific Range has a specific variable as
// input-variable (right side of expression)
func (this *DivRange) HasVarAsInput(varid core.VarId) bool {
	return this.r1.HasVarAsInput(varid) || this.r2.HasVarAsInput(varid)
}

func (this *DivRange) GetValue() *core.IvDomain {
	return this.r1.GetValue().DIvIDE(this.r2.GetValue())
	//panic("DivRange-Panic: NOT IMPLEMENTED YET")
}

func (this *DivRange) Evaluable() ixterm.EvalState {
	if this.r1.Evaluable() == ixterm.EMPTY || this.r2.Evaluable() == ixterm.EMPTY {
		return ixterm.EMPTY
	}

	if this.r1.Evaluable() == ixterm.NOT_EVALUABLE_YET || this.r2.Evaluable() == ixterm.NOT_EVALUABLE_YET {
		return ixterm.NOT_EVALUABLE_YET
	}

	return ixterm.EVALUABLE
}

func (this *DivRange) String() string {
	return fmt.Sprintf("%s /R %s)", this.r1, this.r2)
}
