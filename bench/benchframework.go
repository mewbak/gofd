package main

// Mini bench framework benching single functions. We use the signature of the
// golang test/bench framework, but run the individual tests manually.
// The functions to bench must be called in the main function.
// We can run a single bench with a main function with
// $ go run benchframework.go bench_<name>.go
// in this directory.

// A benchmark function may start with any initialization code
// whose execution time is ignored. The following two lines
//    BenchSetRuns(b.N)
//    b.StartTimer() // benchmark starts here
// must precede the actual code to be benchmarked, which is
// typically iterated in a loop. This almost always forces a third line
//    for i:=0; i < b.N; i++ {
// and the code in the loop then is measured.
// Benchmarks shall not be executed concurrently.
// Do not forget to provide information about the benchmark run stored
// in a key/value map bc of strings. Set at least key "name",
// which will appear in the bench logs as task, name. Use other keys,
// which will all end up under the "task, " prefix and may be used
// to select specific bench runs from the command line.

import (
	"bitbucket.org/gofd/gofd/core"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"runtime/pprof"
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

func insert_separators(s string, every int, separator string) string {
	a := make([]string, (len(s)+(every-1))/every)
	initial_skip := len(s) % every
	if initial_skip == 0 {
		initial_skip = every
	}
	a[0] = s[0:initial_skip]
	for i := 1; i < len(a); i++ {
		a[i] = s[(i-1)*every+initial_skip : i*every+initial_skip]
	}
	return strings.Join(a, separator)
}

const bench_verbose = true

// teardown shall be called after any test
func bench_teardown(br testing.BenchmarkResult) {
	benchmark_config[TIMEPERRUN] = fmt.Sprintf("%d", br.NsPerOp())
	if bench_verbose {
		fmt.Printf("%-22s: ",
			benchmark_config[TASKNAME])
		for key, value := range benchmark_config {
			if strings.HasPrefix(key, TASK) && key != TASKNAME {
				fmt.Printf("%s=%-10v  ", key[len(TASK)+2:], value)
			}
		}
		fmt.Printf("%s: %13s ns\n", TIMEPERRUN,
			insert_separators(benchmark_config[TIMEPERRUN], 3, "_"))
		// fmt.Printf("%v\n", benchmark_config[TIMEPERRUN])
	}
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

// bench_addsystem adds runtime, operating system and
// hardware specific information to the log file
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
	// choose a unique filename (if benchs are not executed
	// concurrently and run longer than a nanosecond)
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

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

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
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
	}
	result := testing.Benchmark(f)
	if *cpuprofile != "" {
		pprof.StopCPUProfile()
	}
	bench_teardown(result)
}

// checkRunIt checks whether first arg is a substring of the
// function name and all args are key=val pairs that as soon as
// the corresponding key exists, the value must be a perfect match
func checkRunIt(task bc, args []string) bool {
	if len(args) == 0 {
		return true // run it if there are no args
	}
	name := task["name"]
	if !strings.Contains(name, args[0]) {
		return false // function name does not fit
	}
	for _, kv := range args[1:] {
		if strings.Contains(kv, "=") {
			kva := strings.SplitN(kv, "=", 2)
			key := kva[0]
			value := kva[1]
			if val, ok := task[key]; ok {
				if val != value {
					return false // key=value did not match exactly
				}
			} // otherwise ignore filter
		} // otherwise ignore filter
	}
	return true // if none of the filters rule it out, run it
}

// run the benchmark if matched
func benchd(f func(b *testing.B), task bc) {
	flag.Parse()
	if checkRunIt(task, flag.Args()) {
		var bcg bc = bc{"": "gofd"}
		var bce bc = bc{}
		bench(f, task, bcg, bce)
	}
}
