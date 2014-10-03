package main

import (
	"strconv"
	"testing"
)

func main() {
	bench_Intervalle_Good()
}

// the driver for everything benching Domain
func bench_Intervalle_Good() {
	name := "Good.ExDomain.Remove(" + strconv.Itoa(getGoodRemoveVal()) + ")"
	benchd(bExRemoveGood, bc{"name": name, "size": "good"})
	name = "Good.IvDomain.Remove(" + strconv.Itoa(getGoodRemoveVal()) + ")"
	benchd(bIvRemoveGood, bc{"name": name, "size": "good"})
	name = "Good.ExDomain.Removes(D([" + getGoodRemovesIntervalAsString() + "]))"
	benchd(bExRemovesGood, bc{"name": name, "size": "good"})
	name = "Good.IvDomain.Removes(D([" + getGoodRemovesIntervalAsString() + "]))"
	benchd(bIvRemovesGood, bc{"name": name, "size": "good"})
	name = "Good.ExDomain.Add(" + strconv.Itoa(getGoodAddVal()) + ")"
	benchd(bExAddGood, bc{"name": name, "size": "good"})
	name = "Good.IvDomain.Add(" + strconv.Itoa(getGoodAddVal()) + ")"
	benchd(bIvAddGood, bc{"name": name, "size": "good"})
	name = "Good.ExDomain.Contains(" + strconv.Itoa(getGoodContainsVal()) + ")"
	benchd(bExContainsGood, bc{"name": name, "size": "good"})
	name = "Good.IvDomain.Contains(" + strconv.Itoa(getGoodContainsVal()) + ")"
	benchd(bIvContainsGood, bc{"name": name, "size": "good"})
	benchd(bExCopyGood, bc{"name": "Good.ExDomain.Copy", "size": "good"})
	benchd(bIvCopyGood, bc{"name": "Good.IvDomain.Copy", "size": "good"})
	benchd(bExMinGood, bc{"name": "Good.ExDomain.GetMin", "size": "good"})
	benchd(bIvMinGood, bc{"name": "Good.IvDomain.GetMin", "size": "good"})
	benchd(bExMaxGood, bc{"name": "Good.ExDomain.GetMax", "size": "good"})
	benchd(bIvMaxGood, bc{"name": "Good.IvDomain.GetMax", "size": "good"})
	benchd(bExIsEmptyGood, bc{"name": "Good.ExDomain.IsEmpty", "size": "good"})
	benchd(bIvIsEmptyGood, bc{"name": "Good.IvDomain.IsEmpty", "size": "good"})
}

func getGoodRemovesIntervalAsString() string {
	vals := getGoodRemovesVal()
	return "(" + strconv.Itoa(vals[0]) + "," + strconv.Itoa(vals[1]) + ")"
}

func getGoodIntervals() [][]int {
	return [][]int{{0, 10000}, {15000, 80000}}
}

func getGoodRemoveVal() int {
	return 5000
}

func getGoodRemovesVal() []int {
	return []int{5000, 10000}
}

func getGoodAddVal() int {
	return 13000
}

func getGoodContainsVal() int {
	return 5000
}

func bExRemoveGood(b *testing.B) { bExRemoveImpl(b, getGoodIntervals(), getGoodRemoveVal()) }
func bIvRemoveGood(b *testing.B) { bIvRemoveImpl(b, getGoodIntervals(), getGoodRemoveVal()) }

func bExRemovesGood(b *testing.B) { bExRemovesImpl(b, getGoodIntervals(), getGoodRemovesVal()) }
func bIvRemovesGood(b *testing.B) { bIvRemovesImpl(b, getGoodIntervals(), getGoodRemovesVal()) }

func bExAddGood(b *testing.B) { bExAddImpl(b, getGoodIntervals(), getGoodAddVal()) }
func bIvAddGood(b *testing.B) { bIvAddImpl(b, getGoodIntervals(), getGoodAddVal()) }

func bExContainsGood(b *testing.B) { bExContainsImpl(b, getGoodIntervals(), getGoodContainsVal()) }
func bIvContainsGood(b *testing.B) { bIvContainsImpl(b, getGoodIntervals(), getGoodContainsVal()) }

func bExCopyGood(b *testing.B) { bExCopyImpl(b, getGoodIntervals()) }
func bIvCopyGood(b *testing.B) { bIvCopyImpl(b, getGoodIntervals()) }

func bExIsEmptyGood(b *testing.B) { bExIsEmptyImpl(b, getGoodIntervals()) }
func bIvIsEmptyGood(b *testing.B) { bIvIsEmptyImpl(b, getGoodIntervals()) }

func bExMinGood(b *testing.B) { bExMinImpl(b, getGoodIntervals()) }
func bIvMinGood(b *testing.B) { bIvMinImpl(b, getGoodIntervals()) }

func bExMaxGood(b *testing.B) { bExMaxImpl(b, getGoodIntervals()) }
func bIvMaxGood(b *testing.B) { bIvMaxImpl(b, getGoodIntervals()) }
