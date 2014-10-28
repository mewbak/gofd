package core

import (
	"testing"
)

func checkVarNamesRegistryStore(t *testing.T, msg, got, want string) {
	if want != got {
		t.Error("RegistryStore %s: got %s, want %s",
			msg, want, got)
	}
}

func Test_GetVarIdToNameMap(t *testing.T) {
	setup()
	defer teardown()
	log("RegistryStoreGetVarIdToNameMap")
	rs := createRegistryStore()
	rs.SetVarName(1, "X")
	rs.SetVarName(2, "Y")
	rs.SetVarName(3, "Z")
	rsIdToName := rs.getVarIdToNameMap()
	want := 3
	if len(rsIdToName) != want {
		t.Errorf("RegistryStore count varnames got %d, want %d",
			len(rsIdToName), want)
	}
	msg := "inserted"
	checkVarNamesRegistryStore(t, msg, rsIdToName[1], "X")
	checkVarNamesRegistryStore(t, msg, rsIdToName[2], "Y")
	checkVarNamesRegistryStore(t, msg, rsIdToName[3], "Z")
}

func Test_GetSetHasName(t *testing.T) {
	setup()
	defer teardown()
	log("RegistryStoreGetSetHasName")
	rs := createRegistryStore()
	rs.SetVarName(1, "X")
	rs.SetVarName(2, "Y")
	rs.SetVarName(3, "Z")
	msg := "SetGetVarName"
	checkVarNamesRegistryStore(t, msg, rs.GetVarName(1), "X")
	checkVarNamesRegistryStore(t, msg, rs.GetVarName(2), "Y")
	checkVarNamesRegistryStore(t, msg, rs.GetVarName(3), "Z")
	v, k := rs.HasVarName(1)
	if !k || v != "X" {
		t.Errorf("RegistryStore.HasName(%s) not found", "X")
	}
}
