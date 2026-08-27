package transforms

import (
	"fmt"
	"testing"

	"github.com/grafana/cog/internal/testutils"
	"github.com/grafana/cog/pkg/ir"
)

func TestAppendCommentObjects(t *testing.T) {
	comment := "+k8s:openapi-gen=true"

	// Prepare test input
	obj := ir.NewObject("sandbox", "AString", ir.String())
	obj.Comments = []string{"This is a string"}
	schema := &ir.Schema{
		Package: "sandbox",
		Objects: testutils.ObjectsMap(obj),
	}

	expectedObj := obj.DeepCopy()
	expectedObj.AddToPassesTrail(fmt.Sprintf("AppendCommentObjects[%s]", comment))
	expectedObj.Comments = []string{"This is a string", comment}
	expected := &ir.Schema{
		Package: "sandbox",
		Objects: testutils.ObjectsMap(expectedObj),
	}

	// Run the compiler pass
	runPassOnSchema(t, &AppendCommentObjects{Comment: comment}, schema, expected)
}
