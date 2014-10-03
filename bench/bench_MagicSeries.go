package main

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"bitbucket.org/gofd/gofd/propagator/explicit"
	"bitbucket.org/gofd/gofd/propagator/indexical"
	"bitbucket.org/gofd/gofd/propagator/interval"
	"bitbucket.org/gofd/gofd/propagator/reification"
	"testing"
)

func main() {
	bench_MagicSeries()
	bench_MagicSeriesWithoutAmong()
}

// the driver for everything benching IntVar
func bench_MagicSeries() {
	benchd(bMagicSeries1, bc{"name": "MagicSeries", "size": "3"})
	benchd(bMagicSeries2, bc{"name": "MagicSeries", "size": "4"})
	benchd(bMagicSeries3, bc{"name": "MagicSeries", "size": "5"})
	benchd(bMagicSeries4, bc{"name": "MagicSeries", "size": "6"})
	benchd(bMagicSeries5, bc{"name": "MagicSeries", "size": "7"})
	benchd(bMagicSeries6, bc{"name": "MagicSeries", "size": "10"})
	benchd(bMagicSeries7, bc{"name": "MagicSeries", "size": "17"})
}

func bMagicSeries1(b *testing.B) { bMagicSeries(b, 3) }
func bMagicSeries2(b *testing.B) { bMagicSeries(b, 4) }
func bMagicSeries3(b *testing.B) { bMagicSeries(b, 5) }
func bMagicSeries4(b *testing.B) { bMagicSeries(b, 6) }
func bMagicSeries5(b *testing.B) { bMagicSeries(b, 7) }
func bMagicSeries6(b *testing.B) { bMagicSeries(b, 10) }
func bMagicSeries7(b *testing.B) { bMagicSeries(b, 17) }

func bMagicSeries(b *testing.B, length int) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	var query *labeling.SearchOneQuery
	for i := 0; i < b.N; i++ {
		store := core.CreateStore()
		// define variables X0,...,Xn
		variables := make([]core.VarId, length+1)
		for i := 0; i < len(variables); i++ {
			variables[i] = core.CreateAuxIntVarExFromTo(store, 0, length)
		}

		// define constraints
		// each value j can occur Xj times
		for i := 0; i < len(variables); i++ {
			store.AddPropagator(explicit.CreateAmong(variables, []int{i}, variables[i]))
		}
		//		println("numberOfPropagators magic among: ", store.GetNumPropagators())
		query = labeling.CreateSearchOneQueryVariableSelect(variables)
		labeling.Labeling(store, query, labeling.VarSelect, labeling.InDomainMin)
	}
	println("among magic:", length, "nodes:", query.GetSearchStatistics().GetNodes())
}

func bench_MagicSeriesWithoutAmong() {
	benchd(bMagicSeriesWithoutAmong1, bc{"name": "MagicSeriesWithoutAmong", "size": "3"})
	benchd(bMagicSeriesWithoutAmong2, bc{"name": "MagicSeriesWithoutAmong", "size": "4"})
	benchd(bMagicSeriesWithoutAmong3, bc{"name": "MagicSeriesWithoutAmong", "size": "5"})
	benchd(bMagicSeriesWithoutAmong4, bc{"name": "MagicSeriesWithoutAmong", "size": "6"})
	benchd(bMagicSeriesWithoutAmong5, bc{"name": "MagicSeriesWithoutAmong", "size": "7"})
	benchd(bMagicSeriesWithoutAmong6, bc{"name": "MagicSeriesWithoutAmong", "size": "10"})
	benchd(bMagicSeriesWithoutAmong7, bc{"name": "MagicSeriesWithoutAmong", "size": "17"})
}

func bMagicSeriesWithoutAmong1(b *testing.B) { bMagicSeriesWithoutAmong(b, 3) }
func bMagicSeriesWithoutAmong2(b *testing.B) { bMagicSeriesWithoutAmong(b, 4) }
func bMagicSeriesWithoutAmong3(b *testing.B) { bMagicSeriesWithoutAmong(b, 5) }
func bMagicSeriesWithoutAmong4(b *testing.B) { bMagicSeriesWithoutAmong(b, 6) }
func bMagicSeriesWithoutAmong5(b *testing.B) { bMagicSeriesWithoutAmong(b, 7) }
func bMagicSeriesWithoutAmong6(b *testing.B) { bMagicSeriesWithoutAmong(b, 10) }
func bMagicSeriesWithoutAmong7(b *testing.B) { bMagicSeriesWithoutAmong(b, 17) }

func bMagicSeriesWithoutAmong(b *testing.B, n int) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	var query *labeling.SearchOneQuery
	for i := 0; i < b.N; i++ {
		store := core.CreateStore()
		//define variables X0,...,Xn
		variables := make([]core.VarId, n+1)
		for i := 0; i < len(variables); i++ {
			variables[i] = core.CreateAuxIntVarIvFromTo(store, 0, n)
		}

		// define constraints
		// each value j can occur Xj times
		for i := 0; i < len(variables); i++ {
			//array for reified constraints
			variables := make([]core.VarId, len(variables))
			for j := 0; j < len(variables); j++ {
				//store in variables[j] whether Xj (variables[j]) takes the value i or not
				variables[j] = core.CreateAuxIntVarIvFromTo(store, 0, 1)
				xeqc := indexical.CreateXeqC(variables[j], i)
				reifiedConstraint := reification.CreateReifiedConstraint(xeqc, variables[j])
				store.AddPropagator(reifiedConstraint)
			}
			// the amount of variables in X0,...,Xn that have taken the value i
			// must correspond to Xi (variables[i])
			store.AddPropagator(interval.CreateSum(store, variables[i], variables))
		}
		query = labeling.CreateSearchOneQueryVariableSelect(variables)
		labeling.Labeling(store, query, labeling.VarSelect, labeling.InDomainMin)
	}
	println("primitive magic:", n, "nodes:", query.GetSearchStatistics().GetNodes())
}
