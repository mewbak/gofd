package ixterm

import (
	"bitbucket.org/gofd/gofd/core"
)

const INFINITY = core.INFINITY
const NEG_INFINITY = core.NEG_INFINITY

//terms are: max, min, val, value

type ITerm interface {
	// returns the (current) value of the term (single-value domain)
	GetValue() *core.IvDomain
	// returns a string representation of the specific term
	String() string
	// returns true, if varid is involved in term
	HasVarId(varid core.VarId) bool
	// returns evaluation state of the specific term (e.g. empty)
	Evaluable() EvalState
	//returns the infimum of the specific term
	GetInf() ITerm
	//returns the supremum of the specific term
	GetSup() ITerm
}
