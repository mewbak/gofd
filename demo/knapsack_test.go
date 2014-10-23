package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"fmt"
	"strings"
	"testing"
)

func testDefinedKnapsack(t *testing.T, name string) {
	ks := Knapsack(name)
	if ks == nil {
		panic(fmt.Sprintf("Problem %s not found", name))
	}
	testKnapsack(t, "knapsack"+ks.name,
		ks.weights, ks.values, ks.capacity, ks.maxValue)
}

func Test_knapsack0a(t *testing.T) {
	setup()
	defer teardown()
	testDefinedKnapsack(t, "0a")
}

func Test_knapsack0b(t *testing.T) {
	setup()
	defer teardown()
	testDefinedKnapsack(t, "0b")
}

func Test_knapsack1a(t *testing.T) {
	setup()
	defer teardown()
	testDefinedKnapsack(t, "1a")
}

func Test_knapsack1b(t *testing.T) {
	setup()
	defer teardown()
	testDefinedKnapsack(t, "1b")
}

func Test_knapsack1c(t *testing.T) {
	setup()
	defer teardown()
	testDefinedKnapsack(t, "1c")
}

func Test_knapsack1d(t *testing.T) {
	setup()
	defer teardown()
	testDefinedKnapsack(t, "1d")
}

func Test_knapsack1e(t *testing.T) {
	setup()
	defer teardown()
	testDefinedKnapsack(t, "1e")
}

func Test_knapsack2a(t *testing.T) {
	setup()
	defer teardown()
	testDefinedKnapsack(t, "2a")
}

func Test_knapsack2b(t *testing.T) {
	setup()
	defer teardown()
	testDefinedKnapsack(t, "2b")
}

func Test_knapsack2c(t *testing.T) {
	setup()
	defer teardown()
	testDefinedKnapsack(t, "2b")
}

/* helper */

func showProblem(name string, weights, values []int, capacity int) {
	log(fmt.Sprintf("%s: capacity %d, %d items",
		name, capacity, len(weights)))
	indexes := make([]int, len(weights))
	for i, _ := range indexes {
		indexes[i] = i
	}
	sindexes := core.IntSliceToStringSliceFormatted(indexes, "%4d")
	log(fmt.Sprintf("index     : %s", strings.Join(sindexes, ", ")))
	sweights := core.IntSliceToStringSliceFormatted(weights, "%4d")
	log(fmt.Sprintf("weights   : %s", strings.Join(sweights, ", ")))
	svalues := core.IntSliceToStringSliceFormatted(values, "%4d")
	log(fmt.Sprintf("values    : %s", strings.Join(svalues, ", ")))
}

func showSolution(solution map[core.VarId]int, vars []core.VarId, weights []int, maxValue int,
	statStore *core.StoreStatistics,
	statSearch *labeling.SearchStatistics) {
	s := make([]string, 0)
	for _, takeVar := range vars {
		if solution[takeVar] == 1 {
			s = append(s, fmt.Sprintf("%4s", store.GetName(takeVar)[1:]))
		}
	}
	msg := "Take items: %s    giving value %d"
	log(fmt.Sprintf(msg, strings.Join(s, ", "), maxValue))
	log("SearchStat: " + statSearch.SearchString())
	log("StoreStat : " + statSearch.StoreString())
}

func testKnapsack(t *testing.T, testname string,
	weights, values []int, capacity, sol int) {
	showProblem(testname, weights, values, capacity)
	vars, objectiveVar := ConstrainKnapsack(store, weights, values, capacity)
	maxVal, solution, searchStat, storeStat := labeling.MaximizeStats(store, objectiveVar)
	equalsInt_test(t, "Test_"+testname, maxVal, sol)
	showSolution(solution, vars, weights, maxVal, storeStat, searchStat)
}
