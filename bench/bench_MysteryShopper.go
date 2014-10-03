package main

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"bitbucket.org/gofd/gofd/propagator/explicit"
	"bitbucket.org/gofd/gofd/propagator/indexical"
	"bitbucket.org/gofd/gofd/propagator/interval"
	"bitbucket.org/gofd/gofd/propagator/reification"
	"fmt"
	"testing"
)

func main() {
	bench_MysteryShopper()
	bench_MysteryShopperWithoutAmong()
}

// the driver for everything benching IntVar
func bench_MysteryShopper() {
	benchd(bMysteryShopper1, bc{"name": "MysteryShopper", "size": "4"})
	benchd(bMysteryShopper2, bc{"name": "MysteryShopper", "size": "6"})
	benchd(bMysteryShopper3, bc{"name": "MysteryShopper", "size": "12"})
}

func bMysteryShopper1(b *testing.B) { bMysteryShopper(b, []int{2, 1}, []int{2, 2}, 2) }
func bMysteryShopper2(b *testing.B) { bMysteryShopper(b, []int{2, 1}, []int{2, 2, 2}, 2) }
func bMysteryShopper3(b *testing.B) { bMysteryShopper(b, []int{2, 1}, []int{4, 4, 4}, 2) }

func bMysteryShopper(b *testing.B, locationsWithSalesladies []int,
	groupsWithShoppers []int, numberOfVisits int) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	var query *labeling.SearchOneQuery
	var numberOfSalesladies int
	var numberOfShoppers int
	for i := 0; i < b.N; i++ {
		store := core.CreateStore()
		numberOfSalesladies = 0
		for i := 0; i < len(locationsWithSalesladies); i++ {
			numberOfSalesladies += locationsWithSalesladies[i]
		}
		numberOfShoppers := 0
		for i := 0; i < len(groupsWithShoppers); i++ {
			numberOfShoppers += groupsWithShoppers[i]
		}
		// create IDs for all salesladies
		salesladyID := make([][]int, numberOfSalesladies)
		for i := 0; i < numberOfSalesladies; i++ {
			salesladyID[i] = []int{i + 1}
		}
		// define visits from shoppers at salesladies
		allVisitsFromAllShoppers := make([][]core.VarId, numberOfShoppers)
		allVisits := make([]core.VarId, 0)
		for shopper := 0; shopper < numberOfShoppers; shopper++ {
			shopperIvisitVSalesladies := make([]core.VarId, numberOfVisits)
			for i := 0; i < numberOfVisits; i++ {
				s := fmt.Sprintf("Shopper%vVisit%v", shopper+1, i+1)
				shopperIvisitVSalesladies[i] =
					core.CreateIntVarExFromTo(s, store, 1, numberOfSalesladies)
				allVisits = append(allVisits, shopperIvisitVSalesladies[i])
			}
			allVisitsFromAllShoppers[shopper] = shopperIvisitVSalesladies
		}

		// determine which saleslady works in which shop:
		// e.g. if locationsWithSalesladies looks like this {2,1},
		// this means that two salesladies work at the first shop and one at
		// the second. it is automatically determined that the salesladies with
		// ID 1 and 2 work at location/shop 1 and saleslady 3 works at shop 2.
		shops := make([][]int, len(locationsWithSalesladies))
		shopIndex := 0
		salesladyLimit := 1
		for _, shopMaximum := range locationsWithSalesladies {
			salesladiesInShop := make([]int, shopMaximum)
			for i := 0; i < shopMaximum; i++ {
				salesladiesInShop[i] = salesladyLimit
				salesladyLimit++
			}
			shops[shopIndex] = salesladiesInShop
			shopIndex += 1
		}
		// define constraints
		// only one visit per location
		for _, shopper := range allVisitsFromAllShoppers {
			// every shopper visits every location at most one time
			for _, shop := range shops {
				store.AddPropagator(explicit.CreateAmong(shopper, shop,
					core.CreateAuxIntVarExFromTo(store, 0, 1)))
			}
			//every shopper visits no saleslady twice
			for _, saleslady := range salesladyID {
				store.AddPropagator(explicit.CreateAmong(shopper, saleslady,
					core.CreateAuxIntVarExFromTo(store, 0, 1)))
			}
		}
		//every saleslady is at least visited twice from shoppers from at least 2 different groups
		currentLimit := 0
		currentGroup := 0
		i := 0
		groups := make([][]core.VarId, len(groupsWithShoppers))
		for _, groupMaximum := range groupsWithShoppers {
			currentLimit += groupMaximum
			j := 0
			visits := make([]core.VarId, groupMaximum*numberOfVisits)
			for ; i < currentLimit; i++ {
				for m := 0; m < numberOfVisits; m++ {
					visits[j] = allVisitsFromAllShoppers[i][m]
					j += 1
				}
			}
			groups[currentGroup] = visits
			currentGroup += 1
		}
		for _, group := range groups {
			for _, saleslady := range salesladyID {
				store.AddPropagator(explicit.CreateAmong(group, saleslady,
					core.CreateAuxIntVarExFromTo(store, 1, numberOfVisits)))
			}
		}
		store.IsConsistent()
		// println("numberOfPropagators mystery among: ", store.GetNumPropagators())
		query = labeling.CreateSearchOneQueryVariableSelect(allVisits)
		labeling.Labeling(store, query, labeling.VarSelect, labeling.InDomainMin)
	}
	println("among mystery:", numberOfShoppers, "nodes:", query.GetSearchStatistics().GetNodes())
}

