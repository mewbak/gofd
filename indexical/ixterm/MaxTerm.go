package ixterm

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

//Single-Val-Terms

// MaxTerm examples
// X in max(Y)
// X in max(Y) .. max(Z) 		(both MaxTerms)
type MaxTerm struct {
	varId core.VarId
	dom   *core.IvDomain
}

// returns the (current) value of the term (can change, cause of underlying
// domain)
func (this *MaxTerm) GetValue() *core.IvDomain {
	max := this.dom.GetMax()
	return core.CreateIvDomainFromTo(max, max)
}

func (this *MaxTerm) String() string {
	return fmt.Sprintf("max(%s)", this.dom)
}

func (this *MaxTerm) HasVarId(varid core.VarId) bool {
	return this.varId == varid
}

func (this *MaxTerm) Evaluable() EvalState {
	if this.dom.IsEmpty() {
		return EMPTY
	}

	return EVALUABLE
}

func (this *MaxTerm) GetInf() ITerm {
	//min
	return CreateMinTerm(this.varId, this.dom)
}

func (this *MaxTerm) GetSup() ITerm {
	//max
	return CreateMaxTerm(this.varId, this.dom)
}

// CreateMaxTerm returns a new MaxTerm
func CreateMaxTerm(varId core.VarId, dom *core.IvDomain) *MaxTerm {
	t := new(MaxTerm)
	t.varId = varId
	t.dom = dom

	return t
}
