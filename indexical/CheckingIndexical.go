package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical/ixrange"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
)

// CheckingIndexical
type CheckingIndexical struct {
	out_varid        core.VarId
	out_var_IvDomain *core.IvDomain
	r                ixrange.IRange
}

func (this *CheckingIndexical) IsEntailed() bool {
	if this.r.Evaluable() == ixterm.EMPTY ||
		this.r.Evaluable() == ixterm.NOT_EVALUABLE_YET {
		return false
	}
	b := this.r.CheckEntail(this.out_var_IvDomain)
	return b
}

func (this *CheckingIndexical) IsDisentailed() bool {
	return this.r.CheckDisentail(this.out_var_IvDomain)
}

func (this *CheckingIndexical) Evaluable() ixterm.EvalState {
	if this.out_var_IvDomain.IsEmpty() {
		return ixterm.EMPTY
	}
	return this.r.Evaluable()
}

func (this *CheckingIndexical) String() string {
	return fmt.Sprintf("CheckingIndexical: %s in %s",
		this.out_var_IvDomain, this.r)
}

// CreateCheckingIndexical creates a CheckingIndexical
// r: checkingRange
func CreateCheckingIndexical(out_varid core.VarId, out_var_IvDomain *core.IvDomain,
	r ixrange.IRange) *CheckingIndexical {
	CheckingIndexical := new(CheckingIndexical)
	CheckingIndexical.out_varid = out_varid
	CheckingIndexical.out_var_IvDomain = out_var_IvDomain
	CheckingIndexical.r = r
	return CheckingIndexical
}
