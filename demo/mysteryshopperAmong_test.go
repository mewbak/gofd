package demo

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/labeling"
	"bitbucket.org/gofd/gofd/propagator/explicit"
	"fmt"
	"testing"
)

// mystery shopper problem modelled with among
// - set of salesladies and set of shoppers
// - each saleslady belongs to a certain location,
//   each shopper belongs to a group
// - every saleslady is visited by more than one shopper
// - each saleslady is visited by shoppers of different groups
// - each shopper visits saleslady of different locations

// testMysteryShopperGAC generates a mystery shopper problem with the given
// amount of salesladies working at shops and shoppers in groups
func testMysteryShopperGAC(t *testing.T, locationsWithSalesladies []int,
	groupsWithShoppers []int, numberOfVisits int, expectedResult bool) {
	numberOfSalesladies := sum_intarray(locationsWithSalesladies)
	numberOfShoppers := sum_intarray(groupsWithShoppers)
	log(fmt.Sprintf("Mystery Shopper with    Among: "+
		"shoppers=%d, groups=%d, salesladies=%d",
		numberOfShoppers, len(groupsWithShoppers), numberOfSalesladies))
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
			ivvar := core.CreateIntVarExFromTo(s, store, 1, numberOfSalesladies)
			shopperIvisitVSalesladies[i] = ivvar
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
	salesladyLimit := 1
	for shopIndex, shopMaximum := range locationsWithSalesladies {
		salesladiesInShop := make([]int, shopMaximum)
		for i := 0; i < shopMaximum; i++ {
			salesladiesInShop[i] = salesladyLimit
			salesladyLimit++
		}
		shops[shopIndex] = salesladiesInShop
	}
	// define constraints
	// only one visit per location
	for _, shopper := range allVisitsFromAllShoppers {
		// every shopper visits every location at most one time
		for _, shop := range shops {
			auxvar := core.CreateAuxIntVarExFromTo(store, 0, 1)
			shop_once := explicit.CreateAmong(shopper, shop, auxvar)
			store.AddPropagator(shop_once)
		}
		// every shopper visits no saleslady twice
		for _, saleslady := range salesladyID {
			auxvar := core.CreateAuxIntVarExFromTo(store, 0, 1)
			not_twice := explicit.CreateAmong(shopper, saleslady, auxvar)
			store.AddPropagator(not_twice)
		}
	}
	// every saleslady is at least visited twice from shoppers
	// from at least two different groups
	currentLimit, currentGroup := 0, 0
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
			auxvar := core.CreateAuxIntVarExFromTo(store, 1, numberOfVisits)
			store.AddPropagator(explicit.CreateAmong(group, saleslady, auxvar))
		}
	}
	query := labeling.CreateSearchOneQueryVariableSelect(allVisits)
	labeling.Labeling(store, query, labeling.VarSelect, labeling.InDomainMin)
	ready := store.IsConsistent()
	log(fmt.Sprintf("ready: %6v,    search nodes=%4d",
		ready, query.GetSearchStatistics().GetNodes()))
	ready_test(t, "Mysteryshopper", ready, expectedResult)
}

func Test_mysteryShopper3Shoppers3Salesladies(t *testing.T) {
	setup()
	defer teardown()
	testMysteryShopperGAC(t, []int{2, 1}, []int{2, 1}, 2, true)
}

func Test_mysteryShopper4Shoppers3Salesladies(t *testing.T) {
	setup()
	defer teardown()
	testMysteryShopperGAC(t, []int{2, 1}, []int{2, 2}, 2, true)
}

func Test_mysteryShopper6Shoppers3Salesladies(t *testing.T) {
	setup()
	defer teardown()
	testMysteryShopperGAC(t, []int{2, 1}, []int{2, 2, 2}, 1, true)
}

func Test_mysteryShopper15Shoppers5Salesladies(t *testing.T) {
	setup()
	defer teardown()
	testMysteryShopperGAC(t, []int{2, 2, 1}, []int{4, 4, 2, 2, 3}, 1, true)
}
