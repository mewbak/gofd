package core

import (
	"fmt"
)

// VarId is a type for unique identifier for IntVars
type VarId int

// VarIdSet is a type for a set of VarIds
type VarIdSet map[VarId]bool

// IntVar is an Integer Variable with a Domain
type IntVar struct {
	ID     VarId
	Domain Domain
}

// CreateIntVarDom creates a new IntVar object with a domain
// and adds it to the store.
func CreateIntVarDom(name string, store *Store, dom Domain) VarId {
	intvar := new(IntVar)
	intvar.Domain = dom
	name = intvar_fixname(name, store)
	return store.registerIntVarAtStore(name, intvar)
}

// CreateIntVarValues creates a new IntVar object with domain values
// and adds it to the store.
func CreateIntVarValues(name string, store *Store, values []int) VarId {
	return CreateIntVarIvValues(name, store, values) // default interval domain
}

// CreateIntVarIvValues creates a new IntVar object with a given int-slice
// as domain values.
func CreateIntVarIvValues(name string, store *Store, values []int) VarId {
	intvar := new(IntVar)
	intvar.Domain = CreateIvDomainFromIntArr(values)
	name = intvar_fixname(name, store)
	return store.registerIntVarAtStore(name, intvar)
}

// CreateIntVarValues creates a new IntVar object with an explicit domain
// from domain values and adds it to the store.
func CreateIntVarExValues(name string, store *Store, values []int) VarId {
	intvar := new(IntVar)
	intvar.Domain = CreateExDomainAdds(values)
	name = intvar_fixname(name, store)
	return store.registerIntVarAtStore(name, intvar)
}

// CreateIntVarFromTo creates a new IntVar object with domain values
// from and to inclusive and adds it to the store.
func CreateIntVarFromTo(name string, store *Store, from, to int) VarId {
	return CreateIntVarIvFromTo(name, store, from, to)
}

// CreateIntVarIvFromTo creates a new IntVar object with interval domain from
// domain values from and to inclusive and adds it to the store.
func CreateIntVarIvFromTo(name string, store *Store, from int, to int) VarId {
	intvar := new(IntVar)
	intvar.Domain = CreateIvDomainFromTo(from, to)
	name = intvar_fixname(name, store)
	return store.registerIntVarAtStore(name, intvar)
}

// CreateIntVarExFromTo creates a new IntVar object with explicit domain from
// domain values from and to inclusive and adds it to the store.
func CreateIntVarExFromTo(name string, store *Store, from, to int) VarId {
	intvar := new(IntVar)
	intvar.Domain = CreateExDomainFromTo(from, to)
	name = intvar_fixname(name, store)
	return store.registerIntVarAtStore(name, intvar)
}

// CreateAuxIntVarValues creates a new auxiliary IntVar object from
// array values and adds it to the store.
func CreateAuxIntVarValues(store *Store, values []int) VarId {
	return CreateAuxIntVarIvValues(store, values)
}

// CreateAuxIntVarValues creates a new auxiliary IntVar object with interval
// domain from array values and adds it to the store.
func CreateAuxIntVarIvValues(store *Store, values []int) VarId {
	intvar := new(IntVar)
	intvar.Domain = CreateIvDomainFromIntArr(values)
	return store.registerAuxIntVarAtStore(intvar)
}

// CreateAuxIntVarValues creates a new auxiliary IntVar object with explicit
// domain from array values and adds it to the store.
func CreateAuxIntVarExValues(store *Store, values []int) VarId {
	intvar := new(IntVar)
	intvar.Domain = CreateExDomainAdds(values)
	return store.registerAuxIntVarAtStore(intvar)
}

// CreateAuxIntVarFromTo creates a new auxiliary IntVar object with domain
// values from and to inclusive and adds it to the store.
func CreateAuxIntVarFromTo(store *Store, from, to int) VarId {
	return CreateAuxIntVarExFromTo(store, from, to) // ToDo: default Iv
}

// CreateAuxIntVarIvFromTo creates a new auxiliary IntVar object with interval
// domain from values from and to inclusive and adds it to the store.
func CreateAuxIntVarIvFromTo(store *Store, from, to int) VarId {
	intvar := new(IntVar)
	intvar.Domain = CreateIvDomainFromTo(from, to)
	return store.registerAuxIntVarAtStore(intvar)
}

// CreateAuxIntVarFromTo creates a new auxiliary IntVar object with explicit
// domain from values from and to inclusive and adds it to the store.
func CreateAuxIntVarExFromTo(store *Store, from, to int) VarId {
	intvar := new(IntVar)
	intvar.Domain = CreateExDomainFromTo(from, to)
	return store.registerAuxIntVarAtStore(intvar)
}

/* more convencience functions */

// CreateIntVarsIvFromTo sets the given IntVars (*VarId). The same interval
// domain for all Intvars is used (from and to inclusive). IntVars will be
// added it to the store. varids must be given as reference, not value.
func CreateIntVarsIvFromTo(varids []*VarId, names []string,
	store *Store, from int, to int) {
	for i, varid := range varids {
		intvar := new(IntVar)
		intvar.Domain = CreateIvDomainFromTo(from, to)
		name := intvar_fixname(names[i], store)
		(*varid) = store.registerIntVarAtStore(name, intvar)
	}
}

// CreateIntVarsIvValues sets the given IntVars (*VarId). The same interval
// domain for all Intvars is used (from and to inclusive). IntVars will be
// added it to the store. varids must be given as reference, not value.
func CreateIntVarsIvValues(varids []*VarId, names []string,
	store *Store, values []int) {
	for i, varid := range varids {
		intvar := new(IntVar)
		intvar.Domain = CreateIvDomainFromIntArr(values)
		name := intvar_fixname(names[i], store)
		(*varid) = store.registerIntVarAtStore(name, intvar)
	}
}

// CreateIntVarIvDomBool creates a new bool IntVar object with interval domain.
func CreateIntVarIvDomBool(name string, store *Store) VarId {
	intvar := new(IntVar)
	intvar.Domain = CreateIvDomainFromTo(0, 1)
	name = intvar_fixname(name, store)
	return store.registerIntVarAtStore(name, intvar)
}

// String representation of an IntVar
func (this *IntVar) String() string {
	return fmt.Sprintf("IntVar{%d,%s}", this.ID, this.Domain.String())
}

// Clone creates a deep copy of this IntVar
func (this *IntVar) Clone() *IntVar {
	newVar := new(IntVar)
	newVar.ID = this.ID
	newVar.Domain = this.Domain.Copy()
	return newVar
}

// IsGround returns true iff the domain holds exactly one value
func (this *IntVar) IsGround() bool {
	return this.Domain.IsGround()
}

/* helper */

// name may not be empty and may not begin with '_' (auxiliary var prefix)
func intvar_fixname(name string, store *Store) string {
	if name == "" || name[0] == '_' {
		currentName := name
		name = store.generateNewVariableName()
		if logger.DoInfo() {
			msg := "CreateIntVar*_name '%s' for variable not valid,"
			msg += " generated new one '%s'"
			logger.If(msg, currentName, name)
		}
	}
	return name
}
