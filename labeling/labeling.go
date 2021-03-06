// package labeling allows to search a solution of a finite domain problem
// by implicit enumeration
package labeling

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

// Labeling searches recursivly assignments for all variables using
// backtracking, as soon as assignment leads to a failed domain.
// Configurations choose enumerations strategy and variable selection
// strategy, where the last configuration counts.
func Labeling(store *core.Store, resultQuery ResultQuery,
	configurations ...interface{}) bool {
	store.IsConsistent()             // finalize propagation
	newStore := store.Clone(nil)     // local copy
	strategy := InDomainRange        // default strategy
	varSelect := SmallestDomainFirst // default variable selection
	for _, configuration := range configurations {
		stmp, instanceofStrategy :=
			configuration.(func(domain core.Domain) chan int)
		vtmp, instanceofVarSelect :=
			configuration.(func(store *core.Store) (core.VarId, bool))
		if instanceofStrategy {
			strategy = stmp
		} else if instanceofVarSelect {
			varSelect = vtmp
		} else {
			msg := "Labeling: config neither strat. nor var. selection %v"
			panic(fmt.Sprintf(msg, configuration))
		}
	}
	statSearch := resultQuery.GetSearchStatistics()
	if store.GetLoggingStat() {
		statSearch.UpdateStoreStatistics(store)
		statSearch.IncNodes()
		statStore := store.GetStat()
		statSearch.setInitialPropagators(statStore.GetActPropagators())
		statSearch.setInitialVars(statStore.GetActVariables())
	}
	fix(newStore, resultQuery, strategy, varSelect)
	logger := core.GetLogger()
	if logger.DoDebug() {
		logger.Dln("LABELING_Finished")
	}
	newStore.Close() // cleanup local copy of the store
	return resultQuery.getResultStatus()
}

// fix fixes given variable to value and continues recursivly with the next
// variable if there is one. Backtracks if assignment leads to failure.
func fix(store *core.Store, resultQuery ResultQuery,
	strategy func(domain core.Domain) chan int,
	varSelect func(store *core.Store) (core.VarId, bool)) bool {
	stat := resultQuery.GetSearchStatistics()
	consistent := store.IsConsistent() // finalize propagation
	if store.GetLoggingStat() {
		stat.UpdateStoreStatistics(store) // update our cumulative statistics
	}
	if consistent { // expecting true for "idle" and false for "failed"
		varId, hasNext := varSelect(store)
		logger := core.GetLogger()
		loggerDoInfo := logger.DoInfo()
		if loggerDoInfo {
			logger.If("LABELING_Store is ready, start/continue labeling...")
			logger.If("LABELING_HasNext: %v VarId: %v", hasNext, varId)
		}
		if hasNext {
			variable, _ := store.GetIntVar(varId)
			// println(fmt.Sprintf("Fix: %v", variable))
			for value := range strategy(variable.Domain) {
				// println(fmt.Sprintf("  Fix to %v", value))
				if loggerDoInfo {
					logger.If("LABELING_FixVariable: %v value: %v", varId,
						value)
				}
				// changes is a ChangeEvent to fix a variable
				changes := FixVariable(store, varId, value)
				// the propagations from FixVariable are collected in fix
				newStore := store.Clone(changes)
				if store.GetLoggingStat() {
					stat.IncNodes()
				}
				fix_result := fix(newStore, resultQuery, strategy, varSelect)
				newStore.Close()
				if fix_result {
					return true
				}
			}
			// node failed, backstep, next value
			stat.IncFailedNodes()
			return false
		} else {
			if loggerDoInfo {
				for _, id := range store.GetVariableIDs() {
					logger.If("LABELING_Variable %s (id=%v) assigned to %v",
						store.GetName(id), id, store.GetDomain(id))
				}
			}
			return resultQuery.onResult(store)
		}
	}
	stat.IncFailedNodes()
	return false
}

// LabelingSplit ...
func LabelingSplit(store *core.Store, resultQuery ResultQuery,
	configurations ...interface{}) bool {
	newStore := store.Clone(nil)
	varSelect := SmallestDomainFirst
	for _, configuration := range configurations {
		vtmp, instanceofVarSelect :=
			configuration.(func(store *core.Store) (core.VarId, bool))
		if instanceofVarSelect {
			varSelect = vtmp
		} else {
			msg := "Labeling: configuration not a variable selection %v"
			panic(fmt.Sprintf(msg, configuration))
		}
	}
	if core.GetLogger().DoInfo() {
		core.GetLogger().If("LABELING_LX start Labeling Divide")
	}
	if store.GetLoggingStat() {
		stat := resultQuery.GetSearchStatistics()
		stat.IncNodes()
	}
	divide(newStore, resultQuery, varSelect)
	newStore.Close()
	return resultQuery.getResultStatus()
}

// divide ...
func divide(store *core.Store, resultQuery ResultQuery,
	varSelect func(store *core.Store) (core.VarId, bool)) bool {
	stat := resultQuery.GetSearchStatistics()
	logInfo := core.GetLogger().DoInfo()
	if store.IsConsistent() { //expecting true ("idle") or false ("failed")
		varId, hasNext := varSelect(store)
		if hasNext {
			domain := store.GetDomain(varId)
			min, max := domain.GetMin(), domain.GetMax()
			middle := ((max - min) / 2) + min
			if logInfo {
				core.GetLogger().If("LABELING_LX var %d %s", varId, domain)
				core.GetLogger().If("LABELING_LX min=%d middle=%d max=%d",
					min, middle, max)
			}
			changesLeft := ResizeVariableDomain(store, varId, min, middle)
			lessThanStore := store.Clone(changesLeft)
			if logInfo {
				core.GetLogger().If("LABELING_LX divide")
			}
			if store.GetLoggingStat() {
				stat.IncNodes()
			}
			if divide(lessThanStore, resultQuery, varSelect) {
				return true
			}
			lessThanStore.Close()
			changesRight := ResizeVariableDomain(store, varId, middle+1, max)
			greaterThanStore := store.Clone(changesRight)
			if logInfo {
				core.GetLogger().If("LABELING_LX divide, left didnt work" +
					"-> try right")
			}
			if store.GetLoggingStat() {
				stat.IncNodes()
			}
			if divide(greaterThanStore, resultQuery, varSelect) {
				return true
			}
			greaterThanStore.Close()
			if logInfo {
				core.GetLogger().If("LABELING_LX divide, non worked" +
					" -> backstep")
			}
			stat.IncFailedNodes()
			return false
		} else {
			return resultQuery.onResult(store)
		}
	}
	if logInfo {
		core.GetLogger().If("LABELING_failed")
	}
	stat.IncFailedNodes()
	return false
}

// FixVariable creates ChangeEvent to fix a variable directly in the store
func FixVariable(store *core.Store,
	varId core.VarId, value int) *core.ChangeEvent {
	change_event := core.CreateChangeEvent()
	change_entry := core.CreateChangeEntry(varId)
	domain := store.GetDomain(varId)
	d := domain.Copy()
	d.Remove(value)
	change_entry.SetValues(d)
	change_event.AddChangeEntry(change_entry)
	return change_event
}

// ResizeVariableDomain creates ChangeEvent to restrict a variable
// to new bounds
func ResizeVariableDomain(store *core.Store,
	varId core.VarId, min, max int) *core.ChangeEvent {
	change_event := core.CreateChangeEvent()
	change_entry := core.CreateChangeEntry(varId)
	domain := store.GetDomain(varId)
	d := domain.GetDomainOutOfBounds(min, max)
	change_entry.SetValues(d)
	change_event.AddChangeEntry(change_entry)
	return change_event
}
