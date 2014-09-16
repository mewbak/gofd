package labeling

import (
	"bitbucket.org/gofd/gofd/core"
	"fmt"
)

/* Search Strategies
 * How are the values within a domain of a variable selected.
 */

var Allvars []core.VarId

func SetAllvars(vars []core.VarId) {
	Allvars = vars
}

func VarSelect(store *core.Store) (core.VarId, bool) {
	for _, v := range Allvars {
		if !store.GetDomain(v).IsGround() {
			return v, true
		}
	}
	return -1, false
}

func createSizedChannel(domain core.Domain) chan int {
	channelSize := domain.Size()
	if channelSize > 17 { // fill in one go for small domain sizes
		channelSize = 17 // limit for huge domains
	}
	return make(chan int, channelSize)
}

// InDomainRange is a search strategy that iterates over the given Domain
// and sends its values on the returned Channel. The order of the provided
// values is undefined.
func InDomainRange(domain core.Domain) chan int {
	channel := createSizedChannel(domain)
	go func() {
		loggerInfo := core.GetLogger().DoInfo()
		if loggerInfo {
			core.GetLogger().If("InDomainRange of %s", domain)
		}
		for _, k := range domain.Values_asSlice() {
			if loggerInfo {
				core.GetLogger().If("InDomainRange choose value %d", k)
			}
			channel <- k
		}
		close(channel)
	}()
	return channel
}

// InDomainMin is a search strategy that iterates over the given Domain
// and puts its values on the returned Channel. The provided values are
// sorted in ascending order.
func InDomainMin(domain core.Domain) chan int {
	channel := createSizedChannel(domain)
	go func() {
		loggerInfo := core.GetLogger().DoInfo()
		if loggerInfo {
			core.GetLogger().If("InDomainMin of %v", domain)
		}
		for _, v := range domain.SortedValues() {
			if loggerInfo {
				core.GetLogger().If("InDomainMin choose value %d", v)
			}
			channel <- v
		}
		close(channel)
	}()
	return channel
}

/* Variable Selection Strategies
 * Which variable is used to be exploited next during search
 */

// GreatestDomainFirst returns true and the id of a non-ground variable
// with maximal domain size if available, false otherwise.
func GreatestDomainFirst(store *core.Store) (core.VarId, bool) {
	maxId := store.SelectVarIdUnfixedDomain(false)
	return maxId, maxId != -1
}

// SmallesDomainFirst returns true and the id of a non-ground variable
// with minimal domain size if available, false otherwise.
func SmallestDomainFirst(store *core.Store) (core.VarId, bool) {
	minId := store.SelectVarIdUnfixedDomain(true)
	return minId, minId != -1
}

/* Statistics */

// SearchStatistics holds the statistic values to report on labeling.
// It also contains the cumulated statistics for the stores that are
// generated during search.
type SearchStatistics struct {
	numNodes              int
	numInitialVars        int
	numInitialPropagators int
	numFailedNode         int
	statStore             *core.StoreStatistics
}

// CreateSearchStatistics returns a new empty statistics instance
func CreateSearchStatistics() *SearchStatistics {
	stats := new(SearchStatistics)
	stats.numNodes = 0
	stats.numInitialVars = 0
	stats.numInitialPropagators = 0
	stats.numFailedNode = 0
	stats.statStore = core.CreateStoreStatistics()
	return stats
}

func (this *SearchStatistics) IncNodes() {
	this.numNodes += 1
}

func (this *SearchStatistics) setInitialVars(n int) {
	this.numInitialVars = n
}

func (this *SearchStatistics) GetInitialVars() int {
	return this.numInitialVars
}

func (this *SearchStatistics) setInitialPropagators(n int) {
	this.numInitialPropagators = n
}

func (this *SearchStatistics) GetInitialPropagators() int {
	return this.numInitialPropagators
}

func (this *SearchStatistics) GetNodes() int {
	return this.numNodes
}

func (this *SearchStatistics) IncFailedNodes() {
	this.numFailedNode += 1
}

func (this *SearchStatistics) GetFailedNodes() int {
	return this.numFailedNode
}

func (this *SearchStatistics) GetStoreStatistics() *core.StoreStatistics {
	return this.statStore
}

func (this *SearchStatistics) UpdateStoreStatistics(store *core.Store) {
	this.statStore.UpdateStoreStatistics(store.GetStat())
}

func (this *SearchStatistics) UpdateSearchStatistics(other *SearchStatistics) {
	this.numInitialVars = other.numInitialVars
	this.numInitialPropagators = other.numInitialPropagators
	this.numFailedNode += other.numFailedNode
	this.numNodes += other.numNodes
	this.statStore.UpdateStoreStatistics(other.statStore)
}

func (this *SearchStatistics) SearchString() string {
	msg := "%d nodes, %d vars, %d propagators "
	return fmt.Sprintf(msg, this.numNodes, this.numInitialVars,
		this.numInitialPropagators)
}

func (this *SearchStatistics) StoreString() string {
	return this.statStore.String()
}

// ResultQuery ...
type ResultQuery interface {
	onResult(store *core.Store) bool
	getResultStatus() bool
	GetResultSet() map[int]map[core.VarId]int
	GetSearchStatistics() *SearchStatistics
}

func CreateSearchOneQuery() *SearchOneQuery {
	query := new(SearchOneQuery)
	query.status = false
	query.searchStats = CreateSearchStatistics()
	return query
}

func CreateSearchOneQueryVariableSelect(variables []core.VarId) *SearchOneQuery {
	Allvars = variables
	return CreateSearchOneQuery()
}

type SearchOneQuery struct {
	status      bool
	results     map[core.VarId]int
	searchStats *SearchStatistics
}

func (this *SearchOneQuery) onResult(store *core.Store) bool {
	this.results = make(map[core.VarId]int)
	varIds := store.GetVariableIDs()
	domains := store.GetDomains(varIds)
	for i, id := range varIds {
		this.results[id] = domains[i].GetAnyElement()
	}
	this.status = true
	return true //one result found, abort
}
func (this *SearchOneQuery) getResultStatus() bool {
	return this.status
}

// GetResultSet returns a map of solutions
// map[int]map[core.VarId]int --> map[solutionnumber]map[varid]solutionvalue
func (this *SearchOneQuery) GetResultSet() map[int]map[core.VarId]int {
	results := make(map[int]map[core.VarId]int)
	results[0] = this.results
	return results
}

func (this *SearchOneQuery) GetSearchStatistics() *SearchStatistics {
	return this.searchStats
}

func CreateSearchAllQuery() *SearchAllQuery {
	query := new(SearchAllQuery)
	query.searchStats = CreateSearchStatistics()
	query.status = false
	return query
}

type SearchAllQuery struct {
	status      bool
	results     map[int]map[core.VarId]int
	searchStats *SearchStatistics
}

func (this *SearchAllQuery) onResult(store *core.Store) bool {
	if this.results == nil {
		this.results = make(map[int]map[core.VarId]int)
	}
	index := len(this.results)
	this.results[index] = make(map[core.VarId]int)
	varIds := store.GetVariableIDs()
	domains := store.GetDomains(varIds)
	for i, id := range varIds {
		this.results[index][id] = domains[i].GetAnyElement()
	}
	this.status = true
	return false //still more solutions possible, do not abort
}

func (this *SearchAllQuery) getResultStatus() bool {
	return this.status
}

func (this *SearchAllQuery) GetResultSet() map[int]map[core.VarId]int {
	return this.results
}

func (this *SearchAllQuery) GetSearchStatistics() *SearchStatistics {
	return this.searchStats
}
