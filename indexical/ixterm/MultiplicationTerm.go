package ixterm

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// MultiplicationTerm t+t examples
// min(X)*min(Y)
// min(X)*3
type MultiplicationTerm struct {
	term1 ITerm
	term2 ITerm
}

// GetValue returns the (current) value of the term (can change, cause of
// underlying domains). Returns result of term1*term2
func (this *MultiplicationTerm) GetValue() *core.IvDomain {
	v := this.term1.GetValue().GetAnyElement() * this.term2.GetValue().GetAnyElement()
	return core.CreateIvDomainFromTo(v, v)
}

func (this *MultiplicationTerm) String() string {
	return fmt.Sprintf("%s *T %s)", this.term1, this.term2)
}

func (this *MultiplicationTerm) HasVarId(varid core.VarId) bool {
	return this.term1.HasVarId(varid) || this.term2.HasVarId(varid)
}

func (this *MultiplicationTerm) Evaluable() EvalState {
	if this.term1.Evaluable() == EMPTY || this.term2.Evaluable() == EMPTY {
		return EMPTY
	}

	if this.term1.Evaluable() == NOT_EVALUABLE_YET || this.term2.Evaluable() == NOT_EVALUABLE_YET {
		return NOT_EVALUABLE_YET
	}

	return EVALUABLE
}

func (this *MultiplicationTerm) GetInf() ITerm {
	panic("not implemented yet")
	return nil
}

func (this *MultiplicationTerm) GetSup() ITerm {
	panic("not implemented yet")
	return nil
}

// CreateDomTerm returns a new DomTerm
func CreateMultiplicationTerm(term1 ITerm, term2 ITerm) *MultiplicationTerm {
	t := new(MultiplicationTerm)

	t.term1 = term1
	t.term2 = term2

	return t
}
