package core

type NameRegistry struct {
	iDToName map[VarId]string // variable names back
}

var nr *NameRegistry

func CreateNameRegistry() {
	nr = &NameRegistry{iDToName: make(map[VarId]string)}
}

func GetNameRegistry() *NameRegistry {
	return nr
}

func (*NameRegistry) Clean() {
	nr = &NameRegistry{iDToName: make(map[VarId]string)}
}

// GetAll returns the whole idToName-map
func (*NameRegistry) GetAll() map[VarId]string {
	return nr.iDToName
}

// SetName sets the human readable name for VarId id
func (*NameRegistry) HasName(id VarId) (string, bool) {
	val, ok := nr.iDToName[id]
	return val, ok
}

// SetName sets the human readable name for VarId id
func (*NameRegistry) SetName(id VarId, name string) {
	if _, exists := nr.iDToName[id]; exists {
		panic("Same varId may not be used twice")
	} else {
		nr.iDToName[id] = name
	}
}

// GetName returns the human readable name of VarId id
func (*NameRegistry) GetName(id VarId) string {
	return nr.iDToName[id]
}
