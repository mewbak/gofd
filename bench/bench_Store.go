package main

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator"
	"runtime"
	"testing"
)

func main() {
	bench_Store()
}

// the driver for everything benching IntVar
func bench_Store() {
	name := "Clone"
	// Too much memory...
	//benchd(bStoreClone1, bc{"name": name, "size": "1"})
	//benchd(bStoreClone10, bc{"name": name, "size": "10"})
	//benchd(bStoreClone100, bc{"name": name, "size": "100"})
	name = "AddSimplePropagator"
	benchd(bStoreAddSimplePropagator1, bc{"name": name, "size": "1"})
	benchd(bStoreAddSimplePropagator10, bc{"name": name, "size": "10"})
	benchd(bStoreAddSimplePropagator100, bc{"name": name, "size": "100"})
	name = "AddComplexPropagator"
	benchd(bStoreAddComplexPropagator1, bc{"name": name, "size": "1"})
	benchd(bStoreAddComplexPropagator10, bc{"name": name, "size": "10"})
	benchd(bStoreAddComplexPropagator100, bc{"name": name, "size": "100"})
}

func bStoreClone1(b *testing.B)   { bStoreClone(b, 1) }
func bStoreClone10(b *testing.B)  { bStoreClone(b, 10) }
func bStoreClone100(b *testing.B) { bStoreClone(b, 100) }

func bStoreClone(b *testing.B, to int) {
	// init
	c := 5
	for j := 0; j < to; j++ {
		XId := core.CreateAuxIntVarFromTo(store, 1, 10)
		xgtc := propagator.CreateXgtC(XId, c)
		store.AddPropagator(xgtc)
	}
	// bench
	BenchSetRuns(b.N)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		store.Clone(nil)
	}
	b.StopTimer()
	// deinit, reclaim memory
	runtime.GC() // not working good, as the store leaks..
}

func bStoreAddSimplePropagator1(b *testing.B) {
	bStoreAddSimplePropagator(b, 1)
}

func bStoreAddSimplePropagator10(b *testing.B) {
	bStoreAddSimplePropagator(b, 10)
}

func bStoreAddSimplePropagator100(b *testing.B) {
	bStoreAddSimplePropagator(b, 100)
}

func bStoreAddSimplePropagator(b *testing.B, to int) {
	//init
	c := 5
	stores := make([]*core.Store, b.N)
	for i := 0; i < b.N; i++ {
		stores[i] = core.CreateStoreWithoutLogging()
	}
	props := make([]core.Constraint, to)
	// bench
	BenchSetRuns(b.N)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		store := stores[i]
		for j := 0; j < to; j++ {
			XId := core.CreateAuxIntVarFromTo(store, 1, 10)
			props[j] = propagator.CreateXgtC(XId, c)
		}
		store.AddPropagators(props...)
		stores[i].IsConsistent()
	}
}

func bStoreAddComplexPropagator1(b *testing.B) {
	bStoreAddComplexPropagator(b, 1)
}

func bStoreAddComplexPropagator10(b *testing.B) {
	bStoreAddComplexPropagator(b, 10)
}

func bStoreAddComplexPropagator100(b *testing.B) {
	bStoreAddComplexPropagator(b, 100)
}

func bStoreAddComplexPropagator(b *testing.B, to int) {
	//init
	stores := make([]*core.Store, b.N)
	for i := 0; i < b.N; i++ {
		stores[i] = core.CreateStoreWithoutLogging()
	}
	props := make([]core.Constraint, to)
	xinit := []int{0, 1, 2, 3, 4}
	yinit := []int{0, 1, 2, 3, 4}
	zinit := []int{6, 8, 9, 16}
	// bench
	BenchSetRuns(b.N)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		store := stores[i]
		for j := 0; j < to; j++ {
			XId := core.CreateAuxIntVarValues(store, xinit)
			YId := core.CreateAuxIntVarValues(store, yinit)
			ZId := core.CreateAuxIntVarValues(store, zinit)
			props[j] = propagator.CreateXmultYeqZ(XId, YId, ZId)
		}
		store.AddPropagators(props...)
	}
	store.IsConsistent()
}
