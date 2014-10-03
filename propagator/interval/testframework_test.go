package interval

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

func setup() {
	if TEST_PARALLEL {
		runtime.GOMAXPROCS(runtime.NumCPU()) // use the cores
	}
	test_counter += 1
	test_time = time.Now()
	if TEST_VERBOSE {
		fmt.Printf("%3d >>> \n", test_counter)
	}
	test_logger = core.GetLogger()
	test_logger.SetLoggingLevel(core.LOG_ERROR)
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
func domainEquals_test(t *testing.T, test string,
	id core.VarId, values []int) {
	want := core.CreateExDomainAdds(values)
	got := store.GetDomain(id)
	if !got.Equals(want) {
		msg := "%s %s: domain got %s, want %s"
		t.Errorf(msg, test, store.GetName(id), got.String(), want.String())
	}
}

func DomainEquals_test(t *testing.T, test string, id core.VarId,
	want core.Domain) {
	got := store.GetDomain(id)
	if !got.Equals(want) {
		msg := "%s %s: domain got %s, want %s"
		t.Errorf(msg, test, store.GetName(id), got.String(), want.String())
	}
}

// equalsInt_test checks two int-values, if they are equal.
func equalsInt_test(t *testing.T, test string, val int, expval int) {
	if val != expval {
		msg := "%s: value got %v, want %v"
		t.Errorf(msg, test, val, expval)
	}
}

// equalsInt_test checks, if the number of results, is equal to the expected
// number of results
func result_count_test(t *testing.T, test string,
	resultSet map[int]map[core.VarId]int, expCount int) {
	if len(resultSet) != expCount {
		msg := "%s: result_count got %v, want %v"
		t.Errorf(msg, test, len(resultSet), expCount)
	}
}

// resultSet_test checks, if the result-row (result_id) in resultSet
// is as expected. Respective, if the fixed value for a specific IntVar
// (var_id) is equal to an expected value (expValue). So it checks,
// if resultSet[result_id][var_id] == expValue
func resultSet_test(t *testing.T, test string,
	resultSet map[int]map[core.VarId]int, result_id int,
	var_id core.VarId, expValue int) {
	if len(resultSet) < result_id {
		msg := "%s: no such result_id (number of results to low)"
		t.Errorf(msg, test)
	}
	if resultSet[result_id][var_id] != expValue {
		t.Errorf("resultSet_test: got %v, want %v",
			resultSet[result_id][var_id], expValue)
	}
}

// ready_test checks if the store answers with the expected result
// (true: solution found, false: propagation failed)
func ready_test(t *testing.T, test string, ready bool, expready bool) {
	if ready != expready {
		msg := "%s: ready got %v, want %v"
		t.Errorf(msg, test, ready, expready)
	}
}

// log forces a log message with TEST out
func log(msgs ...string) {
	for _, msg := range msgs {
		test_logger.P(fmt.Sprintf("TEST: %v\n", msg))
	}
}

// clone_test test the correctness of the clone function of a given constraint
// c1, which contains variables created in a given store store1. It clones
// the original store to store2 and the given constraint to c2, then adds
// c2 to store2 and checks, if the resulting domains in store1 and store2
// are the same after propagating.
func clone_test(t *testing.T, store1 *core.Store, c1 core.Constraint) {
	store2 := store1.Clone(nil)
	c2 := c1.Clone()
	store1.AddPropagator(c1)
	store2.AddPropagator(c2)
	ready1 := store1.IsConsistent()
	ready2 := store2.IsConsistent()
	if ready1 != ready2 {
		msg := "%s: clone test ready got %v, want %v"
		t.Errorf(msg, c1, ready2, ready1)
	}
	varids := store1.GetVariableIDs()
	for _, vid := range varids {
		d1 := store1.GetDomain(vid)
		d2 := store2.GetDomain(vid)
		if !d1.Equals(d2) {
			msg := "%s: clone test domain got %s, want %s"
			t.Errorf(msg, c1, d2.String(), d1.String())
		}
	}
}

func createXYZtestVars(xinit []int, yinit []int,
	zinit []int) (core.VarId, core.VarId, core.VarId) {
	X := core.CreateIntVarIvValues("X", store, xinit)
	Y := core.CreateIntVarIvValues("Y", store, yinit)
	Z := core.CreateIntVarIvValues("Z", store, zinit)

	return X, Y, Z
}

func createXYtestVars(xinit []int,
	yinit []int) (core.VarId, core.VarId) {
	X := core.CreateIntVarIvValues("X", store, xinit)
	Y := core.CreateIntVarIvValues("Y", store, yinit)
	return X, Y
}

func createTestVars(inits [][]int, names []string) []core.VarId {
	varids := make([]core.VarId, len(inits))
	for i, _ := range inits {
		varids[i] = core.CreateIntVarIvValues(names[i], store, inits[i])
	}
	return varids
}
