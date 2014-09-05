package main

import (
	"bitbucket.org/gofd/gofd/core"
	"testing"
)

func main() {
	bench_Domain()
}

// the driver for everything benching Domain
func bench_Domain() {
	name := "CreateExDomainAdds"
	benchd(bCreateExDomainAdds1, bc{"name": name, "size": "1"})
	benchd(bCreateExDomainAdds10, bc{"name": name, "size": "10"})
	benchd(bCreateExDomainAdds100, bc{"name": name, "size": "100"})
	name = "Domain.GetMin"
	benchd(bDomainMin1, bc{"name": name, "size": "1"})
	benchd(bDomainMin10, bc{"name": name, "size": "10"})
	benchd(bDomainMin100, bc{"name": name, "size": "100"})
	name = "Domain.Equals"
	benchd(bDomainEquals1, bc{"name": name, "size": "1"})
	benchd(bDomainEquals10, bc{"name": name, "size": "10"})
	benchd(bDomainEquals100, bc{"name": name, "size": "100"})
}

func bDomainEquals1(b *testing.B)   { bDEquals(b, 1) }
func bDomainEquals10(b *testing.B)  { bDEquals(b, 10) }
func bDomainEquals100(b *testing.B) { bDEquals(b, 100) }

func bDEquals(b *testing.B, to int) {
	// init
	a := make([]int, to)
	for i := 0; i < to; i++ {
		a[i] = i + 1
	}
	d1 := core.CreateExDomainAdds(a)
	d2 := core.CreateExDomainAdds(a)
	// bench
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		d1.Equals(d2)
	}
}

func bDomainMin1(b *testing.B)   { bDMin(b, 1) }
func bDomainMin10(b *testing.B)  { bDMin(b, 10) }
func bDomainMin100(b *testing.B) { bDMin(b, 100) }

func bDMin(b *testing.B, to int) {
	// init
	a := make([]int, to)
	for i := 0; i < to; i++ {
		a[i] = i + 1
	}
	d := core.CreateExDomainAdds(a)
	// bench
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		d.GetMin()
	}
}

func bCreateExDomainAdds1(b *testing.B)   { bCDAdds(b, 1) }
func bCreateExDomainAdds10(b *testing.B)  { bCDAdds(b, 10) }
func bCreateExDomainAdds100(b *testing.B) { bCDAdds(b, 100) }

func bCDAdds(b *testing.B, to int) {
	// init
	a := make([]int, to)
	for i := 0; i < to; i++ {
		a[i] = i + 1
	}
	// bench
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		core.CreateExDomainAdds(a)
	}
}
