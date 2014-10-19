package ixrange

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
)

//examples:
// 1 .. 5
// 1 .. min(X)
// 1 .. max(X)
// -infinity .. max(X)
// min(X) .. max(X)

type FromToRange struct {
	fromTerm ixterm.ITerm //1, min(X)
	toTerm   ixterm.ITerm //10, max(Y)
}

// CreateFromToRange creates a FromToRange with the given from- and toTerm
func CreateFromToRange(fromTerm ixterm.ITerm, toTerm ixterm.ITerm) *FromToRange {
	r := new(FromToRange)

	if fromTerm.Evaluable() == ixterm.EVALUABLE && toTerm.Evaluable() == ixterm.EVALUABLE {
		if fromTerm.GetInf().GetValue().GetAnyElement() > toTerm.GetSup().GetValue().GetAnyElement() {
			panic("FromToRange-Creation-Fail: fromTerm " + fromTerm.GetInf().String() + " must be leq than toTerm " + toTerm.GetSup().String())
		}
	}

	r.fromTerm = fromTerm
	r.toTerm = toTerm

	return r
}

// CreateFromToRangeInts creates a FromToRange with the given from and to
// int-value
func CreateFromToRangeInts(from int, to int) *FromToRange {
	r := new(FromToRange)

	r.fromTerm = ixterm.CreateValueTerm(from)
	r.toTerm = ixterm.CreateValueTerm(to)

	return r
}

func (this *FromToRange) CheckEntail(outDom *core.IvDomain) bool {

	fValue := this.fromTerm.GetValue().GetAnyElement()
	tValue := this.toTerm.GetValue().GetAnyElement()

	if fValue == ixterm.NEG_INFINITY && tValue == ixterm.INFINITY {
		//-inf..inf
		return true
	} else if tValue == ixterm.INFINITY {
		//t..inf
		return outDom.GetMin() >= this.fromTerm.GetSup().GetValue().GetAnyElement()
	} else if fValue == ixterm.NEG_INFINITY {
		//-inf..t
		return outDom.GetMax() <= this.toTerm.GetInf().GetValue().GetAnyElement()
	}
	// else, t1..t2
	return outDom.GetMin() >= this.fromTerm.GetSup().GetValue().GetAnyElement() && outDom.GetMax() <= this.toTerm.GetInf().GetValue().GetAnyElement()
}

func (this *FromToRange) CheckDisentail(outDom *core.IvDomain) bool {

	fValue := this.fromTerm.GetValue().GetAnyElement()
	tValue := this.toTerm.GetValue().GetAnyElement()

	if fValue == ixterm.NEG_INFINITY && tValue == ixterm.INFINITY {
		//-inf..inf
		return false
	} else if tValue == ixterm.INFINITY {
		//t..inf
		return outDom.GetMax() < this.fromTerm.GetInf().GetValue().GetAnyElement()
	} else if fValue == ixterm.NEG_INFINITY {
		//-inf..t
		return outDom.GetMin() > this.toTerm.GetSup().GetValue().GetAnyElement()
	}
	// else, t1..t2
	return outDom.GetMax() < this.fromTerm.GetInf().GetValue().GetAnyElement() || outDom.GetMin() > this.toTerm.GetSup().GetValue().GetAnyElement() || this.fromTerm.GetInf().GetValue().GetAnyElement() > this.toTerm.GetSup().GetValue().GetAnyElement()
}

// Process collectChanges, if it has the changed variable/domain (where some
// values has been removed) as input-variable. Have to return slice, because
// of usage of append
func (this *FromToRange) Process(dom *core.IvDomain) []*core.IvDomPart {
	removingParts := make([]*core.IvDomPart, 0)

	fromVal := this.fromTerm.GetValue().GetAnyElement()
	toVal := this.toTerm.GetValue().GetAnyElement()

	//completely contains
	if toVal == ixterm.INFINITY {
		if fromVal <= dom.GetMin() {
			return removingParts
		} else {
			return []*core.IvDomPart{core.CreateIvDomPart(ixterm.NEG_INFINITY, fromVal-1)}
		}
	}

	//completely contains
	if fromVal == ixterm.NEG_INFINITY {
		if dom.GetMax() <= toVal {
			return removingParts
		} else {
			return []*core.IvDomPart{core.CreateIvDomPart(toVal+1, ixterm.INFINITY)}
		}
	}

	i1 := core.CreateIvDomPart(ixterm.NEG_INFINITY, fromVal-1)
	i2 := core.CreateIvDomPart(toVal+1, ixterm.INFINITY)

	return []*core.IvDomPart{i1, i2}
}

// HasVarAsInput returns, if the specific Range has a specific variable as
// input-variable (right side of expression)
func (this *FromToRange) HasVarAsInput(varid core.VarId) bool {
	if this.fromTerm.HasVarId(varid) || this.toTerm.HasVarId(varid) {
		return true
	}

	return false
}

func (this *FromToRange) String() string {
	return fmt.Sprintf("%s..%s", this.fromTerm, this.toTerm)
}

func (this *FromToRange) GetValue() *core.IvDomain {
	from := this.fromTerm.GetValue().GetAnyElement()
	to := this.toTerm.GetValue().GetAnyElement()
	return core.CreateIvDomainFromTo(from, to)
}

func (this *FromToRange) Evaluable() ixterm.EvalState {
	if this.fromTerm.Evaluable() == ixterm.EMPTY || this.toTerm.Evaluable() == ixterm.EMPTY {
		return ixterm.EMPTY
	}

	if this.fromTerm.Evaluable() == ixterm.NOT_EVALUABLE_YET || this.toTerm.Evaluable() == ixterm.NOT_EVALUABLE_YET {
		return ixterm.NOT_EVALUABLE_YET
	}

	return ixterm.EVALUABLE
}
