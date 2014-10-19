package core

import (
	"testing"
)

func Test_GetVarIdToNameMap(t *testing.T) {
	setup()
	defer teardown()
	log("RegistryStoreGetVarIdToNameMap")
	rs := CreateRegistryStore()
	rs.SetVarName(1, "X")
	rs.SetVarName(2, "Y")
	rs.SetVarName(3, "Z")

	rsIdToName := rs.GetVarIdToNameMap()

	if len(rsIdToName) != 3 {
		t.Error("RegistryStore has not the expected count of varname entries")
	}

	if rsIdToName[1] != "X" || rsIdToName[2] != "Y" || rsIdToName[3] != "Z" {
		t.Error("RegistryStore inserted varname-items incorrectly")
	}
}

func Test_GetSetHasName(t *testing.T) {
	setup()
	defer teardown()
	log("RegistryStoreGetSetHasName")
	rs := CreateRegistryStore()
	rs.SetVarName(1, "X")
	rs.SetVarName(2, "Y")
	rs.SetVarName(3, "Z")

	if rs.GetVarName(1) != "X" || rs.GetVarName(2) != "Y" || rs.GetVarName(3) != "Z" {
		t.Error("RegistryStore.GetName/SetName works incorrectly")
	}

	v, k := rs.HasVarName(1)
	if !k || v != "X" {
		t.Error("RegistryStore.HasName worked incorrectly")
	}
}
