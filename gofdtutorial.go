package main

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"bitbucket.org/gofd/gofd/propagator"
	"bitbucket.org/gofd/gofd/propagator/interval"
	"fmt"
)

func GettingStarted() {
	store := core.CreateStore()
	X := core.CreateIntVarFromTo("X", store, 0, 9)
	Y := core.CreateIntVarFromTo("Y", store, 0, 9)
	eq1 := propagator.CreateC1XplusC2YeqC3(1, X, 1, Y, 9)
	eq2 := propagator.CreateC1XplusC2YeqC3(2, X, 4, Y, 24)
	store.AddPropagator(eq1)
	store.AddPropagator(eq2)

	consistent := store.IsConsistent()
	fmt.Printf("consistent: %v \n", consistent)
	fmt.Printf("X: %s, Y: %s \n", store.GetDomain(X), store.GetDomain(Y))
}

func Labeling() {
	store := core.CreateStore()
	appetizer := core.CreateIntVarValues("appetizer", store, []int{1, 4})
	main := core.CreateIntVarValues("main", store, []int{7, 6})
	dessert := core.CreateIntVarValues("dessert", store, []int{2, 5})
	weight := core.CreateIntVarFromTo("sum", store, 0, 10)
	sum := propagator.CreateSum(store, weight,
		[]core.VarId{appetizer, main, dessert})
	store.AddPropagator(sum)

	query := labeling.CreateSearchAllQuery()
	solutionFound := labeling.Labeling(store, query,
		labeling.SmallestDomainFirst, labeling.InDomainMin)
	fmt.Printf("solutionFound: %v \n", solutionFound)
	if solutionFound {
		resultSet := query.GetResultSet()
		msg := "Solution %v: appetizer=%v, main=%v, dessert=%v, weigth=%v\n"
		for solutionNumber, result := range resultSet {
			fmt.Printf(msg, solutionNumber,
				result[appetizer], result[main],
				result[dessert], result[weight])
		}
	}
}

func EightQueens() {
	store := core.CreateStore()
	n := 8
	queens := make([]core.VarId, n)
	for i := 0; i < n; i++ {
		varname := fmt.Sprintf("Q%d", i)
		queens[i] = core.CreateIntVarIvFromTo(varname, store, 0, n-1)
	}
	prop := propagator.CreateAlldifferent(queens...)
	store.AddPropagators(prop)
	left_offset := make([]int, len(queens))
	right_offset := make([]int, len(queens))
	for i, _ := range queens {
		left_offset[i] = -i
		right_offset[i] = i
	}
	left_prop := interval.CreateAlldifferent_Offset(queens, left_offset)
	store.AddPropagator(left_prop)
	right_prop := interval.CreateAlldifferent_Offset(queens, right_offset)
	store.AddPropagator(right_prop)
	query := labeling.CreateSearchAllQuery()
	solutionFound := labeling.Labeling(store, query,
		labeling.SmallestDomainFirst, labeling.InDomainMin)
	if solutionFound {
		println(n, "queens problem has",
			len(query.GetResultSet()), "solutions.")
	}
}

func main() {
	fmt.Println("Welcome to gofd!")
	fmt.Println("  Only a sample program that runs through the tutorial.")
	fmt.Println("# Getting Started #")
	GettingStarted()
	fmt.Println("# Labeling #")
	Labeling()
	fmt.Println("# Eight Queens #")
	EightQueens()
	fmt.Println("# Done #")
	fmt.Println("https://bitbucket.org/gofd/gofd")
}
