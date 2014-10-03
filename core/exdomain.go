package core

import (
	"fmt"
	"runtime/debug"
)

// ExDomain representation explicitly as a set of integer values based
// on builtin golang hashing.
type ExDomain struct {
	Values   map[int]bool
	Min, Max int
}

// Reset resets the domain to an empty domain.
func (this *ExDomain) Reset() {
	this.Values = make(map[int]bool)
	this.Min, this.Max = 1, 0
}

// CreateExDomain returns a pointer to a new empty ExDomain.
func CreateExDomain() *ExDomain {
	domain := new(ExDomain)
	domain.Values = make(map[int]bool)
	domain.Min, domain.Max = 1, 0
	return domain
}

// CreateExDomainAdds returns a pointer to a new Domain with
// elements from elements.
func CreateExDomainAdds(elements []int) *ExDomain {
	domain := new(ExDomain)
	domain.Values = make(map[int]bool)
	domain.Min, domain.Max = 1, 0
	domain.Adds(elements)
	return domain
}

// CreateExDomainFromIvDomain returns a pointer to a new
// Domain intialised with the values from the given domain.
func CreateExDomainFromIvDomain(dom *IvDomain) *ExDomain {
	return CreateExDomainAdds(dom.GetValues())
}

// CreateExDomainFromTo returns a pointer to a new Domain with
// elements from from to to (inclusive).
func CreateExDomainFromTo(from, to int) *ExDomain {
	domain := new(ExDomain)
	size := (to - from) + 1
	if size < 0 {
		size = 0
	}
	domain.Values = make(map[int]bool, size)
	for i := from; i <= to; i++ {
		domain.Values[i] = true
	}
	domain.Min, domain.Max = from, to
	return domain
}

// GetMin returns the smallest member of the integers in this domain;
// panics on empty domains.
func (this *ExDomain) GetMin() int {
	if this.IsEmpty() {
		logger.If("Domain %s", *this)
		debug.PrintStack()
		panic("GetMin on empty domain")
	}
	return this.Min
}

// GetMax returns the largest member of the integers in this domain;
// panics on empty domains.
func (this *ExDomain) GetMax() int {
	if this.IsEmpty() {
		logger.If("Domain %s", *this)
		debug.PrintStack()
		panic("GetMax on empty domain")
	}
	return this.Max
}

// GetMinAndMax returns the smallest and the largest member of the
// integers in this domain; panics on empty domains.
func (this *ExDomain) GetMinAndMax() (int, int) {
	if this.IsEmpty() {
		logger.If("Domain %s", *this)
		debug.PrintStack()
		panic("GetMinMax on empty domain")
	}
	return this.Min, this.Max
}

// GetAnyElement returns any one of the elements in this domain;
// panics on empty domains.
func (this *ExDomain) GetAnyElement() int {
	if this.Min > this.Max {
		logger.If("Domain %s", *this)
		debug.PrintStack()
		panic("GetAnyElement on empty domain")
	}
	return this.Min // whatever
}

// Contains checks if ele is in Domain.
func (this *ExDomain) Contains(ele int) bool {
	_, exists := this.Values[ele]
	return exists
}

// Equals checks whether this domain contains the same
// elements as the other domain.
func (this *ExDomain) Equals(otherDomain Domain) bool {
	var otherDom *ExDomain
	if other, k := otherDomain.(*ExDomain); k {
		otherDom = other
	} else if other, k := otherDomain.(*IvDomain); k {
		otherDom = CreateExDomainFromIvDomain(other)
	} else {
		panic(fmt.Sprintf("Domain %s is not comparable with Domain %s",
			this, other))
	}
	for k := range this.Values {
		if _, exists := otherDom.Values[k]; !exists { // inline Contains
			return false
		}
	}
	for k := range otherDom.Values {
		if _, exists := this.Values[k]; !exists { // inline Contains
			return false
		}
	}
	return true
}

// Copy (deep) copies a domain.
func (this *ExDomain) Copy() Domain {
	domain := new(ExDomain)
	domain.Values = make(map[int]bool, len(this.Values))
	for k, v := range this.Values {
		if v {
			domain.Values[k] = v
		}
	}
	domain.Min, domain.Max = this.Min, this.Max
	return domain
}

// Add adds ele to this domain.
func (this *ExDomain) Add(ele int) {
	this.Values[ele] = true
	if this.Min > this.Max { // empty
		this.Min, this.Max = ele, ele
	} else if ele < this.Min {
		this.Min = ele
	} else if ele > this.Max {
		this.Max = ele
	}
}

// Adds adds all values in elements to this domain.
func (this *ExDomain) Adds(elements []int) {
	for _, ele := range elements {
		this.Add(ele)
	}
}

// Remove deletes ele from this domain.
func (this *ExDomain) Remove(ele int) {
	if _, exists := this.Values[ele]; !exists {
		return
	}
	delete(this.Values, ele)
	if ele == this.Min {
		if ele == this.Max {
			this.Min, this.Max = 1, 0 // empty domain
		} else {
			min := MaxInt
			for val := range this.Values {
				if val < min {
					min = val
				}
			}
			this.Min = min
		}
	} else if ele == this.Max { // but not min, thus at least two
		max := MinInt
		for val := range this.Values {
			if val > max {
				max = val
			}
		}
		this.Max = max
	}
}

