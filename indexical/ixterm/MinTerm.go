package ixterm

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

//Single-Val-Terms

// MinTerm examples
// X in min(Y)
// X in min(Y) .. min(Z) 		(both MinTerms)
type MinTerm struct {
	varId core.VarId
	dom   *core.IvDomain
}

// returns the (current) value of the term (can change, cause of underlying
// domain)
func (this *MinTerm) GetValue() *core.IvDomain {
	min := this.dom.GetMin()
	return core.CreateIvDomainFromTo(min, min)
}

func (this *MinTerm) HasVarId(varid core.VarId) bool {
	return this.varId == varid
}

func (this *MinTerm) String() string {
	return fmt.Sprintf("min(%s)", this.dom)
}

func (this *MinTerm) Evaluable() EvalState {
	if this.dom.IsEmpty() {
		return EMPTY
	}

	return EVALUABLE
}

func (this *MinTerm) GetInf() ITerm {
	//min
	return CreateMinTerm(this.varId, this.dom)
}

func (this *MinTerm) GetSup() ITerm {
	//max
	return CreateMaxTerm(this.varId, this.dom)
}

// CreateMinTerm returns a new MinTerm
func CreateMinTerm(varId core.VarId, dom *core.IvDomain) *MinTerm {
	t := new(MinTerm)
	t.varId = varId
	t.dom = dom

	return t
}
