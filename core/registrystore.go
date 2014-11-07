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

func CreateRegistryStore() *RegistryStore {
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
	//close all channels
	for constraintData, _ := range this.constraintsToVarIds {
		close(constraintData.channel)
	}
	//set memory free for garbage collection
	this.constraintsToVarIds = nil
	this.varIdsToConstraints = nil
	this.constraints = nil
}

// ConnectVarIdWithConstraint connect Constraint with related varids
// - a constraint is interested in varids v1, v2, ..., vx (list can be
//   reduced while propagation/labeling)
//   (see AddConstraintToVarIds)
// - changes for the domain of a varid can be "observed" by several constraints
//   (see AddVarIdsToConstraint)
func (this *RegistryStore) RegisterVarIdWithConstraint(propId PropId,
	writeChannel chan *ChangeEntry, varIds []VarId,
	interestedInVarIds []VarId) {
	constraintData := CreateConstraintData(propId,
		this.constraints[propId], writeChannel)

	this.addVarIdsToConstraint(varIds, constraintData)
	this.addConstraintToVarIds(constraintData, interestedInVarIds)
}

// addVarIdsToConstraint stores the relation between several varids and the
// specific constraint, e.g. constraint C1: X+Y=Z with id c1 and
// writechannel w1 is given. C1, c1 and w1 are store in a ConstraintData
// instance cd1. Then for each variable X,Y and Z cd1 is stored. So, when
// a change for X is incoming, then all cdx will be informed.
func (this *RegistryStore) addVarIdsToConstraint(varIds []VarId,
	cD *ConstraintData) {
	for _, varId := range varIds {
		if _, exists := this.varIdsToConstraints[varId]; !exists {
			this.varIdsToConstraints[varId] = make([]*ConstraintData, 0)
		}
		this.varIdsToConstraints[varId] =
			append(this.varIdsToConstraints[varId], cD)
	}
}

// addConstraintToVarIds stores the relation between a specific constraint
// and several varids, e.g. constraint C1: X+Y=Z with id c1 and
// writechannel w1 is given. C1, c1 and w1 are store in a ConstraintData
// instance cd1. Then cd1 is stored with the slice of variables (X,Y and Z),
// for which he wants to get change information.
func (this *RegistryStore) addConstraintToVarIds(cD *ConstraintData,
	varids []VarId) {
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
