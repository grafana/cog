package compiler

import (
	"fmt"
	"testing"

	"github.com/grafana/cog/internal/ast"
	"github.com/grafana/cog/internal/testutils"
)

func TestAppendCommentObjects(t *testing.T) {
	comment := "+k8s:openapi-gen=true"

	// Prepare test input
	obj := ast.NewObject("sandbox", "AString", ast.String())
	obj.Comments = []string{"This is a string"}
	schema := &ast.Schema{
		Package: "sandbox",
		Objects: testutils.ObjectsMap(obj),
	}

	expectedObj := obj.DeepCopy()
	expectedObj.AddToPassesTrail(fmt.Sprintf("AppendCommentObjects[%s]", comment))
	expectedObj.Comments = []string{"This is a string", comment}
	expected := &ast.Schema{
		Package: "sandbox",
		Objects: testutils.ObjectsMap(expectedObj),
	}

	// Run the compiler pass
	runPassOnSchema(t, &AppendCommentObjects{Comment: comment}, schema, expected)
}
