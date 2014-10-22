package main

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
	"strconv"
	"testing"
)

func main() {
	bench_Interval_Bad()
	bench_Interval_Good()
	bench_Interval_Trend()
}

// constants for the tests
const removeVal int = 5000
const addVal = 13000
const containsVal = 5000

var removeVals = []int{5000, 10000}
var badIntervals [][]int
var goodIntervals [][]int

// make sure all "constants" are initialised
func init() {
	badIntervals = make([][]int, 0)
	for i := 0; i <= 80000; i += 2 {
		badIntervals = append(badIntervals, []int{i, i})
	}
	goodIntervals = make([][]int, 0)
	goodIntervals = append(goodIntervals, []int{0, 10000})
	goodIntervals = append(goodIntervals, []int{15000, 80000})
}

// the driver for everything benching Interval_Bad
func bench_Interval_Bad() {
	name := fmt.Sprintf("Remove(%d)", removeVal)
	benchd(bExRemoveBad, bc{"name": name, "type": "bad", "domain": "Ex"})
	name = fmt.Sprintf("Remove(%d)", removeVal)
	benchd(bIvRemoveBad, bc{"name": name, "type": "bad", "domain": "Iv"})
	name = fmt.Sprintf("Removes(%v)", removeVals)
	benchd(bExRemovesBad, bc{"name": name, "type": "bad", "domain": "Ex"})
	name = fmt.Sprintf("Removes(%v)", removeVals)
	benchd(bIvRemovesBad, bc{"name": name, "type": "bad", "domain": "Iv"})
	name = fmt.Sprintf("Add(%d)", addVal)
	benchd(bExAddBad, bc{"name": name, "type": "bad", "domain": "Ex"})
	name = fmt.Sprintf("Add(%d)", addVal)
	benchd(bIvAddBad, bc{"name": name, "type": "bad", "domain": "Iv"})
	name = fmt.Sprintf("Contains(%d)", containsVal)
	benchd(bExContainsBad, bc{"name": name, "type": "bad", "domain": "Ex"})
	name = fmt.Sprintf("Contains(%d", containsVal)
	benchd(bIvContainsBad, bc{"name": name, "type": "bad", "domain": "Iv"})
	benchd(bExCopyBad, bc{"name": "Copy", "type": "bad", "domain": "Ex"})
	benchd(bIvCopyBad, bc{"name": "Copy", "type": "bad", "domain": "Iv"})
	benchd(bExMinBad, bc{"name": "GetMin", "type": "bad", "domain": "Ex"})
	benchd(bIvMinBad, bc{"name": "GetMin", "type": "bad", "domain": "Iv"})
	benchd(bExMaxBad, bc{"name": "GetMax", "type": "bad", "domain": "Ex"})
	benchd(bIvMaxBad, bc{"name": "GetMax", "type": "bad", "domain": "Iv"})
	benchd(bExIsEmptyBad,
		bc{"name": "IsEmpty", "type": "bad", "domain": "Ex"})
	benchd(bIvIsEmptyBad,
		bc{"name": "IsEmpty", "type": "bad", "domain": "Iv"})
}

func bExRemoveBad(b *testing.B) {
	bExRemoveImpl(b, badIntervals, removeVal)
}
func bIvRemoveBad(b *testing.B) {
	bIvRemoveImpl(b, badIntervals, removeVal)
}

func bExRemovesBad(b *testing.B) {
	bExRemovesImpl(b, badIntervals, removeVals)
}
func bIvRemovesBad(b *testing.B) {
	bIvRemovesImpl(b, badIntervals, removeVals)
}

func bExAddBad(b *testing.B) {
	bExAddImpl(b, badIntervals, addVal)
}
func bIvAddBad(b *testing.B) {
	bIvAddImpl(b, badIntervals, addVal)
}

func bExContainsBad(b *testing.B) {
	bExContainsImpl(b, badIntervals, containsVal)
}
func bIvContainsBad(b *testing.B) {
	bIvContainsImpl(b, badIntervals, containsVal)
}

