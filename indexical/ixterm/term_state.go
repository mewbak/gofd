package ixterm

type EvalState int

const (
	//states of term, range, indexical
	EVALUABLE         = 0
	NOT_EVALUABLE_YET = 1
	EMPTY             = 2
)
