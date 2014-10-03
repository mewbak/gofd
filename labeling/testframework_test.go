package labeling

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

var store *core.Store
var logger *core.Logger
var test_time = time.Now()
var test_counter int = 0

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
	logger = core.GetLogger()
	logger.SetLoggingLevel(core.LOG_ERROR)
	store = core.CreateStore()
}

func teardown() {
	dur := time.Now().Sub(test_time)
	if TEST_VERBOSE {
		fmt.Printf("%3d <<< %s\n", test_counter, dur.String())
	}
}

func searchStat(searchStatistics *SearchStatistics) {
	log(searchStatistics.SearchString())
	log(searchStatistics.StoreString())
}

// log forces a log message with TEST out
func log(msgs ...string) {
	for _, msg := range msgs {
		logger.P(fmt.Sprintf("TEST: %v\n", msg))
	}
}
