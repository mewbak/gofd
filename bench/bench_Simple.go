package main

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/explicit"
	"testing"
)

func main() {
	bench_Simple()
}

func bench_Simple() {
	name := "Simple"
	benchd(bSimple1, bc{"name": name, "size": "1"})
}

func bSimple1(b *testing.B) { bSimple(b, 1) }

func bSimple(b *testing.B, to int) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		store := core.CreateStore()
		X := core.CreateIntVarFromTo("X", store, 0, 9)
		Y := core.CreateIntVarFromTo("Y", store, 0, 9)
		prop1 := explicit.CreateC1XplusC2YeqC3(1, X, 1, Y, 9)
		store.AddPropagator(prop1)
		prop2 := explicit.CreateC1XplusC2YeqC3(2, X, 4, Y, 24)
		store.AddPropagator(prop2)
		store.IsConsistent()
	}
}
