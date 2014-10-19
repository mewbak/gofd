package explicit

import (
	"bitbucket.org/gofd/gofd/core"
)

// CreateWeightedSumBounds creates a weighted sum-constraint.
// it uses X*C=Y and GSum-Constraint
// Example: 			1*X + 2*Y + 3*Q + 4*R = SUM
// with X*C=Y -->       H1  + H2  + H3  + H4 = SUM
func CreateWeightedSumBounds(store *core.Store, sum core.VarId,
	weights []int, intVars []core.VarId) []core.Constraint {
	if core.GetLogger().DoDebug() {
		core.GetLogger().Dln("CreateWeightedSumBounds")
	}
	if len(weights) != len(intVars) {
		msg := "CreateWeightedSumBounds: # weights != # vars"
		panic(msg)
	}
	prop_list := make([]core.Constraint, len(intVars))
	helperVarList := make([]core.VarId, len(intVars))
	// from X*w1 + Y*w2 + ... + Q*wn  (=Z)
	// to 	H1  + H2 + ...   + Hn    (=Z)
	for i, X := range intVars {
		// for Hi = wi*Xi only multiples of Xi may be values of Hi
		H := core.CreateAuxIntVarExValues(store,
			core.ScalarSlice(weights[i], store.GetDomain(X).Values_asSlice()))
		helperVarList[i] = H
		prop := CreateCXeqYBounds(weights[i], X, H)
		prop_list[i] = prop
	}
	// sum = H1 + H2 + ... + Hn
	for _, prop := range CreateSumBounds(store, sum, helperVarList) {
		prop_list = append(prop_list, prop)
	}
	return prop_list
}
