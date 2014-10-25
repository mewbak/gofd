package core

type RegistryStore struct {
	idToName map[VarId]string // variable names back
	nameToId map[string]VarId // variable names
}

func CreateRegistryStore() *RegistryStore {
	return &RegistryStore{
		idToName: make(map[VarId]string),
		nameToId: make(map[string]VarId)}
}

// GetVarIdToNameMap returns the whole idToName-map
func (this *RegistryStore) Copy() *RegistryStore {
	rs := new(RegistryStore)
	rs.idToName = this.idToName
	rs.nameToId = this.nameToId
	return rs
}

// GetVarIdToNameMap returns the whole idToName-map
func (this *RegistryStore) GetVarIdToNameMap() map[VarId]string {
	return this.idToName
}

// GetNameToVarIdMap returns the whole nameToID-map
func (this *RegistryStore) GetNameToVarIdMap() map[string]VarId {
	return this.nameToId
}

// HasVarName returns true, iff a varname exists in RegistryStore
func (this *RegistryStore) HasVarName(id VarId) (string, bool) {
	val, ok := this.idToName[id]
	return val, ok
}

// SetVarName sets the human readable varName for VarId id
func (this *RegistryStore) SetVarName(id VarId, varName string) {
	if _, exists := this.idToName[id]; exists {
		panic("Same varId may not be used twice")
	} else if _, exists := this.nameToId[varName]; exists {
		panic("Same varName may not be used twice")
	} else {
		this.idToName[id] = varName
		this.nameToId[varName] = id
	}
}

// GetVarName returns the human readable name of VarId id
func (this *RegistryStore) GetVarName(id VarId) string {
	return this.idToName[id]
}

// GetVarName returns the human readable name of VarId id
func (this *RegistryStore) GetVarId(varName string) VarId {
	return this.nameToId[varName]
}
