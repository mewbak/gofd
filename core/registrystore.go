package core

type RegistryStore struct {
	idToName map[VarId]string // variable names back
	nameToId map[string]VarId // variable names
	// constraintsInterestedInVarids, varids of unfixed variables per propagator
	constraintsInterestedInVarids map[*ConstraintData][]VarId
	varidsConnectedToConstraints map[VarId][]*ConstraintData
	constraints map[PropId]Constraint
}

// ---------- general functions ---------- 

func CreateRegistryStore() *RegistryStore {
	return &RegistryStore{
		idToName: make(map[VarId]string),
		nameToId: make(map[string]VarId),
		constraintsInterestedInVarids: make(map[*ConstraintData][]VarId),
		varidsConnectedToConstraints: make(map[VarId][]*ConstraintData),
		constraints: make(map[PropId]Constraint)}
}

// Copy copies the registry store
func (this *RegistryStore) Clone() (*RegistryStore, []Constraint) {
	rs := new(RegistryStore)
	rs.idToName = this.idToName
	rs.nameToId = this.nameToId
	rs.constraintsInterestedInVarids = make(map[*ConstraintData][]VarId, len(this.constraintsInterestedInVarids))
	rs.varidsConnectedToConstraints = make(map[VarId][]*ConstraintData, len(this.varidsConnectedToConstraints))
	rs.constraints = make(map[PropId]Constraint, len(this.constraints))

	clonedConstraints := make([]Constraint,len(this.constraints))

	i:=0
	for _, prop := range this.constraints {
		clonedConstraints[i] = prop.Clone()
		clonedConstraints[i].SetID(0) // temporary value. Will be overridden
		// for each propagator in function "Addconstraints"
		i+=1
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

//---------- propagator <-> varid mappings ----------

type ConstraintData struct {
	propId PropId
	constraint Constraint
	channel chan *ChangeEntry	
}

func CreateConstraintData(propId PropId, 
		constraint Constraint, channel chan *ChangeEntry) *ConstraintData{
	cd := new(ConstraintData)
	cd.propId = propId
	cd.constraint = constraint
	cd.channel = channel
	return cd
}

func (this *RegistryStore) Close(){
	
	varids:=make([]VarId,len(this.varidsConnectedToConstraints))
	constraintDatas := make([]*ConstraintData,len(this.constraintsInterestedInVarids))
	i:=0
	for varid,_ := range this.varidsConnectedToConstraints {
		varids[i] = varid
		i+=1
	}
	
	i=0
	for constraintData,_ := range this.constraintsInterestedInVarids {
		constraintDatas[i]=constraintData
		i+=1
	}
	
	for _,constraintData := range constraintDatas {
		this.constraintsInterestedInVarids[constraintData] = nil
		close(constraintData.channel)
		//deprecated
		/*constraintData.channel = nil
		constraintData.constraint = nil
		constraintData.propId = 0*/
	}
	
	this.constraintsInterestedInVarids = nil
	this.varidsConnectedToConstraints = nil
	this.constraints = nil
//	println("JOU", len(this.constraints))
}

func (this *RegistryStore) AddVaridsConntectedToConstraint(varIds []VarId, constraintData *ConstraintData) {
	for _, varId := range varIds {
		if _, exists := this.varidsConnectedToConstraints[varId]; !exists {
			this.varidsConnectedToConstraints[varId] = make([]*ConstraintData,0) 
		}
		this.varidsConnectedToConstraints[varId] = append(this.varidsConnectedToConstraints[varId],constraintData)
	}
}

func (this *RegistryStore) AddConstraintInterestedInVarids(constraintData *ConstraintData, varids []VarId) {
	if _, exists := this.constraintsInterestedInVarids[constraintData]; !exists {
		this.constraintsInterestedInVarids[constraintData] = make([]VarId,0) 
	}
	
	this.constraintsInterestedInVarids[constraintData] = varids
}

// RemoveFixedRelations removes relations between varid and propId, if 
// domain of varid is ground
func (this *RegistryStore) RemoveFixedRelations(varid VarId) int {
	constraintData := this.varidsConnectedToConstraints[varid]

	removedConstraints:=0

	for _,constraintD := range constraintData {
		//entfernen	
		indexToRemove:=-1
		for i,vid := range this.constraintsInterestedInVarids[constraintD] {
			if vid==varid{
				indexToRemove = i
				break
			}	
		}
	
		this.constraintsInterestedInVarids[constraintD] = append(this.constraintsInterestedInVarids[constraintD][:indexToRemove], this.constraintsInterestedInVarids[constraintD][indexToRemove+1:]...)	
		if len(this.constraintsInterestedInVarids[constraintD]) == 0 {
			close(constraintD.channel)
/*			constraintD.channel = nil
			constraintD.constraint = nil
			constraintD.propId = 0*/
			// remove dangling reference to propagator, allow gc
			delete(this.constraints, constraintD.propId)
			delete(this.constraintsInterestedInVarids, constraintD)
			removedConstraints+=1
		}
	}
	
	this.varidsConnectedToConstraints[varid] = nil

	return removedConstraints
}

func (this *RegistryStore) ConnectedConstraints(varid VarId) []*ConstraintData{
	return this.varidsConnectedToConstraints[varid]
}