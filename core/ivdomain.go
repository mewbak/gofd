package core

import (
	"fmt"
	"runtime/debug"
	"sort"
)

// IvDomain represents a Domain as list of intervals, where each interval
// is a set of consecutive values represented as tuple (from, to) where
// from <= to must always hold. Example: [[0,5], [10,50], [100,200]]
type IvDomain struct {
	partList *SortedIvDomPartList
}

// CreateIvDomainFromDomain creates an IvDomain from a Domain in
// explicit representation.
func CreateIvDomainFromDomain(dom *ExDomain) *IvDomain {
	return CreateIvDomainFromIntArr(dom.Values_asSlice())
}

// CreateIvDomainFromTos creates an interval domain from given
// fromTo pairs.
func CreateIvDomainFromTos(fromTos [][]int) *IvDomain {
	domain := CreateIvDomain()
	domain.partList = CreateSortedIvDomPartListWithSortedParts(
		CreateIvDomParts(fromTos))
	return domain
}

// CreateIvDomain creates an interval domain.
func CreateIvDomain() *IvDomain {
	domain := new(IvDomain)
	domain.partList = CreateSortedIvDomPartList()
	return domain
}

// CreateIvDomainUnion returns a domain merged from parts with help of
// a set union operation applied to all intervals. If intervals overlap,
// then they are joined together.
func CreateIvDomainUnion(parts []*IvDomPart) *IvDomain {
	domain := CreateIvDomain()
	domain.partList = CreateSortedIvDomPartList()
	domain.partList.UnionWithParts(parts)
	return domain
}

// CreateIvDomainFromTo creates a new IvDomain (domain with internal interval
// representation of values) from given from and to values (only one IvDomPart)
func CreateIvDomainFromTo(from, to int) *IvDomain {
	domain := CreateIvDomain()
	part := CreateIvDomPart(from, to)
	domain.partList = CreateSortedIvDomPartListWithSortedParts([]*IvDomPart{part})
	return domain
}

// CreateIvDomainDomParts initialises a new domain with the given parts
func CreateIvDomainDomParts(parts []*IvDomPart) *IvDomain {
	dom := CreateIvDomain()
	dom.partList = CreateSortedIvDomPartListWithSortedParts(parts)
	return dom
}

// CreateIvDomainFromIntArr creates an IvDomain and inits its IvDomParts with
// the given int slice. The int-slice will be sorted first.
func CreateIvDomainFromIntArr(eles []int) *IvDomain {
	domain := CreateIvDomain()
	domain.partList = CreateSortedIvDomPartListFromIntArr(eles)
	return domain
}

// CreateIvDomainDomPartsWithSort init a new domain with the given parts,
// but the parts firt will be sorted.
func CreateIvDomainDomPartsWithSort(parts []*IvDomPart) *IvDomain {
	dom := CreateIvDomain()
	p := SortableDomPart(parts)
	sort.Sort(p)
	dom.partList = CreateSortedIvDomPartListWithSortedParts(p)
	return dom
}

// GetParts returns domain-parts (intervals).
func (this *IvDomain) GetParts() []*IvDomPart {
	return this.partList.GetParts()
}

// Equals checks whether this domain contains the IvDomParts
// elements as the other domain.
func (this *IvDomain) Equals(otherDomain Domain) bool {
	var otherDom *IvDomain
	if other, k := otherDomain.(*IvDomain); k {
		otherDom = other
	} else if other, k := otherDomain.(*ExDomain); k {
		otherDom = CreateIvDomainFromDomain(other)
	} else {
		msg := "Domain %s is not comparable with Domain %s"
		panic(fmt.Sprintf(msg, this, other))
	}
	return this.partList.Equals(otherDom.partList)
}

// Copy (deep) copies a domain.
func (this *IvDomain) Copy() Domain {
	domain := new(IvDomain)
	domain.partList = this.partList.Copy()
	return domain
}

// IsEmpty returns true iff there is no element in this domain.
func (this *IvDomain) IsEmpty() bool {
	return this.partList.IsEmpty()
}

// IsGround returns true iff there is one element in this domain.
func (this *IvDomain) IsGround() bool {
	return (len(this.partList.GetParts()) == 1) &&
		(this.partList.GetParts()[0].IsGround())
}

// GetMax returns the largest member of the integers in this domain;
// panics on empty domains.
func (this *IvDomain) GetMax() int {
	if this.IsEmpty() {
		logger.If("Domain %s", *this)
		debug.PrintStack()
		panic("GetMax on empty domain")
	}
	return this.partList.GetMax()
}

