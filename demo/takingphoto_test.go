package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"bitbucket.org/gofd/gofd/propagator/indexical"
	"bitbucket.org/gofd/gofd/propagator/reification"
	"fmt"
	"testing"
)

/* http://www.ps.uni-saarland.de/alice/manual/cptutorial/node44.html

Betty, Chris, Donald, Fred, Gary, Mary, and Paul want to align in
one row for taking a photo. Some of them have preferences next
to whom they want to stand:
- Betty wants to stand next to Gary and Mary.
- Chris wants to stand next to Betty and Gary.
- Fred wants to stand next to Mary and Donald.
- Paul wants to stand next to Fred and Donald.

necessary Constraints
* vars...
* dist-Constraints
  - dBG = CreateRICDist(B,G,1, S0)
  - dBM = CreateRICDist(B,M,1, S1)
  - ...
ok- iccardinality
  - CreateCardinality([S0,S1,...,S7],S)
  - Sum... S0+S1+S2+S3+...= S
ok- icAlldifferent
Results:
number of satisfied constraints 1-6: ok
number of satisfied constraints 7-8: no solution

Demonstrates the need and usefulness of reification
*/

func Test_takingphoto0a(t *testing.T) {
	setup()
	defer teardown()
	log("taking_photo: B,C,D in 0..2, 2 Constraints true --> ok")
	DoAlignInOneRowVERYSIMPLE(t, 2, true)
}

func Test_takingphoto0b(t *testing.T) {
	setup()
	defer teardown()
	log("taking_photo: B,C,D in 0..2, 2 Constraints true --> ok")
	DoAlignInOneRowSIMPLE(t, 2, true)
}

func Test_takingphoto0c(t *testing.T) {
	setup()
	defer teardown()
	log("taking_photo: B,C,D in 0..2, 3 Constraints true --> ok")
	DoAlignInOneRowSIMPLE(t, 3, false)
}

func Test_takingphoto1a(t *testing.T) {
	setup()
	defer teardown()
	log("taking_photo: B,C,D,F,G,M,P in 0..6, 3 Constraints true --> ok")
	DoAlignInOneRow(t, 3, true)
}

func Test_takingphoto1b(t *testing.T) {
	setup()
	defer teardown()
	log("taking_photo: B,C,D,F,G,M,P in 0..6, 6 Constraints true --> ok")
	DoAlignInOneRow(t, 6, true)
}

func Test_takingphoto2(t *testing.T) {
	setup()
	defer teardown()
	log("taking_photo: B,C,D,F,G,M,P in 0..6, 7 Constraints true --> fail")
	DoAlignInOneRow(t, 7, false)
}

func DoAlignInOneRow(t *testing.T, minSatisfyingConstraints int, expReady bool) {
	var B, C, D, F, G, M, P core.VarId
	core.CreateIntVarsFromTo(
		[]*core.VarId{&B, &C, &D, &F, &G, &M, &P},
		[]string{"B", "C", "D", "F", "G", "M", "P"}, store, 0, 6)
	var S0, S1, S2, S3, S4, S5, S6, S7 core.VarId
	var Z0, Z1, Z2, Z3, Z4, Z5, Z6, Z7 core.VarId
	core.CreateIntVarsFromTo(
		[]*core.VarId{&S0, &S1, &S2, &S3, &S4, &S5, &S6, &S7},
		[]string{"S0", "S1", "S2", "S3", "S4", "S5", "S6", "S7"},
		store, 0, 1)
	core.CreateIntVarsValues(
		[]*core.VarId{&Z0, &Z1, &Z2, &Z3, &Z4, &Z5, &Z6, &Z7},
		[]string{"Z0", "Z1", "Z2", "Z3", "Z4", "Z5", "Z6", "Z7"},
		store, []int{-1, 1})
	// |(B-G)|=1, |(B-M)|=1, ...
	dBG := indexical.CreateXplusYeqZ(B, Z0, G) // Betty next to Gary
	store.AddPropagator(reification.CreateReifiedConstraint(dBG, S0))

	dBM := indexical.CreateXplusYeqZ(B, Z1, M) // Betty next to Mary
	store.AddPropagator(reification.CreateReifiedConstraint(dBM, S1))

	dCB := indexical.CreateXplusYeqZ(C, Z2, B) // Chris next to Betty
	store.AddPropagator(reification.CreateReifiedConstraint(dCB, S2))

	dCG := indexical.CreateXplusYeqZ(C, Z3, G) // Chris next to Gary
	store.AddPropagator(reification.CreateReifiedConstraint(dCG, S3))

	dFM := indexical.CreateXplusYeqZ(F, Z4, M) //Fred next to Mary
	store.AddPropagator(reification.CreateReifiedConstraint(dFM, S4))

	dFD := indexical.CreateXplusYeqZ(F, Z5, D) // Fred next to Donald
	store.AddPropagator(reification.CreateReifiedConstraint(dFD, S5))

	dPF := indexical.CreateXplusYeqZ(P, Z6, F) // Paul next to Fred
	store.AddPropagator(reification.CreateReifiedConstraint(dPF, S6))

	dPD := indexical.CreateXplusYeqZ(P, Z7, D) // Paul next to Donald
	store.AddPropagator(reification.CreateReifiedConstraint(dPD, S7))

	S := core.CreateIntVarFromTo("S", store, minSatisfyingConstraints, 8)
	card := reification.CreateCardinality(store, S,
		[]core.VarId{S0, S1, S2, S3, S4, S5, S6, S7})
	store.AddPropagator(card)

	alldiff_prop := indexical.CreateAlldifferent(B, C, D, F, G, M, P)
	store.AddPropagators(alldiff_prop)

	labeling.SetAllvars([]core.VarId{S0, S1, S2, S3, S4, S5, S6, S7,
		B, C, D, F, G, M, P, Z0, Z1, Z2, Z3, Z4, Z5, Z6, Z7})
	query := labeling.CreateSearchOneQuery()
	result := labeling.Labeling(store, query,
		labeling.InDomainMin, labeling.VarSelect)
	set := query.GetResultSet()
	ready_test(t, "taking_photo", result, expReady)
	searchStat(query.GetSearchStatistics())
	if logger.GetLoggingLevel() >= core.LOG_INFO {
		for varid, value := range set[0] {
			if varid == B || varid == C || varid == D || varid == F ||
				varid == G || varid == M || varid == P {
				logger.Iln(fmt.Sprintf("%s: %v",
					store.GetName(varid), value))
			}
		}
	}
}

