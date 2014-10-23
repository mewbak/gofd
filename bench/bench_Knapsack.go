package main

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/demo"
	"bitbucket.org/gofd/gofd/labeling"
	"testing"
)

func main() {
	benchd(bKnapsack0a, bc{"name": "Knapsack", "variant": "0a"})
	benchd(bKnapsack0b, bc{"name": "Knapsack", "variant": "0b"})
	benchd(bKnapsack1a, bc{"name": "Knapsack", "variant": "1a"})
	benchd(bKnapsack1b, bc{"name": "Knapsack", "variant": "1b"})
	benchd(bKnapsack1c, bc{"name": "Knapsack", "variant": "1c"})
	benchd(bKnapsack1d, bc{"name": "Knapsack", "variant": "1d"})
	benchd(bKnapsack1e, bc{"name": "Knapsack", "variant": "1e"})
	benchd(bKnapsack2a, bc{"name": "Knapsack", "variant": "2a"})
	benchd(bKnapsack2b, bc{"name": "Knapsack", "variant": "2b"})
	benchd(bKnapsack2c, bc{"name": "Knapsack", "variant": "2c"})
}

func bKnapsack0a(b *testing.B) { bKnapsack(b, "0a") }
func bKnapsack0b(b *testing.B) { bKnapsack(b, "0b") }
func bKnapsack1a(b *testing.B) { bKnapsack(b, "1a") }
func bKnapsack1b(b *testing.B) { bKnapsack(b, "1b") }
func bKnapsack1c(b *testing.B) { bKnapsack(b, "1c") }
func bKnapsack1d(b *testing.B) { bKnapsack(b, "1d") }
func bKnapsack1e(b *testing.B) { bKnapsack(b, "1e") }
func bKnapsack2a(b *testing.B) { bKnapsack(b, "2a") }
func bKnapsack2b(b *testing.B) { bKnapsack(b, "2b") }
func bKnapsack2c(b *testing.B) { bKnapsack(b, "2c") }

func bKnapsack(b *testing.B, name string) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		store := core.CreateStoreWithoutLogging()
		ks := demo.Knapsack(name)
		_, objective := demo.ConstrainKnapsack(store,
			ks.Weights, ks.Values, ks.Capacity)
		// need variant with fixed order of variable selection
		labeling.Maximize(store, objective)
	}
}
