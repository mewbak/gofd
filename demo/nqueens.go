package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator"
	"bitbucket.org/gofd/gofd/propagator/interval"
	"fmt"
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

// Create an NQueens problem parametrized by size and
// functions that generate the constraints for
// different columns and different diagonals
func nqueensConstrain(store *core.Store, N int,
	diffcols func(queens []core.VarId, store *core.Store),
	diffdiag func(queens []core.VarId, store *core.Store)) []core.VarId {
	queens := createQueens(N, store) // create the finite domain vars
	diffcols(queens, store)          // add constraints different columns
	diffdiag(queens, store)          // add constraints different diags
	return queens
}

// create the variables Qi for 0 <= i <= N-1
// identical for all models
func createQueens(n int, store *core.Store) []core.VarId {
	queens := make([]core.VarId, n)
	for i := 0; i < n; i++ {
		varname := fmt.Sprintf("Q%d", i)
		queens[i] = core.CreateIntVarIvFromTo(varname, store, 0, n-1)
	}
	return queens
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
func NQueensPrim(store *core.Store, N int) []core.VarId {
	return nqueensConstrain(store, N,
		differentColsPrim, differentDiagPrim)
}

/* Second: We replace the different columns constraint with a
 * dedicated AllDifferent one. This propagator generates the
 * same propagators internally, but uses them in one coroutine
 * and thus just needs one copy of the variables.
 */

func differentColsAllDiff(queens []core.VarId, store *core.Store) {
	prop := propagator.CreateAlldifferent(queens...)
	store.AddPropagators(prop)
}

// nqueens with alldifferent on columns but primitive on diagonals
func NQueensAllDiffCols(store *core.Store, N int) []core.VarId {
	return nqueensConstrain(store, N,
		differentColsAllDiff, differentDiagPrim)
}

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

func NQueensAllDiffBoth(store *core.Store, N int) []core.VarId {
	return nqueensConstrain(store, N,
		differentColsAllDiff,    // reuse the columns constraint
		differentDiagAuxAllDiff) // new one for diagonals
}

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

func NQueensAllDiffOnly(store *core.Store, N int) []core.VarId {
	return nqueensConstrain(store, N,
		differentColsAllDiff, // reuse column constraints
		differentDiagAllDiff, // new dedicated offset ones for diagonals
	)
}
