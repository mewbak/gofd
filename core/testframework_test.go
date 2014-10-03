package core

// Mini test framework providing a setup and teardown
// DO NOT COPY TO OTHER PACKAGES. Copy from e.g. propagator.
// Two lines at the beginning of every test
//    setup()
//    defer teardown()
// Use log("message", ...) to force log messages at info level

import (
	"fmt"
	"time"
)

var store *Store
var test_logger *Logger
var test_time = time.Now()
var test_counter int = 0

const TEST_VERBOSE = true

// setup shall be called before any test
func setup() {
	test_counter += 1
	test_time = time.Now()
	if TEST_VERBOSE {
		fmt.Printf("%3d >>> \n", test_counter)
	}
	test_logger = GetLogger()
	test_logger.SetLoggingLevel(LOG_ERROR)
	store = CreateStore()
}

// teardown shall be called after any test
func teardown() {
	dur := time.Now().Sub(test_time)
	if TEST_VERBOSE {
		fmt.Printf("%3d <<< %s\n", test_counter, dur.String())
	}
}

// log forces a log message with TEST out
func log(msgs ...string) {
	for _, msg := range msgs {
		test_logger.P(fmt.Sprintf("TEST: %v\n", msg))
	}
}
