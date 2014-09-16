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

// car sequencing problem modelled without among

// testCarSequencingWithoutAmong generates a car sequencing problem.
func testCarSequencingWithoutAmong(t *testing.T, carsPerClass []int,
	carsWithOptions [][]int, consecutiveCars []int,
	howManyOfconsecutiveCarsAtLeast []int,
	howManyOfconsecutiveCarsAtMost []int,
	numberOfCars int, expectedResult bool) {
	log(fmt.Sprintf("car sequencing without Among: cars=%2d, options=%2d",
		numberOfCars, len(carsWithOptions)))
	numberOfDiffClasses := len(carsPerClass)
	// define car variables
	cars := make([]core.VarId, numberOfCars)
	for i := 0; i < len(cars); i++ {
		cars[i] = core.CreateAuxIntVarIvFromTo(store, 1, numberOfDiffClasses)
	}
	// define constraints
	// every car belongs to one of six classes (all cars in a class have
	// the same set of options)
	for i := 0; i < numberOfDiffClasses; i++ {
		// array to store variables with a {0,1} domain
		// which are needed for reification
		variables := make([]core.VarId, len(cars))
		for j := 0; j < len(cars); j++ {
			variables[j] = core.CreateAuxIntVarIvFromTo(store, 0, 1)
			xeqc := indexical.CreateXeqC(cars[j], i+1)
			reif := reification.CreateReifiedConstraint(xeqc, variables[j])
			store.AddPropagator(reif)
		}
		store.AddPropagator(interval.CreateSumBounds(store,
			core.CreateAuxIntVarIvFromTo(store, carsPerClass[i],
				carsPerClass[i]), variables))
	}
	// automated constraint modelling of problem
	// j is used as index for arrays concerning the options
	for j := 0; j < len(carsWithOptions); j++ {
		for i := 0; i < len(cars)-consecutiveCars[j]+1; i++ {
			// create array to store the results of reification
			// its size must be the number of consecutive cars that
			// can have the option
			// multiplied with the number of options the class has
			currentCarsWithOptions := carsWithOptions[j]
			variables := make([]core.VarId,
				consecutiveCars[j]*(len(currentCarsWithOptions)))
			variablesIndex := 0
			for n := 0; n < consecutiveCars[j]; n++ {
				for opt := 0; opt < len(currentCarsWithOptions); opt++ {
					aux := core.CreateAuxIntVarIvFromTo(store, 0, 1)
					variables[variablesIndex] = aux
					xeqc := indexical.CreateXeqC(cars[i+n],
						currentCarsWithOptions[opt])
					reifiedConstraint :=
						reification.CreateReifiedConstraint(xeqc,
							variables[variablesIndex])
					store.AddPropagator(reifiedConstraint)
					variablesIndex += 1
				}
			}
			// constrain how many consecutive cars can have the current
			// option
			store.AddPropagator(interval.CreateSumBounds(store,
				core.CreateAuxIntVarIvFromTo(store,
					howManyOfconsecutiveCarsAtLeast[j],
					howManyOfconsecutiveCarsAtMost[j]), variables))
		}
	}
	numberOfPropagators := store.GetNumPropagators()
	query := labeling.CreateSearchOneQueryVariableSelect(cars)
	labeling.Labeling(store, query, labeling.VarSelect, labeling.InDomainMin)
	ready := store.IsConsistent()
	log(fmt.Sprintf("propagators=%3d, ready=%5v, nodes=%3d",
		numberOfPropagators, ready, query.GetSearchStatistics().GetNodes()))
	ready_test(t, "Car Sequencing Without Among", ready, expectedResult)
}

func Test_carSequencingWithoutAmongA(t *testing.T) {
	setup()
	defer teardown()

	//10 cars divided in six classes
	carsPerClass := []int{1, 1, 2, 2, 2, 2}

	//5 options
	carsWithOptions := make([][]int, 5)
	//cars of the classes 1, 5 and 6 need option 1
	carsWithOptions[0] = []int{1, 5, 6}
	//cars of the classes 3, 4 and 6 need option 2 ...
	carsWithOptions[1] = []int{3, 4, 6}
	carsWithOptions[2] = []int{1, 5}
	carsWithOptions[3] = []int{1, 2, 4}
	carsWithOptions[4] = []int{3}

	//2 consecutive cars are concerned with option 1
	//3 consecutive cars are concerned with option 2, ...
	consecutiveCarsPerOption := []int{2, 3, 3, 5, 5}

	//at least 0 of 2 consecutive cars can have option 1
	//at least 0 of 3 consecutive cars can have option 2, ...
	howManyOfconsecutiveCarsAtLeast := []int{0, 0, 0, 0, 0}

	//at most 1 of 2 consecutive cars can have option 1
	//at most 2 of 3 consecutive cars can have option 2, ...
	howManyOfconsecutiveCarsAtMost := []int{1, 2, 1, 2, 1}

	numberOfCars := sum_intarray(carsPerClass)
	testCarSequencingWithoutAmong(t, carsPerClass, carsWithOptions,
		consecutiveCarsPerOption, howManyOfconsecutiveCarsAtLeast,
		howManyOfconsecutiveCarsAtMost, numberOfCars, true)
}

func Test_carSequencingWithoutAmongB(t *testing.T) {
	setup()
	defer teardown()
	carsPerClass := []int{3, 1, 2, 2, 2, 3, 2}
	carsWithOptions := make([][]int, 5)
	carsWithOptions[0] = []int{2, 3, 7}
	carsWithOptions[1] = []int{1, 3, 5}
	carsWithOptions[2] = []int{2, 7}
	carsWithOptions[3] = []int{4, 7}
	carsWithOptions[4] = []int{2, 5}
	consecutiveCarsPerOption := []int{2, 3, 3, 5, 5}
	howManyOfconsecutiveCarsAtLeast := []int{0, 0, 0, 0, 0}
	howManyOfconsecutiveCarsAtMost := []int{1, 2, 1, 2, 1}
	numberOfCars := sum_intarray(carsPerClass)
	testCarSequencingWithoutAmong(t, carsPerClass, carsWithOptions,
		consecutiveCarsPerOption, howManyOfconsecutiveCarsAtLeast,
		howManyOfconsecutiveCarsAtMost, numberOfCars, true)
}

func Test_carSequencingWithoutAmongC(t *testing.T) {
	setup()
	defer teardown()
	carsPerClass := []int{1, 1, 2, 2, 2, 2, 3, 3, 4}
	carsWithOptions := make([][]int, 5)
	carsWithOptions[0] = []int{1, 5, 8}
	carsWithOptions[1] = []int{3, 4, 6}
	carsWithOptions[2] = []int{1, 5, 9}
	carsWithOptions[3] = []int{1, 5, 7}
	carsWithOptions[4] = []int{3, 4}
	consecutiveCarsPerOption := []int{4, 3, 3, 7, 5}
	howManyOfconsecutiveCarsAtLeast := []int{0, 0, 0, 0, 0}
	howManyOfconsecutiveCarsAtMost := []int{2, 2, 4, 2, 5}
	numberOfCars := sum_intarray(carsPerClass)
	testCarSequencingWithoutAmong(t, carsPerClass, carsWithOptions,
		consecutiveCarsPerOption, howManyOfconsecutiveCarsAtLeast,
		howManyOfconsecutiveCarsAtMost, numberOfCars, true)
}
