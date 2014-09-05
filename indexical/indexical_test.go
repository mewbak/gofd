package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical/ixrange"
	"bitbucket.org/gofd/gofd/indexical/ixterm"
	"testing"
)

//ToDo: more tests needed

func Test_Indexical(t *testing.T) {
	setup()
	defer teardown()
	log("Indexical")
	// 1..100 in [50..100] + 5
	outVarD := core.CreateIvDomainFromTo(1, 100)
	inVarD := core.CreateIvDomainFromTo(50, 100)
	expRemoving := core.CreateIvDomainFromTo(1, 54)
	expDResult := core.CreateIvDomainFromTo(55, 100)
	c := 5
	domR := ixrange.CreateDomRange(1, inVarD)
	valueT := ixterm.CreateValueTerm(c)
	valueR := ixrange.CreateSingleValueRange(valueT)
	addR := ixrange.CreateAddRange(domR, valueR)
	indexical_test(t, outVarD, addR, expRemoving)
	indexical_withRemoving_test(t, outVarD, addR, expDResult)
}

func indexical_test(t *testing.T, outD *core.IvDomain,
	r ixrange.IRange, expRemoving *core.IvDomain) {
	i := CreateIndexical(1, outD, r)
	evt := core.CreateChangeEvent()
	i.Process(evt, nil, false)
	evt2 := core.CreateChangeEvent()
	chEntry := core.CreateChangeEntryWithValues(1, expRemoving)
	evt2.AddChangeEntry(chEntry)
	if !evt.Equals(evt2) {
		t.Errorf("Indexical-InitProcessing failed: created evts are not "+
			"equal. Evt1 %s, evt2 %s", evt, evt2)
	}
}

func indexical_withRemoving_test(t *testing.T, outD *core.IvDomain,
	r ixrange.IRange, expD *core.IvDomain) {
	i := CreateIndexical(1, outD, r)
	evt := core.CreateChangeEvent()
	i.Process(evt, nil, true)
	if !outD.Equals(expD) {
		t.Errorf("Indexical-InitProcessing (with removing) failed")
	}
}
