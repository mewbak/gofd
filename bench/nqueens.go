package main

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"bitbucket.org/gofd/gofd/propagator/explicit"
	"bitbucket.org/gofd/gofd/propagator/interval"
	"strconv"
	"testing"
)

//-- Helper-class --

//------------ NQueens IvDomain --------------

func bIvDNQueensImpl(b *testing.B, queensCount int) {
	bIvDNDameAllDiff2(b, queensCount)
}

func bIvDNDameAllDiff2(b *testing.B, queensCount int) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		store := core.CreateStoreWithoutLogging()
		queens := createIvDQueens(queensCount, store)
		ivd_differentColsCTRAllDiff2(queens, store)
		ivd_differentDiagCTR(queens, store)
		labeling.SetAllvars(queens)
		query := labeling.CreateSearchAllQuery()
		labeling.Labeling(store, query, labeling.VarSelect, labeling.InDomainMin)
	}
}

func createIvDQueens(countDames int, store *core.Store) []core.VarId {
	dames := make([]core.VarId, countDames)
	for i := 0; i < countDames; i++ {
		varname := "Q" + strconv.Itoa(i)
		dames[i] = core.CreateIntVarIvFromTo(varname, store, 0, countDames-1)
	}
	return dames
}

func ivd_differentColsCTRAllDiff2(queens []core.VarId, store *core.Store) {
	prop := interval.CreateAlldifferent(queens...)
	store.AddPropagator(prop)
}

func ivd_differentDiagCTR(queens []core.VarId, store *core.Store) {
	for i := 0; i < len(queens)-1; i++ {
		remaining := queens[i:]
		for offset := 1; offset < len(remaining); offset++ {
			ivd_checkOffset(remaining, store, offset)
			ivd_checkOffset(remaining, store, -offset)
		}
	}
}

func ivd_checkOffset(queens []core.VarId, store *core.Store, offset int) {
	headQueen := queens[0]
	index := offset
	if index < 0 {
		index = -index
	}
	prop := interval.CreateXplusCneqY(headQueen, offset, queens[index])
	store.AddPropagator(prop)
}

//------------ NQueens Domain --------------

func bNQueensImpl(b *testing.B, nDame int) {
	bNDameAllDiff2(b, nDame)
}

func bNDameAllDiff2(b *testing.B, to int) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		store := core.CreateStoreWithoutLogging()
		queens := createQueens(to, store)
		differentColsCTRAllDiff2(queens, store)
		differentDiagCTR(queens, store)
		labeling.SetAllvars(queens)
		query := labeling.CreateSearchAllQuery()
		labeling.Labeling(store, query, labeling.VarSelect, labeling.InDomainMin)
	}
}

func createQueens(countDames int, store *core.Store) []core.VarId {
	dames := make([]core.VarId, countDames)
	for i := 0; i < countDames; i++ {
		varname := "Q" + strconv.Itoa(i)
		dames[i] = core.CreateIntVarFromTo(varname, store, 0, countDames-1)
	}
	return dames
}

func differentColsCTRAllDiff2(queens []core.VarId, store *core.Store) {
	prop := explicit.CreateAlldifferent_Primitives(queens...)
	store.AddPropagator(prop)
}

func differentDiagCTR(queens []core.VarId, store *core.Store) {
	for i := 0; i < len(queens)-1; i++ {
		remaining := queens[i:]
		for offset := 1; offset < len(remaining); offset++ {
			checkOffset(remaining, store, offset)
			checkOffset(remaining, store, -offset)
		}
	}
}

func checkOffset(queens []core.VarId, store *core.Store, offset int) {
	headQueen := queens[0]
	index := offset
	if index < 0 {
		index = -index
	}
	prop := explicit.CreateXplusCneqY(headQueen, offset, queens[index])
	store.AddPropagator(prop)
}
