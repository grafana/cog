package ast

const (
	// HintDisjunctionOfScalars indicates that the struct was previously
	// represented in the IR by a disjunction of scalars (+ array), the
	// original definition of which is associated to this hint.
	HintDisjunctionOfScalars = "disjunction_of_scalars"

	// HintDiscriminatedDisjunctionOfRefs indicates that the struct was
	// previously represented in the IR by a disjunction of a fixed list of
	// references to structs, the original definition of which is associated
	// to this hint.
	HintDiscriminatedDisjunctionOfRefs = "disjunction_of_refs"
)

const DiscriminatorCatchAll = "cog_discriminator_catch_all"
