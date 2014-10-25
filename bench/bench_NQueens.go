package main

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/demo"
	"bitbucket.org/gofd/gofd/labeling"
	"testing"
)

func main() {
	name := "NQueens"
	typ := "prim"
	benchd(b6QueensPrim, bc{"name": name, "type": typ, "size": "6"})
	benchd(b7QueensPrim, bc{"name": name, "type": typ, "size": "7"})
	benchd(b8QueensPrim, bc{"name": name, "type": typ, "size": "8"})
	typ = "ad_col"
	benchd(b7QueensADCol, bc{"name": name, "type": typ, "size": "7"})
	benchd(b8QueensADCol, bc{"name": name, "type": typ, "size": "8"})
	benchd(b9QueensADCol, bc{"name": name, "type": typ, "size": "9"})
	typ = "ad_both"
	benchd(b7QueensADBoth, bc{"name": name, "type": typ, "size": "7"})
	benchd(b8QueensADBoth, bc{"name": name, "type": typ, "size": "8"})
	benchd(b9QueensADBoth, bc{"name": name, "type": typ, "size": "9"})
	typ = "ad_only"
	benchd(b7QueensADOnly, bc{"name": name, "type": typ, "size": "7"})
	benchd(b8QueensADOnly, bc{"name": name, "type": typ, "size": "8"})
	benchd(b9QueensADOnly, bc{"name": name, "type": typ, "size": "9"})
}

func b6QueensPrim(b *testing.B) { bNQueens(b, demo.NQueensPrim, 6) }
func b7QueensPrim(b *testing.B) { bNQueens(b, demo.NQueensPrim, 7) }
func b8QueensPrim(b *testing.B) { bNQueens(b, demo.NQueensPrim, 8) }

func b7QueensADCol(b *testing.B) { bNQueens(b, demo.NQueensAllDiffCols, 7) }
func b8QueensADCol(b *testing.B) { bNQueens(b, demo.NQueensAllDiffCols, 8) }
func b9QueensADCol(b *testing.B) { bNQueens(b, demo.NQueensAllDiffCols, 9) }

func b7QueensADBoth(b *testing.B) { bNQueens(b, demo.NQueensAllDiffBoth, 7) }
func b8QueensADBoth(b *testing.B) { bNQueens(b, demo.NQueensAllDiffBoth, 8) }
func b9QueensADBoth(b *testing.B) { bNQueens(b, demo.NQueensAllDiffBoth, 9) }

func b7QueensADOnly(b *testing.B) { bNQueens(b, demo.NQueensAllDiffOnly, 7) }
func b8QueensADOnly(b *testing.B) { bNQueens(b, demo.NQueensAllDiffOnly, 8) }
func b9QueensADOnly(b *testing.B) { bNQueens(b, demo.NQueensAllDiffOnly, 9) }

func bNQueens(b *testing.B,
	constrain func(store *core.Store, N int) []core.VarId,
	N int) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		store := core.CreateStoreWithoutLogging()
		queens := constrain(store, N)
		labeling.SetAllvars(queens)
		query := labeling.CreateSearchAllQuery()
		labeling.Labeling(store, query,
			labeling.VarSelect, labeling.InDomainMin)
	}
}
