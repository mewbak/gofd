package core

import (
	"fmt"
	"strings"
)

// SortedIvDomPartList is a list with sorted IvDomParts
// implementation with SliceTricks.
// https://code.google.com/p/go-wiki/wiki/SliceTricks
type SortedIvDomPartList struct {
	parts []*IvDomPart
}

// --- create functions ---

// CreateSortedIvDomPartList returns an empty list with initialized
// min/max values (min greater max)
func CreateSortedIvDomPartList() *SortedIvDomPartList {
	l := new(SortedIvDomPartList)
	return l
}

// CreateSortedIvDomPartListWithSortedParts creates an IvFomPart array
// and directly puts (no copy) the parts there by reference.
func CreateSortedIvDomPartListWithSortedParts(parts []*IvDomPart) *SortedIvDomPartList {
	l := new(SortedIvDomPartList)
	l.parts = parts
	return l
}

// CreateSortedIvDomPartListFromIntArr creates an IvDomain and initialzes
// it with the given int-slice. The int-slice will be sorted first.
func CreateSortedIvDomPartListFromIntArr(eles []int) *SortedIvDomPartList {
	if len(eles) == 0 {
		return CreateSortedIvDomPartList()
	}
	SortIntArr(eles)
	l := CreateSortedIvDomPartList()
	begin := eles[0]
	end := eles[0]
	next := eles[0] + 1
	for _, v := range eles[1:] {
		if v == next {
			next += 1
		} else {
			end = next - 1
			l.parts = append(l.parts, CreateIvDomPart(begin, end))
			begin = v
			end = v
			next = v + 1
		}
	}
	end = next - 1
	l.parts = append(l.parts, CreateIvDomPart(begin, end))
	return l
}

//	--- manipulation functions ---

// UnionWithParts needs a comment, ToDo
func (this *SortedIvDomPartList) UnionWithParts(parts []*IvDomPart) {
	for _, part := range parts {
		this.AddAnyPart(part)
	}
}

// AddAnyPart adds a part to the current domain
// in contrast to other add/insert functions, this function
// doesn't rely on conventions. You can add any part, you want
// (even in incorrect order).
// useful for union (see CreateIvDomainUnion)
// Examples:
// d:[1,4]  [10,15]
// addParts:
// [2,12] --> merge		[1,15]
// [5,9] --> merge		[1,15]
// [4,10] --> merge		[1,15]
// [0,12] --> From-change+Merge	[0,15]
// [0,17] --> From, To-change+Merge [0,17]
func (this *SortedIvDomPartList) AddAnyPart(part *IvDomPart) {
	// extend part or merge two parts
	splitParts := DIFFERENCE_DomParts(part, this.GetParts()...)
	match := false
	for _, sP := range splitParts {
		// println(sP.String())
		for i := 0; i < len(this.parts); i++ {
			if i+1 < len(this.parts) &&
				this.parts[i].ContainsInt(sP.From-1) &&
				this.parts[i+1].ContainsInt(sP.To+1) {
				// merge
				// print("merge"+ this.parts[i].String()+","
				// println(this.parts[i+1].String())
				this.parts[i].To = this.parts[i+1].To
				this.removeDomPart(i + 1)
				match = true
				break
			} else {
				if !this.parts[i].ContainsInt(sP.From-1) &&
					this.parts[i].ContainsInt(sP.To+1) {
					// extend
					// println("extend from"+ this.parts[i].String())
					this.parts[i].From = sP.From
					match = true
					break
				}
				if this.parts[i].ContainsInt(sP.From-1) &&
					!this.parts[i].ContainsInt(sP.To+1) {
					// extend
					// println("extend to"+ this.parts[i].String())
					this.parts[i].To = sP.To
					match = true
					break
				}
			}
		}
		if !match { //insert new part
			if len(this.parts) == 0 {
				// println("insert new")
				this.Append(CreateIvDomPart(sP.From, sP.To))
				break
			}
			if this.parts[0].GT_DP(sP) {
				// min: insert at position 0
				// println("insert at pos 0")
				this.insertDomPart(0, CreateIvDomPart(sP.From, sP.To))
			} else if this.parts[len(this.parts)-1].LT_DP(sP) {
				// max
				// println("insert at pos max")
				this.insertDomPart(len(this.parts),
					CreateIvDomPart(sP.From, sP.To))
			} else { // search position
				var lpart, rpart *IvDomPart
				for i := 0; i < len(this.parts)-1; i++ {
					lpart = this.parts[i]
					rpart = this.parts[i+1]
					if lpart.LT_DP(sP) && rpart.GT_DP(sP) {
						// println("insert at pos ",i+1)
						this.insertDomPart(i+1,
							CreateIvDomPart(sP.From, sP.To))
						break
					}
				}
			}
		}
	}
}

// Append appends a part at the end of the current list
// important: part must be greater than the other with
// difference>1.
func (this *SortedIvDomPartList) Append(part *IvDomPart) {
	this.parts = append(this.parts, part)
}

