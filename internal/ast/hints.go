package ast

type JennyHint string

const (
	HintDisjunctionOfScalars              JennyHint = "disjunction_of_scalars"
	HintDiscriminatedDisjunctionOfStructs JennyHint = "disjunction_of_structs"
)
