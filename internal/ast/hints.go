package ast

type JennyHint string

const (
	// HintDisjunctionOfScalars indicates that the struct was previously
	// represented in the IR by a disjunction of scalars (+ array), the
	// original definition of which is associated to this hint.
	HintDisjunctionOfScalars JennyHint = "disjunction_of_scalars"

	// HintDiscriminatedDisjunctionOfRefs indicates that the struct was
	// previously represented in the IR by a disjunction of a fixed list of
	// references to structs, the original definition of which is associated
	// to this hint.
	HintDiscriminatedDisjunctionOfRefs JennyHint = "disjunction_of_refs"
)
