package main

import (
	"strconv"
	"testing"
)

func main() {
	bench_Intervalle_Bad()
}

// the driver for everything benching Domain
func bench_Intervalle_Bad() {
	name := "Bad.ExDomain.Remove(" + strconv.Itoa(getBadRemoveVal()) + ")"
	benchd(bExRemoveBad, bc{"name": name, "size": "bad"})
	name = "Bad.Ivomain.Remove(" + strconv.Itoa(getBadRemoveVal()) + ")"
	benchd(bIvRemoveBad, bc{"name": name, "size": "bad"})
	name = "Bad.ExplDomain.Removes(D([" + getBadRemovesIntervalAsString() + "]))"
	benchd(bExRemovesBad, bc{"name": name, "size": "bad"})
	name = "Bad.Ivomain.Removes(D([" + getBadRemovesIntervalAsString() + "]))"
	benchd(bIvRemovesBad, bc{"name": name, "size": "bad"})
	name = "Bad.ExplDomain.Add(" + strconv.Itoa(getBadAddVal()) + ")"
	benchd(bExAddBad, bc{"name": name, "size": "bad"})
	name = "Bad.Ivomain.Add(" + strconv.Itoa(getBadAddVal()) + ")"
	benchd(bIvAddBad, bc{"name": name, "size": "bad"})
	name = "Bad.ExplDomain.Contains(" + strconv.Itoa(getBadContainsVal()) + ")"
	benchd(bExContainsBad, bc{"name": name, "size": "bad"})
	name = "Bad.Ivomain.Contains(" + strconv.Itoa(getBadContainsVal()) + ")"
	benchd(bIvContainsBad, bc{"name": name, "size": "bad"})
	benchd(bExCopyBad, bc{"name": "Bad.ExplDomain.Copy", "size": "bad"})
	benchd(bIvCopyBad, bc{"name": "Bad.Ivomain.Copy", "size": "bad"})
	benchd(bExMinBad, bc{"name": "Bad.ExplDomain.GetMin", "size": "bad"})
	benchd(bIvMinBad, bc{"name": "Bad.Ivomain.GetMin", "size": "bad"})
	benchd(bExMaxBad, bc{"name": "Bad.ExplDomain.GetMax", "size": "bad"})
	benchd(bIvMaxBad, bc{"name": "Bad.Ivomain.GetMax", "size": "bad"})
	benchd(bExIsEmptyBad, bc{"name": "Bad.ExplDomain.IsEmpty", "size": "bad"})
	benchd(bIvIsEmptyBad, bc{"name": "Bad.Ivomain.IsEmpty", "size": "bad"})
}

func getBadRemovesIntervalAsString() string {
	vals := getBadRemovesVal()
	return "(" + strconv.Itoa(vals[0]) + "," + strconv.Itoa(vals[1]) + ")"
}

func getBadIntervals() [][]int {
	vals := make([][]int, 0)
	for i := 0; i <= 80000; i += 2 {
		vals = append(vals, []int{i, i})
	}
	return vals
}

func getBadRemoveVal() int {
	return 5000
}

func getBadRemovesVal() []int {
	return []int{5000, 10000}
}

func getBadAddVal() int {
	return 13000
}

func getBadContainsVal() int {
	return 5000
}

func bExRemoveBad(b *testing.B) { bExRemoveImpl(b, getBadIntervals(), getBadRemoveVal()) }
func bIvRemoveBad(b *testing.B) { bIvRemoveImpl(b, getBadIntervals(), getBadRemoveVal()) }

func bExRemovesBad(b *testing.B) { bExRemovesImpl(b, getBadIntervals(), getBadRemovesVal()) }
func bIvRemovesBad(b *testing.B) { bIvRemovesImpl(b, getBadIntervals(), getBadRemovesVal()) }

func bExAddBad(b *testing.B) { bExAddImpl(b, getBadIntervals(), getBadAddVal()) }
func bIvAddBad(b *testing.B) { bIvAddImpl(b, getBadIntervals(), getBadAddVal()) }

func bExContainsBad(b *testing.B) { bExContainsImpl(b, getBadIntervals(), getBadContainsVal()) }
func bIvContainsBad(b *testing.B) { bIvContainsImpl(b, getBadIntervals(), getBadContainsVal()) }

func bExCopyBad(b *testing.B) { bExCopyImpl(b, getBadIntervals()) }
func bIvCopyBad(b *testing.B) { bIvCopyImpl(b, getBadIntervals()) }

func bExIsEmptyBad(b *testing.B) { bExIsEmptyImpl(b, getBadIntervals()) }
func bIvIsEmptyBad(b *testing.B) { bIvIsEmptyImpl(b, getBadIntervals()) }

func bExMinBad(b *testing.B) { bExMinImpl(b, getBadIntervals()) }
func bIvMinBad(b *testing.B) { bIvMinImpl(b, getBadIntervals()) }

func bExMaxBad(b *testing.B) { bExMaxImpl(b, getBadIntervals()) }
func bIvMaxBad(b *testing.B) { bIvMaxImpl(b, getBadIntervals()) }
