package demo

// Mini test framework providing a setup and teardown.
// Two lines at the beginning of every test
//    setup()
//    defer teardown()

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"fmt"
	"runtime"
	"testing"
	"time"
)

var store *core.Store
var store_ready chan bool
var logger *core.Logger
var test_time = time.Now()
var test_counter int = 0

const TEST_VERBOSE = true
const TEST_PARALLEL = true

func init() {
	logger = core.GetLogger()
	logger.SetLoggingLevel(core.LOG_ERROR)
}

// provides in addition to book keeping a logger and running store
func setup() {
	if TEST_PARALLEL {
		runtime.GOMAXPROCS(runtime.NumCPU()) // use the cores
	}
	test_counter += 1
	test_time = time.Now()
	if TEST_VERBOSE {
		fmt.Printf("%3d >>> \n", test_counter)
	}
	store = core.CreateStore()
}

func setupWithoutStore() {
	if TEST_PARALLEL {
		runtime.GOMAXPROCS(runtime.NumCPU()) // use the cores
	}
	test_counter += 1
	test_time = time.Now()
	if TEST_VERBOSE {
		fmt.Printf("%3d >>> \n", test_counter)
	}
	logger = core.GetLogger()
	logger.SetLoggingLevel(core.LOG_ERROR)
}

func teardown() {
	dur := time.Now().Sub(test_time)
	if TEST_VERBOSE {
		fmt.Printf("%3d <<< %s\n", test_counter, dur.String())
	}
}

// log forces a log message with TEST out
func log(msgs ...string) {
	for _, msg := range msgs {
		logger.P(fmt.Sprintf("TEST: %v\n", msg))
	}
}

func searchStat(searchStatistics *labeling.SearchStatistics) {
	log(searchStatistics.SearchString())
	log(searchStatistics.StoreString())
}

func propStat() {
	log(store.GetStat().String())
}

// domainEquals_test checks whether the domain of a variable (id) is the same
// as a given set of values
func domainEquals_test(t *testing.T,
	test string, id core.VarId, values []int) {
	want := core.CreateExDomainAdds(values)
	got := store.GetDomain(id)
	if !got.Equals(want) {
		msg := "%s %s: Domain calculated = %v, want %v"
		t.Errorf(msg, test, core.GetNameRegistry().GetName(id), got, want)
	}
}

// equalsInt_test checks two int-values, if they are equal.
func equalsInt_test(t *testing.T, test string, val int, expval int) {
	if val != expval {
		msg := "%s: value calculated = %v, want %v"
		t.Errorf(msg, test, val, expval)
	}
}

// result_count_test checks, if the number of results, is equal to
// the expected number of results
func result_count_test(t *testing.T,
	test string, resultSet map[int]map[core.VarId]int, expCount int) {
	if len(resultSet) != expCount {
		msg := "%s: result_count = %v, want %v"
		t.Errorf(msg, test, len(resultSet), expCount)
	}
}

// resultSet_test checks, if the result-row (result_id) in resultSet is as
// expected. Respective, if the fixed value for a specific IntVar (var_id)
// is equal to an expected value (expValue). So it checks,
// if resultSet[result_id][var_id] == expValue
func resultSet_test(t *testing.T,
	test string, resultSet map[int]map[core.VarId]int,
	result_id int, var_id core.VarId, expValue int) {
	if len(resultSet) < result_id {
		msg := "%s: no such result_id (number of results to low)"
		t.Errorf(msg, test)
	}
	if resultSet[result_id][var_id] != expValue {
		t.Errorf("chicken and rabbit: result = %v, want %v",
			resultSet[result_id][var_id], expValue)
	}
}

// ready_test checks if the store answers with the expected result
// (true: solution found, false: propagation failed)
func ready_test(t *testing.T, test string, ready bool, expready bool) bool {
	if ready != expready {
		msg := "%s: ready = %v, want %v"
		t.Errorf(msg, test, ready, expready)
		return false
	}
	return true
}

func rangestep(from, to, step int) []int {
	noelems := ((to - from) + step) / step
	steps := make([]int, noelems)
	i := 0
	for value := from; value <= to; value += step {
		steps[i] = value
		i += 1
	}
	return steps
}

/// Example
//func Test_Aaa(t *testing.T) {
//	setup()
//	defer teardown()
//	log("Testframework-example-test")
//}
