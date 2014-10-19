package indexical

import (
	"bitbucket.org/gofd/gofd/core"
)

func InitProcessConstraint(caller IIndexicalConstraint, updating bool) {
	var evt *core.ChangeEvent
	core.LogInitConsistency(caller)
	iColl := caller.GetIndexicalCollection()
	evt = ProcessIndexicals(iColl, nil, updating)
	core.SendChangesToStore(evt, caller)
}

func ProcessConstraint(caller IIndexicalConstraint, updating bool) {
	var evt *core.ChangeEvent
	inCh := caller.GetInCh()
	for changeEntry := range inCh {
		RemoveValues(caller, changeEntry)
		iColl := caller.GetIndexicalCollection()
		evt = ProcessIndexicals(iColl, changeEntry, updating)
		core.SendChangesToStore(evt, caller)
	}
}

func RemoveValues(caller core.Constraint, changeEntry *core.ChangeEntry) {
	varIds := caller.GetVarIds()
	domains := caller.GetDomains()
	for i, varid := range varIds {
		if varid == changeEntry.GetID() {
			domains[i].Removes(changeEntry.GetValues())
		}
	}
}
