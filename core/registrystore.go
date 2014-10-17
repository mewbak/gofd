package core

type RegistryStore struct {
	iDToName map[VarId]string // variable names back
}

func CreateRegistryStore() *RegistryStore{
	return &RegistryStore{iDToName: make(map[VarId]string)}
}

// getVarIdToNameMap returns the whole idToName-map
func (this *RegistryStore) GetVarIdToNameMap() map[VarId]string {
	return this.iDToName
}

// setName sets the human readable name for VarId id
func (this *RegistryStore) HasVarName(id VarId) (string, bool) {
	val, ok := this.iDToName[id]
	return val, ok
}

// setName sets the human readable name for VarId id
func (this *RegistryStore) SetVarName(id VarId, name string) {
	if _, exists := this.iDToName[id]; exists{
		panic("Same varId may not be used twice") 	
	}else{
		this.iDToName[id] = name
	}
}

// getName returns the human readable name of VarId id
func (this *RegistryStore) GetVarName(id VarId) string {
	return this.iDToName[id]
}
