package propagator

import (
	"bitbucket.org/gofd/gofd/core"
	"bitbucket.org/gofd/gofd/propagator/interval"
)

// CreateXltC creates the constraint X < C
func CreateXltC(x core.VarId, c int) *interval.XltC {
	return interval.CreateXltC(x, c)
}