func bExCopyBad(b *testing.B) {
	bExCopyImpl(b, badIntervals)
}
func bIvCopyBad(b *testing.B) {
	bIvCopyImpl(b, badIntervals)
}

func bExIsEmptyBad(b *testing.B) {
	bExIsEmptyImpl(b, badIntervals)
}
func bIvIsEmptyBad(b *testing.B) {
	bIvIsEmptyImpl(b, badIntervals)
}

func bExMinBad(b *testing.B) {
	bExMinImpl(b, badIntervals)
}
func bIvMinBad(b *testing.B) {
	bIvMinImpl(b, badIntervals)
}

func bExMaxBad(b *testing.B) {
	bExMaxImpl(b, badIntervals)
}
func bIvMaxBad(b *testing.B) {
	bIvMaxImpl(b, badIntervals)
}

// the driver for everything benching Interval_Good
func bench_Interval_Good() {
	name := fmt.Sprintf("Remove(%d)", removeVal)
	benchd(bExRemoveGood, bc{"name": name, "type": "good", "domain": "Ex"})
	benchd(bIvRemoveGood, bc{"name": name, "type": "good", "domain": "Iv"})
	name = fmt.Sprintf("Removes(%v)", removeVals)
	benchd(bExRemovesGood, bc{"name": name, "type": "good", "domain": "Ex"})
	benchd(bIvRemovesGood, bc{"name": name, "type": "good", "domain": "Iv"})
	name = fmt.Sprintf("Add(%d)", addVal)
	benchd(bExAddGood, bc{"name": name, "type": "good", "domain": "Ex"})
	benchd(bIvAddGood, bc{"name": name, "type": "good", "domain": "Iv"})
	name = fmt.Sprintf("Contains(%d)", containsVal)
	benchd(bExContainsGood, bc{"name": name, "type": "good", "domain": "Ex"})
	benchd(bIvContainsGood, bc{"name": name, "type": "good", "domain": "Iv"})
	name = "Copy"
	benchd(bExCopyGood, bc{"name": name, "type": "good", "domain": "Ex"})
	benchd(bIvCopyGood, bc{"name": name, "type": "good", "domain": "Iv"})
	name = "GetMin"
	benchd(bExMinGood, bc{"name": name, "type": "good", "domain": "Ex"})
	benchd(bIvMinGood, bc{"name": name, "type": "good", "domain": "Iv"})
	name = "GetMax"
	benchd(bExMaxGood, bc{"name": name, "type": "good", "domain": "Ex"})
	benchd(bIvMaxGood, bc{"name": name, "type": "good", "domain": "Iv"})
	name = "IsEmpty"
	benchd(bExIsEmptyGood, bc{"name": name, "type": "good", "domain": "Ex"})
	benchd(bIvIsEmptyGood, bc{"name": name, "type": "good", "domain": "Iv"})
}

func bExRemoveGood(b *testing.B) {
	bExRemoveImpl(b, goodIntervals, removeVal)
}
func bIvRemoveGood(b *testing.B) {
	bIvRemoveImpl(b, goodIntervals, removeVal)
}

func bExRemovesGood(b *testing.B) {
	bExRemovesImpl(b, goodIntervals, removeVals)
}
func bIvRemovesGood(b *testing.B) {
	bIvRemovesImpl(b, goodIntervals, removeVals)
}

func bExAddGood(b *testing.B) {
	bExAddImpl(b, goodIntervals, addVal)
}
func bIvAddGood(b *testing.B) {
	bIvAddImpl(b, goodIntervals, addVal)
}

func bExContainsGood(b *testing.B) {
	bExContainsImpl(b, goodIntervals, containsVal)
}
func bIvContainsGood(b *testing.B) {
	bIvContainsImpl(b, goodIntervals, containsVal)
}

