package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"fmt"
	"testing"
)

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
// parametrized by size and a function that generates the variables as
// well as the constraints that there must be different columns and
// different diagonals.
func nqueensall(t *testing.T, testname string,
	constrain func(store *core.Store, N int) []core.VarId,
	N int) {
	msg := fmt.Sprintf("nqueens%d%s", N, testname)
	log(msg)
	// create the finite domain vars and generate constraints in the store
	constrain(store, N) // return value, varids queens not needed
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
	if logger.GetLoggingLevel() <= core.LOG_NONE {
		searchStat(query.GetSearchStatistics())
		if logger.GetLoggingLevel() <= core.LOG_INFO {
			show_nqueens_results(store, query.GetResultSet())
		}
	}
}

// helper to show the solutions
func show_nqueens_results(store *core.Store,
	resultSet map[int]map[core.VarId]int) {
	for i, result := range resultSet {
		fmt.Printf("Sol %d:", i)
		for _, varId := range core.SortedKeys_MapVarIdToInt(result) {
			fmt.Printf(" %s=%d",
				store.GetName(varId), result[varId])
		}
		fmt.Printf("\n")
	}
}

// helper to show one set of variables
func show_nqueens(store *core.Store,
	queens []core.VarId) {
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

// tests with primitive constraints
func Test_nqueens5Prim(t *testing.T) {
	setup()
	defer teardown()
	nqueensall(t, "Prim", NQueensPrim, 5)
}

func Test_nqueens7Prim(t *testing.T) {
	setup()
	defer teardown()
	nqueensall(t, "Prim", NQueensPrim, 7)
}

func Test_nqueens8Prim(t *testing.T) {
	setup()
	defer teardown()
	nqueensall(t, "Prim", NQueensPrim, 8)
}

// tests with alldifferent constraint in the columns
func Test_nqueens7AllDiffCols(t *testing.T) {
	setup()
	defer teardown()
	nqueensall(t, "AllDiffCols", NQueensAllDiffCols, 7)
}

func Test_nqueens8AllDiffCols(t *testing.T) {
	setup()
	defer teardown()
	nqueensall(t, "AllDiffCols", NQueensAllDiffCols, 8)
}

// tests with alldifferent constaint in both column and diagonal
func Test_nqueens7AllDiffBoth(t *testing.T) {
	setup()
	defer teardown()
	nqueensall(t, "AllDiffBoth", NQueensAllDiffBoth, 7)
}

func Test_nqueens8AllDiffBoth(t *testing.T) {
	setup()
	defer teardown()
	nqueensall(t, "AllDiffBoth", NQueensAllDiffBoth, 8)
}

// tests with only alldifferent constraints
func Test_nqueens7AllDiffOnly(t *testing.T) {
	setup()
	defer teardown()
	nqueensall(t, "AllDiffOnly", NQueensAllDiffOnly, 7)
}

func Test_nqueens8AllDiffOnly(t *testing.T) {
	setup()
	defer teardown()
	nqueensall(t, "AllDiffOnly", NQueensAllDiffOnly, 8)
}

func Test_nqueens9AllDiffOnly(t *testing.T) {
	setup()
	defer teardown()
	nqueensall(t, "AllDiffOnly", NQueensAllDiffOnly, 9)
}

// Todo: tests with Alldistinct instead of Alldifferent