// removesIvDomain from the given values of an IvDomain from local
// data storage.
func (this *ExDomain) removesIvDomain(ivdom *IvDomain,
	modifyingOtherDomain bool) {
	if modifyingOtherDomain {
		dels := make([]int, 0)
		for _, part := range ivdom.GetParts() {
			for v := part.From; v <= part.To; v++ {
				if this.Contains(v) {
					this.Remove(v)
				} else {
					dels = append(dels, v)
				}
			}
		}
		if len(dels) != 0 {
			ivdom.Removes(CreateIvDomainFromIntArr(dels))
		}
	} else {
		for _, part := range ivdom.GetParts() {
			for v := part.From; v <= part.To; v++ {
				this.Remove(v)
			}
		}
	}
}

func (this *ExDomain) removesDomain(dom *ExDomain, modifyingOtherDomain bool) {
	if modifyingOtherDomain {
		for ele := range dom.GetValues() { // a map
			if this.Contains(ele) {
				this.Remove(ele)
			} else {
				dom.Remove(ele)
			}
		}
	} else {
		for ele := range dom.GetValues() {
			this.Remove(ele)
		}
	}
	// ToDo: setting max/min after removing all (only remembering,
	// if max/min has changed)
}

// RemovesWithOther, ToDo: explain
func (this *ExDomain) RemovesWithOther(eles Domain) {
	if d, ok := eles.(*IvDomain); ok {
		this.removesIvDomain(d, true)
	} else {
		d := eles.(*ExDomain)
		this.removesDomain(d, true)
	}
}

// Removes deletes eles from this domain.
func (this *ExDomain) Removes(eles Domain) {
	if d, ok := eles.(*ExDomain); ok {
		this.removesDomain(d, false)
	} else {
		d := eles.(*IvDomain)
		this.removesIvDomain(d, false)
	}
}

// Values_asSlice returns a slice of values of the domain.
func (this *ExDomain) Values_asSlice() []int {
	return Keys_MapIntToBool(this.GetValues())
}

// GetValues returns the values of this domain.
func (this *ExDomain) GetValues() map[int]bool {
	return this.Values
}

// Values_asMap just returns the values.
func (this *ExDomain) Values_asMap() map[int]bool {
	return this.GetValues()
}

// Size returns the number of elements of this domain.
func (this *ExDomain) Size() int {
	return len(this.Values) // assumes no mappings to false
}

// IsEmpty returns true iff there is no element in this domain.
func (this *ExDomain) IsEmpty() bool {
	return this.Min > this.Max
}

// IsGround returns true iff there is one element in this domain.
func (this *ExDomain) IsGround() bool {
	return this.Min == this.Max
}

// SortedValues provides a slice of values in ascending order.
func (this *ExDomain) SortedValues() []int {
	return SortedKeys_MapIntToBool(this.Values)
}

// String returns a sorted string representation of this domain.
func (this *ExDomain) String() string {
	return BeautifulOutput("", this.Values)
}

// GetSubDomainBounds returns a new domain with all values
// within the given boundaries min to max inclusive.
func (this *ExDomain) GetSubDomainBounds(min, max int) Domain {
	d := CreateExDomain()
	for _, v := range this.SortedValues() {
		if v >= min && v <= max {
			d.Add(v)
		}
	}
	return d
}

// GetDomainOutOfBounds returns a new domain with all values
// outside of the given boundaries.
func (this *ExDomain) GetDomainOutOfBounds(min, max int) Domain {
	d := CreateExDomain()
	for _, v := range this.SortedValues() {
		if v < min || v > max {
			d.Add(v)
		}
	}
	return d
}

// Intersection returns a new domain containing all values from
// the intersection of this and domain
func (this *ExDomain) Intersection(domain Domain) Domain {
	newDomain := CreateExDomain()
	domainValues := domain.Values_asMap()
	for key := range this.Values {
		if domainValues[key] {
			newDomain.Add(key)
		}
	}
	return newDomain
}

// Union returns a new domain containing all values from the union of
// this and domain
func (this *ExDomain) Union(domain Domain) Domain {
	newDomain := CreateExDomain()
	newDomainValues := this.Values_asMap()
	for k := range newDomainValues {
		newDomain.Add(k)
	}
	domainValues := domain.Values_asMap()
	for key := range domainValues {
		if !newDomainValues[key] {
			newDomain.Add(key)
		}
	}
	return newDomain
}

// Difference returns a new domain containing all values that this and
// domain have not in common
func (this *ExDomain) Difference(domain Domain) Domain {
	newDomain := CreateExDomain()
	domainValues := domain.Values_asMap()
	for key := range this.Values {
		if !domainValues[key] {
			newDomain.Add(key)
		}
	}

	for key := range domainValues {
		if !this.Values[key] {
			newDomain.Add(key)
		}
	}
	return newDomain
}

// IsSubset returns a boolean indicating if this is a subset of domain.
func (this *ExDomain) IsSubset(domain Domain) bool {
	domainValues := domain.Values_asMap()
	for key := range this.Values {
		if !domainValues[key] {
			return false
		}
	}
	return true
}
