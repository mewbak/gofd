package main

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

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
