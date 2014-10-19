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

// mystery shopper problem modelled without among
// - set of salesladies and set of shoppers
// - each saleslady belongs to a certain location,
//   each shopper belongs to a group
// - every saleslady is visited by more than one shopper
// - each saleslady is visited by shoppers of different groups
// - each shopper visits saleslady of different locations

// testMysteryShopperWithoutAmong generates a mystery shopper problem
// with the given amount of salesladies working at shops and shoppers
// in groups
func testMysteryShopperWithoutAmong(t *testing.T,
	locationsWithSalesladies []int, groupsWithShoppers []int,
	numberOfVisits int, expectedResult bool) {
	numberOfSalesladies := sum_intarray(locationsWithSalesladies)
	numberOfShoppers := sum_intarray(groupsWithShoppers)
	log(fmt.Sprintf("Mystery Shopper without Among: "+
		"shoppers=%d, groups=%d, salesladies=%d",
		numberOfShoppers, len(groupsWithShoppers), numberOfSalesladies))
	// create IDs for all salesladies
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
			ivvar := core.CreateIntVarIvFromTo(s, store, 1, numberOfSalesladies)
			shopperIvisitVSalesladies[i] = ivvar
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
			// array to store results of the reified constraints per visit
			ShopperXVisitISalesladies := make([]core.VarId, numberOfSalesladies)
			for j := 0; j < numberOfSalesladies; j++ {
				// create variable to store the results of the reified
				// constraint per saleslady
				auxvar := core.CreateAuxIntVarIvFromTo(store, 0, 1)
				ShopperXVisitISalesladyJ := auxvar
				ShopperXVisitISalesladies[j] = ShopperXVisitISalesladyJ
				// only one visit per location
				// shopper X visits location Y during visit J
				xeqc := indexical.CreateXeqC(shopperIvisitVSalesladies[i], j+1)
				reif := reification.CreateReifiedConstraint(xeqc,
					ShopperXVisitISalesladyJ)
				store.AddPropagator(reif)
			}
			// during visit I shopper X can only visit one of all the
			// locations
			auxvar := core.CreateAuxIntVarIvFromTo(store, 1, numberOfLocations)
			ShopperXVisitILocationY[i] = auxvar
			store.AddPropagator(interval.CreateWeightedSum(store,
				ShopperXVisitILocationY[i], weightedLocations,
				ShopperXVisitISalesladies...))
		}
		// shoppers visit always different locations
		alldiffs := interval.CreateAlldifferent(ShopperXVisitILocationY...)
		store.AddPropagator(alldiffs)
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
			SalesladyIGroupXShopperYVisitJ := make([]core.VarId,
				groupsWithShoppers[y]*numberOfVisits)
			for ; x < len(allVisitsFromAllShoppers); x++ {
				// every shopper in a group
				shopperIvisitVSalesladies := allVisitsFromAllShoppers[x]
				for i := 0; i < numberOfVisits; i++ {
					// visits the saleslady at most one time
					SJGMSXVI := core.CreateAuxIntVarIvFromTo(store, 0, 1)
					xeqc := indexical.CreateXeqC(shopperIvisitVSalesladies[i],
						j+1)
					reif := reification.CreateReifiedConstraint(xeqc, SJGMSXVI)
					store.AddPropagator(reif)
					SalesladyIGroupXShopperYVisitJ[g] = SJGMSXVI
					g += 1
				}
				if g == groupsWithShoppers[y]*numberOfVisits {
					break
				}
			}
			V1G2 := core.CreateAuxIntVarIvFromTo(store,
				1, groupsWithShoppers[y])
			store.AddPropagator(interval.CreateSum(store,
				V1G2, SalesladyIGroupXShopperYVisitJ))
		}
	}
	query := labeling.CreateSearchOneQueryVariableSelect(allVisits)
	labeling.Labeling(store, query, labeling.VarSelect, labeling.InDomainMin)
	ready := store.IsConsistent()
	log(fmt.Sprintf("ready: %6v,    search nodes=%4d",
		ready, query.GetSearchStatistics().GetNodes()))
	ready_test(t, "Mysteryshopper without Among", ready, expectedResult)
}

func Test_mysteryShopperWithoutAmong3Shoppers3Salesladies(t *testing.T) {
	setup()
	defer teardown()
	testMysteryShopperWithoutAmong(t, []int{2, 1}, []int{2, 1}, 2, true)
}

func Test_mysteryShopperWithoutAmong4Shoppers3Salesladies(t *testing.T) {
	setup()
	defer teardown()
	testMysteryShopperWithoutAmong(t, []int{2, 1}, []int{2, 2}, 2, true)
}

func Test_mysteryShopperWithoutAmong6Shoppers3Salesladies(t *testing.T) {
	setup()
	defer teardown()
	testMysteryShopperWithoutAmong(t, []int{2, 1}, []int{2, 2, 2}, 1, true)
}

// Takes rather long
//func Test_mysteryShopperWithoutAmong15Shoppers5Salesladies(t *testing.T) {
//	setup()
//	defer teardown()
//	testMysteryShopperWithoutAmong(t, []int{2, 2, 1},
//		[]int{4, 4, 2, 2, 3}, 1, true)
//}
