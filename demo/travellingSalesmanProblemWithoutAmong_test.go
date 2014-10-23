package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"bitbucket.org/gofd/gofd/propagator/indexical"
	"bitbucket.org/gofd/gofd/propagator/interval"
	"bitbucket.org/gofd/gofd/propagator/reification"
	"fmt"
	"testing"
)

// Travelling salesman problem modelled without among

// testTravellingSalesmanWithoutAmong generates a TSP with given nodes
// implementation meant for TSPs with complete graphs.
func testTravellingSalesmanWithoutAmong(t *testing.T, nodes []core.VarId,
	expectedResult bool) {
	numberOfNodes := len(nodes)
	// avoid subtours by always allowing only n-1 nodes to take values
	// corresponding to their indexes
	// get all needed subsets
	nodeSubsets := make([][]core.VarId, 0)
	getAllSets(nodes, &nodeSubsets)
	// define constraints
	// all nodes must have different successors
	store.AddPropagator(interval.CreateAlldifferent(nodes...))
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
		reifiedVariables := make([]core.VarId, setSize*len(domain))
		i := 0
		for _, value := range domain {
			for _, node := range set {
				reifiedVariables[i] = core.CreateAuxIntVarFromTo(store, 0, 1)
				xeqc := indexical.CreateXeqC(node, value)
				reified := reification.CreateReifiedConstraint(xeqc,
					reifiedVariables[i])
				store.AddPropagator(reified)
				i++
			}
		}
		// impose that at most setSize-1 variables in the set can take values
		// from domain
		store.AddPropagator(interval.CreateSum(store,
			core.CreateAuxIntVarFromTo(store, 0, setSize-1),
			reifiedVariables))
	}

	query := labeling.CreateSearchOneQueryVariableSelect(nodes)
	labeling.Labeling(store, query, labeling.VarSelect, labeling.InDomainMin)
	ready := store.IsConsistent()
	log(fmt.Sprintf("propagators=%3d, ready=%5v, nodes=%3d",
		store.GetNumPropagators(), ready,
		query.GetSearchStatistics().GetNodes()))
	ready_test(t, "Travelling Salesman WithoutAmong", ready, expectedResult)
}

// generateTSPWithSize generates the variables for a TSP with a given number
// of nodes, it generates TSPs with complete graphs,
// the method returns an array containing the variables' VarIds.
func generateTSPWithSizeWithoutAmong(numberOfNodes int) []core.VarId {
	// every node has its possible successors in its domain, but not itself
	// node 1 is replaced by 5 in the domains
	log(fmt.Sprintf("TSP with %d cities, complete graph", numberOfNodes))
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
// index 1 should be replaced by (number of nodes + 1)
func generateTSPFromSliceWithoutAmong(nodes [][]int) []core.VarId {
	// every node has its possible successors in its domain, but not itself
	// node 1 is replaced by 5 in the domains
	log(fmt.Sprintf("TSP with %d cities, from slice", len(nodes)))
	variables := make([]core.VarId, len(nodes))
	for i := 0; i < len(nodes); i++ {
		variables[i] = core.CreateAuxIntVarValues(store, nodes[i])
	}
	return variables
}

func Test_travellingSalesmanWithoutAmongA(t *testing.T) {
	setup()
	defer teardown()
	nodes := generateTSPWithSizeWithoutAmong(2)
	testTravellingSalesmanWithoutAmong(t, nodes, true)
}

func Test_travellingSalesmanWithoutAmongB(t *testing.T) {
	setup()
	defer teardown()
	nodes := generateTSPWithSizeWithoutAmong(3)
	testTravellingSalesmanWithoutAmong(t, nodes, true)
}

func Test_travellingSalesmanWithoutAmongC(t *testing.T) {
	setup()
	defer teardown()
	nodes := generateTSPWithSizeWithoutAmong(4)
	testTravellingSalesmanWithoutAmong(t, nodes, true)
}

func Test_travellingSalesmanWithoutAmongD(t *testing.T) {
	setup()
	defer teardown()
	nodes := generateTSPWithSizeWithoutAmong(5)
	testTravellingSalesmanWithoutAmong(t, nodes, true)
}

func Test_travellingSalesmanWithoutAmongE(t *testing.T) {
	setup()
	defer teardown()
	variables := make([][]int, 4)
	variables[0] = []int{2, 3}
	variables[1] = []int{3, 5}
	variables[2] = []int{4}
	variables[3] = []int{2, 5}
	nodes := generateTSPFromSliceWithoutAmong(variables)
	testTravellingSalesmanWithoutAmong(t, nodes, true)
}

func Test_travellingSalesmanWithoutAmongF(t *testing.T) {
	setup()
	defer teardown()
	variables := make([][]int, 5)
	variables[0] = []int{2, 4}
	variables[1] = []int{4, 5}
	variables[2] = []int{2, 6}
	variables[3] = []int{3, 5}
	variables[4] = []int{3, 6}
	nodes := generateTSPFromSliceWithoutAmong(variables)
	testTravellingSalesmanWithoutAmong(t, nodes, true)
}
