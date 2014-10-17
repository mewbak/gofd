package main

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"bitbucket.org/gofd/gofd/propagator/explicit"
	"bitbucket.org/gofd/gofd/propagator/indexical"
	"bitbucket.org/gofd/gofd/propagator/interval"
	"bitbucket.org/gofd/gofd/propagator/reification"
	"testing"
)

func main() {
	bench_CarSequencing()
	bench_CarSequencingWithoutAmong()
}

// the driver for everything benching CarSequencing
func bench_CarSequencing() {
	benchd(bCarSequencing1, bc{"name": "CarSequencing", "size": "10"})
	benchd(bCarSequencing2, bc{"name": "CarSequencing", "size": "15"})
	benchd(bCarSequencing3, bc{"name": "CarSequencing", "size": "20"})
}

func bCarSequencing1(b *testing.B) {
	carsPerClass := []int{1, 1, 2, 2, 2, 2}
	carsWithOptions := make([][]int, 5)
	carsWithOptions[0] = []int{1, 5, 6}
	carsWithOptions[1] = []int{3, 4, 6}
	carsWithOptions[2] = []int{1, 5}
	carsWithOptions[3] = []int{1, 2, 4}
	carsWithOptions[4] = []int{3}
	consecutiveCarsPerOption := []int{2, 3, 3, 5, 5}
	howManyOfconsecutiveCarsAtLeast := []int{0, 0, 0, 0, 0}
	howManyOfconsecutiveCarsAtMost := []int{1, 2, 1, 2, 1}
	numberOfCars := 10
	bCarSequencing(b, carsPerClass, carsWithOptions,
		consecutiveCarsPerOption, howManyOfconsecutiveCarsAtLeast,
		howManyOfconsecutiveCarsAtMost, numberOfCars)
}

func bCarSequencing2(b *testing.B) {
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
	numberOfCars := 15
	bCarSequencing(b, carsPerClass, carsWithOptions,
		consecutiveCarsPerOption, howManyOfconsecutiveCarsAtLeast,
		howManyOfconsecutiveCarsAtMost, numberOfCars)
}

func bCarSequencing3(b *testing.B) {
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
	numberOfCars := 20
	bCarSequencing(b, carsPerClass, carsWithOptions,
		consecutiveCarsPerOption, howManyOfconsecutiveCarsAtLeast,
		howManyOfconsecutiveCarsAtMost, numberOfCars)
}

func bCarSequencing(b *testing.B, carsPerClass []int, carsWithOptions [][]int,
	consecutiveCars []int, howManyOfconsecutiveCarsAtLeast []int,
	howManyOfconsecutiveCarsAtMost []int, numberOfCars int) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	var query *labeling.SearchOneQuery
	for i := 0; i < b.N; i++ {
		store := core.CreateStore()
		numberOfDifferentClasses := len(carsPerClass)

		cars := make([]core.VarId, numberOfCars)
		for i := 0; i < len(cars); i++ {
			cars[i] = core.CreateAuxIntVarExFromTo(store, 1, numberOfDifferentClasses)
		}

		// define constraints
		// every car belongs to one of a number of classes
		// (all cars in a class have the same set of options)
		for i := 0; i < numberOfDifferentClasses; i++ {
			store.AddPropagator(explicit.CreateAmong(cars, []int{i + 1},
				core.CreateAuxIntVarExFromTo(store, carsPerClass[i], carsPerClass[i])))
		}

		// automated constraint modelling of problem
		// j is used as index for arrays concerning the options
		for j := 0; j < len(carsWithOptions); j++ {
			// create intvar with domain containing cars which need the current
			// option
			option := carsWithOptions[j]
			for i := 0; i < len(cars)-consecutiveCars[j]+1; i++ {
				// create array containing the consecutive cars that can have
				// the current options
				curCars := make([]core.VarId, consecutiveCars[j])
				for n := 0; n < len(curCars); n++ {
					curCars[n] = cars[i+n]
				}
				// constrain how many consecutive cars can have the current
				// option
				store.AddPropagator(explicit.CreateAmong(curCars, option,
					core.CreateAuxIntVarExFromTo(store,
						howManyOfconsecutiveCarsAtLeast[j],
						howManyOfconsecutiveCarsAtMost[j])))
			}
		}

		//		println("numberOfPropagators car sequencing: ", store.GetNumPropagators())
		query = labeling.CreateSearchOneQueryVariableSelect(cars)
		labeling.Labeling(store, query, labeling.VarSelect, labeling.InDomainMin)
	}
	println("among cars:", numberOfCars, "nodes:", query.GetSearchStatistics().GetNodes())
}

func bench_CarSequencingWithoutAmong() {
	benchd(bCarSequencingWithoutAmong1, bc{"name": "CarSequencingWithoutAmong", "size": "10"})
	benchd(bCarSequencingWithoutAmong2, bc{"name": "CarSequencingWithoutAmong", "size": "15"})
	benchd(bCarSequencingWithoutAmong3, bc{"name": "CarSequencingWithoutAmong", "size": "20"})
}

