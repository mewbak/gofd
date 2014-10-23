package main

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/demo"
	"bitbucket.org/gofd/gofd/labeling"
	"testing"
)

func main() {
	benchd(bSendMoreMoneyPlain,
		bc{"name": "SendMoreMoney", "type": "plain"})
	benchd(bSendMoreMoneyIndexical,
		bc{"name": "SendMoreMoney", "type": "indexical"})
}

func bSendMoreMoneyPlain(b *testing.B) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		store := core.CreateStoreWithoutLogging()
		vars := demo.ConstrainSendMoreMoney(store, false)
		labeling.SetAllvars(vars)
		query := labeling.CreateSearchOneQuery()
		labeling.Labeling(store, query,
			labeling.InDomainMin, labeling.VarSelect)
	}
}

func bSendMoreMoneyIndexical(b *testing.B) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	for i := 0; i < b.N; i++ {
		store := core.CreateStoreWithoutLogging()
		vars := demo.ConstrainSendMoreMoney(store, true)
		labeling.SetAllvars(vars)
		query := labeling.CreateSearchOneQuery()
		labeling.Labeling(store, query,
			labeling.InDomainMin, labeling.VarSelect)
	}
}
