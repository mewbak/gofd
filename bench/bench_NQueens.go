package main

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"bitbucket.org/gofd/gofd/propagator"
	"bitbucket.org/gofd/gofd/propagator/interval"
	"fmt"
	"testing"
)

func main() {
	name := "NQueens"
	benchd(b8Queens, bc{"name": name, "size": "8"})
	benchd(b9Queens, bc{"name": name, "size": "9"})
}

func b8Queens(b *testing.B) { bNQueensImpl(b, 8) }
func b9Queens(b *testing.B) { bNQueensImpl(b, 9) }

func bNQueensImpl(b *testing.B, queensCount int) {
	bNQueensAllDiff(b, queensCount)
}

func bNQueensAllDiff(b *testing.B, queensCount int) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		store := core.CreateStoreWithoutLogging()
		queens := createQueens(queensCount, store)
		differentColsAllDiff(queens, store)
		differentDiagAllDiff(queens, store)
		labeling.SetAllvars(queens)
		query := labeling.CreateSearchAllQuery()
		labeling.Labeling(store, query, labeling.VarSelect, labeling.InDomainMin)
	}
}

func createQueens(n int, store *core.Store) []core.VarId {
	queens := make([]core.VarId, n)
	for i := 0; i < n; i++ {
		varname := fmt.Sprintf("Q%d", i)
		queens[i] = core.CreateIntVarFromTo(varname, store, 0, n-1)
	}
	return queens
}

func differentColsAllDiff(queens []core.VarId, store *core.Store) {
	prop := propagator.CreateAlldifferent(queens...)
	store.AddPropagators(prop)
}

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
