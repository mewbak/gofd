package core

import (
	"fmt"
	"strings"
)

// ChangeEvent signaling to a store to delete domain values from domains.
// Contains multiple ChangeEntry objects. It may contain several ChangeEntry
// objects for one changing variable/domain
type ChangeEvent struct {
	changes []*ChangeEntry
}

// IsEmpty returns true, iff this has no changeEntry
func (this *ChangeEvent) IsEmpty() bool {
	return len(this.changes) == 0
}

// Clone creates a new ChangeEvent with the values given by "this" (deepcopy)
func (this *ChangeEvent) Clone() *ChangeEvent {
	changeEvent := CreateChangeEvent()
	// try to be fast for empty ChangeEvents
	for _, ch := range this.changes {
		changeEvent.AddChangeEntry(ch.Clone())
	}
	return changeEvent
}

// Equals equals two ChangeEvents
func (this *ChangeEvent) Equals(other *ChangeEvent) bool {
	if len(this.changes) != len(other.changes) {
		return false
	}
	// try to be fast for empty ChangeEvents
	for i := 0; i < len(this.changes); i++ {
		if !this.changes[i].Equals(other.changes[i]) {
			return false
		}
	}
	return true
}

// CreateChangeEvent creates a new empty ChangeEvent
func CreateChangeEvent() *ChangeEvent {
	changeEvent := new(ChangeEvent)
	// try to be fast for empty ChangeEvents
	changeEvent.changes = make([]*ChangeEntry, 0, 0)
	return changeEvent
}

// AddChangeEntry adds a change entry to a change event
func (this *ChangeEvent) AddChangeEntry(change *ChangeEntry) {
	this.changes = append(this.changes, change)
}

// String returns a readable string-representation of this ChangeEvent
func (this *ChangeEvent) String() string {
	evtsString := make([]string, len(this.changes))
	for k, v := range this.changes {
		evtsString[k] = fmt.Sprintf("%s", v)
	}
	return "ChangeEvent[" + strings.Join(evtsString, ", ") + "]"
}

// ChangeEntry domain values to be removed from one domain of a variable
type ChangeEntry struct {
	varId        VarId  // which variable has changed
	domainvalues Domain // which values to remove
}

// Equals equals two ChangeEntries
func (this *ChangeEntry) Equals(other *ChangeEntry) bool {
	if this.varId != other.varId {
		return false
	}
	return this.domainvalues.Equals(other.domainvalues)
}

// Clone creates a new ChangeEntry with the values given by "this"
func (this *ChangeEntry) Clone() *ChangeEntry {
	entry := CreateChangeEntry(this.varId)
	entry.domainvalues = this.domainvalues.Copy()
	return entry
}

// CreateChangeEntry creates a new instance of ChangeEntry, which can
// be added to a ChangeEvent; one ChangeEntry per VarId.
// Parameter varId of the Variable to delete the value from.
func CreateChangeEntry(varId VarId) *ChangeEntry {
	entry := new(ChangeEntry)
	entry.varId = varId
	entry.domainvalues = CreateIvDomain()
	return entry
}

// CreateChangeEntryWithValues creates a new instance of ChangeEntry,
// which can be added to a ChangeEvent; one ChangeEntry per VarId.
// Parameter varId of the variable to delete the value from.
func CreateChangeEntryWithValues(varId VarId, values Domain) *ChangeEntry {
	entry := new(ChangeEntry)
	entry.varId = varId
	entry.domainvalues = values
	return entry
}

// CreateChangeEntryWithIntValues creates a new instance of ChangeEntry,
// which can be added to a ChangeEvent; one ChangeEntry per VarId.
// Parameter varId of the Variable to delete the value from.
func CreateChangeEntryWithIntValues(varId VarId, values []int) *ChangeEntry {
	entry := new(ChangeEntry)
	entry.varId = varId
	entry.domainvalues = CreateIvDomainFromIntArr(values)
	return entry
}

// CreateChangeEntryWithIntValue creates a new instance of ChangeEntry,
// which can be added to a ChangeEvent; one ChangeEntry per VarId.
// Parameter varId of the Variable to delete the value from.
func CreateChangeEntryWithIntValue(varId VarId, value int) *ChangeEntry {
	entry := new(ChangeEntry)
	entry.varId = varId
	entry.domainvalues = CreateIvDomainFromTo(value, value)
	return entry
}

// GetID returns the varId of the variable of the affected domain
func (this *ChangeEntry) GetID() VarId {
	return this.varId
}

// IsEmpty returns true, iff no values are to be removed
func (this *ChangeEntry) IsEmpty() bool {
	return this.domainvalues.IsEmpty()
}

// GetValues returns the values to be removed
func (this *ChangeEntry) GetValues() Domain {
	return this.domainvalues
}

// SetValues sets the sorted values/subdomain to be removed
func (this *ChangeEntry) SetValues(values Domain) {
	this.domainvalues = values
}

// SetValuesByIntArr sets the unsorted or sorted values to be removed
func (this *ChangeEntry) SetValuesByIntArr(values []int) {
	this.domainvalues = CreateIvDomainFromIntArr(values)
}

// Remove do not remove that value due to this ChangeEntry
func (this *ChangeEntry) Remove(value int) {
	this.domainvalues.Remove(value)
}

// Add adds a new value to a ChangeEntry
func (this *ChangeEntry) Add(val int) {
	this.domainvalues.Add(val)
}

// Add adds a new value to a ChangeEntry
func (this *ChangeEntry) AddValues(values map[int]bool) {
	for val, _ := range values {
		this.domainvalues.Add(val)
	}
}

// String returns a readable string-representation of this ChangeEntry
func (this *ChangeEntry) String() string {
	return fmt.Sprintf("ChangeEntry{%d, %s}",
		this.varId, this.domainvalues)
}
