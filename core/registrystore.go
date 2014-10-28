package core

// RegistryStore represents a central part of the data
// in a store connecting varIds to names and propagators
type RegistryStore struct {
	idToName map[VarId]string // variable names back
	nameToId map[string]VarId // variable names
	// varids of unfixed variables per propagator
	constraintsToVarIds map[*ConstraintData][]VarId
	// pointer to connected constraints per varid
	varIdsToConstraints map[VarId][]*ConstraintData
	// collection of all constraints/propagators
	constraints map[PropId]Constraint
}

// ---------- general functions ----------

func createRegistryStore() *RegistryStore {
	return &RegistryStore{
		idToName:            make(map[VarId]string),
		nameToId:            make(map[string]VarId),
		constraintsToVarIds: make(map[*ConstraintData][]VarId),
		varIdsToConstraints: make(map[VarId][]*ConstraintData),
		constraints:         make(map[PropId]Constraint)}
}

// Clone clones the registry store
func (this *RegistryStore) Clone() (*RegistryStore, []Constraint) {
	rs := new(RegistryStore)
	rs.idToName = this.idToName
	rs.nameToId = this.nameToId
	rs.constraintsToVarIds = make(map[*ConstraintData][]VarId,
		len(this.constraintsToVarIds))
	rs.varIdsToConstraints = make(map[VarId][]*ConstraintData,
		len(this.varIdsToConstraints))
	rs.constraints = make(map[PropId]Constraint, len(this.constraints))
	clonedConstraints := make([]Constraint, len(this.constraints))
	i := 0
	for _, prop := range this.constraints {
		clonedConstraints[i] = prop.Clone()
		// Set temporary id. Will be overridden for
		// each propagator in Addconstraints
		clonedConstraints[i].SetID(0)
		i += 1
	}
	return rs, clonedConstraints
}

//---------- varid to name and vice versa ----------

// getVarIdToNameMap returns the whole idToName-map
func (this *RegistryStore) getVarIdToNameMap() map[VarId]string {
	return this.idToName
}

// getNameToVarIdMap returns the whole nameToID-map
func (this *RegistryStore) getNameToVarIdMap() map[string]VarId {
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

func (this *RegistryStore) numberOfConstraints() int {
	return len(this.constraints)
}

func (this *RegistryStore) numberOfActiveVarIds() int {
	return len(this.varIdsToConstraints)
}

//---------- propagator <-> varid mappings ----------

// ConstraintData represents all information about a propagator
type ConstraintData struct {
	propId     PropId
	constraint Constraint
	channel    chan *ChangeEntry
}

func CreateConstraintData(propId PropId, constraint Constraint,
	channel chan *ChangeEntry) *ConstraintData {
	cd := new(ConstraintData)
	cd.propId = propId
	cd.constraint = constraint
	cd.channel = channel
	return cd
}

func (this *RegistryStore) Close() {
	varids := make([]VarId, len(this.varIdsToConstraints))
	constraintDatas := make([]*ConstraintData,
		len(this.constraintsToVarIds))
	i := 0
	for varid, _ := range this.varIdsToConstraints {
		varids[i] = varid
		i += 1
	}
	i = 0
	for constraintData, _ := range this.constraintsToVarIds {
		constraintDatas[i] = constraintData
		i += 1
	}
	for _, constraintData := range constraintDatas {
		this.constraintsToVarIds[constraintData] = nil
		close(constraintData.channel)
	}
	this.constraintsToVarIds = nil
	this.varIdsToConstraints = nil
	this.constraints = nil
}

func (this *RegistryStore) AddVarIdsToConstraint(varIds []VarId,
	cD *ConstraintData) {
	for _, varId := range varIds {
		if _, exists := this.varIdsToConstraints[varId]; !exists {
			this.varIdsToConstraints[varId] = make([]*ConstraintData, 0)
		}
		this.varIdsToConstraints[varId] =
			append(this.varIdsToConstraints[varId], cD)
	}
}

func (this *RegistryStore) AddConstraintToVarIds(cD *ConstraintData,
	varids []VarId) {
	if _, exists := this.constraintsToVarIds[cD]; !exists {
		this.constraintsToVarIds[cD] = make([]VarId, 0)
	}
	// TODO: that does not make sense, it overrides?
	this.constraintsToVarIds[cD] = varids
}

// RemoveFixedRelations removes relations between varid and propId, if
// domain of varid is ground
func (this *RegistryStore) RemoveFixedRelations(varid VarId) int {
	constraintData := this.varIdsToConstraints[varid]
	removedConstraints := 0
	for _, constraintD := range constraintData {
		// remove
		indexToRemove := -1
		for i, vid := range this.constraintsToVarIds[constraintD] {
			if vid == varid {
				indexToRemove = i
				break
			}
		}
		this.constraintsToVarIds[constraintD] =
			append(this.constraintsToVarIds[constraintD][:indexToRemove],
				this.constraintsToVarIds[constraintD][indexToRemove+1:]...)
		if len(this.constraintsToVarIds[constraintD]) == 0 {
			close(constraintD.channel)
			// remove dangling reference to constraints, allow gc
			delete(this.constraints, constraintD.propId)
			delete(this.constraintsToVarIds, constraintD)
			removedConstraints += 1
		}
	}
	this.varIdsToConstraints[varid] = nil
	return removedConstraints
}

func (this *RegistryStore) Constraints(varid VarId) []*ConstraintData {
	return this.varIdsToConstraints[varid]
}
