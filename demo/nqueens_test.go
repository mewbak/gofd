package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"bitbucket.org/gofd/gofd/propagator"
	"bitbucket.org/gofd/gofd/propagator/interval"
	"fmt"
	//"runtime"
	"testing"
)

// The infamous N queens problem:
//   Place N queens on an NxN chess board without any queen
//   attacking any other queen.
// We model each queen as finite domain variable from 0 to N-1
// and impose the constraints, that there are no two queens in
// the same column, the same left diagonal, the same right
// diagonal. There are different models ranging from using
// primitive constraints only to various degrees of alldifferent
// propagators.

/* Helper functions */

// helper to show the solutions
func show_nqueens_results(resultSet map[int]map[core.VarId]int) {
	for i, result := range resultSet {
		fmt.Printf("Sol %d:", i)
		for _, varId := range core.SortedKeys_MapVarIdToInt(result) {
			fmt.Printf(" %s=%d",
				core.GetNameRegistry().GetName(varId), result[varId])
		}
		fmt.Printf("\n")
	}
}

// helper to show one set of variables
func show_nqueens(queens []core.VarId) {
	for i, vId := range queens {
		v, _ := store.GetIntVar(vId)
		fmt.Printf("%2d: ", i)
		if v.IsGround() {
			for i := 0; i < v.Domain.GetMin(); i += 1 {
				fmt.Printf(" ")
			}
			fmt.Printf("*\n")
		} else {
			fmt.Printf("(open: %v)\n", v.Domain)
		}
	}
}

// create the variables Qi for 0 <= i <= N-1
// identical for all models
func createQueens(n int) []core.VarId {
	queens := make([]core.VarId, n)
	for i := 0; i < n; i++ {
		varname := fmt.Sprintf("Q%d", i)
		queens[i] = core.CreateIntVarIvFromTo(varname, store, 0, n-1)
	}
	return queens
}

// number of solutions for some N-Queens instances
var nqueensNoSols []int = []int{
	0,    // 0
	1,    // 1
	0,    // 2
	0,    // 3
	2,    // 4
	10,   // 5
	4,    // 6
	40,   // 7
	92,   // 8
	352,  // 9
	724,  // 10
	2680, // 11
}

// nqueensall computes and checks all solutions for the nqueens problems
// parametrized by functions to generate the constraints that there
// must be different columns and different diagonals.
func nqueensall(t *testing.T, testname string, N int,
	diffcols func(queens []core.VarId, store *core.Store),
	diffdiag func(queens []core.VarId, store *core.Store)) {
	msg := fmt.Sprintf("nqueens%d%s", N, testname)
	log(msg)
	queens := createQueens(N) // create the finite domain vars
	diffcols(queens, store)   // add constraints different columns
	diffdiag(queens, store)   // add constraints different queens
	// compute all solutions
	query := labeling.CreateSearchAllQuery()
	result := labeling.Labeling(store, query,
		labeling.SmallestDomainFirst, labeling.InDomainMin)
	noSols := nqueensNoSols[N]
	ready_test(t, msg, result, noSols > 0)
	if logger.DoInfo() {
		logger.If("No results: %v", query.GetResultSet())
	}
	if len(query.GetResultSet()) != noSols {
		t.Errorf("nqueens(%d) number of solutions = %d, want %d",
			N, len(query.GetResultSet()), noSols)
	}
	// show_nqueens(queens)
	if logger.GetLoggingLevel() <= core.LOG_NONE {
		searchStat(query.GetSearchStatistics())
		if logger.GetLoggingLevel() <= core.LOG_INFO {
			show_nqueens_results(query.GetResultSet())
		}
	}
}

/* First: model with primitive constraints only resulting in many
 * independent propagators.
 */

// no two queens in the same column; Qi != Qj for 0 <= i < j <= N-1
func differentColsPrim(queens []core.VarId, store *core.Store) {
	for i := 0; i < len(queens); i++ {
		for j := i + 1; j < len(queens); j++ {
			prop := propagator.CreateXneqY(queens[j], queens[i])
			store.AddPropagator(prop)
		}
	}
}

// helper for Qi +/- o != Qj
func checkOffset(queens []core.VarId, store *core.Store, offset int) {
	headQueen := queens[0]
	prop := propagator.CreateXplusCneqY(headQueen,
		offset, queens[core.AbsInt(offset)])
	store.AddPropagator(prop)
}

// no two queens in the same left and right diagonal
// Qi +/- o != Qj for 0 <= i < j <= N-1, for 1 <= o <= j-i
func differentDiagPrim(queens []core.VarId, store *core.Store) {
	for i := 0; i < len(queens)-1; i++ {
		remaining := queens[i:]
		for offset := 1; offset < len(remaining); offset++ {
			checkOffset(remaining, store, offset)
			checkOffset(remaining, store, -offset)
		}
	}
}

// nqueens with primitive constraints only
func nqueensPrim(N int, t *testing.T) {
	nqueensall(t, "Prim", N, differentColsPrim, differentDiagPrim)
}

// run some tests
func Test_nqueens5Prim(t *testing.T) {
	setup()
	defer teardown()
	nqueensPrim(5, t)
}

