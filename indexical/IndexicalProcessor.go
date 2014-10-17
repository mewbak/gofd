package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
)

// ProcessIndexicals executes the indexical in the desired order
// indexicals: [[bounds_i1,bounds_i2][arc_i1,arc_i2]] means:
// first execute bounds (if changes appear, repeat execution). Then, if
// no changes can be collected from higher priority, go deeper to next
// indexical-level (in this case: from bound- to arc-consistency)
func ProcessIndexicals(iColl *IndexicalCollection,
	changeEntry *core.ChangeEntry,
	updating bool) *core.ChangeEvent {
	evt := core.CreateChangeEvent()
	for prio := range GetIndexicalPriosHighestToLowest() {
		for _, i := range iColl.GetIndexicalsAtPrio(prio) {
			if changeEntry == nil || i.HasVarAsInput(changeEntry.GetID()) {
				evalRes := i.Evaluable()
				if evalRes == ixterm.EVALUABLE {
					i.Process(evt, changeEntry, updating)
				} else if evalRes == ixterm.NOT_EVALUABLE_YET {
					//not executing current indexical
				} else if evalRes == ixterm.EMPTY {
					// stop indexical processing. A domain is empty or so...
					// send collected changes from previous indexicals
					// immediately to store
					return evt
				}
			}
		}
	}
	return evt
}
