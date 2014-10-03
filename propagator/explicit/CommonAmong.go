package explicit

import (
	"bitbucket.org/gofd/gofd/core"
)

// implementation of the Common constraint with Among
// signature: Count({X1,...,Xi}, {Y1,...,Yj}, N, M)
// exactly N variables in Xi take values in Yj and
// exactly M variables in Yj take values in Xi

// CreateCommonAmong creates a Common constraint modelled with Among
func CreateCommonAmong(xi []core.VarId, yj []core.VarId, n core.VarId,
	m core.VarId, store *core.Store) []core.Constraint {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateCommon_propagator")
	}

	//make unions of the domains so the unions can be used as K
	xiUnion := createUnionFromDomains(xi, store)
	yjUnion := createUnionFromDomains(yj, store)

	//create two Among constraints
	return []core.Constraint{CreateAmong(xi, yjUnion, n),
		CreateAmong(yj, xiUnion, m)}
}

// createUnionFromDomains creates the union of the variables' domains
// and returns the result as an int-slice
func createUnionFromDomains(variables []core.VarId, store *core.Store) []int {
	domain := core.CreateExDomain()
	for _, y := range variables {
		varDomain := store.GetDomain(y)
		temp := domain.Union(varDomain)
		domain = temp.(*core.ExDomain)
	}

	domainValues := make([]int, len(domain.Values))
	i := 0
	for value, _ := range domain.Values {
		domainValues[i] = value
		i += 1
	}

	return domainValues
}
