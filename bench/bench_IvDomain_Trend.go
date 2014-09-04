package main

import (
	"bitbucket.org/gofd/gofd/core"
	"strconv"
	"testing"
)

func main() {
	bench_Intervalle_trend()
}

var curEx, curIv, z int
var dIv *core.IvDomain
var dEx *core.ExDomain

func bench_Intervalle_trend() {
	dIv = core.CreateIvDomainFromTo(0, 200000)
	dEx = core.CreateExDomainFromTo(0, 200000)
	curIv = 1
	curEx = 1
	for z = 1; z <= 100001; {
		name := strconv.Itoa(z) + ".ExDomain.Removes(D([" + getTrendRemovesIntervalAsString() + "]))"
		benchd(bExRemovesTrend, bc{"name": name, "size": "Trend"})
		name = strconv.Itoa(z) + ".IvDomain.Removes(D([" + getTrendRemovesIntervalAsString() + "]))"
		benchd(bIvRemovesTrend, bc{"name": name, "size": "Trend"})
		benchd(bExCopyTrend, bc{"name": strconv.Itoa(z) + ".ExDomain.Copy", "size": "Trend"})
		benchd(bIvCopyTrend, bc{"name": strconv.Itoa(z) + ".IvDomain.Copy", "size": "Trend"})
		z = z * 10
	}
}

func getTrendRemovesIntervalAsString() string {
	vals := getTrendRemovesIvVal()
	return "(" + strconv.Itoa(vals.GetParts()[0].From) + "," + strconv.Itoa(vals.GetParts()[0].To) + ")"
}

/*
X in [(0,0),(2,2),(4,4),(6,6),(8,8),(10,10)]
  --> add 1, add 9
X in [(0,2),(4,4),(6,6),(8,10)]
  --> add 3, add 7
X in [(0,4),(6,10)]
  --> add 5
X in [(0,10)]
[0,80000]
80000/1 = 80000 (remove all 80000 one value)
80000/10 = 8000
80000/100 = 800
*/
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

func getTrendRemoveVal() int {
	return 5000
}

func getTrendRemovesIvVal() *core.IvDomain {
	return core.CreateIvDomainFromTo(5000, 10000)
}

func getTrendRemovesExVal() *core.ExDomain {
	return core.CreateExDomainFromTo(5000, 10000)
}

func getTrendAddVal() int {
	return 13000
}

func getTrendContainsVal() int {
	return 5000
}

func bExRemovesTrend(b *testing.B) { removesCheck(b, getTrendExDomain(), getTrendRemovesExVal()) }
func bIvRemovesTrend(b *testing.B) { removesCheck(b, getTrendIvDomain(), getTrendRemovesIvVal()) }

func bExCopyTrend(b *testing.B) { copyCheck(b, getTrendExDomain()) }
func bIvCopyTrend(b *testing.B) { copyCheck(b, getTrendIvDomain()) }
