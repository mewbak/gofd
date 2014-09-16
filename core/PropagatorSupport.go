package core

// SendChangesToStore sends changeevent to store and logs if logging == true.
func SendChangesToStore(evt *ChangeEvent, caller Constraint) {
	if GetLogger().DoDebug() {
		msg := "%s: propagate_'communicate change, evt-value: %s'"
		GetLogger().Df(msg, caller, evt)
	}
	caller.GetOutCh() <- evt
}

// LogIncomingChange logs incoming changes if logging == true.
func LogIncomingChange(caller Constraint, store *Store,
	changeEntry *ChangeEntry) {
	if GetLogger().DoDebug() {
		msg := "%s_Start_'Incoming Change for %s'"
		GetLogger().Df(msg, caller, store.GetName(changeEntry.GetID()))
	}
}

// LogInitConsistency logs text "initial consistency check"
func LogInitConsistency(caller Constraint) {
	if GetLogger().DoDebug() {
		GetLogger().Df("%s: %s", caller, "initial consistency check")
	}
}

// LogEntailmentDetected logs, that entailment is detected
func LogEntailmentDetected(caller Constraint, c Constraint) {
	if GetLogger().DoDebug() {
		GetLogger().Df("%s - %s - %s", caller, "entailment detected on", c)
	}
}

// LogPropagationOfConstraint logs propagation
func LogPropagationOfConstraint(c Constraint) {
	if GetLogger().DoDebug() {
		GetLogger().Df("propagate constraint %s", c)
	}
}

// GetVaridToIntervalDomains from any domain implementation (no new instance)
func GetVaridToIntervalDomains(idoms map[VarId]Domain) map[VarId]*IvDomain {
	intervalDoms := make(map[VarId]*IvDomain, len(idoms))
	for varid, idom := range idoms {
		if dom, k := idom.(*IvDomain); k {
			intervalDoms[varid] = dom
		} else if dom, k := idom.(*ExDomain); k {
			intervalDoms[varid] = CreateIvDomainFromDomain(dom)
		}
	}
	return intervalDoms
}

// GetVaridToExplicitDomainsMap from any domain implementation
// (no new instance)
func GetVaridToExplicitDomainsMap(idoms map[VarId]Domain) map[VarId]*ExDomain {
	explDoms := make(map[VarId]*ExDomain, len(idoms))
	for varid, idom := range idoms {
		if dom, k := idom.(*ExDomain); k {
			explDoms[varid] = dom
		} else if dom, k := idom.(*IvDomain); k {
			explDoms[varid] = CreateExDomainFromIvDomain(dom)
		}
	}
	return explDoms
}

// GetVaridToExplicitDomainsSlice from any domain implementation
// (no new instance).
func GetVaridToExplicitDomainsSlice(idoms []Domain) []*ExDomain {
	explDoms := make([]*ExDomain, len(idoms))
	for i, idom := range idoms {
		if dom, k := idom.(*ExDomain); k {
			explDoms[i] = dom
		} else if dom, k := idom.(*IvDomain); k {
			explDoms[i] = CreateExDomainFromIvDomain(dom)
		}
	}
	return explDoms
}
