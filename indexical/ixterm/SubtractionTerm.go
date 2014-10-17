package ixterm

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// SubtractionTerm t-t examples
// min(X)-min(Y)
// max(X)-min(Y)
type SubtractionTerm struct {
	term1 ITerm
	term2 ITerm
}

// GetValue returns the (current) value of the term (can change, cause of
// underlying domains). Returns result of term1-term2
func (this *SubtractionTerm) GetValue() *core.IvDomain {
	v := this.term1.GetValue().GetAnyElement() - this.term2.GetValue().GetAnyElement()
	return core.CreateIvDomainFromTo(v, v)
}

func (this *SubtractionTerm) String() string {
	return fmt.Sprintf("%s -T %s", this.term1, this.term2)
}

func (this *SubtractionTerm) HasVarId(varid core.VarId) bool {
	return this.term1.HasVarId(varid) || this.term2.HasVarId(varid)
}

func (this *SubtractionTerm) Evaluable() EvalState {
	if this.term1.Evaluable() == EMPTY || this.term2.Evaluable() == EMPTY {
		return EMPTY
	}

	if this.term1.Evaluable() == NOT_EVALUABLE_YET || this.term2.Evaluable() == NOT_EVALUABLE_YET {
		return NOT_EVALUABLE_YET
	}

	return EVALUABLE
}

func (this *SubtractionTerm) GetInf() ITerm {
	//inf-sup
	return CreateSubtractionTerm(this.term1.GetInf(), this.term2.GetSup())
}

func (this *SubtractionTerm) GetSup() ITerm {
	//sup-inf
	return CreateSubtractionTerm(this.term1.GetSup(), this.term2.GetInf())
}

// CreateDomTerm returns a new DomTerm
func CreateSubtractionTerm(term1 ITerm, term2 ITerm) *SubtractionTerm {
	t := new(SubtractionTerm)

	t.term1 = term1
	t.term2 = term2

	return t
}
