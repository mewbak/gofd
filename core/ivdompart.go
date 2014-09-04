package core

import (
	"fmt"
	"math"
)

// IvDomPart is an interval/part, using only the lower (From) and
// upper (To) bound to represent all values (including the bounds)
// in between.
type IvDomPart struct {
	From int
	To   int
}

// CreateIvDomPart creates new IvDomPart from given from and to values
func CreateIvDomPart(from, to int) *IvDomPart {
	if from > to {
		panic("CreateIvDomPart: to must be greater than from")
	}
	dp := new(IvDomPart)
	dp.From = from
	dp.To = to
	return dp
}

// DIFFERENCE_MIN_MAX returns a new part (difference between this and
// a range from min to max)
func (this *IvDomPart) DIFFERENCE_MIN_MAX(min, max int) []*IvDomPart {
	_, p := this.DIFFERENCE_NEW(CreateIvDomPart(min, max))
	return p
}

// CreateIvDomParts creates new IvDomParts from given from and to values.
func CreateIvDomParts(fromTos [][]int) []*IvDomPart {
	parts := make([]*IvDomPart, len(fromTos))
	partsPos := 0
	for _, fromTo := range fromTos {
		partLen := len(fromTo)
		if partLen == 0 {
			// hack, who would like to instance an empty dompart? Useless
			parts = make([]*IvDomPart, len(parts)-1)
			continue
		}
		if partLen > 2 {
			panic("CreateIvDomParts: Wrong param format of fromTo pairs")
		}
		from := fromTo[0]
		to := fromTo[0]
		if partLen == 2 {
			to = fromTo[1]
		}
		parts[partsPos] = CreateIvDomPart(from, to)
		partsPos += 1
	}
	if len(parts) == 0 {
		return nil
	}
	return parts
}

// DIFFERENCE_Int computes the set difference
func (this *IvDomPart) DIFFERENCE_Int(e int) (int, *IvDomPart) {
	if this.ContainsInt(e) {
		state := this.RelationCheckInt(e)
		if state == CONTAINS_SAME_FROM {
			// println("Change_bound_from", p.String(), e)
			this.From = e + 1
			return MODIFIED_PART, nil
		} else if state == CONTAINS_SAME_TO {
			// println("Change_bound_to", p.String(), e)
			this.To = e - 1
			return MODIFIED_PART, nil
		} else if state == SAME {
			// println("remove", p.String(), e)
			return REMOVE_PART, nil
		} else if state == CONTAINS_NOT_SAME_FROM_OR_TO {
			// println("split", p.String(), e)
			tmp := this.To
			this.To = e - 1
			return INSERT_PART, CreateIvDomPart(e+1, tmp)
		}
	}
	return NOTHING, nil
}

// DIFFERENCE_Ints creates a list of *IvDomPart (splitted from the
// given IvDomPart) with given eles (which must be contained by DomPart).
// Note: eles must be sorted
// Its Difference-Calculating on p with given eles (returning difference).
func DIFFERENCE_Ints(p *IvDomPart, eles ...int) []*IvDomPart {
	var splits []*IvDomPart
	var dp1, dp2 *IvDomPart
	part := p.Copy()
	for _, ele := range eles {
		if part.ContainsInt(ele) {
			state := part.RelationCheckInt(ele)
			if state == CONTAINS_SAME_FROM {
				// println("Change_bound_from", part.String(), ele)
				part.From = ele + 1
			} else if state == CONTAINS_SAME_TO {
				// println("Change_bound_to", part.String(), ele)
				part.To = ele - 1
			} else if state == SAME {
				// println("remove", part.String(), ele)
				return splits
			} else if state == CONTAINS_NOT_SAME_FROM_OR_TO {
				// println("split", part.String(), ele)
				if splits == nil {
					splits = make([]*IvDomPart, 1)
					dp1, dp2 = part.removeWithTwoSplits(ele)
					part = dp2
					splits[0] = dp1
				} else {
					dp1, dp2 = part.removeWithTwoSplits(ele)
					part = dp2
					splits = append(splits, dp1)
				}
			}
		}
	}
	splits = append(splits, part)
	return splits
}

