package indexical

import (
	"bitbucket.org/gofd/gofd/core"
)

// IIndexicalConstraint interface to be implemented by every IndexicalConstraint,
type IIndexicalConstraint interface {
	core.Constraint
	GetIndexicalCollection() *IndexicalCollection
}
