package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"bitbucket.org/gofd/gofd/propagator/interval"
	"fmt"
	"strings"
	"testing"
)

func showSolutions(sol []bool, weights []int, maxValue int,
	statStore *core.StoreStatistics,
	statSearch *labeling.SearchStatistics) {
	s := make([]string, 0)
	for i, val := range sol {
		if val {
			s = append(s, fmt.Sprintf("%4d", i))
		}
	}
	msg := "Take items: %s    giving value %d"
	log(fmt.Sprintf(msg, strings.Join(s, ", "), maxValue))
	log("SearchStat: " + statSearch.SearchString())
	log("StoreStat : " + statSearch.StoreString())
}

func Test_knapsack0a(t *testing.T) {
	setupWithoutStore()
	defer teardown()
	weights := []int{4, 3}
	values := []int{2, 1}
	capacity := 3
	testKnapsack(t, "knapsack0a", weights, values, capacity, 1)
}

func Test_knapsack0b(t *testing.T) {
	setupWithoutStore()
	defer teardown()
	weights := []int{4, 3}
	values := []int{2, 1}
	capacity := 4
	testKnapsack(t, "knapsack0b", weights, values, capacity, 2)
}

func Test_knapsack1a(t *testing.T) {
	setupWithoutStore()
	defer teardown()
	weights := []int{4, 3, 2}
	values := []int{15, 10, 7}
	capacity := 9
	testKnapsack(t, "knapsack1a", weights, values, capacity, 32)
}

func Test_knapsack1b(t *testing.T) {
	setupWithoutStore()
	defer teardown()
	weights := []int{4, 3, 2}
	values := []int{15, 10, 7}
	capacity := 7
	testKnapsack(t, "knapsack1b", weights, values, capacity, 25)
}

func Test_knapsack1c(t *testing.T) {
	setupWithoutStore()
	defer teardown()
	weights := []int{4, 3, 2}
	values := []int{15, 10, 7}
	capacity := 6
	testKnapsack(t, "knapsack1c", weights, values, capacity, 22)
}

func Test_knapsack1d(t *testing.T) {
	setupWithoutStore()
	defer teardown()
	weights := []int{4, 3, 2}
	values := []int{15, 10, 7}
	capacity := 5
	testKnapsack(t, "knapsack1d", weights, values, capacity, 17)
}

func Test_knapsack1e(t *testing.T) {
	setupWithoutStore()
	defer teardown()
	weights := []int{4, 3, 2}
	values := []int{15, 10, 7}
	capacity := 4
	testKnapsack(t, "knapsack1e", weights, values, capacity, 15)
}

func Test_knapsack2a(t *testing.T) {
	setupWithoutStore()
	defer teardown()
	weights := []int{580, 1616, 1906, 1942, 50, 294}
	values := []int{874, 620, 345, 369, 360, 470}
	capacity := 2000
	testKnapsack(t, "knapsack2a", weights, values, capacity, 1704)
}

func Test_knapsack2b(t *testing.T) {
	setupWithoutStore()
	defer teardown()
	weights := []int{946, 859, 147, 512, 656, 328, 625, 585, 145, 754}
	values := []int{986, 292, 867, 720, 734, 888, 585, 541, 300, 475}
	capacity := 1000
	testKnapsack(t, "knapsack2b", weights, values, capacity, 2475)
}

func Test_knapsack2c(t *testing.T) {
	setupWithoutStore()
	defer teardown()
	weights := []int{61, 425, 95, 464, 470, 155, 17, 309, 124, 121}
	values := []int{683, 996, 61, 562, 652, 826, 833, 784, 539, 46}
	capacity := 1000
	testKnapsack(t, "knapsack2c", weights, values, capacity, 4122)
}

func testKnapsack(t *testing.T, testname string,
	weights, values []int, capacity, sol int) {
	lastSol := make([]bool, len(values))
	maxValue, statSearch, statStore :=
		doKnapsack(t, testname, weights, values, capacity, lastSol)
	equalsInt_test(t, "Test_"+testname, maxValue, sol)
	showSolutions(lastSol, weights, maxValue, statStore, statSearch)
}

func doKnapsack(t *testing.T, testname string,
	weights, values []int, capacity int,
	lastSol []bool) (int, *labeling.SearchStatistics, *core.StoreStatistics) {
	first := true
	result := true
	curValue := -1
	searchStats := labeling.CreateSearchStatistics()
	for result { // increase bound incrementally until unsatisfiable
		store := core.CreateStore()
		// fmt.Printf("knapsack with bound %d\n", curValue)
		take := setupKnapsack(store, weights, values, capacity)
		store.AddPropagator(interval.CreateXgtC(take[len(take)-1], curValue))
		if first {
			log(fmt.Sprintf("%s: capacity %d, %d items",
				testname, capacity, len(weights)))
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
			first = false
		}
		// fmt.Printf("  setup\n")
		query := labeling.CreateSearchOneQuery()
		result = labeling.Labeling(store, query)
		stats := query.GetSearchStatistics()
		searchStats.UpdateSearchStatistics(stats)
		// fmt.Printf("  %s\n", searchStats.GetStoreStatistics().String())
		resultSet := query.GetResultSet()
		if result {
			resetLastSol(lastSol)
			lentakem1 := len(take) - 1
			solValue := 0
			for i := 0; i < lentakem1; i++ {
				if resultSet[0][take[i]] == 1 {
					lastSol[i] = true
					solValue += values[i]
				}
			}
			curValue = solValue // lift to value of current assignment
			// fmt.Printf("  solution value: %d\n", curValue)
		}
	}
	return curValue, searchStats, searchStats.GetStoreStatistics()
}

func resetLastSol(lastSol []bool) {
	for i, _ := range lastSol {
		lastSol[i] = false
	}
}

func setupKnapsack(store *core.Store, weights, values []int,
	capacity int) []core.VarId {
	lenValues := len(values)
	take := make([]core.VarId, lenValues+1)
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
	// fmt.Printf("    takes\n")
	take[len(take)-1] = core.CreateIntVarIvFromTo("sumValues",
		store, 0, maxSumValues)
	// fmt.Printf("    var sumvalues\n")
	sumprop := interval.CreateSumBounds(store,
		take[len(take)-1], sumValues)
	// fmt.Printf("    prop sumvalues\n")
	store.AddPropagators(sumprop)
	// fmt.Printf("    added sumvalues\n")
	sumWeight := core.CreateIntVarIvFromTo("sumWeights", store, 0,
		maxSumWeights)
	// fmt.Printf("    var sumweights\n")
	sumprop = interval.CreateSumBounds(store, sumWeight, sumWeights)
	// fmt.Printf("    prop sumweights\n")
	store.AddPropagators(sumprop)
	// fmt.Printf("    added sumweights\n")
	store.AddPropagator(interval.CreateXlteqC(sumWeight, capacity))
	// for _, vid := range store.GetVariableIDs() {
	// 	 println(vid, store.GetDomain(vid).String())
	// }
	// println("DUMP VARS")
	return take
}