func DoAlignInOneRowVERYSIMPLE(t *testing.T, minSatisfied int, expReady bool) {
	var B, C core.VarId
	core.CreateIntVarsFromTo(
		[]*core.VarId{&B, &C},
		[]string{"B", "C"}, store, 0, 1)
	// C wants to be beside B
	// B wants to be beside C
	var S0, S1, Z0, Z1 core.VarId
	core.CreateIntVarsFromTo(
		[]*core.VarId{&S0, &S1},
		[]string{"S0", "S1"}, store, 0, 1)
	core.CreateIntVarsValues(
		[]*core.VarId{&Z0, &Z1},
		[]string{"Z0", "Z1"}, store, []int{-1, 1})
	// Distance as X-Y = Z, Z = {-1, 1}
	// C-B = Z0 --> C = B + Z0
	dCB := indexical.CreateXplusYeqZ(B, Z0, C)
	store.AddPropagator(reification.CreateReifiedConstraint(dCB, S0))
	dBC := indexical.CreateXplusYeqZ(C, Z1, B)
	store.AddPropagator(reification.CreateReifiedConstraint(dBC, S1))
	// C<->B<->D : 1, 0, 1
	// B<->C<->D : 1, 1, 0
	S := core.CreateIntVarFromTo("S", store, minSatisfied, 2)
	store.AddPropagator(
		reification.CreateCardinality(store, S, []core.VarId{S0, S1}))
	alldiff_prop := indexical.CreateAlldifferent(B, C)
	store.AddPropagators(alldiff_prop)
	query := labeling.CreateSearchOneQuery()
	result := labeling.Labeling(store, query)
	set := query.GetResultSet()
	ready_test(t, "taking_photo", result, expReady)
	if logger.GetLoggingLevel() >= core.LOG_INFO {
		for varid, value := range set[0] {
			if varid == B || varid == C {
				logger.Iln(fmt.Sprintf("%s: %v",
					store.GetName(varid), value))
			}
		}
	}
}

func DoAlignInOneRowSIMPLE(t *testing.T, minSatisfied int, expReady bool) {
	var B, C, D core.VarId
	core.CreateIntVarsFromTo(
		[]*core.VarId{&B, &C, &D},
		[]string{"B", "C", "D"}, store, 0, 2)
	// C wants to be beside B
	// B wants to be beside C
	var S0, S1, S2 core.VarId
	var Z0, Z1, Z2 core.VarId
	core.CreateIntVarsFromTo(
		[]*core.VarId{&S0, &S1, &S2},
		[]string{"S0", "S1", "S2"}, store, 0, 1)
	core.CreateIntVarsValues(
		[]*core.VarId{&Z0, &Z1, &Z2},
		[]string{"Z0", "Z1", "Z2"}, store, []int{-1, 1})
	// distance as X - Y = Z, Z = {-1, 1}
	// C - B = Z0 --> C = B + Z0
	dCB := indexical.CreateXplusYeqZ(B, Z0, C)
	store.AddPropagator(reification.CreateReifiedConstraint(dCB, S0))
	dCD := indexical.CreateXplusYeqZ(D, Z1, C)
	store.AddPropagator(reification.CreateReifiedConstraint(dCD, S1))
	dBD := indexical.CreateXplusYeqZ(D, Z2, B)
	store.AddPropagator(reification.CreateReifiedConstraint(dBD, S2))
	// C<->B<->D : 1, 0, 1
	// B<->C<->D : 1, 1, 0
	S := core.CreateIntVarFromTo("S", store, minSatisfied, 3)
	store.AddPropagator(
		reification.CreateCardinality(store, S, []core.VarId{S0, S1, S2}))
	alldiff_prop := indexical.CreateAlldifferent(B, C, D)
	store.AddPropagators(alldiff_prop)
	query := labeling.CreateSearchOneQuery()
	result := labeling.Labeling(store, query)
	set := query.GetResultSet()
	ready_test(t, "taking_photo", result, expReady)
	if logger.GetLoggingLevel() >= core.LOG_INFO {
		for varid, value := range set[0] {
			if varid == B || varid == C || varid == D {
				logger.Iln(fmt.Sprintf("%s: %v",
					store.GetName(varid), value))
			}
		}
	}
}
