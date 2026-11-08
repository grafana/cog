package simplecue

import (
	"cuelang.org/go/cue"
	cueast "cuelang.org/go/cue/ast"
	"cuelang.org/go/cue/token"
	"github.com/grafana/cog/pkg/ir"
)

// maxIntegerLookupDepth bounds how far references and nested conjunctions
// are followed while looking for a declared integer type.
const maxIntegerLookupDepth = 32

// integerRange describes one of CUE's predeclared fixed-width integer types.
// In CUE, these are not distinct types but bounds: `int32` is defined as
// `int & >=-2147483648 & <=2147483647`.
type integerRange struct {
	kind ir.ScalarKind
	bits int
	min  string
	max  string
}

func integerRanges() []integerRange {
	return []integerRange{
		{kind: ir.KindInt8, bits: 8, min: "-128", max: "127"},
		{kind: ir.KindUint8, bits: 8, min: "0", max: "255"},
		{kind: ir.KindInt16, bits: 16, min: "-32768", max: "32767"},
		{kind: ir.KindUint16, bits: 16, min: "0", max: "65535"},
		{kind: ir.KindInt32, bits: 32, min: "-2147483648", max: "2147483647"},
		{kind: ir.KindUint32, bits: 32, min: "0", max: "4294967295"},
		{kind: ir.KindInt64, bits: 64, min: "-9223372036854775808", max: "9223372036854775807"},
		{kind: ir.KindUint64, bits: 64, min: "0", max: "18446744073709551615"},
	}
}

func isIntegerKind(kind ir.ScalarKind) bool {
	switch kind {
	case ir.KindInt8, ir.KindInt16, ir.KindInt32, ir.KindInt64,
		ir.KindUint8, ir.KindUint16, ir.KindUint32, ir.KindUint64:
		return true
	default:
		return false
	}
}

func integerRangeOf(kind ir.ScalarKind) (integerRange, bool) {
	for _, r := range integerRanges() {
		if r.kind == kind {
			return r, true
		}
	}

	return integerRange{}, false
}

// integerConjuncts holds the `int` type and the inclusive bounds that an
// integer value was built from, before CUE merged them.
type integerConjuncts struct {
	sawInt bool
	min    map[string]bool
	max    map[string]bool
}

// declaredIntegerKind returns the width of an integer value, derived from the
// predeclared types (int8, uint32, int64, ...) it was declared with.
//
// CUE's simplified export syntax, used to infer number types, rewrites bounds
// into predeclared type names: `int & >0 & <=2147483647` becomes `int32 & >0`.
// Relying on it narrows the width of an `int` (or `int64`) whenever a bound
// happens to fit a smaller type. Instead, this walks the value's expression
// tree, where a predeclared type still appears as its pair of bounds.
//
// The second return value is false when the expression couldn't be inspected,
// and the caller should fall back to the type inferred from the syntax.
func declaredIntegerKind(v cue.Value, inferred ir.ScalarKind) (ir.ScalarKind, bool) {
	conjuncts := integerConjuncts{min: map[string]bool{}, max: map[string]bool{}}
	conjuncts.collect(v, 0)

	if !conjuncts.sawInt {
		return "", false
	}

	var declared *integerRange
	for _, r := range integerRanges() {
		if !conjuncts.min[r.min] || !conjuncts.max[r.max] {
			continue
		}

		// Disjunctions (`int8 | int16`) need the widest type.
		if declared == nil || r.bits > declared.bits {
			declared = &r
		}
	}

	if declared != nil {
		return declared.kind, true
	}

	// `uint` is `int & >=0`, which is indistinguishable from an `int` written
	// with the same bound. Keep the unsigned type CUE's syntax gave it.
	if inferred == ir.KindUint64 && conjuncts.min["0"] {
		return ir.KindUint64, true
	}

	return ir.KindInt64, true
}

func (c *integerConjuncts) collect(v cue.Value, depth int) {
	if depth > maxIntegerLookupDepth {
		return
	}

	// Follow references to the definition they point to: evaluating them
	// would merge the bounds of the type they were declared with.
	if root, path := v.ReferencePath(); len(path.Selectors()) != 0 {
		c.collect(root.LookupPath(path), depth+1)
		return
	}

	op, args := v.Expr()
	switch op {
	case cue.AndOp, cue.OrOp:
		for _, arg := range args {
			c.collect(arg, depth+1)
		}
		return
	default:
		c.collectSyntax(v.Syntax(cue.Raw()))
	}
}

func (c *integerConjuncts) collectSyntax(node cueast.Node) {
	switch x := node.(type) {
	case *cueast.Ident:
		if x.Name == "int" {
			c.sawInt = true
		}
	case *cueast.ParenExpr:
		c.collectSyntax(x.X)
	case *cueast.BinaryExpr:
		switch x.Op {
		case token.AND:
			c.collectSyntax(x.X)
			c.collectSyntax(x.Y)
		case token.OR:
			if !isDefault(x.X) {
				c.collectSyntax(x.X)
			}
			if !isDefault(x.Y) {
				c.collectSyntax(x.Y)
			}
		}
	case *cueast.UnaryExpr:
		raw, ok := extractNumber(x.X)
		if !ok {
			return
		}

		switch x.Op {
		case token.GEQ:
			c.min[raw] = true
		case token.LEQ:
			c.max[raw] = true
		}
	}
}

// impliedByIntegerKind reports whether a bound is one of the limits of the
// given integer type, and is therefore already conveyed by the type itself.
func impliedByIntegerKind(kind ir.ScalarKind, op token.Token, raw string) bool {
	r, ok := integerRangeOf(kind)
	if !ok {
		return false
	}

	switch op {
	case token.GEQ:
		return raw == r.min
	case token.LEQ:
		return raw == r.max
	default:
		return false
	}
}