// GetMin returns the smallest member of the integers in this domain;
// panics on empty domains.
func (this *IvDomain) GetMin() int {
	if this.IsEmpty() {
		logger.If("Domain %s", *this)
		debug.PrintStack()
		panic("GetMin on empty domain")
	}
	return this.partList.GetMin()
}

// GetMinAndMax returns the smallest and the largest member of the
// integers in this domain; panics on empty domains
func (this *IvDomain) GetMinAndMax() (int, int) {
	return this.GetMin(), this.GetMax()
}

// GetAnyElement returns any one of the elements in this domain;
// panics on empty domains.
func (this *IvDomain) GetAnyElement() int {
	if this.IsEmpty() {
		logger.If("Domain %s", *this)
		debug.PrintStack()
		panic("GetAnyElement on empty domain")
	}
	return this.GetMin()
}

// Contains checks if ele is in Domain.
func (this *IvDomain) Contains(ele int) bool {
	return this.partList.Contains(ele)
}

// Reset resets the interval domain to the empty Domain.
func (this *IvDomain) Reset() {
	this.partList = CreateSortedIvDomPartList()
}

// Removes removes the parts of a domain from the current domain
// in contrast to Difference, it manipulates the curren domain instead of
// returning a new domain; will be called from Propagators.
func (this *IvDomain) removesSubIvDomain(other *IvDomain,
	modifyingOtherDomain bool) {
	// remove difference or in other words: take the intersection
	var thisCopy *IvDomain
	if modifyingOtherDomain {
		thisCopy = this.Copy().(*IvDomain)
	}
	this.partList.Removes(other.partList)
	if modifyingOtherDomain {
		other.partList = thisCopy.partList.Intersection(other.partList)
	}
}

// will be called from store
func (this *IvDomain) removesSubDomain(other *ExDomain,
	modifyingOtherDomain bool) {
	//remove difference or in other words: take the intersection
	var thisCopy *IvDomain
	if modifyingOtherDomain {
		thisCopy = this.Copy().(*IvDomain)
	}
	this.partList.RemovesSortedInts(other.SortedValues())
	if modifyingOtherDomain {
		intersectionPerPart := thisCopy.IntersectionWithDomainThisAssoc(other)
		other.Reset()
		for _, v := range intersectionPerPart {
			other.Adds(v)
		}
	}
}

// RemovesWithOther
// ToDo: Explain
func (this *IvDomain) RemovesWithOther(eles Domain) {
	// i.e. [[0, 5], [10, 50], [55, 100]]
	// eles: 4, 5, 8, 10, 45
	if d, ok := eles.(*IvDomain); ok {
		this.removesSubIvDomain(d, true)
	} else {
		d := eles.(*ExDomain)
		this.removesSubDomain(d, true)
	}
}

// Removes removes eles (Domain) from local data storage
// modifyingOtherDomain specifies if the eles domain should be modified,
// if some values are not in the current Domain
func (this *IvDomain) Removes(eles Domain) {
	// i.e. [[0, 5], [10, 50], [55, 100]]
	// eles: 4, 5, 8, 10, 45
	if d, ok := eles.(*IvDomain); ok {
		this.removesSubIvDomain(d, false)
	} else {
		d := eles.(*ExDomain)
		this.removesSubDomain(d, false)
	}
}

func getIntersectionAsSortedDomPartList(intersectionMap map[int][]*IvDomPart) []*IvDomPart {
	k := 0
	i := 0
	parts := make([]*IvDomPart, 0)
	for k != len(intersectionMap) {
		if intersecPart, ok := intersectionMap[i]; ok {
			parts = append(parts, intersecPart...)
			k += 1
		}
		i += 1
	}
	return parts
}

// IntersectionWithDomainThisAssoc calculates the intersection between this
// and dom. It returns a map: key is the part-index of this.parts, value
// is a list of intersections in the specific part (at part-index).
// intersection between this and dom
// example:
// 	this([1..12])
//	dom([1..4],[10..15])
//  this function returns [0:[1,4],1:[10,12]], the intersection
func (this *IvDomain) IntersectionWithDomainThisAssoc(dom *ExDomain) map[int][]int {
	eles := dom.SortedValues()
	matchingEles := make(map[int][]int)
	difference := make([]int, 0)
	var match bool
	// take ele and look, in which part it belongs to
	for _, ele := range eles {
		match = false
		for i, part := range this.partList.GetParts() {
			if part.ContainsInt(ele) {
				match = true
				matchingEles[i] = append(matchingEles[i], ele)
				break
			} else if part.GT(ele) { // cause of sorted parts
				break
			}
		}
		// we use index "-1" in the map for all values, which should be removed
		// from the otherDomain "dom"
		if !match {
			difference = append(difference, ele)
		}
	}
	return matchingEles
}