func Test_nqueens7Prim(t *testing.T) {
	setup()
	defer teardown()
	nqueensPrim(7, t)
}

func Test_nqueens8Prim(t *testing.T) {
	setup()
	defer teardown()
	nqueensPrim(8, t)
}

//func Test_nqueens9Prim(t *testing.T) {
//	setup()
//  defer teardown()
//  nqueensPrim(9, t)
//}

/* Second: We replace the different columns constraint with a
 * dedicated AllDifferent one. This propagator generates the
 * same propagators internally, but uses them in one coroutine
 * and thus just needs one copy of the variables.
 */

func differentColsAllDiff(queens []core.VarId, store *core.Store) {
	prop := propagator.CreateAlldifferent(queens...)
	store.AddPropagators(prop)
}

func nqueensAllDiffCols(N int, t *testing.T) {
	nqueensall(t, "AllDiffCols", N, differentColsAllDiff, differentDiagPrim)
}

// run some tests
func Test_nqueens7AllDiffCols(t *testing.T) {
	setup()
	defer teardown()
	nqueensAllDiffCols(7, t)
}

func Test_nqueens8AllDiffCols(t *testing.T) {
	setup()
	defer teardown()
	nqueensAllDiffCols(8, t)
}

//func Test_nqueens9AllDiffCols(t *testing.T) {
//	setup()
//  defer teardown()
//  nqueensAllDiffCols(9, t)
//}

/* Third: We use the better dedicated AllDifferent constraint for the
 * diagonals as well. To that end we need variables expressing the values
 * on the diagonal.
 */

// no two queens in the same left and right diagonal
// Qi +/- o != Qj for 0 <= i < j <= N-1, for 1 <= o <= j-i
func differentDiagAuxAllDiff(queens []core.VarId, store *core.Store) {
	lenqueens := len(queens)
	// one more to ease indexing
	negoffqueens := make([]core.VarId, lenqueens)
	posoffqueens := make([]core.VarId, lenqueens)
	min, max := store.GetMinMaxDomain(queens[0]) // any queen
	for i := 1; i < len(queens); i++ {
		negoffqueens[i] = core.CreateAuxIntVarFromTo(store, min-i, max-i)
		propneg := propagator.CreateXplusCeqY(queens[i], -i, negoffqueens[i])
		store.AddPropagator(propneg)
		posoffqueens[i] = core.CreateAuxIntVarFromTo(store, min+i, max+i)
		proppos := propagator.CreateXplusCeqY(queens[i], i, posoffqueens[i])
		store.AddPropagator(proppos)
	}
	negoffqueens[0] = queens[0]
	posoffqueens[0] = queens[0]
	negprop := propagator.CreateAlldifferent(negoffqueens...)
	store.AddPropagator(negprop)
	posprop := propagator.CreateAlldifferent(posoffqueens...)
	store.AddPropagator(posprop)
}

func nqueensAllDiffBoth(N int, t *testing.T) {
	nqueensall(t, "AllDiffBoth", N,
		differentColsAllDiff,    // reuse the columns constraint
		differentDiagAuxAllDiff) // new one for diagonals
}

// run some tests
func Test_nqueens7AllDiffBoth(t *testing.T) {
	setup()
	defer teardown()
	nqueensAllDiffBoth(7, t)
}

func Test_nqueens8AllDiffBoth(t *testing.T) {
	setup()
	defer teardown()
	nqueensAllDiffBoth(8, t)
}

//func Test_nqueens9AllDiffBoth(t *testing.T) {
//	setup()
//  defer teardown()
//  nqueensAllDiffBoth(9, t)
//}

/* Fourth: We replace the diagonal constraints by a dedicated AllDifferent
 * constraint, that allows an offset per variable. As with the columns, that
 * constraint has only one propagator and thus reduces the amount of copying
 * and therefore the communication, which in turn reduces the runtime for
 * these small problems with relatively simple propagators.
 */

func differentDiagAllDiff(queens []core.VarId, store *core.Store) {
	left_offset := make([]int, len(queens))
	right_offset := make([]int, len(queens))
	for i, _ := range queens {
		left_offset[i] = -i
		right_offset[i] = i
	}
	left_prop := interval.CreateAlldifferent_Offset(queens, left_offset)
	store.AddPropagator(left_prop)
	right_prop := interval.CreateAlldifferent_Offset(queens, right_offset)
	store.AddPropagator(right_prop)
}

func nqueensAllDiffOnly(N int, t *testing.T) {
	nqueensall(t, "OnlyAllDiffs", N,
		differentColsAllDiff, // reuse column constraints
		differentDiagAllDiff, // new dedicated offset ones for diagonals
	)
}

// run some tests
func Test_nqueens7AllDiffOnly(t *testing.T) {
	setup()
	defer teardown()
	nqueensAllDiffOnly(7, t)
}

func Test_nqueens8AllDiffOnly(t *testing.T) {
	setup()
	defer teardown()
	nqueensAllDiffOnly(8, t)
}

func Test_nqueens9AllDiffOnly(t *testing.T) {
	setup()
	defer teardown()
	nqueensAllDiffOnly(9, t)
}

/* Todo: Alldistinct for all Alldifferent variants */
