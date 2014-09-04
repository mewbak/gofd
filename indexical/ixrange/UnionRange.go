package ixrange

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
	"strings"
)

//examples:
// X in dom(Y),dom(Z),5,1..20

type UnionRange struct {
	ranges []IRange //ranges
}

// CreateFromToRange creates a FromToRange with the given from- and toTerm
func CreateUnionRange(ranges ...IRange) *UnionRange {
	newr := new(UnionRange)
	newr.ranges = ranges

	return newr
}

func (this *UnionRange) CheckEntail(outDom *core.IvDomain) bool {
	f, t := GetMinMaxTermsWithoutVarID(this.GetValue())
	ftR := CreateFromToRange(f, t)
	return ftR.CheckEntail(outDom)
}

func (this *UnionRange) CheckDisentail(outDom *core.IvDomain) bool {
	f, t := GetMinMaxTermsWithoutVarID(this.GetValue())
	ftR := CreateFromToRange(f, t)
	return ftR.CheckDisentail(outDom)
}

// Process collectChanges, if it has the changed variable/domain (where some
// values has been removed) as input-variable. Have to return slice, because
// of usage of append
// dom: output-variable/domain
func (this *UnionRange) Process(dom *core.IvDomain) []*core.IvDomPart {
	return dom.DifferenceWithIvDomain(this.GetValue()).GetParts()
}

// HasVarAsInput returns, if the specific Range has a specific variable as
// input-variable (right side of expression)
func (this *UnionRange) HasVarAsInput(varid core.VarId) bool {
	for _, r := range this.ranges {
		if r.HasVarAsInput(varid) {
			return true
		}
	}
	return false
}

func (this *UnionRange) GetValue() *core.IvDomain {
	parts := make([]*core.IvDomPart, 0)

	for _, r := range this.ranges {
		parts = append(parts, r.GetValue().GetParts()...)
	}

	return core.CreateIvDomainUnion(parts)
}

func (this *UnionRange) Evaluable() ixterm.EvalState {
	for _, r := range this.ranges {
		rstate := r.Evaluable()
		if rstate == ixterm.EMPTY {
			return ixterm.EMPTY
		} else if rstate == ixterm.NOT_EVALUABLE_YET {
			return ixterm.NOT_EVALUABLE_YET
		}
	}

	return ixterm.EVALUABLE
}

func (this *UnionRange) String() string {

	s := make([]string, len(this.ranges))
	for i, r := range this.ranges {
		s[i] = r.String()
	}
	return fmt.Sprintf("%s)", strings.Join(s, " v "))
}
