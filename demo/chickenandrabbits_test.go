package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"bitbucket.org/gofd/gofd/propagator"
	"fmt"
	"testing"
)

// Constraint Programming -- The B-Prolog Project; clpfd_tutorial.pdf
// In a farmyard, there are only chickens and rabbits.
// Its is known that there are 18 heads and 58 feet.
// How many chickens and rabbits are there?
func Test_chickenAndRabbits(t *testing.T) {
	setup()
	defer teardown()
	log("Chicken and Rabbits")
	heads, feet, chicken, rabbits := 18, 58, 7, 11
	v := setupChickenAndRabbits(store, heads, feet)
	query := labeling.CreateSearchOneQuery()
	result := labeling.Labeling(store, query)
	ready_test(t, "Test_chickenAndRabbits", result, true)
	resultSet := query.GetResultSet()
	resultSet_test(t, "Test_chickenAndRabbits", resultSet, 0, v[0], chicken)
	resultSet_test(t, "Test_chickenAndRabbits", resultSet, 0, v[1], rabbits)
	msg := "With %d heads and %d feet there are %d chicken and %d rabbits"
	log(fmt.Sprintf(msg, heads, feet, chicken, rabbits))
	searchStat(query.GetSearchStatistics())
}

func setupChickenAndRabbits(store *core.Store, heads, feet int) []core.VarId {
	v := make([]core.VarId, 4)
	v[0] = core.CreateIntVarFromTo("chicken", store, 0, heads)
	v[1] = core.CreateIntVarFromTo("rabbit", store, 0, heads)
	store.AddPropagator(propagator.CreateXplusYeqC(v[0], v[1], heads))
	v[2] = core.CreateIntVarFromTo("chicken feet", store, 0, feet)
	v[3] = core.CreateIntVarFromTo("rabbit feet", store, 0, feet)
	store.AddPropagator(propagator.CreateXplusYeqC(v[2], v[3], feet))
	// Two feet for one chicken
	store.AddPropagator(propagator.CreateXmultCeqY(v[0], 2, v[2]))
	// Four feet for one rabbit
	store.AddPropagator(propagator.CreateXmultCeqY(v[1], 4, v[3]))
	return v
}
