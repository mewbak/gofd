package ixterm

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

//Single-Val-Terms (waiting)

// ValTerm examples
// X in val(Y)
type ValTerm struct {
	varId core.VarId
	dom   *core.IvDomain
}

// returns the (current) value of the term (can change, cause of underlying
// domain)
func (this *ValTerm) GetValue() *core.IvDomain {
	if !this.dom.IsGround() {
		panic("val-Term not ground! You have to call Evaluable-function before" +
			"calling GetValue-Function to avoid this panicing")
	}
	val := this.dom.GetMin()
	return core.CreateIvDomainFromTo(val, val)
}

func (this *ValTerm) String() string {
	return fmt.Sprintf("val(%s)", this.dom)
}

func (this *ValTerm) HasVarId(varid core.VarId) bool {
	return this.varId == varid
}

func (this *ValTerm) Evaluable() EvalState {
	if this.dom.IsEmpty() {
		return EMPTY
	}

	if !this.dom.IsGround() {
		return NOT_EVALUABLE_YET
	}

	return EVALUABLE
}

func (this *ValTerm) GetInf() ITerm {
	//val
	return CreateValTerm(this.varId, this.dom)
}

func (this *ValTerm) GetSup() ITerm {
	//val
	return CreateValTerm(this.varId, this.dom)
}

// CreateMinTerm returns a new MinTerm
func CreateValTerm(varId core.VarId, dom *core.IvDomain) *ValTerm {
	t := new(ValTerm)
	t.varId = varId
	t.dom = dom

	return t
}