func bExCopyGood(b *testing.B) { bExCopyImpl(b, goodIntervals) }
func bIvCopyGood(b *testing.B) { bIvCopyImpl(b, goodIntervals) }

func bExIsEmptyGood(b *testing.B) { bExIsEmptyImpl(b, goodIntervals) }
func bIvIsEmptyGood(b *testing.B) { bIvIsEmptyImpl(b, goodIntervals) }

func bExMinGood(b *testing.B) { bExMinImpl(b, goodIntervals) }
func bIvMinGood(b *testing.B) { bIvMinImpl(b, goodIntervals) }

func bExMaxGood(b *testing.B) { bExMaxImpl(b, goodIntervals) }
func bIvMaxGood(b *testing.B) { bIvMaxImpl(b, goodIntervals) }

var curEx, curIv, z int
var dIv *core.IvDomain
var dEx *core.ExDomain

func bench_Interval_Trend() {
	dIv = core.CreateIvDomainFromTo(0, 200000)
	dEx = core.CreateExDomainFromTo(0, 200000)
	curIv = 1
	curEx = 1
	for z = 1; z <= 100001; z = z * 10 {
		zs := strconv.Itoa(z)
		name := fmt.Sprintf("Removes(%v)", removeVals)
		benchd(bExRemovesTrend,
			bc{"name": name, "type": "trend", "domain": "Ex", "size": zs})
		benchd(bIvRemovesTrend,
			bc{"name": name, "type": "trend", "domain": "Iv", "size": zs})
		name = "Copy"
		benchd(bExCopyTrend,
			bc{"name": name, "type": "trend", "domain": "Ex", "size": zs})
		benchd(bIvCopyTrend,
			bc{"name": name, "type": "trend", "domain": "Iv", "size": zs})
	}
}

func getTrendIvDomain() core.Domain {
	if curIv == z {
		return dIv.Copy()
	} // else
	curIv = z
	step := 200000 / curIv
	dIv = core.CreateIvDomainFromIntArr(makeVals(step))
	return dIv.Copy()
}

func getTrendExDomain() core.Domain {
	if curEx == z {
		return dEx.Copy()
	} // else
	curEx = z
	step := 200000 / curEx
	dEx = core.CreateExDomainAdds(makeVals(step))
	return dEx.Copy()
}

func makeVals(step int) []int {
	vs := make([]int, 0)
	for i := 0; i < 200000; i++ {
		if step == 1 { // worst
			if (i % 2) != 0 {
				vs = append(vs, i)
			}
		} else if (i % step) != 0 {
			vs = append(vs, i)
		}
	}
	return vs
}

func getTrendRemovesIvVal() *core.IvDomain {
	return core.CreateIvDomainFromTo(5000, 10000)
}

func getTrendRemovesExVal() *core.ExDomain {
	return core.CreateExDomainFromTo(5000, 10000)
}

func bExRemovesTrend(b *testing.B) {
	removesCheck(b, getTrendExDomain(), getTrendRemovesExVal())
}
func bIvRemovesTrend(b *testing.B) {
	removesCheck(b, getTrendIvDomain(), getTrendRemovesIvVal())
}

func bExCopyTrend(b *testing.B) { copyCheck(b, getTrendExDomain()) }
func bIvCopyTrend(b *testing.B) { copyCheck(b, getTrendIvDomain()) }

//-- Helper-class --

func bExRemoveImpl(b *testing.B, fromTos [][]int, removeVal int) {
	// init
	d := makeDomain(fromTos)
	removeCheck(b, d, removeVal)
}

func bIvRemoveImpl(b *testing.B, fromTos [][]int, removeVal int) {
	// init
	d := core.CreateIvDomainFromTos(fromTos)
	removeCheck(b, d, removeVal)
}

func removeCheck(b *testing.B, d core.Domain, removeVal int) {
	// bench
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		d.Remove(removeVal)
	}
}

