package main

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/explicit"
	"bitbucket.org/gofd/gofd/propagator/interval"
	"testing"
)

func main() {
	bench_XmultCeqY()
}

// the driver for everything benching XmultCeqY
func bench_XmultCeqY() {
	benchd(b_XmultCeqYex, bc{"name": "explicit.XmultCeqY", "size": "1"})
	benchd(b_XmultCeqYiv, bc{"name": "interval.XmultCeqY", "size": "1"})
}

func b_XmultCeqYex(b *testing.B) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		store := core.CreateStore()
		X := core.CreateIntVarFromTo("X", store, 100, 10000)
		Y := core.CreateIntVarFromTo("Y", store, 100, 1000)
		c := 10
		store.AddPropagator(explicit.CreateXmultCeqY(X, c, Y))
		store.IsConsistent()
	}
}

func b_XmultCeqYiv(b *testing.B) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		store := core.CreateStore()
		X := core.CreateIntVarFromTo("X", store, 100, 10000)
		Y := core.CreateIntVarFromTo("Y", store, 100, 1000)
		c := 10
		store.AddPropagator(interval.CreateXmultCeqY(X, c, Y))
		store.IsConsistent()
	}
}