// IntersectionWithIvDomainThisAssoc calculates the intersection between this
// and dom. It returns a map: key is the part-index of this.parts, value
// is a list of intersections in the specific part (at part-index).
// intersection between this and dom
// this.IntersectionWithIvDomainThisAssoc(dom) returns the same values as
// dom.IntersectionWithIvDomainThisAssoc(this), but with different map-indexes.
// example:
// 	this([1..12])
//	dom([1..4], [10..15])
//  this function returns [0:[1, 4], 0:[10, 12]], the intersection
func (this *IvDomain) IntersectionWithIvDomainThisAssoc(dom *IvDomain) map[int][]*IvDomPart {
	domparts := dom.GetParts()
	matching := make(map[int][]*IvDomPart)
	for _, dpart := range domparts {
		for i, part := range this.partList.GetParts() {
			if part.GT_DP(dpart) { //cause of sorted parts
				// out of boundaries
				break
			} else if part.LT_DP(dpart) {
				continue
			} else {
				newp := part.INTERSECTION(dpart)
				if newp != nil {
					matching[i] = append(matching[i], newp)
				}
			}
		}
	}
	return matching
}

// --- set operations ---

// IntersectionIvDomain returns a domain, which represents the intersection
// between this and dom
func (this *IvDomain) IntersectionIvDomain(dom *IvDomain) *IvDomain {
	d := CreateIvDomain()
	d.partList = this.partList.Intersection(dom.partList)
	return d
}

// DifferenceWithIvDomain calculates the differences between each part of
// this.parts and dom.parts. It returns the difference
// example
//   this([1..3][6..9][15..20])
//   dom([1..2][5..10] [12..12] [16..19])
//	 result: newdom([3..3][15..15][20..20])
func (this *IvDomain) DifferenceWithIvDomain(dom *IvDomain) *IvDomain {
	thisCopy := this.Copy().(*IvDomain)
	thisCopy.removesSubIvDomain(dom, false)
	return thisCopy
}

// --- end set operations ---

// SetParts updates the local storage with the given parts (parts and
// min/max updating)
func (this *IvDomain) SetParts(parts []*IvDomPart) {
	this.partList = CreateSortedIvDomPartListWithSortedParts(parts)
}

// Remove removes ele from local data storage
// (CAUTION: only for success from external)
func (this *IvDomain) Remove(ele int) {
	this.partList.Remove(ele)
}

// String returns a sorted string representation of this domain.
func (this *IvDomain) String() string {
	return fmt.Sprintf("%s", this.partList)
}

// Size returns the number of elements of this domain.
func (this *IvDomain) Size() int {
	return this.partList.Size()
}

// Values_asMap returns a map representation of an IvDomain
func (this *IvDomain) Values_asMap() map[int]bool {
	return this.partList.Values_asMap()
}

// GetValues returns the IvDomain as a slice.
func (this *IvDomain) GetValues() []int {
	return this.Values_asSlice()
}

// Values_asSlice returns the IvDomain as a slice.
func (this *IvDomain) Values_asSlice() []int {
	return this.partList.Values_asSlice()
}

// SortedValues returns the IvDomain as a sorted slice.
func (this *IvDomain) SortedValues() []int {
	return this.Values_asSlice()
}

// ToDo: Test
// Append appends a part to the current part list. Attention: the given
// part must be greater than the parts in this-Domain
func (this *IvDomain) Append(part *IvDomPart) {
	lenDom := len(this.GetParts())
	if part.GT_DP(this.GetParts()[lenDom-1]) && this.GetParts()[lenDom-1].To+1 < part.From {
		this.partList.Append(part)
		return
	}
	panic("IvDomain.Append error: part have to be greater than the greatest" +
		"part in current IvDomain (from > greatestPart.To+1)")
}

// ToDo: Test
// Append appends a part to the current part list. Attention: the given
// parts must be sorted and the first one must be greater than the
// parts in this-Domain
func (this *IvDomain) Appends(parts []*IvDomPart) {
	lenDom := len(this.GetParts())
	if len(parts) > 0 {
		if parts[0].GT_DP(this.GetParts()[lenDom-1]) {
			this.partList.Appends(parts)
			return
		}
		panic("IvDomain.Append error: part have to be greater than the" +
			"parts in current IvDomain (this)")
	}
}