// DIFFERENCE_DomParts creates a list of *IvDomPart splitted from the current
// IvDomPart with the given splitParts. splitParts must be contained by
// part (Note: splitParts must be sorted!)
// Its Difference-Calculating on p with given splitPart (returning difference)
func DIFFERENCE_DomParts(p *IvDomPart, splitPart ...*IvDomPart) []*IvDomPart {
	var splits []*IvDomPart
	var dp1, dp2 *IvDomPart
	part := p.Copy()
	for _, sp := range splitPart {
		if part.ContainsOrSame(sp) {
			state := part.RelationCheck(sp)
			if state == CONTAINS_SAME_FROM {
				// println("Change_bound_from")
				part.From = sp.To + 1 // 1..10, 1..4 --> 5..10
			} else if state == CONTAINS_SAME_TO {
				// println("Change_bound_to")
				part.To = sp.From - 1 // 1..10, 5..10 --> 1..4
			} else if state == SAME ||
				state == IS_CONTAINED_NOT_SAME_FROM_OR_TO ||
				state == IS_CONTAINED_SAME_FROM ||
				state == IS_CONTAINED_SAME_TO {
				// println("remove")
				return nil
			} else if state == CONTAINS_NOT_SAME_FROM_OR_TO {
				// println("split")
				if splits == nil {
					splits = make([]*IvDomPart, 1)
					dp1, dp2 = part.removeWithTwoSplitsIvDom(sp)
					part = dp2
					splits[0] = dp1
				} else {
					dp1, dp2 = part.removeWithTwoSplitsIvDom(sp)
					part = dp2
					splits = append(splits, dp1)
				}
			}
		}
	}
	splits = append(splits, part)
	return splits
}

// IsGround on IvDomPart returns true iff there is one element in this part.
func (this *IvDomPart) IsGround() bool {
	return this.From == this.To
}

// Copy deep-copies/clones an IvDomPart
func (this *IvDomPart) Copy() *IvDomPart {
	return CreateIvDomPart(this.From, this.To)
}

// Equals checks whether this domain contains the same
// elements as the other domain.
func (this *IvDomPart) Equals(other *IvDomPart) bool {
	if this.From != other.From || this.To != other.To {
		return false
	}
	return true
}

// removeWithTwoSplits returns splitted parts from a specific int value
// 1. return-param: part with lower bounds
// 2. return-param: part with higher bounds
func (this *IvDomPart) removeWithTwoSplits(ele int) (*IvDomPart, *IvDomPart) {
	dpBefore := CreateIvDomPart(this.From, ele-1)
	dpAfter := CreateIvDomPart(ele+1, this.To)
	return dpBefore, dpAfter
}

// removeWithTwoSplitsIvDom returns splitted parts from a specific part
// 1. return-param: part with lower bounds
// 2. return-param: part with higher bounds
// i.e. this is [1,10] and part is [3,8] --> return-values are [1,2] [9,10]
func (this *IvDomPart) removeWithTwoSplitsIvDom(part *IvDomPart) (*IvDomPart, *IvDomPart) {
	dpBefore := CreateIvDomPart(this.From, part.From-1)
	dpAfter := CreateIvDomPart(part.To+1, this.To)
	return dpBefore, dpAfter
}

// ContainsInt returns true iff ele is in the IvDomPart
func (this *IvDomPart) ContainsInt(ele int) bool {
	return (this.From <= ele) && (ele <= this.To)
}

// INTERSECTION calculates the intersection between this and part
func (this *IvDomPart) INTERSECTION(part *IvDomPart) *IvDomPart {
	// println("INTER-this:",this.String())
	// println("INTER-part:",part.String())
	res := this.RelationCheck(part)
	// i: intersect
	// this i domPart
	if res == SAME ||
		res == IS_CONTAINED_SAME_FROM ||
		res == IS_CONTAINED_SAME_TO ||
		res == IS_CONTAINED_NOT_SAME_FROM_OR_TO {
		// [1..3] i [1..3] --> [1..3]
		// do nothing
		// println("WHOLE")
		return CreateIvDomPart(this.From, this.To)
	} else if res == CONTAINS_NOT_SAME_FROM_OR_TO {
		// [1..8] i [2..6] --> [2..6]
		// println("CONTAINS")
		return CreateIvDomPart(part.From, part.To)
	} else if res == NOT_IN_LOWER || res == NOT_IN_HIGHER {
		// [1..8] i [10..15] --> []
		// println("NIL")
		return nil
	} else if res == CONTAINS_SAME_FROM {
		// [1..10] i [1..9] --> [1..9]
		// println("CONTAINS_SAME_FROM")
		return CreateIvDomPart(this.From, part.To)
	} else if res == CONTAINS_SAME_TO {
		// [1..10] i [2..10] --> [2..10]
		// println("CONTAINS_SAME_TO")
		return CreateIvDomPart(part.From, part.To)
	} else if res == OVERLAPS_LOWER {
		// [1..10] i [5..15] --> [5..10]
		// println("OVER_LOWER")
		return CreateIvDomPart(part.From, this.To)
	} else if res == OVERLAPS_HIGHER {
		// [5..15] i [1..10] --> [5..10]
		// println("OVER_HIGHER")
		return CreateIvDomPart(this.From, part.To)
	}
	panic("INTERSECTION: case not matched")
}

