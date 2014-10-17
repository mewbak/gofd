package ixterm

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// DivisionTerm t/t examples
// min(X)/min(Y)
// min(X)/3
type DivisionTerm struct {
	term1   ITerm
	term2   ITerm
	roundUp bool //true: ceiling. false: floor
}

// GetValue returns the (current) value of the term (can change, cause of
// underlying domains). Returns result of term1/term2
func (this *DivisionTerm) GetValue() *core.IvDomain {

	t1v := this.term1.GetValue().GetAnyElement()
	t2v := this.term2.GetValue().GetAnyElement()
	var v int
	if t1v%t2v != 0 {
		if this.roundUp {
			v = (t1v / t2v) + 1
			return core.CreateIvDomainFromTo(v, v)
		}
	}

	v = t1v / t2v
	return core.CreateIvDomainFromTo(v, v)
}

func (this *DivisionTerm) String() string {
	return fmt.Sprintf("%s /T %s", this.term1, this.term2)
}

func (this *DivisionTerm) HasVarId(varid core.VarId) bool {
	return this.term1.HasVarId(varid) || this.term2.HasVarId(varid)
}

func (this *DivisionTerm) Evaluable() EvalState {
	if this.term1.Evaluable() == EMPTY || this.term2.Evaluable() == EMPTY {
		return EMPTY
	}

	if this.term1.Evaluable() == NOT_EVALUABLE_YET || this.term2.Evaluable() == NOT_EVALUABLE_YET {
		return NOT_EVALUABLE_YET
	}

	return EVALUABLE
}

func (this *DivisionTerm) GetInf() ITerm {
	panic("not implemented yet")
	return nil
}

func (this *DivisionTerm) GetSup() ITerm {
	panic("not implemented yet")
	return nil
}

// CreateDomTerm returns a new DomTerm
func CreateDivisionTerm(term1 ITerm, term2 ITerm, roundUp bool) *DivisionTerm {
	t := new(DivisionTerm)

	t.term1 = term1
	t.term2 = term2
	t.roundUp = roundUp

	return t
}