func bExRemovesImpl(b *testing.B, fromTos [][]int, removeVals []int) {
	// init
	d := makeDomain(fromTos)
	removeD := core.CreateExDomainFromTo(removeVals[0], removeVals[1])
	removesCheck(b, d, removeD)
}

func bIvRemovesImpl(b *testing.B, fromTos [][]int, removeVals []int) {
	// init
	d := core.CreateIvDomainFromTos(fromTos)
	removeD := core.CreateIvDomainFromTo(removeVals[0], removeVals[1])
	removesCheck(b, d, removeD)
}

func removesCheck(b *testing.B, d core.Domain, removeVals core.Domain) {
	// bench
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		d.Removes(removeVals)
	}
}

func bExAddImpl(b *testing.B, fromTos [][]int, addVal int) {
	// init
	d := makeDomain(fromTos)
	addCheck(b, d, addVal)
}

func bIvAddImpl(b *testing.B, fromTos [][]int, addVal int) {
	// init
	d := core.CreateIvDomainFromTos(fromTos)
	addCheck(b, d, addVal)
}

func addCheck(b *testing.B, d core.Domain, addVal int) {
	// bench
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		d.Add(addVal)
	}
}

func bExContainsImpl(b *testing.B, fromTos [][]int, val int) {
	// init
	d := makeDomain(fromTos)
	containsCheck(b, d, val)
}

func bIvContainsImpl(b *testing.B, fromTos [][]int, val int) {
	// init
	d := core.CreateIvDomainFromTos(fromTos)
	containsCheck(b, d, val)
}

func containsCheck(b *testing.B, d core.Domain, val int) {
	// bench
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		d.Contains(val)
	}
}

func bExCopyImpl(b *testing.B, fromTos [][]int) {
	// init
	d := makeDomain(fromTos)
	copyCheck(b, d)
}

func bIvCopyImpl(b *testing.B, fromTos [][]int) {
	// init
	d := core.CreateIvDomainFromTos(fromTos)
	copyCheck(b, d)
}

func copyCheck(b *testing.B, d core.Domain) {
	// bench
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		d.Copy()
	}
}

func bExIsEmptyImpl(b *testing.B, fromTos [][]int) {
	// init
	d := makeDomain(fromTos)
	isEmptyCheck(b, d)
}

func bIvIsEmptyImpl(b *testing.B, fromTos [][]int) {
	// init
	d := core.CreateIvDomainFromTos(fromTos)
	isEmptyCheck(b, d)
}

func isEmptyCheck(b *testing.B, d core.Domain) {
	// bench
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		d.IsEmpty()
	}
}

func bExMinImpl(b *testing.B, fromTos [][]int) {
	// init
	d := makeDomain(fromTos)
	minCheck(b, d)
}

func bIvMinImpl(b *testing.B, fromTos [][]int) {
	// init
	d := core.CreateIvDomainFromTos(fromTos)
	minCheck(b, d)
}

func minCheck(b *testing.B, d core.Domain) {
	// bench
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		d.GetMin()
	}
}

func bExMaxImpl(b *testing.B, fromTos [][]int) {
	// init
	d := makeDomain(fromTos)
	maxCheck(b, d)
}

func bIvMaxImpl(b *testing.B, fromTos [][]int) {
	// init
	d := core.CreateIvDomainFromTos(fromTos)
	maxCheck(b, d)
}

func maxCheck(b *testing.B, d core.Domain) {
	// bench
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		d.GetMax()
	}
}

func makeDomain(fromTos [][]int) *core.ExDomain {
	size := 0
	for _, fromto := range fromTos {
		from := fromto[0]
		to := fromto[1]
		size += (to - from) + 1
	}
	a := make([]int, size)
	j := 0
	for _, fromto := range fromTos {
		for i := fromto[0]; i <= fromto[1]; i++ {
			a[j] = i
			j += 1
		}
	}
	d := core.CreateExDomainAdds(a)
	return d
}