// ToDo: Test
// Appends appends parts at the end of the current list
// important: parts must be sorted and part[0] must be greater than
// greatest part of this (difference>1).
func (this *SortedIvDomPartList) Appends(parts []*IvDomPart) {
	this.parts = append(this.parts, parts...)
}

// Add adds an element to the current IvDomain
func (this *SortedIvDomPartList) Add(ele int) {
	// extend part or merge two parts
	for i := 0; i < len(this.parts); i++ {
		if i+1 < len(this.parts) &&
			this.parts[i].ContainsInt(ele-1) &&
			this.parts[i+1].ContainsInt(ele+1) {
			// merge
			this.parts[i].To = this.parts[i+1].To
			this.removeDomPart(i + 1)
			return
		} else if this.parts[i].From == ele+1 {
			// extend
			this.parts[i].From -= 1
			return
		} else if this.parts[i].To == ele-1 {
			// extend
			this.parts[i].To += 1
			return
		}
	}
	// insert new part
	if len(this.parts) == 0 {
		this.Append(CreateIvDomPart(ele, ele))
		return
	}
	if this.parts[0].GT(ele + 1) {
		//min: insert at position 0
		this.insertDomPart(0, CreateIvDomPart(ele, ele))
	} else if this.parts[len(this.parts)-1].LT(ele - 1) {
		// max
		this.insertDomPart(len(this.parts), CreateIvDomPart(ele, ele))
	} else {
		// search position
		var lpart, rpart *IvDomPart
		for i := 0; i < len(this.parts)-1; i++ {
			lpart = this.parts[i]
			rpart = this.parts[i+1]
			if lpart.LT(ele-1) && rpart.GT(ele+1) {
				this.insertDomPart(i+1, CreateIvDomPart(ele, ele))
				return
			}
		}
	}
}

// insertDomPart inserts a new IvDomPart at a given index to the given
// data storage with SliceTricks
func (this *SortedIvDomPartList) insertDomPart(index int, part *IvDomPart) {
	this.parts = append(this.parts, nil)
	copy(this.parts[index+1:], this.parts[index:])
	this.parts[index] = part
}

// removeDomPart removes a part at the given index from the local storage
// (Note: too expensive for frequently usage!)
func (this *SortedIvDomPartList) removeDomPart(index int) {
	copy(this.parts[index:], this.parts[index+1:])
	this.parts[len(this.parts)-1] = nil // or the zero value of T
	this.parts = this.parts[:len(this.parts)-1]
}

// Removes removes a sorted list of domain parts
func (this *SortedIvDomPartList) Removes(other *SortedIvDomPartList) {
	i := 0
	j := 0
	for i < len(this.GetParts()) && j < len(other.GetParts()) {
		// println(len(this.parts), " - ", this.String())
		// println(len(other.parts), " - ", other.String())
		// println(i, " - ", j)
		if this.parts[i].LT_DP(other.parts[j]) {
			i += 1
			// println("weiter i")
		} else if other.parts[j].LT_DP(this.parts[i]) {
			// println("weiter j")
			j += 1
		} else {
			state, p := this.parts[i].DIFFERENCE(other.parts[j])
			if state == REMOVE_PART {
				// SliceTrick
				// println("remove")
				this.removeDomPart(i)
				continue
			} else if state == MODIFIED_PART {
				// do nothing
				// println("do nothing")
			} else if state == INSERT_PART {
				// SliceTrick
				// println("insert")
				this.insertDomPart(i+1, p)
			} else {
				// println("nothing... State:", state)
			}
			if this.parts[i].To > other.parts[j].To {
				j += 1
			} else {
				i += 1
			}
		}
	}
}

// RemovesSortedInts removes the given integer values.
func (this *SortedIvDomPartList) RemovesSortedInts(sortedEles []int) {
	i := 0
	j := 0
	for i < len(this.parts) && j < len(sortedEles) {
		// println(len(this.parts)," - ",this.String(), " :::::: "
		// println(i, " - ", j)
		if this.parts[i].LT(sortedEles[j]) {
			i += 1
			// println("weiter i")
		} else if this.parts[i].GT(sortedEles[j]) {
			j += 1
			// println("weiter j")
		} else {
			state, p := this.parts[i].DIFFERENCE_Int(sortedEles[j])
			if state == REMOVE_PART {
				// SliceTrick
				this.removeDomPart(i)
				// println("remove")
				continue
			} else if state == MODIFIED_PART {
				// do nothing
				// println("do nothing")
			} else if state == INSERT_PART {
				// SliceTrick
				this.insertDomPart(i+1, p)
				// println("insert")
			}
			if this.parts[i].To > sortedEles[j] {
				j += 1
			} else {
				i += 1
			}
		}
	}
}

