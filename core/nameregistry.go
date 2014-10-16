package core

type NameRegistry struct {
	iDToName map[VarId]string // variable names back
}

var nm = &NameRegistry{iDToName: make(map[VarId]string)}

func GetNameRegistry() *NameRegistry {
	return nm
}

// GetAll returns the whole idToName-map
func (*NameRegistry) GetAll() map[VarId]string {
	return nm.iDToName
}

// SetName sets the human readable name for VarId id
func (*NameRegistry) HasName(id VarId) (string, bool) {
	val, ok := nm.iDToName[id]
	return val, ok
}

// SetName sets the human readable name for VarId id
func (*NameRegistry) SetName(id VarId, name string) {
	nm.iDToName[id] = name
}

// GetName returns the human readable name of VarId id
func (*NameRegistry) GetName(id VarId) string {
	return nm.iDToName[id]
}
