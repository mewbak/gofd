package core

import (
	"fmt"
	"sort"
	"strings"
)

const MaxInt = int(^uint(0) >> 1)
const MinInt = -MaxInt - 1
const NEG_INFINITY = MinInt
const INFINITY = MaxInt

// ScalarSlice generates a new slice of values each multiplied with a weight
// e.g.: weight=10, values=[3, 6, 8] -> [30, 60, 80]
func ScalarSlice(weight int, values []int) []int {
	a := make([]int, len(values))
	for i, value := range values {
		a[i] = weight * value
	}
	return a
}

func makeSlice(start, end int) []int {
	values := make([]int, (end-start)+1)
	for v := start; v <= end; v++ {
		values[v-start] = v
	}
	return values
}

// AbsInt computes the absolute value of an int.
func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// BeautifulOutput computes a sorted representation of an explicit Domain
// represented by a map of int to bool, where consecutive values are
// collapsed into a range similar to interval representation.
func BeautifulOutput(what string, values map[int]bool) string {
	if len(values) == 0 {
		return fmt.Sprintf("%s[]", what)
	}
	vals := SortedKeys_MapIntToBool(values)
	beginning, ending := vals[0], vals[0]
	var s []string
	for i := 1; i < len(vals); i++ {
		v := vals[i]
		if v == ending+1 { // incremented
			ending = v
		} else { // end of continuous region
			s = append(s, makeEntry(beginning, ending))
			beginning, ending = v, v
		}
	}
	s = append(s, makeEntry(beginning, ending))
	return fmt.Sprintf("%s[%s]", what, strings.Join(s, ","))
}

// returns a string of a domain-beginning and -ending for beautiful output.
func makeEntry(beginning, ending int) string {
	if beginning == ending {
		return fmt.Sprintf("%d", beginning)
	}
	return fmt.Sprintf("%d..%d", beginning, ending)
}

// SortedKeys_MapIntToBool provides a slice of values in ascending order.
func SortedKeys_MapIntToBool(vals map[int]bool) []int {
	keys := Keys_MapIntToBool(vals)
	SortIntArr(keys)
	return keys
}

func SortIntArr(values []int) {
	sort.Sort(SortableIntSlice(values))
}

// ValuesOfMapVarIdToIvDomain returns a slice of domains from a
// given map if IvDomains.
func ValuesOfMapVarIdToIvDomain(varids []VarId,
	m map[VarId]*IvDomain) []Domain {
	doms := make([]Domain, len(varids))
	for i, varid := range varids {
		doms[i] = m[varid]
	}
	return doms
}

// ValuesOfMapVarIdToDomain returns a slice of domains from explicit domains.
func ValuesOfMapVarIdToExDomain(varids []VarId, m map[VarId]*ExDomain) []Domain {
	doms := make([]Domain, len(varids))
	for i, varid := range varids {
		doms[i] = m[varid]
	}
	return doms
}

// Keys_MapIntToBool returns a slice of keys of the given map[int]bool
func Keys_MapIntToBool(values map[int]bool) []int {
	keys := make([]int, len(values))
	i := 0
	for k, v := range values {
		if v { // only the ones that are true (from the domain)
			keys[i] = k
			i += 1
		}
	}
	return keys
}

// Keys_MapIntToBool returns a slice of keys of the given map[int]bool.
func Keys_MapVarIdsToBool(values map[VarId]Domain) []VarId {
	keys := make([]VarId, len(values))
	i := 0
	for k, _ := range values {
		keys[i] = k
		i += 1
	}
	return keys
}

// SliceToKeys_MapIntToBool returns a map with the values of the given slice.
func SliceToKeys_MapIntToBool(values []int) map[int]bool {
	m := make(map[int]bool, len(values))
	for _, v := range values {
		m[v] = true
	}
	return m
}

// SortedKeys_MapVarIdToInt provides a slice of values in ascending order.
func SortedKeys_MapVarIdToInt(vals map[VarId]int) []VarId {
	keys := Keys_MapVarIdToInt(vals)
	ikeys := make([]int, len(keys))
	for i, key := range keys { // copy as ints
		ikeys[i] = int(key)
	}
	sort.Sort(SortableIntSlice(ikeys)) // sorts ascending
	for i, ikey := range ikeys {       // copy as VarId
		keys[i] = VarId(ikey)
	}
	return keys
}

// Keys_MapVarIdToInt provides a slice of values.
func Keys_MapVarIdToInt(values map[VarId]int) []VarId {
	keys := make([]VarId, len(values))
	i := 0
	for k, _ := range values {
		keys[i] = k
		i += 1
	}
	return keys
}

// makes a one-dim slice of two-dim-interval
// e.g. from {{1,10}} --> {1,2,3,4,5,6,7,8,9,10}
func makeTwoDim_OneDim(values [][]int) []int {
	valuesRes := make([]int, 0)
	for _, vslice := range values {
		if len(vslice) == 1 {
			v := vslice[0]
			valuesRes = append(valuesRes, v)
		} else if len(vslice) == 2 {
			for v := vslice[0]; v <= vslice[1]; v++ {
				valuesRes = append(valuesRes, v)
			}
		}
	}
	return valuesRes
}

func IntSliceToStringSlice(is []int) []string {
	return IntSliceToStringSliceFormatted(is, "%d")
}

func IntSliceToStringSliceFormatted(is []int, format string) []string {
	s := make([]string, len(is))
	for i, val := range is {
		s[i] = fmt.Sprintf(format, val)
	}
	return s
}

// type to implement sort.Interface for an int slice
type SortableIntSlice []int

func (this SortableIntSlice) Len() int {
	return len(this)
}
func (this SortableIntSlice) Less(i, j int) bool {
	return this[i] < this[j]
}
func (this SortableIntSlice) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
