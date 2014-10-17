// Package core provides the central constraint store and the finite domain
// variables as interval domains and explicit domains.
package core

// PropId is a type for a unique identifier for Propagators
type PropId int

// Constraint interface to be implemented by every Propagator,
// for execution, registration, and copying in store
type Constraint interface {
	// Start starts propagator listening for events and propagating
	Start(store *Store)
	// Register registers itself at store and gets needed data from store
	// (channel, domains)
	Register(store *Store)
	// Clone clones the propagator (in particular for labeling)
	Clone() Constraint
	// SetID sets own Propagator-ID, only to be called initially by Store
	SetID(propID PropId)
	// GetID gets own Propagator-ID
	GetID() PropId
	// String returns a humanreadable string representation of the propagator
	String() string
	// GetVarIds returns the involved varids
	GetVarIds() []VarId
	// GetDomains returns the involved domains related with involved varids
	GetDomains() []Domain
	// GetInCh returns the inChannel of the propagator
	GetInCh() <-chan *ChangeEntry
	//GetOutCh returns the outChannel of the propagator
	GetOutCh() chan<- *ChangeEvent
}

// Domain interface to be implemented by every Domain implementation
type Domain interface {
	// GetMin returns the minimum value of a domain
	// panics on empty domains
	GetMin() int
	// GetMax returns the maximum value of a domain
	// panics on empty domains
	GetMax() int
	// GetMinAndMax returns min and max of a domain
	// panics on empty domains
	GetMinAndMax() (int, int)
	// GetAnyElement returns any one of the elements in this domain;
	// panics on empty domains
	GetAnyElement() int
	// Contains checks if ele is in a domain
	Contains(ele int) bool
	// Equals checks whether this domain contains the same
	// elements as the other domain
	Equals(other Domain) bool
	// Copy (deep) copies a domain
	Copy() Domain
	// Add adds ele to this domain
	Add(ele int)
	// Remove deletes ele from this domain
	Remove(ele int)
	// Removes deletes eles from this domain
	Removes(eles Domain)
	// RemovesWithOther deletes eles from this domain and from Other for
	// all elements which have no intersection with elements of this
	RemovesWithOther(eles Domain)
	// Values_asSlice returns a slice of values of the domain
	Values_asSlice() []int
	// Values_asMap returns a map of values of the domain
	Values_asMap() map[int]bool
	// Size returns the number of elements of this domain
	Size() int
	// IsEmpty returns true iff there is no element in this domain
	IsEmpty() bool
	// IsGround returns true iff there is one element in this domain
	IsGround() bool
	// SortedValues provides a slice of values in ascending order
	SortedValues() []int
	// String returns a human readable string representation of the propagator
	String() string
	// GetSubDomain() returns a Domain, which is a subdomain of current domain
	//GetSubDomainBounds(min, max int) Domain
	// GetDomainOutOfBounds returns a Domain, which contains values lower than
	// min and higher than max
	GetDomainOutOfBounds(min, max int) Domain
	// Intersection returns a new domain containing all values from the
	// intersection of this and domain
	Intersection(domain Domain) Domain
	// Union returns a new domain containing all values from the union
	// of this and domain
	Union(domain Domain) Domain
	// Difference returns a new domain containing all values that this
	// and domain have not in common
	Difference(domain Domain) Domain
	// IsSubset returns a boolean indicating if this is a subset of domain
	IsSubset(domain Domain) bool
}
