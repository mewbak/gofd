package ixterm

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// AdditionTerm t+t examples
// min(X)+min(Y)
// min(X)+max(Y)

type AdditionTerm struct {
	term1 ITerm
	term2 ITerm
}

// GetValue returns the (current) value of the term (can change, cause of
// underlying domains). Returns result of term1+term2
func (this *AdditionTerm) GetValue() *core.IvDomain {
	v := this.term1.GetValue().GetAnyElement() + this.term2.GetValue().GetAnyElement()
	return core.CreateIvDomainFromTo(v, v)
}

func (this *AdditionTerm) String() string {
	return fmt.Sprintf("%s +T %s", this.term1, this.term2)
}

func (this *AdditionTerm) HasVarId(varid core.VarId) bool {
	return this.term1.HasVarId(varid) || this.term2.HasVarId(varid)
}

func (this *AdditionTerm) Evaluable() EvalState {
	if this.term1.Evaluable() == EMPTY || this.term2.Evaluable() == EMPTY {
		return EMPTY
	}

	if this.term1.Evaluable() == NOT_EVALUABLE_YET || this.term2.Evaluable() == NOT_EVALUABLE_YET {
		return NOT_EVALUABLE_YET
	}

	return EVALUABLE
}

func (this *AdditionTerm) GetInf() ITerm {
	//inf+inf
	return CreateAdditionTerm(this.term1.GetInf(), this.term2.GetInf())
}

func (this *AdditionTerm) GetSup() ITerm {
	//sup+sup
	return CreateAdditionTerm(this.term1.GetSup(), this.term2.GetSup())
}

// CreateDomTerm returns a new DomTerm
func CreateAdditionTerm(term1 ITerm, term2 ITerm) *AdditionTerm {
	t := new(AdditionTerm)

	t.term1 = term1
	t.term2 = term2

	return t
}
