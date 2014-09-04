package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"bitbucket.org/gofd/gofd/propagator/explicit"
	"fmt"
	"testing"
)

// Travelling salesman problem modelled with among

// testTravellingSalesman generates a TSP with given nodes,
// implementation meant for TSPs with complete graphs
func testTravellingSalesman(t *testing.T, nodes []core.VarId, expectedResult bool) {
	numberOfNodes := len(nodes)
	// create all needed subsets for the no subtour among-constraints
	nodeSubsets := make([][]core.VarId, 0)
	getAllSets(nodes, &nodeSubsets)
	// define constraints
	// all nodes must have different successors
	store.AddPropagator(explicit.CreateAlldifferent_Primitives(nodes...))
	for _, set := range nodeSubsets {
		// generate the values which the variables should take
		// those values correspond to the indexes of the variables in the
		// current set
		domain := make([]int, 0)
		setSize := len(set)
		for _, node := range set {
			if int(node) == 1 {
				domain = append(domain, numberOfNodes+1)
			} else {
				domain = append(domain, int(node))
			}
		}
		// Avoid subtours by always allowing only setSize-1 nodes to take values
		// corresponding to their indexes:
		// Impose that at most setSize-1 variables in the set can take values
		// from domain.
		store.AddPropagator(explicit.CreateAmong(set, domain,
			core.CreateAuxIntVarExFromTo(store, 0, setSize-1)))
	}
	query := labeling.CreateSearchOneQueryVariableSelect(nodes)
	labeling.Labeling(store, query, labeling.VarSelect, labeling.InDomainMin)
	ready := store.IsConsistent()
	log(fmt.Sprintf("propagators=%3d, ready=%5v, nodes=%3d",
		store.GetNumPropagators(), ready,
		query.GetSearchStatistics().GetNodes()))
	ready_test(t, "Travelling Salesman", ready, expectedResult)
}

// getAllSets generates subsets (with minimum 2 elements up to n-1 elements)
// of a list S and stores the results in ls.
func getAllSets(varList []core.VarId, resultList *[][]core.VarId) {
	if len(varList) < 3 {
		return
	}
	for _, element := range varList {
		sc := make([]core.VarId, 0)
		//make a copy of S not containing the current element
		for _, variable := range varList {
			if variable != element {
				sc = append(sc, variable)
			}
		}
		if !sliceContainsSlice(sc, *resultList) {
			*resultList = append(*resultList, sc)
		}
		getAllSets(sc, resultList)
	}
}

// slicesIdentical checks if two slices are identical
// and returns a boolean containing the check result.
func slicesIdentical(slice1 []core.VarId, slice2 []core.VarId) bool {
	slice1Length := len(slice1)
	if slice1Length != len(slice2) {
		return false
	}
	countIdenticalValues := 0
	for _, element1 := range slice1 {
		elementIdentical := false
		for _, element2 := range slice2 {
			if element1 == element2 {
				countIdenticalValues += 1
				elementIdentical = true
				break
			}
		}
		if !elementIdentical {
			countIdenticalValues = 0
		}
	}
	if slice1Length == countIdenticalValues {
		return true
	}
	return false
}

// sliceContainsSlice checks if a two dimensional slice already contains a
// given one dimensional slice and returns a boolean containing the result.
func sliceContainsSlice(smallSlice []core.VarId, bigSlice [][]core.VarId) bool {
	for _, slice := range bigSlice {
		if slicesIdentical(slice, smallSlice) {
			return true
		}
	}
	return false
}

// generateTSPWithSize generates the variables for a TSP with a given number
// of nodes, it generates TSPs with complete graphs,
// returns an array containing the variables' VarIds.
func generateTSPWithSize(numberOfNodes int) []core.VarId {
	// every node has its possible successors in its domain, but not itself
	// node 1 is replaced by 5 in the domains
	log(fmt.Sprintf("TSP with %d cities, complete graph, Among", numberOfNodes))
	nodes := make([]core.VarId, numberOfNodes)
	for i := 0; i < len(nodes); i++ {
		domain := make([]int, 0)
		for j := 0; j < len(nodes); j++ {
			if i != 0 && j == 0 {
				domain = append(domain, len(nodes)+1)
			} else if i != j {
				domain = append(domain, j+1)
			}
		}
		nodes[i] = core.CreateAuxIntVarValues(store, domain)
	}
	return nodes
}

// generateTSPFromSlice generates the variables for a TSP from a given slice,
// the slice contains the indexes of the nodes' successors,
// domains begin at value 2, no variable should contain its own index,
// index 1 should be replaced by (number of nodes + 1).
func generateTSPFromSlice(nodes [][]int) []core.VarId {
	// every node has its possible successors in its domain, but not itself
	// node 1 is replaced by 5 in the domains
	log(fmt.Sprintf("TSP with %d cities, from slice,     Among", len(nodes)))
	variables := make([]core.VarId, len(nodes))
	for i := 0; i < len(nodes); i++ {
		variables[i] = core.CreateAuxIntVarValues(store, nodes[i])
	}
	return variables
}

func Test_travellingSalesman2(t *testing.T) {
	setup()
	defer teardown()
	nodes := generateTSPWithSize(2)
	testTravellingSalesman(t, nodes, true)
}

func Test_travellingSalesman3(t *testing.T) {
	setup()
	defer teardown()
	nodes := generateTSPWithSize(3)
	testTravellingSalesman(t, nodes, true)
}

func Test_travellingSalesman4(t *testing.T) {
	setup()
	defer teardown()
	nodes := generateTSPWithSize(4)
	testTravellingSalesman(t, nodes, true)
}

func Test_travellingSalesman5(t *testing.T) {
	setup()
	defer teardown()
	nodes := generateTSPWithSize(5)
	testTravellingSalesman(t, nodes, true)
}

func Test_travellingSalesmans4(t *testing.T) {
	setup()
	defer teardown()
	variables := make([][]int, 4)
	variables[0] = []int{2, 3}
	variables[1] = []int{3, 5}
	variables[2] = []int{4}
	variables[3] = []int{2, 5}
	nodes := generateTSPFromSlice(variables)
	testTravellingSalesman(t, nodes, true)
}

func Test_travellingSalesmanF(t *testing.T) {
	setup()
	defer teardown()
	variables := make([][]int, 5)
	variables[0] = []int{2, 4}
	variables[1] = []int{4, 5}
	variables[2] = []int{2, 6}
	variables[3] = []int{3, 5}
	variables[4] = []int{3, 6}
	nodes := generateTSPFromSlice(variables)
	testTravellingSalesman(t, nodes, true)
}