func bCarSequencingWithoutAmong1(b *testing.B) {
	carsPerClass := []int{1, 1, 2, 2, 2, 2}
	carsWithOptions := make([][]int, 5)
	carsWithOptions[0] = []int{1, 5, 6}
	carsWithOptions[1] = []int{3, 4, 6}
	carsWithOptions[2] = []int{1, 5}
	carsWithOptions[3] = []int{1, 2, 4}
	carsWithOptions[4] = []int{3}
	consecutiveCarsPerOption := []int{2, 3, 3, 5, 5}
	howManyOfconsecutiveCarsAtLeast := []int{0, 0, 0, 0, 0}
	howManyOfconsecutiveCarsAtMost := []int{1, 2, 1, 2, 1}
	numberOfCars := 10
	bCarSequencingWithoutAmong(b, carsPerClass, carsWithOptions,
		consecutiveCarsPerOption, howManyOfconsecutiveCarsAtLeast,
		howManyOfconsecutiveCarsAtMost, numberOfCars)
}

func bCarSequencingWithoutAmong2(b *testing.B) {
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
	numberOfCars := 15
	bCarSequencingWithoutAmong(b, carsPerClass, carsWithOptions,
		consecutiveCarsPerOption, howManyOfconsecutiveCarsAtLeast,
		howManyOfconsecutiveCarsAtMost, numberOfCars)
}

func bCarSequencingWithoutAmong3(b *testing.B) {
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
	numberOfCars := 20
	bCarSequencingWithoutAmong(b, carsPerClass, carsWithOptions,
		consecutiveCarsPerOption, howManyOfconsecutiveCarsAtLeast,
		howManyOfconsecutiveCarsAtMost, numberOfCars)
}

func bCarSequencingWithoutAmong(b *testing.B, carsPerClass []int, carsWithOptions [][]int, consecutiveCars []int, howManyOfconsecutiveCarsAtLeast []int, howManyOfconsecutiveCarsAtMost []int, numberOfCars int) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	var query *labeling.SearchOneQuery
	for i := 0; i < b.N; i++ {
		store := core.CreateStore()
		numberOfDifferentClasses := len(carsPerClass)

		// define car variables
		cars := make([]core.VarId, numberOfCars)
		for i := 0; i < len(cars); i++ {
			cars[i] = core.CreateAuxIntVarIvFromTo(store, 1, numberOfDifferentClasses)
		}

		// define constraints
		// every car belongs to one of six classes (all cars in a class have the same set of options)
		for i := 0; i < numberOfDifferentClasses; i++ {
			//array to store variables with a {0,1} domain which are needed for reification
			variables := make([]core.VarId, len(cars))
			for j := 0; j < len(cars); j++ {
				variables[j] = core.CreateAuxIntVarIvFromTo(store, 0, 1)
				xeqc := indexical.CreateXeqC(cars[j], i+1)
				reifiedConstraint := reification.CreateReifiedConstraint(xeqc, variables[j])
				store.AddPropagator(reifiedConstraint)
			}
			store.AddPropagator(interval.CreateSum(store,
				core.CreateAuxIntVarIvFromTo(store, carsPerClass[i], carsPerClass[i]), variables))
		}

		// automated constraint modelling of problem
		// j is used as index for arrays concerning the options
		for j := 0; j < len(carsWithOptions); j++ {
			for i := 0; i < len(cars)-consecutiveCars[j]+1; i++ {
				// create array to store the results of reification
				// its size must be the number of consecutive cars that can have the option
				// multiplied with the number of options the class has
				currentCarsWithOptions := carsWithOptions[j]
				variables := make([]core.VarId, consecutiveCars[j]*(len(currentCarsWithOptions)))
				variablesIndex := 0
				for n := 0; n < consecutiveCars[j]; n++ {
					for opt := 0; opt < len(currentCarsWithOptions); opt++ {
						variables[variablesIndex] = core.CreateAuxIntVarIvFromTo(store, 0, 1)
						xeqc := indexical.CreateXeqC(cars[i+n], currentCarsWithOptions[opt])
						reified := reification.CreateReifiedConstraint(xeqc, variables[variablesIndex])
						store.AddPropagator(reified)
						variablesIndex += 1
					}
				}
				// constrain how many consecutive cars can have the current
				// option
				store.AddPropagator(interval.CreateSum(store, core.CreateAuxIntVarIvFromTo(store,
					howManyOfconsecutiveCarsAtLeast[j],
					howManyOfconsecutiveCarsAtMost[j]), variables))
			}
		}
		query = labeling.CreateSearchOneQueryVariableSelect(cars)
		labeling.Labeling(store, query, labeling.VarSelect, labeling.InDomainMin)
	}
	println("primitive cars:", numberOfCars, "nodes:", query.GetSearchStatistics().GetNodes())
}
