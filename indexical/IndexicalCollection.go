package indexical

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/indexical/ixrange"
	"strings"
)

//ToDo: Test

type IndexicalCollection struct {
	indexicals [][]*Indexical
}

func (this *IndexicalCollection) String() string {
	indexicals := this.GetIndexicalsAtPrio(HIGHEST)
	s := []string{"HIGHEST"}
	s = append(s, getIndexicalsAsString(indexicals)...)
	indexicals = this.GetIndexicalsAtPrio(HIGH)
	s = append(s, "HIGH")
	s = append(s, getIndexicalsAsString(indexicals)...)
	indexicals = this.GetIndexicalsAtPrio(LOW)
	s = append(s, "LOW")
	s = append(s, getIndexicalsAsString(indexicals)...)
	indexicals = this.GetIndexicalsAtPrio(LOWEST)
	s = append(s, "LOWEST")
	s = append(s, getIndexicalsAsString(indexicals)...)
	return strings.Join(s, "\r\n")
}

func getIndexicalsAsString(indexicals []*Indexical) []string {
	sh := make([]string, len(indexicals))
	for i, indexical := range indexicals {
		sh[i] = indexical.String()
	}
	return sh
}

// CreateIndxicalCollection creates ....
func CreateIndexicalCollection() *IndexicalCollection {
	indColl := new(IndexicalCollection)
	indColl.indexicals = make([][]*Indexical, 4)
	return indColl
}

func (this *IndexicalCollection) CreateAndAddIndexical(varid core.VarId,
	d *core.IvDomain, r ixrange.IRange, prio int) {
	indexical := CreateIndexical(varid, d, r)
	this.indexicals[prio] = append(this.indexicals[prio], indexical)
}

func (this *IndexicalCollection) AddIndexicalsAtPrio(indexicals []*Indexical, prio int) {
	this.indexicals[prio] = append(this.indexicals[prio], indexicals...)
}

func (this *IndexicalCollection) GetIndexicalsAtPrio(prio int) []*Indexical {
	if prio > LOWEST {
		panic("indexical-prio undefined")
	}
	return this.indexicals[prio]
}
