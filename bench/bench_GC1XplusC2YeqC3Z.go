package main

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator"
	"testing"
)

func main() {
	bench_GC1XplusC2YeqC3Z()
}

// the driver for everything benching IntVar
func bench_GC1XplusC2YeqC3Z() {
	benchd(bGC1XplusC2YeqC3Z1, bc{"name": "GC1XplusC2YeqC3Z", "size": "1"})
}

func bGC1XplusC2YeqC3Z1(b *testing.B) { bGC1XplusC2YeqC3Z(b, 1) }

func bGC1XplusC2YeqC3Z(b *testing.B, to int) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		store := core.CreateStore()
		xinit := []int{0, 1, 2, 3, 4}
		yinit := []int{0, 1, 2, 3, 4}
		zinit := []int{6, 8, 9}
		X := core.CreateIntVarValues("X", store, xinit)
		Y := core.CreateIntVarValues("Y", store, yinit)
		Z := core.CreateIntVarValues("Z", store, zinit)
		c1, c2, c3 := 1, 1, 1
		store.AddPropagator(propagator.CreateC1XplusC2YeqC3ZBounds(c1, X, c2, Y, c3, Z))
		store.IsConsistent()
	}
}
