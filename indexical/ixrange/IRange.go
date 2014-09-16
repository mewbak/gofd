package ixrange

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
)

type IRange interface {
	// Process collectChanges, if it has the changed variable/domain (where some
	// values has been removed) as input-variable. Have to return slice, because
	// of usage of append
	// dom: output-variable e.g. "X" of "X in dom(y)+3"
	Process(dom *core.IvDomain) []*core.IvDomPart

	// HasVarAsInput returns, if the specific Range has a specific variable as
	// input-variable (right side of expression)
	HasVarAsInput(varid core.VarId) bool

	// String returns a string representation of the specific range
	String() string

	// GetValue returns the containing value (*IvDomain)
	// can be expensive for IRanges like AddRange, MultRange, ... and great
	// involved domains
	GetValue() *core.IvDomain

	// returns, if a range is evaluable (yes, no e.g. containing "val-term",
	// containing term is empty)
	Evaluable() ixterm.EvalState

	CheckEntail(outDom *core.IvDomain) bool
	CheckDisentail(outDom *core.IvDomain) bool
}

func GetMinMaxTermsWithoutVarID(dom *core.IvDomain) (ixterm.ITerm, ixterm.ITerm) {
	return ixterm.CreateMinTerm(-1, dom), ixterm.CreateMaxTerm(-1, dom)
}