// DIFFERENCE calculates the difference between this and part.
func (this *IvDomPart) DIFFERENCE(part *IvDomPart) (int, *IvDomPart) {
	if part.ContainsOrSame(this) {
		// is_contained
		return REMOVE_PART, nil
	} else {
		res := this.RelationCheck(part)
		if res == CONTAINS_NOT_SAME_FROM_OR_TO {
			// contains
			oldThisTo := this.To
			this.To = part.From - 1
			return INSERT_PART, CreateIvDomPart(part.To+1, oldThisTo)
		} else if res == NOT_IN_LOWER || res == NOT_IN_HIGHER {
			// not_in
			return NOTHING, nil
		} else if res == CONTAINS_SAME_FROM {
			// this([1,10]), p([1,9])
			this.From = part.To + 1
			return MODIFIED_PART, nil
		} else if res == CONTAINS_SAME_TO {
			// this([1,10]), p([2,10])
			this.To = part.From - 1
			return MODIFIED_PART, nil
		} else if res == OVERLAPS_LOWER {
			// this([1,5]), p([3,10])
			this.To = part.From - 1
			return MODIFIED_PART, nil
		} else if res == OVERLAPS_HIGHER {
			// this([3,10]), p([1,5])
			this.From = part.To + 1
			return MODIFIED_PART, nil
		}
	}
	panic("IvDomPart.DIFFERENCE failed")
}

// DIFFERENCE_NEW returns a new part (difference between this and part)
func (this *IvDomPart) DIFFERENCE_NEW(part *IvDomPart) (int, []*IvDomPart) {
	// println(this.String())
	// println(part.String())
	if part.ContainsOrSame(this) {
		// is_contained
		return REMOVE_PART, nil
	} else {
		res := this.RelationCheck(part)
		// println("",res)
		if res == CONTAINS_NOT_SAME_FROM_OR_TO {
			return INSERT_PART, []*IvDomPart{
				CreateIvDomPart(this.From, part.From-1),
				CreateIvDomPart(part.To+1, this.To)}
		} else if res == NOT_IN_LOWER || res == NOT_IN_HIGHER {
			// not_in
			return NOTHING, []*IvDomPart{CreateIvDomPart(this.From, this.To)}
		} else if res == CONTAINS_SAME_FROM {
			// this([1,10]), p([1,9])
			return MODIFIED_PART, []*IvDomPart{
				CreateIvDomPart(part.To+1, this.To)}
		} else if res == CONTAINS_SAME_TO {
			// this([1,10]), p([2,10])
			return MODIFIED_PART, []*IvDomPart{
				CreateIvDomPart(this.From, part.From-1)}
		} else if res == OVERLAPS_LOWER {
			// this([1,5]), p([3,10])
			return MODIFIED_PART, []*IvDomPart{
				CreateIvDomPart(this.From, part.From-1)}
		} else if res == OVERLAPS_HIGHER {
			// this([3,10]), p([1,5])
			this.From = part.To - 1
			return MODIFIED_PART, []*IvDomPart{
				CreateIvDomPart(part.To+1, this.To)}
		}
	}
	panic("IvDomPart.DIFFERENCE failed")
}

// ContainsOrSame returns true iff the current IvDomPart
// contains the given part.
func (this *IvDomPart) ContainsOrSame(part *IvDomPart) bool {
	return (this.From <= part.From) && (part.To <= this.To)
}

// String returns a string representation of this IvDomPart.
func (this *IvDomPart) String() string {
	if this.From == this.To {
		return fmt.Sprintf("[%d]", this.From)
	}
	return fmt.Sprintf("[%d..%d]", this.From, this.To)
}

// Size returns the size of the dompart interval.
func (this *IvDomPart) Size() int {
	return (this.To - this.From) + 1
}

// makeSlice creates a slice of the given IvDomPart.
func (this *IvDomPart) makeSlice() []int {
	return makeSlice(this.From, this.To)
}

