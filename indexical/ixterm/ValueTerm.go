package ixterm

import (
	"bitbucket.org/gofd/gofd/core"

	"fmt"
)

//Single-Val-Terms

// ValueTerm for indexical-definitions like
// X in 5
// X in min(Y) .. 10 		(where 10 is a ValueTerm)
type ValueTerm struct {
	value int
}

// returns the value of the term
func (this *ValueTerm) GetValue() *core.IvDomain {
	return core.CreateIvDomainFromTo(this.value, this.value)
}

func (this *ValueTerm) String() string {
	return fmt.Sprintf("%v", this.value)
}

func (this *ValueTerm) HasVarId(varid core.VarId) bool {
	return false
}

func (this *ValueTerm) Evaluable() EvalState {
	return EVALUABLE
}

func (this *ValueTerm) GetInf() ITerm {
	return CreateValueTerm(this.value)
}

func (this *ValueTerm) GetSup() ITerm {
	return CreateValueTerm(this.value)
}

// CreateValueTerm returns a new ValueTerm
func CreateValueTerm(v int) *ValueTerm {
	t := new(ValueTerm)
	t.value = v
	return t
}
