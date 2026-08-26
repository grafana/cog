package transforms

import (
	"fmt"

	"github.com/grafana/cog/internal/ir"
	"github.com/grafana/cog/internal/tools"
)

var _ Transform = (*SanitizeEnumMemberNames)(nil)

type SanitizeEnumMemberNames struct {
}

func (pass *SanitizeEnumMemberNames) Process(schemas ir.Schemas) (ir.Schemas, error) {
	visitor := &Visitor{
		OnEnum: pass.processEnum,
	}

	return visitor.VisitSchemas(schemas)
}

func (pass *SanitizeEnumMemberNames) processEnum(_ *Visitor, _ *ir.Schema, def ir.Type) (ir.Type, error) {
	def.Enum.Values = tools.Map(def.Enum.Values, pass.sanitizeEnumMember)

	return def, nil
}

func (pass *SanitizeEnumMemberNames) sanitizeEnumMember(member ir.EnumValue) ir.EnumValue {
	if member.Type.Scalar.ScalarKind == ir.KindString && member.Name == "" && member.Value.(string) == "" {
		member.Name = "None"
	}

	if member.Name[0] == '-' {
		member.Name = tools.UpperCamelCase(fmt.Sprintf("negative%s", member.Name[1:]))
	}
	if member.Name[0] == '+' {
		member.Name = tools.UpperCamelCase(fmt.Sprintf("positive%s", member.Name[1:]))
	}

	return member
}
