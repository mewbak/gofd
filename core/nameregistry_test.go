package core

import (
	"testing"
)

func Test_GetNameRegistry(t *testing.T) {
	setup()
	defer teardown()
	log("NameRegistryGetNameRegistry")
	if GetNameRegistry() == nil {
		t.Error("NameRegistry must be available at start")
	}
}

func Test_GetAll(t *testing.T) {
	setup()
	defer teardown()
	log("NameRegistryGetAll")
	nr := GetNameRegistry()
	nr.SetName(1, "X")
	nr.SetName(2, "Y")
	nr.SetName(3, "Z")

	nrContent := nr.GetAll()

	if len(nrContent) != 3 {
		t.Error("NameRegistry has not the expected count of values")
	}

	if nrContent[1] != "X" || nrContent[2] != "Y" || nrContent[3] != "Z" {
		t.Error("NameRegistry inserted items incorrectly")
	}
}

func Test_GetSetHasName(t *testing.T) {
	setup()
	defer teardown()
	log("NameRegistryGetSetHasName")
	nr := GetNameRegistry()
	nr.SetName(1, "X")
	nr.SetName(2, "Y")
	nr.SetName(3, "Z")

	if nr.GetName(1) != "X" || nr.GetName(2) != "Y" || nr.GetName(3) != "Z" {
		t.Error("NameRegistry.GetName/SetName works wrong")
	}

	v, k := nr.HasName(1)
	if !k || v != "X" {
		t.Error("NameRegistry.HasName worked wrong")
	}
}
