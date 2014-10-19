package ixterm

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"strings"
)

// SumTerm t1+t2+...t3 examples
// min(X)+min(Y)+max(Y)
// min(X)+max(Y)

type SumTerm struct {
	terms []ITerm
}

// GetValue returns the (current) value of the term (can change, cause of
// underlying domains). Returns result of term1+term2
func (this *SumTerm) GetValue() *core.IvDomain {
	v := 0
	for _, t := range this.terms {
		v += t.GetValue().GetAnyElement()
	}

	return core.CreateIvDomainFromTo(v, v)
}

func (this *SumTerm) String() string {
	s := make([]string, len(this.terms))
	for i, t := range this.terms {
		s[i] = fmt.Sprintf("%d", t.GetValue().GetAnyElement())
	}

	return fmt.Sprintf("%s", strings.Join(s, "+"))
}

func (this *SumTerm) HasVarId(varid core.VarId) bool {
	for _, t := range this.terms {
		if t.HasVarId(varid) {
			return true
		}
	}
	return false
}

func (this *SumTerm) Evaluable() EvalState {
	for _, t := range this.terms {
		rstate := t.Evaluable()
		if rstate == EMPTY {
			return EMPTY
		} else if rstate == NOT_EVALUABLE_YET {
			return NOT_EVALUABLE_YET
		}
	}

	return EVALUABLE
}

func (this *SumTerm) GetInf() ITerm {
	//inf+inf+inf+...
	ts := make([]ITerm, len(this.terms))

	for i, t := range this.terms {
		ts[i] = t.GetInf()
	}
	return CreateSumTerm(ts...)
}

func (this *SumTerm) GetSup() ITerm {
	//sup+sup+sup+...
	ts := make([]ITerm, len(this.terms))

	for i, t := range this.terms {
		ts[i] = t.GetSup()
	}
	return CreateSumTerm(ts...)
}

// CreateDomTerm returns a new DomTerm
func CreateSumTerm(terms ...ITerm) *SumTerm {
	t := new(SumTerm)

	t.terms = terms

	return t
}