// AddAnyPart adds a part to the current domain
// in contrast to other add/insert functions, this function
// doesn't rely on conventions. You can add any part, you want
// (even in incorrect order).
// useful for union (see CreateIvDomainUnion)
func (this *IvDomain) AddAnyPart(part *IvDomPart) {
	this.partList.AddAnyPart(part)
}

// Add adds an element to the current IvDomain
func (this *IvDomain) Add(ele int) {
	this.partList.Add(ele)
}

// GetDomainOutOfBounds returns all values lower min or upper max
// and contained in this-IvDomain.
func (this *IvDomain) GetDomainOutOfBounds(min, max int) Domain {
	d := CreateIvDomainFromTo(min, max)
	return this.DifferenceWithIvDomain(d)
}

// ADD returns a new domain with len(this.parts)*len(dom.parts) parts. Every
// part of this.parts will be added with every part of dom.parts. The result
// will be added to the new domain. For the case, that some intervals have
// intersections, a union of the intervals will be created.
// Useful for A+B=C (Domain A plus Domain B). For this example, there must be
// calculated all potential intervals for C (through A+B) and
// the given function ADD do this. It returns all potential intervals for C,
// if this == A and dom == B
// e.g. Domain([1,2],[4,6]) + [3,3] = Domain([4,5],[7,9])
func (this *IvDomain) ADD(dom *IvDomain) *IvDomain {
	newparts := make([]*IvDomPart,
		len(this.partList.GetParts())*len(dom.partList.GetParts()))
	jlen := len(dom.GetParts())
	for i, part1 := range this.GetParts() {
		for j, part2 := range dom.GetParts() {
			res := part1.ADD(part2)
			newparts[(i*jlen)+j] = res
		}
	}
	// union, cause of potential overlapping parts
	return CreateIvDomainUnion(newparts)
}

// SUBTRACT returns a new domain with len(this.parts)*len(dom.parts) parts.
// Every part of this.parts will be subtracted with every part of dom.parts.
// The result will be added to the new domain. For the case, that some
// intervals have intersections, a union of the intervals will be created.
// Useful for C-B=A (Domain C minus Domain B). For the given example, there
// must be calculated all potential intervals for A (through C-B) and
// the given function SUBTRACT do this. It returns all potential intervals for
// A, if this == C and dom == B
func (this *IvDomain) SUBTRACT(dom *IvDomain) *IvDomain {
	newparts := make([]*IvDomPart,
		len(this.partList.GetParts())*len(dom.partList.GetParts()))
	jlen := len(dom.GetParts())
	for i, part1 := range this.GetParts() {
		for j, part2 := range dom.GetParts() {
			newparts[(i*jlen)+j] = part1.SUBTRACT(part2)
		}
	}
	//union, cause of potential overlapping parts
	return CreateIvDomainUnion(newparts)
}

// NEGATE returns a domain with negated parts.
func (this *IvDomain) NEGATE() *IvDomain {
	thisParts := this.GetParts()
	dom := CreateIvDomain()
	parts := make([]*IvDomPart, len(thisParts))
	for i := len(thisParts) - 1; i >= 0; i-- {
		parts[AbsInt(i-(len(thisParts)-1))] = thisParts[i].NEG()
	}
	dom.SetParts(parts)
	return dom
}

// ABS returns a domain with positive parts (all froms and tos are positive).
func (this *IvDomain) ABS() *IvDomain {
	thisParts := this.GetParts()
	parts := make([]*IvDomPart, len(thisParts))
	for i, part := range this.partList.GetParts() {
		if part.From < 0 {
			parts[i] = part.ABS()
		} else {
			parts[i] = part.Copy()
		}
	}
	// Union necessary, example: abs([(-5..2), (4,8)]) --> (0,8)
	return CreateIvDomainUnion(parts)
}

// NOT returns a domain except the elements in current domain.
func (this *IvDomain) NOT() *IvDomain {
	p := CreateIvDomPart(NEG_INFINITY, INFINITY)
	notD := CreateIvDomainDomParts(DIFFERENCE_DomParts(p, this.GetParts()...))
	return notD
}

// --- compatiblity with explicit domain representation ---

