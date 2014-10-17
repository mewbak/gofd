package labeling

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator"
	"bitbucket.org/gofd/gofd/propagator/explicit"
	"bitbucket.org/gofd/gofd/propagator/interval"
	"testing"
)

func Test_simple1_explicit(t *testing.T) {
	setup()
	defer teardown()
	log("simple1_explicit")

	X := core.CreateIntVarExFromTo("X", store, 0, 9)
	Y := core.CreateIntVarExFromTo("Y", store, 0, 9)
	prop := explicit.CreateC1XplusC2YeqC3(1, X, 1, Y, 12)
	store.AddPropagator(prop)
	query := CreateSearchAllQuery()
	result := Labeling(store, query, InDomainMin, GreatestDomainFirst)
	if !result {
		t.Errorf("labeling_test: labeling_test_result = %v, want %v",
			result, !result)
	}
	resultSet := query.GetResultSet()
	logger.If("no results: %v", len(resultSet))
	for _, result := range resultSet {
		sum := result[X] + result[Y]
		if sum != 12 {
			t.Errorf("labeling_test: labeling_test_result X + Y = %v, want %v",
				sum, 12)
		}
	}
	searchStat(query.GetSearchStatistics())
}

func Test_simple1_interval(t *testing.T) {
	setup()
	defer teardown()
	log("simple1_interval")

	X := core.CreateIntVarIvFromTo("X", store, 0, 9)
	Y := core.CreateIntVarIvFromTo("Y", store, 0, 9)
	prop := interval.CreateC1XplusC2YeqC3(1, X, 1, Y, 12)
	store.AddPropagator(prop)
	query := CreateSearchAllQuery()
	result := Labeling(store, query, InDomainMin, GreatestDomainFirst)
	if !result {
		t.Errorf("labeling_test: labeling_test_result = %v, want %v",
			result, !result)
	}
	resultSet := query.GetResultSet()
	logger.If("no results: %v", len(resultSet))
	for _, result := range resultSet {
		sum := result[X] + result[Y]
		if sum != 12 {
			t.Errorf("labeling_test: labeling_test_result X + Y = %v, want %v",
				sum, 12)
		}
	}
	searchStat(query.GetSearchStatistics())
}

func Test_simple1split(t *testing.T) {
	setup()
	defer teardown()
	log("simple1split")

	X := core.CreateIntVarFromTo("X", store, 0, 9)
	Y := core.CreateIntVarFromTo("Y", store, 0, 9)
	prop := propagator.CreateC1XplusC2YeqC3(1, X, 1, Y, 12)
	store.AddPropagator(prop)
	query := CreateSearchAllQuery()
	result := LabelingSplit(store, query, GreatestDomainFirst)
	if !result {
		t.Errorf("labeling_test: labeling_test_result = %v, want %v",
			result, !result)
	}
	resultSet := query.GetResultSet()
	logger.If("no results: %v", len(resultSet))

	for _, result := range resultSet {
		logger.If("x:%v,y:%v", result[X], result[Y])
		sum := result[X] + result[Y]
		if sum != 12 {
			t.Errorf("labeling_test: labeling_test_result X + Y = %v, want %v",
				sum, 12)
		}
	}
	searchStat(query.GetSearchStatistics())
}

func Test_simple2(t *testing.T) {
	setup()
	defer teardown()
	//	logger.SetLoggingLevel(core.LOG_INFO)
	log("simple2")

	X := core.CreateIntVarFromTo("X", store, 0, 9)
	Y := core.CreateIntVarFromTo("Y", store, 0, 9)
	prop1 := propagator.CreateC1XplusC2YeqC3(1, X, 1, Y, 9)
	store.AddPropagator(prop1)
	prop2 := propagator.CreateC1XplusC2YeqC3(2, X, 4, Y, 24)
	store.AddPropagator(prop2)

	query := CreateSearchOneQuery()
	result := Labeling(store, query)
	if !result {
		t.Errorf("labeling_test: labeling_test_result = %v, want %v",
			result, !result)
	}
	resultSet := query.GetResultSet()
	if resultSet[0][X] != 6 {
		t.Errorf("labeling_test: labeling_test_result X = %v, want %v",
			resultSet[0][X], 6)
	}
	if resultSet[0][Y] != 3 {
		t.Errorf("labeling_test: labeling_test_result Y = %v, want %v",
			resultSet[0][Y], 3)
	}
	searchStat(query.GetSearchStatistics())
}
