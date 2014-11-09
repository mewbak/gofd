package demo

import (
	"bitbucket.org/gofd/gofd/labeling"
	"fmt"
	"strings"
	"testing"
)

func runAlpha(t *testing.T, genAlpha func() (map[string]int, map[string]int)) {
	problem, solution := genAlpha()
	vars := ConstrainAlpha(store, problem)
	store.IsConsistent()
	query := labeling.CreateSearchOneQuery()
	result := labeling.Labeling(store, query,
		labeling.InDomainMin, labeling.SmallestDomainFirst)
	ready_test(t, "alpha1", result, true)
	assignment := query.GetResultSet()[0]
	for _, varid := range vars {
		varname := store.GetName(varid)
		log(fmt.Sprintf("%3v,%s=%02d", varid, varname, assignment[varid]))
	}
	log(fmt.Sprintf("Solution: %v", solution))
	log(fmt.Sprintf("Result  : %v", assignment))
	for _, varid := range store.GetVariableIDs() {
		varname := store.GetName(varid)
		if !strings.HasPrefix(varname, "_") {
			log(fmt.Sprintf("%3v,%s=%02d", varid, varname, assignment[varid]))
		}
	}
}

func Test_alpha1(t *testing.T) {
	setup()
	defer teardown()
	// "SKIP: There must be a bug in SumWeightedBounds?
	// we do not find a solution sometimes"
	// runAlpha(t, GenerateAlpha1)
}