// addToMap adds the values of an IvDomPart to a given map
func (this *IvDomPart) addToMap(m map[int]bool) {
	for v := this.From; v <= this.To; v++ {
		m[v] = true
	}
}

// SortableDomPart is a sortable list of interval domain parts
type SortableDomPart []*IvDomPart

// Len computes the length of the SortableDomPart
func (this SortableDomPart) Len() int {
	return len(this)
}

// Less imposes an order over domparts given indexes.
func (this SortableDomPart) Less(i, j int) bool {
	return this[i].From < this[j].From
}

// Swap exchanges inplace
func (this SortableDomPart) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}

// RelationCheck checks the relation between this and the given IvDomPart
// and returns a constant, expressing the relation.
func (this *IvDomPart) RelationCheck(p *IvDomPart) int {
	// this([1,10]), p([1,10])
	if (this.From == p.From) && (p.To == this.To) {
		return SAME
	}
	// this([1,10]), p([2,9])
	if (this.From < p.From) && (p.To < this.To) {
		return CONTAINS_NOT_SAME_FROM_OR_TO
	}
	// this([2,9]), p([1,10])
	if (p.From < this.From) && (this.To < p.To) {
		return IS_CONTAINED_NOT_SAME_FROM_OR_TO
	}
	// this([1,3]), p([5,9])
	if this.To < p.From {
		return NOT_IN_LOWER
	}
	// this([5,9]), p([1,3])
	if p.To < this.From {
		return NOT_IN_HIGHER
	}
	// this([1,10]), p([1,9])
	if (this.From == p.From) && (p.To < this.To) {
		return CONTAINS_SAME_FROM
	}
	// this([1,10]), p([2,10])
	if (this.From < p.From) && (p.To == this.To) {
		return CONTAINS_SAME_TO
	}
	// this([1,9]), p([1,10])
	if (p.From == this.From) && (this.To < p.To) {
		return IS_CONTAINED_SAME_FROM
	}
	// this([2,10]), p([1,10])
	if (p.From < this.From) && (this.To == p.To) {
		return IS_CONTAINED_SAME_TO
	}
	// this([1,5]), p([3,10])
	if (this.From < p.From) && (p.From <= this.To) && (this.To < p.To) {
		return OVERLAPS_LOWER
	}
	// this([3,10]), p([1,5])
	if (p.From < this.From) && (this.From <= p.To) && (p.To < this.To) {
		return OVERLAPS_HIGHER
	}
	panic("IvDomPart.RelationCheck failed... should not happen")
}

// RelationCheckInt checks the relation between this and the given int value
// and returns a constant, epxressing the relation.
func (this *IvDomPart) RelationCheckInt(ele int) int {
	// --- parts-size-changing

	// this([10,10]), ele=10
	if (this.From == ele) && (this.To == ele) {
		return SAME
	}
	// this([1,10]), ele=9
	if (this.From < ele) && (ele < this.To) {
		return CONTAINS_NOT_SAME_FROM_OR_TO
	}
	// this([1,10]), ele=0
	if this.To < ele {
		return NOT_IN_LOWER
	}
	// this([1,10]), ele=15
	if ele < this.From {
		return NOT_IN_HIGHER
	}
	// this([1,10]), ele=1
	if this.From == ele {
		return CONTAINS_SAME_FROM
	}
	// this([1,10]), ele=10
	if this.To == ele {
		return CONTAINS_SAME_TO
	}
	// IS_CONTAINED, IS_CONTAINED_SAME_FROM, IS_CONTAINED_SAME_TO not possible
	panic("IvDomPart.RelationCheckInt failed... should not happen")
}

// LT checks, if a given part is less than ele.
func (this *IvDomPart) LT(ele int) bool {
	return this.To < ele
}

// GT checks, if a given part is greater than ele.
func (this *IvDomPart) GT(ele int) bool {
	return this.From > ele
}

// LT_EQ checks, if a given part is less or equal than ele.
func (this *IvDomPart) LT_EQ(ele int) bool {
	return this.To <= ele
}

// GT_EQ checks, if a given part is greater or equal than ele.
func (this *IvDomPart) GT_EQ(ele int) bool {
	return this.From >= ele
}

// LT_DP checks, if a given IvDomPart is less than "this".
func (this *IvDomPart) LT_DP(dom *IvDomPart) bool {
	return this.To < dom.From
}

