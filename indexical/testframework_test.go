package indexical

// Mini test framework providing a setup and teardown
// Should be copied to every package.
// Two lines at the beginning of every test
//    setup()
//    defer teardown()

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"runtime"
	"testing"
	"time"
)

var test_logger *core.Logger
var store *core.Store
var test_time = time.Now()
var test_counter int = 0

var X_ForCloneTest core.VarId
var Y_ForCloneTest core.VarId

const TEST_VERBOSE = true
const TEST_PARALLEL = false

func init() {
	test_logger = core.GetLogger()
	test_logger.SetLoggingLevel(core.LOG_ERROR)
}

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

func teardown() {
	dur := time.Now().Sub(test_time)
	if TEST_VERBOSE {
		log(store.GetStat().String())
		fmt.Printf("%3d <<< %s\n", test_counter, dur.String())
	}
}

// domainEquals_test checks whether the domain of a variable (id) is the same
// as a given set of values
func domainEquals_test(t *testing.T, test string, id core.VarId, values []int) {
	want := core.CreateIvDomainFromIntArr(values)
	got := store.GetDomain(id)
	if !got.Equals(want) {
		msg := "%s %s: Domain calculated = %v, want %v"
		t.Errorf(msg, test, store.GetName(id), got.String(), want.String())
	}
}

func DomainEquals_test(t *testing.T, test string, id core.VarId,
	want core.Domain) {
	got := store.GetDomain(id)
	if !got.Equals(want) {
		msg := "%s %s: Domain calculated = %v, want %v"
		t.Errorf(msg, test, store.GetName(id), got.String(), want.String())
	}
}

// equalsInt_test checks two int-values, if they are equal.
func equalsInt_test(t *testing.T, test string, val int, expval int) {
	if val != expval {
		msg := "%s: value calculated = %v, want %v"
		t.Errorf(msg, test, val, expval)
	}
}

// equalsInt_test checks, if the number of results, is equal to the expected
// number of results
func result_count_test(t *testing.T, test string, resultSet map[int]map[core.VarId]int, expCount int) {
	if len(resultSet) != expCount {
		msg := "%s: result_count = %v, want %v"
		t.Errorf(msg, test, len(resultSet), expCount)
	}
}

// resultSet_test checks, if the result-row (result_id) in resultSet
// is as expected. Respective, if the fixed value for a specific IntVar (var_id)
// is equal to an expected value (expValue). So it checks,
// if resultSet[result_id][var_id] == expValue
func resultSet_test(t *testing.T, test string, resultSet map[int]map[core.VarId]int, result_id int, var_id core.VarId, expValue int) {
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
func ready_test(t *testing.T, test string, ready bool, expready bool) {
	if ready != expready {
		msg := "%s: ready = %v, want %v"

		t.Errorf(msg, test, ready, expready)
	}
}

// log forces a log message with TEST out
func log(msgs ...string) {
	for _, msg := range msgs {
		test_logger.P(fmt.Sprintf("TEST: %v\n", msg))
	}
}
