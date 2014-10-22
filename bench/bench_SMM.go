package main

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"bitbucket.org/gofd/gofd/propagator/indexical"
	"bitbucket.org/gofd/gofd/propagator/interval"
	"testing"
)

func main() {
	benchd(bSMM,
		bc{"name": "SendMoreMoney", "type": "normal"})
	benchd(bSMMIndexical,
		bc{"name": "SendMoreMoney", "type": "indexical"})
}

func bSMM(b *testing.B) { SMMImpl(b) }

func SMMImpl(b *testing.B) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		store := core.CreateStoreWithoutLogging()
		S := core.CreateIntVarFromTo("S", store, 1, 9)
		E := core.CreateIntVarFromTo("E", store, 0, 9)
		N := core.CreateIntVarFromTo("N", store, 0, 9)
		D := core.CreateIntVarFromTo("D", store, 0, 9)
		M := core.CreateIntVarFromTo("M", store, 1, 9)
		O := core.CreateIntVarFromTo("O", store, 0, 9)
		R := core.CreateIntVarFromTo("R", store, 0, 9)
		Y := core.CreateIntVarFromTo("Y", store, 0, 9)
		alldiff_prop := interval.CreateAlldifferent(S, E, N, D, M, O, R, Y)
		store.AddPropagators(alldiff_prop)
		max := 1000*store.GetDomain(S).GetMax() +
			100*store.GetDomain(E).GetMax() +
			10*store.GetDomain(N).GetMax() +
			1*store.GetDomain(D).GetMax() +
			1000*store.GetDomain(M).GetMax() +
			100*store.GetDomain(O).GetMax() +
			10*store.GetDomain(R).GetMax() +
			1*store.GetDomain(E).GetMax()
		sum := core.CreateIntVarFromTo("sum", store, 0, max)
		weights := []int{1000, 100, 10, 1, 1000, 100, 10, 1}
		vars := []core.VarId{S, E, N, D, M, O, R, E}
		sendmore_props := interval.CreateWeightedSumBounds(store,
			sum, weights, vars...)
		store.AddPropagators(sendmore_props)
		weights = []int{10000, 1000, 100, 10, 1}
		vars = []core.VarId{M, O, N, E, Y}
		money_props := interval.CreateWeightedSumBounds(store,
			sum, weights, vars...)
		store.AddPropagators(money_props)
		labeling.SetAllvars([]core.VarId{S, E, N, D, M, O, R, Y})
		query := labeling.CreateSearchOneQuery()
		labeling.Labeling(store, query,
			labeling.InDomainMin, labeling.VarSelect)
	}
}

func bSMMIndexical(b *testing.B) { SMMIndexicalImpl(b) }

func SMMIndexicalImpl(b *testing.B) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		store := core.CreateStoreWithoutLogging()
		S := core.CreateIntVarFromTo("S", store, 1, 9)
		E := core.CreateIntVarFromTo("E", store, 0, 9)
		N := core.CreateIntVarFromTo("N", store, 0, 9)
		D := core.CreateIntVarFromTo("D", store, 0, 9)
		M := core.CreateIntVarFromTo("M", store, 1, 9)
		O := core.CreateIntVarFromTo("O", store, 0, 9)
		R := core.CreateIntVarFromTo("R", store, 0, 9)
		Y := core.CreateIntVarFromTo("Y", store, 0, 9)
		alldiff_prop := indexical.CreateAlldifferent(S, E, N, D, M, O, R, Y)
		store.AddPropagators(alldiff_prop)
		max := 1000*store.GetDomain(S).GetMax() +
			100*store.GetDomain(E).GetMax() +
			10*store.GetDomain(N).GetMax() +
			1*store.GetDomain(D).GetMax() +
			1000*store.GetDomain(M).GetMax() +
			100*store.GetDomain(O).GetMax() +
			10*store.GetDomain(R).GetMax() +
			1*store.GetDomain(E).GetMax()
		sum := core.CreateIntVarFromTo("sum", store, 0, max)
		weights := []int{1000, 100, 10, 1, 1000, 100, 10, 1}
		vars := []core.VarId{S, E, N, D, M, O, R, E}
		sendmore_props := indexical.CreateWeightedSum(store,
			sum, weights, vars...)
		store.AddPropagators(sendmore_props)
		weights = []int{10000, 1000, 100, 10, 1}
		vars = []core.VarId{M, O, N, E, Y}
		money_props := indexical.CreateWeightedSum(store, sum, weights, vars...)
		store.AddPropagators(money_props)
		labeling.SetAllvars([]core.VarId{S, E, N, D, M, O, R, Y})
		query := labeling.CreateSearchOneQuery()
		labeling.Labeling(store, query,
			labeling.InDomainMin, labeling.VarSelect)
	}
}