// GT_DP checks, if the a given IvDomPart is less than "this".
func (this *IvDomPart) GT_DP(dom *IvDomPart) bool {
	return this.From > dom.To
}

// internal constants
const (
	NOT_IN_LOWER                     = 0
	NOT_IN_HIGHER                    = 1
	SAME                             = 2
	CONTAINS_NOT_SAME_FROM_OR_TO     = 3
	CONTAINS_SAME_FROM               = 4
	CONTAINS_SAME_TO                 = 5
	IS_CONTAINED_SAME_FROM           = 6
	IS_CONTAINED_SAME_TO             = 7
	IS_CONTAINED_NOT_SAME_FROM_OR_TO = 8
	OVERLAPS_LOWER                   = 9
	OVERLAPS_HIGHER                  = 10
)

// internal constants
const (
	INSERT_PART   = 0 // first modified, second insert
	REMOVE_PART   = 1
	MODIFIED_PART = 2
	NOTHING       = 3
)

// internal constants
const (
	SPLIT  = 0 // do not modify anything, only say: you have to split!
	MODIFY = 1
	REMOVE = 2
)

// SUBTRACT subtracts the current part with the given dompart with the
// greatest potential resulting interval.
// cases:
// (1,5)-(1,5)	= (-4,4)
// (1,5)-(-1,5) = (-4,6)
// (1,5)-(-5,-1) = (2,10)
// (1,5)-(-5,1) = (0,10)
func (this *IvDomPart) SUBTRACT(dompart *IvDomPart) *IvDomPart {
	return CreateIvDomPart(this.From-dompart.To, this.To-dompart.From)
}

// ADD adds the current part with the given dompart.
// Example: [3,4] + [6,7] = [9,11]
func (this *IvDomPart) ADD(dom *IvDomPart) *IvDomPart {
	return CreateIvDomPart(this.From+dom.From, this.To+dom.To)
}

// NEG negates the current part.
func (this *IvDomPart) NEG() *IvDomPart {
	return CreateIvDomPart(-this.To, -this.From)
}

// ABS turns the part-bounds to positive
func (this *IvDomPart) ABS() *IvDomPart {
	// -6,-3
	if this.From < 0 && this.To < 0 {
		return CreateIvDomPart(-this.To, -this.From)
	} else if this.From < 0 && this.To >= 0 {
		// -5, 5 --> 0..5
		// -4, 5 --> 0..5
		if -this.From <= this.To {
			return CreateIvDomPart(0, this.To)
		}
		// -6, 5 --> 0..6
		return CreateIvDomPart(0, -this.From)
	}
	// 3,6
	return CreateIvDomPart(this.From, this.To)
}

// ToDo: test

// MULTIPLY multiplies the current part with the given dompart
// Example: [3,4] * [6,7] = [18,28]
// multiplication with negativ number results in switch
// 1,5 * 3,6, --> 3,30
// 1,5 * -6,-3 --> -30,-3
// 1,5 * -3,2 --> -15,2
// -5,-1 * -3,2 --> -2,15
// -5,-1 * -3,-2 --> 2,15
func (this *IvDomPart) MULTIPLY(dom *IvDomPart) *IvDomPart {
	if this.From < 0 || this.To < 0 {
		if dom.From < 0 || dom.To < 0 {
			// double switch
			return CreateIvDomPart(this.To*dom.To, this.From*dom.From)
		} else {
			// single switch
			return CreateIvDomPart(this.From*dom.To, this.To*dom.From)
		}
	} else {
		if dom.From < 0 || dom.To < 0 {
			// single switch
			return CreateIvDomPart(this.To*dom.From, this.From*dom.To)
		} else {
			// no switch
			return CreateIvDomPart(this.From*dom.From, this.To*dom.To)
		}
	}
	// next line is identified as unreachable code, but needed for
	// go 1.0 compatibility
	panic("IvDomPart.MULTIPLY: Nothing matched")
}

// ToDo: test, switch like mulitply for negative values

// DIvIDE divides the current part with the given dompart
// (rounds to natural next included number)
// Example: [4,10]/[2,3]=[2,5] --> 5/3  and  9/2
func (this *IvDomPart) DIvIDE(dom *IvDomPart) *IvDomPart {
	tf := float64(this.From)
	dt := float64(dom.From)
	tt := float64(this.To)
	df := float64(dom.From)
	return CreateIvDomPart(int(math.Ceil(tf/dt)), int(math.Ceil(tt/df)))
}
