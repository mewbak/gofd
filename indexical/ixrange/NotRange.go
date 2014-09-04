package ixrange

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
)

//ToDo: Tests

// examples:
// X in not(dom(Y))	--> X not in dom(Y)
// 1,2,3 not in 3,4,5 --> 1,2

type NotRange struct {
	r1 IRange
}

func CreateNotRange(r1 IRange) *NotRange {
	newr := new(NotRange)
	newr.r1 = r1

	return newr
}

func (this *NotRange) CheckEntail(outDom *core.IvDomain) bool {
	// important: r instead of not not(r)!
	f, t := GetMinMaxTermsWithoutVarID(this.r1.GetValue())

	ftR := CreateFromToRange(f, t)
	return ftR.CheckDisentail(outDom)
}

func (this *NotRange) CheckDisentail(outDom *core.IvDomain) bool {
	f, t := GetMinMaxTermsWithoutVarID(this.GetValue())
	ftR := CreateFromToRange(f, t)
	return ftR.CheckDisentail(outDom)
}

// Process collectChanges, if it has the changed variable/domain (where some
// values has been removed) as input-variable.
// dom: output-variable/domain
func (this *NotRange) Process(dom *core.IvDomain) []*core.IvDomPart {

	removingDom := dom.IntersectionIvDomain(this.r1.GetValue())
	return removingDom.GetParts()
}

// HasVarAsInput returns, if the specific Range has a specific variable as
// input-variable (right side of expression)
func (this *NotRange) HasVarAsInput(varid core.VarId) bool {
	return this.r1.HasVarAsInput(varid)
}

//[0,10] not in [1,20] --> [0,10] in  [[-inf, 0] v [21,inf]]
func (this *NotRange) GetValue() *core.IvDomain {
	return this.r1.GetValue().NOT()
}

func (this *NotRange) Evaluable() ixterm.EvalState {
	return this.r1.Evaluable()
}

func (this *NotRange) String() string {
	return fmt.Sprintf("not(%s))", this.r1)
}
