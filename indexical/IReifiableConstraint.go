package indexical

import (
	"bitbucket.org/gofd/gofd/core"
)

// IReifiableConstraint interface to be implemented by every
// IndexicalConstraint, which can be reified.
type IReifiableConstraint interface {
	IIndexicalConstraint
	// IsEntailed returns true iff the ReifiableConstraint is entailed
	IsEntailed() bool
	// Init to initialize a ReifiableConstraint. A ReifiableConstraint doesn't
	// register itself at the store. The ReifiedConstraints sets Domains and
	// store with this function.
	Init(store *core.Store, domains map[core.VarId]*core.IvDomain)
	// GetNegation returns the negated Constraints for propagating
	GetNegation() IReifiableConstraint
}
