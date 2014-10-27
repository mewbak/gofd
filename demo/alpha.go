package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator"
	"fmt"
)

// ConstrainAlpha imposes an alphabetic puzzle on a store
// http://www.picat-lang.org/bprolog/examples/clpfd/alpha.pl
// Each letter represents a number. The sum of the letters
// of a word is associated to a number. Find the value of
// each letter. Given the word/value combination it returns
// the extracted letter to variable id association in the store.
func ConstrainAlpha(store *core.Store,
	words map[string]int) []core.VarId {
	vars := make(map[string]core.VarId)
	avars := make([]core.VarId, 26)
	for i := 0; i < len(avars); i++ {
		varname := fmt.Sprintf("%c", 'A'+i)
		varid := core.CreateIntVarFromTo(varname, store, 1, len(avars))
		vars[varname] = varid
		avars[i] = varid
	}
	alldiff := propagator.CreateAlldifferent(avars...)
	store.AddPropagator(alldiff)
	for word, value := range words {
		wordvarsmap := make(map[core.VarId]int)
		for _, rune := range word {
			letter := fmt.Sprintf("%c", rune)
			if _, ok := vars[letter]; !ok {
				panic(fmt.Sprintf("letter %v not found", letter))
			}
			varid := vars[letter]
			if _, ok := wordvarsmap[varid]; !ok {
				wordvarsmap[varid] = 0
			}
			// collect the number of occurences
			wordvarsmap[varid] = wordvarsmap[varid] + 1
			// ToDo: Sum propagator should find equal variables on its own
		}
		wordvars := make([]core.VarId, len(wordvarsmap))
		wordoccs := make([]int, len(wordvarsmap))
		i := 0
		for varid, noocc := range wordvarsmap {
			wordvars[i] = varid
			wordoccs[i] = noocc
			i++
		}
		zvar := core.CreateIntVarFromTo(word, store, 0, value)
		sumprop := propagator.CreateWeightedSumBounds(store,
			zvar, wordoccs, wordvars...)
		//for idx, varid := range wordvars {
		//	fmt.Printf("%d*%s ", wordoccs[idx], store.GetName(varid))
		//}
		//fmt.Printf("\n")
		store.AddPropagator(sumprop)
		store.AddPropagator(propagator.CreateXeqC(zvar, value)) // fix late
	}
	return avars // sorted ascending letters
}

// Example:
//    BALLET  45     GLEE  66     POLKA      59     SONG     61
//    CELLO   43     JAZZ  58     QUARTET    50     SOPRANO  82
//    CONCERT 74     LYRE  47     SAXOPHONE 134     THEME    72
//    FLUTE   30     OBOE  53     SCALE      51     VIOLIN  100
//    FUGUE   50     OPERA 65     SOLO       37     WALTZ    34
//
// Solution:
//  [A, B,C, D, E,F, G, H, I, J, K,L,M, N, O, P,Q, R, S,T,U, V,W, X, Y, Z]
//  [5,13,9,16,20,4,24,21,25,17,23,2,8,12,10,19,7,11,15,3,1,26,6,22,14,18]
func GenerateAlpha1() (problem map[string]int, solution map[string]int) {
	problem = make(map[string]int)
	solution = make(map[string]int)
	problem["BALLET"] = 45
	problem["GLEE"] = 66
	problem["POLKA"] = 59
	problem["SONG"] = 61
	problem["CELLO"] = 43
	problem["JAZZ"] = 58
	problem["QUARTET"] = 50
	problem["SOPRANO"] = 82
	problem["CONCERT"] = 74
	problem["LYRE"] = 47
	problem["SAXOPOHONE"] = 134
	problem["THEME"] = 72
	problem["FLUTE"] = 30
	problem["OBOE"] = 53
	problem["SCALE"] = 51
	problem["VIOLIN"] = 100
	problem["FUGUE"] = 50
	problem["OPERA"] = 65
	problem["SOLO"] = 37
	problem["WALTZ"] = 34
	sols := []int{5, 13, 9, 16, 20, 4, 24, 21, 25, 17, 23, 2, 8, 12, 10,
		19, 7, 11, 15, 3, 1, 26, 6, 22, 14, 18}
	for i, val := range sols {
		solution[fmt.Sprintf("%c", 'A'+i)] = val
	}
	return
}