// IntersectionInts calculates the intersection between the current part list
// (this) and the given int array. Linear time algorithm
func (this *SortedIvDomPartList) IntersectionInts(sortedEles []int) *SortedIvDomPartList {
	newInts := make([]int, 0)
	i := 0
	j := 0
	for i < len(this.parts) && j < len(sortedEles) {
		if this.parts[i].ContainsInt(sortedEles[j]) {
			newInts = append(newInts, sortedEles[j])
		}
		if this.parts[i].LT(sortedEles[j]) {
			i += 1
		} else if this.parts[j].GT(sortedEles[j]) {
			j += 1
		}
	}
	return CreateSortedIvDomPartListFromIntArr(newInts)
}

// Intersection calculates the intersection between the current part list
// (this) and the given part list (other). Linear time algorithm
func (this *SortedIvDomPartList) Intersection(other *SortedIvDomPartList) *SortedIvDomPartList {
	newParts := make([]*IvDomPart, 0)
	i, j := 0, 0
	for i < len(this.parts) && j < len(other.parts) {
		if this.parts[i].LT_DP(other.parts[j]) {
			i += 1
		} else if other.parts[j].LT_DP(this.parts[i]) {
			j += 1
		} else {
			newParts = append(newParts,
				this.parts[i].INTERSECTION(other.parts[j]))
			if this.parts[i].To > other.parts[j].To {
				j += 1
			} else {
				i += 1
			}
		}
	}
	return CreateSortedIvDomPartListWithSortedParts(newParts)
}

// Remove removes ele from local data storage
// ToDo: meaning? (CAUTION: only for success from external)
func (this *SortedIvDomPartList) Remove(ele int) {
	index := 0
	var newPart *IvDomPart
	operation := NOTHING
	for i, part := range this.GetParts() {
		if part.LT(ele) {
			continue
		} else if part.GT(ele) {
			return
		} else if part.ContainsInt(ele) {
			index = i
			operation, newPart = part.DIFFERENCE_Int(ele)
			break
		}
	}
	if operation == REMOVE_PART {
		this.removeDomPart(index)
	}
	// SliceTricks
	if operation == INSERT_PART {
		index += 1
		this.insertDomPart(index, newPart)
	}
}

// --- representation functions ---

// GetParts returns domain parts (intervals)
func (this *SortedIvDomPartList) GetParts() []*IvDomPart {
	return this.parts
}

// Values_asMap returns a map representation of an IvDomain.
func (this *SortedIvDomPartList) Values_asMap() map[int]bool {
	values := make(map[int]bool)
	for _, part := range this.parts {
		part.addToMap(values)
	}
	return values
}

// Values_asSlice returns the IvDomain as a slice
func (this *SortedIvDomPartList) Values_asSlice() []int {
	var slice []int
	for _, part := range this.parts {
		slice = append(slice, part.makeSlice()...)
	}
	return slice
}

// --- standard functions ---

// Equals checks, if the given list is equal to another list
func (this *SortedIvDomPartList) Equals(other *SortedIvDomPartList) bool {
	if len(this.parts) != len(other.parts) {
		return false
	}
	for i, _ := range this.parts {
		if !this.parts[i].Equals(other.parts[i]) {
			return false
		}
	}
	return true
}

// Copy copies the current list.
func (this *SortedIvDomPartList) Copy() *SortedIvDomPartList {
	l := CreateSortedIvDomPartList()
	l.parts = make([]*IvDomPart, len(this.parts))
	for index, dp := range this.parts {
		l.parts[index] = dp.Copy()
	}
	return l
}

// IsEmpty returns true iff there is no element in this domain.
func (this *SortedIvDomPartList) IsEmpty() bool {
	return this.parts == nil || len(this.parts) == 0
}

// GetMax returns the largest value in this list;
// panics on empty domains.
func (this *SortedIvDomPartList) GetMax() int {
	return this.parts[len(this.parts)-1].To
}

// GetMin returns the smallest value in this list;
// panics on empty domains.
func (this *SortedIvDomPartList) GetMin() int {
	return this.parts[0].From
}

// Contains checks if ele is in List.
func (this *SortedIvDomPartList) Contains(ele int) bool {
	if this.IsEmpty() {
		return false
	}
	for _, part := range this.parts {
		if part.GT(ele) { // cause: sorted parts
			return false
		}
		if part.ContainsInt(ele) {
			return true
		}
	}
	return false
}

// String returns a sorted string representation of this domain.
func (this *SortedIvDomPartList) String() string {
	aParts := make([]string, len(this.parts))
	for i, part := range this.parts {
		if part.From == part.To {
			aParts[i] = fmt.Sprintf("%d", part.From)
		} else {
			aParts[i] = fmt.Sprintf("%d..%d", part.From, part.To)
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(aParts, ","))
}

// Size returns the number of elements of this domain.
func (this *SortedIvDomPartList) Size() int {
	size := 0
	for _, part := range this.parts {
		size += (part.To - part.From) + 1
	}
	return size
}
