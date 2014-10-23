package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/interval"
	"fmt"
)

// KnapsackProblem represents a knapsack problem
type KnapsackProblem struct {
	name     string
	weights  []int
	values   []int
	capacity int
	maxValue int
}

var KnapsackProblems []KnapsackProblem = []KnapsackProblem{
	KnapsackProblem{"0a", []int{4, 3}, []int{2, 1}, 3, 1},
	KnapsackProblem{"0b", []int{4, 3}, []int{2, 1}, 4, 2},
	KnapsackProblem{"1a", []int{4, 3, 2}, []int{15, 10, 7}, 9, 32},
	KnapsackProblem{"1b", []int{4, 3, 2}, []int{15, 10, 7}, 7, 25},
	KnapsackProblem{"1c", []int{4, 3, 2}, []int{15, 10, 7}, 6, 22},
	KnapsackProblem{"1d", []int{4, 3, 2}, []int{15, 10, 7}, 5, 17},
	KnapsackProblem{"1e", []int{4, 3, 2}, []int{15, 10, 7}, 4, 15},
	KnapsackProblem{"2a",
		[]int{580, 1616, 1906, 1942, 50, 294}, // weights
		[]int{874, 620, 345, 369, 360, 470},   // values
		2000,  // capacity
		1704}, // maxValue
	KnapsackProblem{"2b",
		[]int{946, 859, 147, 512, 656, 328, 625, 585, 145, 754},
		[]int{986, 292, 867, 720, 734, 888, 585, 541, 300, 475},
		1000,
		2475},
	KnapsackProblem{"2c",
		[]int{61, 425, 95, 464, 470, 155, 17, 309, 124, 121},
		[]int{683, 996, 61, 562, 652, 826, 833, 784, 539, 46},
		1000,
		4122},
}

// Knapsack returns one of the predefined knapsack problems
// the one matching the name exactly
func Knapsack(name string) *KnapsackProblem {
	for _, ks := range KnapsackProblems {
		if ks.name == name {
			return &ks
		}
	}
	return nil
}

// ConstrainKnapsack create a knapsack problem and
// returns a slice of varIds that represents for each
// item whether it should be take or not and the
// sum of values which can be used to maximize
func ConstrainKnapsack(store *core.Store,
	weights, values []int,
	capacity int) ([]core.VarId, core.VarId) {
	lenValues := len(values)
	if len(weights) != lenValues {
		msg := "len(weights) = %d != %d = len(values)"
		panic(fmt.Sprintf(msg, len(weights), lenValues))
	}
	take := make([]core.VarId, lenValues)
	sumWeights := make([]core.VarId, lenValues)
	sumValues := make([]core.VarId, lenValues)
	maxSumWeights := sum_intarray(weights)
	maxSumValues := sum_intarray(values)
	propsWeights := make([]core.Constraint, lenValues)
	propsValues := make([]core.Constraint, lenValues)
	for i := 0; i < lenValues; i++ {
		stri := fmt.Sprintf("T%d", i)
		take[i] = core.CreateIntVarIvFromTo(stri, store, 0, 1)
		sumWeights[i] = core.CreateIntVarIvValues("sumweights_"+stri,
			store, []int{0, weights[i]})
		propsWeights[i] = interval.CreateXmultCeqY(take[i],
			weights[i], sumWeights[i])
		sumValues[i] = core.CreateIntVarIvValues("sumvalues_"+stri,
			store, []int{0, values[i]})
		propsValues[i] = interval.CreateXmultCeqY(take[i],
			values[i], sumValues[i])
	}
	store.AddPropagators(propsWeights...)
	store.AddPropagators(propsValues...)
	value := core.CreateIntVarIvFromTo("sumValues",
		store, 0, maxSumValues)
	sumprop := interval.CreateSumBounds(store, value, sumValues)
	store.AddPropagators(sumprop)
	sumWeight := core.CreateIntVarIvFromTo("sumWeights", store, 0,
		maxSumWeights)
	sumprop = interval.CreateSumBounds(store, sumWeight, sumWeights)
	store.AddPropagators(sumprop)
	store.AddPropagator(interval.CreateXlteqC(sumWeight, capacity))
	return take, value
}
