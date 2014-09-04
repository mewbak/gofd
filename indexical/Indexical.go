package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical/ixrange"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"fmt"
)

// INFO:
// ranges with lower positions in slice have greater priority than
// ranges with higher positions in slice
// X>Y --> Y in -inf...(max(X)-1), --> X in (min(Y)+1)...+inf

// Indexical represents an indexial
type Indexical struct {
	out_varid        core.VarId
	out_var_IvDomain *core.IvDomain
	r                ixrange.IRange
}

// collectChanges collect changes from this specific indexical
// one indexical per out-variable. So each indexical creates a ChangeEntry
// for the outvar and adds it to the ChangeEvent-parameter
func (this *Indexical) collectChanges(evt *core.ChangeEvent) *core.IvDomain {
	d := this.r.GetValue()
	removingD := this.out_var_IvDomain.DifferenceWithIvDomain(d)
	if len(removingD.GetParts()) != 0 {
		return removingD
	}
	return nil
}

func (this *Indexical) GetCheckingIndexical() *CheckingIndexical {
	return CreateCheckingIndexical(this.out_varid,
		this.out_var_IvDomain, this.r)
}

// Process collectChanges, if it has the changed variable/domain
// (where some values has been removed) as input-variable
func (this *Indexical) Process(evt *core.ChangeEvent,
	changeEntry *core.ChangeEntry, updating bool) {
	chEntry := core.CreateChangeEntry(this.out_varid)
	changeDom := this.collectChanges(evt)
	if changeDom != nil {
		// updating ggf. in IndexicalProcessor reinziehen
		if updating {
			this.out_var_IvDomain.Removes(changeDom)
		}
		chEntry.SetValues(changeDom)
		evt.AddChangeEntry(chEntry)
	}
}

func (this *Indexical) Evaluable() ixterm.EvalState {
	if this.out_var_IvDomain.IsEmpty() {
		return ixterm.EMPTY
	}
	return this.r.Evaluable()
}

// HasVarAsInput returns, if the specific Indexical has a specific variable as
// input-variable (right side of expression). Indexicals asks each range, it
// contains.
func (this *Indexical) HasVarAsInput(varid core.VarId) bool {
	if this.r.HasVarAsInput(varid) {
		return true
	}
	return false
}

func (this *Indexical) String() string {
	return fmt.Sprintf("%s in %s", this.out_var_IvDomain, this.r)
}

// CreateIndexical creates a new indexical
func CreateIndexical(out_varid core.VarId, out_var_IvDomain *core.IvDomain,
	r ixrange.IRange) *Indexical {
	indexical := new(Indexical)
	indexical.out_varid = out_varid
	indexical.out_var_IvDomain = out_var_IvDomain
	indexical.r = r
	return indexical
}

func GetIndexicalPriosHighestToLowest() []int {
	return []int{HIGHEST, HIGH, LOW, LOWEST}
}

const (
	HIGHEST = 0
	HIGH    = 1
	LOW     = 2
	LOWEST  = 3
)