func bench_MysteryShopperWithoutAmong() {
	benchd(bMysteryShopperWithoutAmong1, bc{"name": "MysteryShopperWithoutAmong", "size": "4"})
	benchd(bMysteryShopperWithoutAmong2, bc{"name": "MysteryShopperWithoutAmong", "size": "6"})
	benchd(bMysteryShopperWithoutAmong3, bc{"name": "MysteryShopperWithoutAmong", "size": "12"})
}

func bMysteryShopperWithoutAmong1(b *testing.B) {
	bMysteryShopperWithoutAmong(b, []int{2, 1}, []int{2, 2}, 2)
}

func bMysteryShopperWithoutAmong2(b *testing.B) {
	bMysteryShopperWithoutAmong(b, []int{2, 1}, []int{2, 2, 2}, 2)
}

func bMysteryShopperWithoutAmong3(b *testing.B) {
	bMysteryShopperWithoutAmong(b, []int{2, 1}, []int{4, 4, 4}, 2)
}

func bMysteryShopperWithoutAmong(b *testing.B, locationsWithSalesladies []int,
	groupsWithShoppers []int, numberOfVisits int) {
	BenchSetRuns(b.N)
	b.StartTimer() // benchmark starts here
	var query *labeling.SearchOneQuery
	var numberOfSalesladies int
	var numberOfShoppers int
	for i := 0; i < b.N; i++ {
		store := core.CreateStore()
		numberOfSalesladies = 0
		for i := 0; i < len(locationsWithSalesladies); i++ {
			numberOfSalesladies += locationsWithSalesladies[i]
		}
		numberOfShoppers := 0
		for i := 0; i < len(groupsWithShoppers); i++ {
			numberOfShoppers += groupsWithShoppers[i]
		}
		//create IDs for all salesladies
		salesladyID := make([]core.VarId, numberOfSalesladies)
		for i := 0; i < numberOfSalesladies; i++ {
			salesladyID[i] = core.CreateAuxIntVarIvFromTo(store, i+1, i+1)
		}
		// define visits from shoppers at salesladies
		allVisitsFromAllShoppers := make([][]core.VarId, numberOfShoppers)
		allVisits := make([]core.VarId, 0)
		for j := 0; j < numberOfShoppers; j++ {
			shopperIvisitVSalesladies := make([]core.VarId, numberOfVisits)
			for i := 0; i < numberOfVisits; i++ {
				s := fmt.Sprintf("Shopper%vVisit%v", j+1, i+1)
				shopperIvisitVSalesladies[i] = core.CreateIntVarIvFromTo(s, store, 1, numberOfSalesladies)
				allVisits = append(allVisits, shopperIvisitVSalesladies[i])
			}
			allVisitsFromAllShoppers[j] = shopperIvisitVSalesladies
		}
		numberOfLocations := len(locationsWithSalesladies)
		// weightedLocations is needed for the weighted sum constraint later in
		// order to prevent a shopper from visiting the same location twice
		weightedLocations := make([]int, numberOfSalesladies)
		m := 0
		n := 1
		// create weightedlocations (add for every saleslady in a shop a weight
		// n in the array corresponding to the shop she is working in)
		for _, shopMaximum := range locationsWithSalesladies {
			for i := 0; i < shopMaximum; i++ {
				weightedLocations[m] = n
				m++
			}
			n++
		}
		// define constraints
		// go through the arrays containing all visits per shopper
		for _, shopperIvisitVSalesladies := range allVisitsFromAllShoppers {
			ShopperXVisitILocationY := make([]core.VarId, numberOfVisits)
			for i := 0; i < numberOfVisits; i++ {
				// array to store the results of the reified constraints per
				// visit
				ShopperXVisitISalesladies := make([]core.VarId, numberOfSalesladies)
				for j := 0; j < numberOfSalesladies; j++ {
					// create variable to store the results of the reified
					// constraint per saleslady
					ShopperXVisitISalesladyJ := core.CreateAuxIntVarIvFromTo(store, 0, 1)
					ShopperXVisitISalesladies[j] = ShopperXVisitISalesladyJ
					// only one visit per location
					// shopper X visits location Y during visit J
					xeqc := indexical.CreateXeqC(shopperIvisitVSalesladies[i], j+1)
					reifiedConstraint := reification.CreateReifiedConstraint(xeqc, ShopperXVisitISalesladyJ)
					store.AddPropagator(reifiedConstraint)
				}
				// during visit I shopper X can only visit one of all the
				// locations
				ShopperXVisitILocationY[i] = core.CreateAuxIntVarIvFromTo(store, 1, numberOfLocations)
				store.AddPropagator(interval.CreateWeightedSum(store, ShopperXVisitILocationY[i],
					weightedLocations, ShopperXVisitISalesladies...))
			}
			// shoppers visit always different locations
			store.AddPropagator(interval.CreateAlldifferent(ShopperXVisitILocationY...))
		}
		// every saleslady is at least visited twice from shoppers from at least
		// 2 different groups
		// example: shopper 1 and 2 from group 1 must have visited saleslady 1
		// at least 1 time and at most 2 times
		// for every saleslady
		for j := 0; j < numberOfSalesladies; j++ {
			x := 0
			for y := 0; y < len(groupsWithShoppers); y++ {
				g := 0
				SalesladyIGroupXShopperYVisitJ := make([]core.VarId, groupsWithShoppers[y]*numberOfVisits)
				for ; x < len(allVisitsFromAllShoppers); x++ {
					// every shopper in a group
					shopperIvisitVSalesladies := allVisitsFromAllShoppers[x]
					for i := 0; i < numberOfVisits; i++ {
						// visits the saleslady at most one time
						SJGMSXVI := core.CreateAuxIntVarIvFromTo(store, 0, 1)
						xeqc := indexical.CreateXeqC(shopperIvisitVSalesladies[i], j+1)
						reifiedConstraint := reification.CreateReifiedConstraint(xeqc, SJGMSXVI)
						store.AddPropagator(reifiedConstraint)
						SalesladyIGroupXShopperYVisitJ[g] = SJGMSXVI
						g += 1
					}
					if g == groupsWithShoppers[y]*numberOfVisits {
						break
					}
				}
				V1G2 := core.CreateAuxIntVarIvFromTo(store, 1, groupsWithShoppers[y])
				store.AddPropagator(interval.CreateSum(store, V1G2, SalesladyIGroupXShopperYVisitJ))
			}
		}
		store.IsConsistent()
		// println("numberOfPropagators mystery primitive: ", store.GetNumPropagators())
		query = labeling.CreateSearchOneQueryVariableSelect(allVisits)
		labeling.Labeling(store, query, labeling.VarSelect, labeling.InDomainMin)
	}
	println("primitive mystery:", numberOfShoppers, "nodes:", query.GetSearchStatistics().GetNodes())
}