// DifferenceWithDomain calculates the differences between each part of
// this.parts and dom.parts. It returns the difference
// example
//   this([1..3][6..9][15..20])
//   dom([1..2][5..10] [12..12] [16..19])
//   result: newdom([3..3][15..15][20..20])
func (this *IvDomain) DifferenceWithDomain(dom *ExDomain) *IvDomain {
	intersectionPerPart := this.IntersectionWithDomainThisAssoc(dom)
	diffParts := this.DifferenceWithIntersectionIntsParts(intersectionPerPart)
	return CreateIvDomainDomParts(diffParts)
}

// DifferenceWithIntersectionIntsParts has a way too long name, ToDo
func (this *IvDomain) DifferenceWithIntersectionIntsParts(intersectionPerPart map[int][]int) []*IvDomPart {
	var parts []*IvDomPart
	if len(intersectionPerPart) > 0 {
		newParts := make([]*IvDomPart, 0)
		for i := 0; i < len(this.partList.GetParts()); i++ {
			if _, ok := intersectionPerPart[i]; ok {
				eles := intersectionPerPart[i]
				newP := DIFFERENCE_Ints(this.partList.GetParts()[i], eles...)
				if newP != nil {
					newParts = append(newParts, newP...)
				}
			} else {
				newParts = append(newParts, this.partList.GetParts()[i])
			}
		}
		parts = newParts
	} else {
		parts = this.Copy().(*IvDomain).GetParts()
	}
	return parts
}

// ToDo: test

// MULTIPLY returns a new domain with len(this.parts)*len(dom.parts) parts.
// Every part of this.parts will be multiplied with every part of dom.parts.
// The result will be added to the new domain. For the case, that some
// intervals have intersections, a union of the intervals will be created at
// the end.
//
// Useful for A*B=C (Domain A multiplied with Domain B). For this example,
// there must be calculated all potential intervals for C (through A*B) and
// the given function MULTIPLY do this. It returns all potential intervals for
// C, if this == A and dom == B
func (this *IvDomain) MULTIPLY(dom *IvDomain) *IvDomain {
	newparts := make([]*IvDomPart,
		len(this.partList.GetParts())*len(dom.GetParts()))
	jlen := len(dom.GetParts())
	for i, part1 := range this.GetParts() {
		for j, part2 := range dom.GetParts() {
			res := part1.MULTIPLY(part2)
			newparts[(i*jlen)+j] = res
		}
	}
	// union, cause of potential overlapping parts
	return CreateIvDomainUnion(newparts)
}

// ToDo: test

// DIvIDE divides the parts of the first domain with ones of the
// second domain and returns a new domain.
func (this *IvDomain) DIvIDE(dom *IvDomain) *IvDomain {
	newparts := make([]*IvDomPart,
		len(this.partList.GetParts())*len(dom.GetParts()))
	jlen := len(dom.GetParts())
	for i, part1 := range this.GetParts() {
		for j, part2 := range dom.GetParts() {
			res := part1.DIvIDE(part2)
			newparts[(i*jlen)+j] = res
		}
	}
	// union, cause of potential overlapping parts
	return CreateIvDomainUnion(newparts)
}

// Intersection returns a new domain containing all values from
// the intersection of this and domain.
func (this *IvDomain) Intersection(domain Domain) Domain {
	newDomain := CreateIvDomain()
	domainValues := domain.Values_asMap()
	for key, _ := range this.Values_asMap() {
		if domainValues[key] {
			newDomain.Add(key)
		}
	}
	return newDomain
}

// Union returns a new domain containing all values from the union of this
// and domain.
func (this *IvDomain) Union(domain Domain) Domain {
	newDomain := CreateIvDomain()
	newDomainValues := this.Values_asMap()
	for k, _ := range newDomainValues {
		newDomain.Add(k)
	}
	domainValues := domain.Values_asMap()
	for key, _ := range domainValues {
		if !newDomainValues[key] {
			newDomain.Add(key)
		}
	}
	return newDomain
}

// Difference returns a new domain containing all values that this and
// domain have not in common.
func (this *IvDomain) Difference(domain Domain) Domain {
	newDomain := CreateIvDomain()
	domainValues := domain.Values_asMap()
	for key, _ := range this.Values_asMap() {
		if !domainValues[key] {
			newDomain.Add(key)
		}
	}
	for key, _ := range domainValues {
		if !this.Values_asMap()[key] {
			newDomain.Add(key)
		}
	}
	return newDomain
}

// IsSubset returns a boolean indicating if this is a subset of domain
func (this *IvDomain) IsSubset(domain Domain) bool {
	domainValues := domain.Values_asMap()
	for key, _ := range this.Values_asMap() {
		if !domainValues[key] {
			return false
		}
	}
	return true
}
