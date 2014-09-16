package main

// Sample benchmark to evaluate the performance of primitives
import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func main() {
	bench_IntVar()
}

// The driver for everything benching IntVar.
func bench_IntVar() {
	// We run 5 benchmarks in different configurations
	// Each run is represented by one function
	// First, two Clone benchmark runs
	name := "IntVar.Clone"
	benchd(bIntVarClone10, bc{"name": name, "size": "10"})
	benchd(bIntVarClone100, bc{"name": name, "size": "100"})
	// Second, three IsGround benchmark runs
	name = "IntVar.IsGround"
	benchd(bIntVarIsGround1, bc{"name": name, "size": "1"})
	benchd(bIntVarIsGround10, bc{"name": name, "size": "10"})
	benchd(bIntVarIsGround100, bc{"name": name, "size": "100"})
}

// First, two Clone benchmark runs
func bIntVarClone10(b *testing.B)  { bIntVarClone(b, 10) }  // with 10
func bIntVarClone100(b *testing.B) { bIntVarClone(b, 100) } // with 100
// values in the finite domain. Cloning should grow linearly with
// the domain size

func bIntVarClone(b *testing.B, to int) {
	// init, executed once per benchmark run
	store := core.CreateStoreWithoutLogging()
	XId := core.CreateIntVarFromTo("X", store, 1, to)
	X, _ := store.GetIntVar(XId)
	// bench, executed b.N times per benchmark run to fill up time
	BenchSetRuns(b.N)
	b.StartTimer()             // benchmark (and timing!) starts here
	for i := 0; i < b.N; i++ { // the loop you must provide
		X.Clone() // the code you want to bench
	}
}

// Second, three IsGround benchmark runs
func bIntVarIsGround1(b *testing.B)   { bIntVarIsGround(b, 1) }
func bIntVarIsGround10(b *testing.B)  { bIntVarIsGround(b, 10) }
func bIntVarIsGround100(b *testing.B) { bIntVarIsGround(b, 100) }

func bIntVarIsGround(b *testing.B, to int) {
	// init
	store := core.CreateStoreWithoutLogging()
	XId := core.CreateIntVarFromTo("X", store, 1, to)
	X, _ := store.GetIntVar(XId)
	// bench
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		X.IsGround()
	}
}
