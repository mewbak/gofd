package core

type RegistryStore struct {
	iDToName map[VarId]string // variable names back
}

func CreateRegistryStore() *RegistryStore {
	return &RegistryStore{iDToName: make(map[VarId]string)}
}

// GetVarIdToNameMap returns the whole idToName-map
func (this *RegistryStore) Copy() *RegistryStore {
	rs := new(RegistryStore)
	rs.iDToName = this.iDToName
	return rs
}

// GetVarIdToNameMap returns the whole idToName-map
func (this *RegistryStore) GetVarIdToNameMap() map[VarId]string {
	return this.iDToName
}

// HasVarName returns true, iff a varname exists in RegistryStore
func (this *RegistryStore) HasVarName(id VarId) (string, bool) {
	val, ok := this.iDToName[id]
	return val, ok
}

// SetVarName sets the human readable name for VarId id
func (this *RegistryStore) SetVarName(id VarId, name string) {
	if _, exists := this.iDToName[id]; exists {
		panic("Same varId may not be used twice")
	} else {
		this.iDToName[id] = name
	}
}

// GetVarName returns the human readable name of VarId id
func (this *RegistryStore) GetVarName(id VarId) string {
	return this.iDToName[id]
}
