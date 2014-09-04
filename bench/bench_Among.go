package main

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"bitbucket.org/gofd/gofd/propagator/explicit"
	"testing"
)

// the driver for everything benching IntVar
func bench_Among() {
	benchd(bAmong1, bc{"name": "Among", "size": "1"})
}

func bAmong1(b *testing.B) {
	xi := make([][]int, 4)
	xi[0] = []int{2, 3}
	xi[1] = []int{4, 5}
	xi[2] = []int{4, 5}
	xi[3] = []int{5, 6}

	expxi := make([][]int, 4)
	expxi[0] = []int{2, 3}
	expxi[1] = []int{5}
	expxi[2] = []int{5}
	expxi[3] = []int{5, 6}

	k := []int{2, 3, 4}
	n := []int{1}

	bAmong(b, xi, k, n, expxi, n)
}

func bAmong(b *testing.B, xiinit [][]int, kinit []int,
	ninit []int, expxi [][]int, expn []int) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		store := core.CreateStore()
		Xi := make([]core.VarId, len(xiinit))
		for i := 0; i < len(xiinit); i++ {
			Xi[i] = core.CreateAuxIntVarValues(store, xiinit[i])
		}
		N := core.CreateIntVarValues("N", store, ninit)
		store.AddPropagators(explicit.CreateAmong(Xi, kinit, N))
		query := labeling.CreateSearchOneQuery()
		labeling.LabelingSplit(store, query, labeling.SmallestDomainFirst)
	}
}
