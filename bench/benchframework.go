package main

// Mini bench framework benching single functions. We use the signature of the
// golang test/bench framework, but run run the individual tests ourselves.
// The functions to bench are organized in bench.go (the main/driver). We have
// run everything with go run *.go in this directory.

// A benchmark function may start with any initialization code whose execution
// time is ignored. The following two lines
//    BenchSetRuns(b.N)
//    b.StartTimer() // benchmark starts here
// must precede the actual code to be benchmarked, which is typically iterated #
// in a loop. This almost always forces a third line
//    for i:=0; i < b.N; i++ {
// and the code in the loop then is measured. Benchmarks shall not be executed
// concurrently.
// Do not forget to provide information about the benchmark run to be stored
// through a key/value map bc of strings. Set at least key "name", which will
// appear in the bench logs as task, name. Use other keys which will all end up
// under the "task, " prefix.

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"sort"
	"strings"
	"testing"
	"time"
)

const DIR = "logs"
const CREATED = "created"
const RUNS = "runs"
const TIMEPERRUN = "timePerRun"
const TASK = "task"
const TASKNAME = "task, name"
const IMPL = "impl"

type bc map[string]string // bench config

var store *core.Store
var bench_logger *core.Logger
var bench_counter int = 0

// global variable, but ok since benchmarks are run sequentially
var benchmark_config bc

// bench_setup shall be called before any test
func bench_setup(c map[string]string) {
	bench_counter += 1
	bench_logger = core.GetLogger()
	bench_logger.SetLoggingLevel(core.LOG_ERROR)

	store = core.CreateStoreWithoutLogging()
	benchmark_config = bc{} // new
	for k, v := range c {   // copy
		benchmark_config[k] = v
	}
	benchmark_config["runs"] = "1" // default, also means not available
}

// teardown shall be called after any test
func bench_teardown(br testing.BenchmarkResult) {
	benchmark_config[TIMEPERRUN] = fmt.Sprintf("%d", br.NsPerOp())
	bench_write(benchmark_config)
}

func bench_writeremove(f *os.File, c bc, key string) {
	// the key itself is part of the map
	if value, ok := c[key]; ok {
		f.WriteString(fmt.Sprintf("%s: %s\n", key, value))
		delete(c, key)
		return
	}
	// sort the keys
	keys := make([]string, len(c))
	i := 0
	for k, _ := range c {
		keys[i] = k
		i += 1
	}
	sort.Strings(keys)
	// print all matching keys in alphabetical order
	found := false
	for _, k := range keys {
		if strings.HasPrefix(k, key) {
			value := c[k]
			f.WriteString(fmt.Sprintf("%s: %s\n", k, value))
			delete(c, k)
			found = true
		}
	}
	if found || key == "" {
		return
	}
	panic(fmt.Sprintf("Did not find mandatory key %v in %v", key, c))
}

// adds system specific information
func bench_addsystem(c bc) {
	c["lang, name"] = "golang"
	c["lang, compiler"] = runtime.Compiler
	c["lang, version"] = runtime.Version()
	c["hardware, arch"] = runtime.GOARCH
	c["hardware, cpus"] = fmt.Sprintf("%d", runtime.NumCPU())
	c["os, name"] = runtime.GOOS
	proccpuinfo, err := ioutil.ReadFile("/proc/cpuinfo")
	if err == nil { // ok, can get more info about the cpu
		lines := strings.Split(string(proccpuinfo), "\n")
		for _, line := range lines {
			info := strings.Split(line, ":")
			if len(info) > 1 {
				info[0] = strings.TrimSpace(info[0])
				info[1] = strings.TrimSpace(info[1])
				if info[0] == "model" {
					c["cpu, model"] = info[1]
				}
				if strings.HasPrefix(info[0], "stepping") {
					c["cpu, stepping"] = info[1]
				}
				if strings.HasPrefix(info[0], "cpu family") {
					c["cpu, family"] = info[1]
				}
				if strings.HasPrefix(info[0], "cpu MHz") {
					c["cpu, MHz"] = info[1]
				}
				if strings.HasPrefix(info[0], "model name") {
					c["cpu, name"] = info[1]
				}
				if strings.HasPrefix(info[0], "bogomips") {
					c["cpu, bogomips"] = info[1]
				}
				if strings.HasPrefix(info[0], "vendor_id") {
					c["cpu, vendor_id"] = info[1]
				}
			}
		}
	}
	meminfo, err := ioutil.ReadFile("/proc/meminfo")
	if err == nil { // ok, can get more info about memory
		lines := strings.Split(string(meminfo), "\n")
		for _, line := range lines {
			info := strings.Split(line, ":")
			if len(info) > 1 {
				info[0] = strings.TrimSpace(info[0])
				info[1] = strings.TrimSpace(info[1])
				if info[0] == "MemTotal" {
					c["hardware, memory"] = info[1]
				}
			}
		}
	}
	versioninfo, err := ioutil.ReadFile("/proc/version_signature")
	if err == nil { // ok
		c["os, version"] = strings.TrimSpace(string(versioninfo))
	}
}

// http://stackoverflow.com/questions/10510691/how-to-check-whether-a-file-or-directory-denoted-by-a-path-exists-in-golang
func exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func bench_write(c bc) {
	bench_time := time.Now().UnixNano()
	c[CREATED] = fmt.Sprintf("%d", bench_time/1000000)
	bench_addsystem(c) // adds system information
	if !exists(DIR) {
		err := os.Mkdir(DIR, 0775)
		if err != nil {
			panic(err)
		}
	}
	filename := fmt.Sprintf("%s.%d.bench", c[TASKNAME], bench_time)
	benchfile, err := os.Create(path.Join(DIR, filename))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := benchfile.Close(); err != nil {
			panic(err)
		}
	}()
	bench_writeremove(benchfile, c, CREATED)
	bench_writeremove(benchfile, c, RUNS)
	bench_writeremove(benchfile, c, TIMEPERRUN)
	bench_writeremove(benchfile, c, TASKNAME)
	bench_writeremove(benchfile, c, TASK)
	bench_writeremove(benchfile, c, IMPL)
	bench_writeremove(benchfile, c, "") // rest
}

// BenchSetRuns stores the runs, which is unfortunately not
// available in BenchmarkResult. Thus, this function has to
// be called in the benchmark function itself.
func BenchSetRuns(n int) {
	benchmark_config[RUNS] = fmt.Sprintf("%d", n)
}

// run function with config for task, impl and more
func bench(f func(b *testing.B), task bc, impl bc, more bc) {
	c := make(bc)
	for k, v := range task {
		if !strings.HasPrefix(k, TASK) {
			k = TASK + ", " + k
		}
		c[k] = v
	}
	for k, v := range impl {
		if !strings.HasPrefix(k, IMPL) {
			if len(k) > 0 {
				k = IMPL + ", " + k
			} else {
				k = IMPL
			}
		}
		c[k] = v
	}
	for k, v := range more {
		c[k] = v
	}
	bench_setup(c)
	bench_teardown(testing.Benchmark(f))
}

// default for impl and more
func benchd(f func(b *testing.B), task bc) {
	var bcg bc = bc{"": "gofd"}
	var bce bc = bc{}
	bench(f, task, bcg, bce)
}
