package parser

import (
	"github.com/the-xlang/x/ast"
	"github.com/the-xlang/x/pkg/x"
)

var builtinFunctions = []*function{
	{
		Name: "out",
		ReturnType: ast.TypeAST{
			Type: x.Void,
		},
		Params: []ast.ParameterAST{{
			Name: "v",
			Type: ast.TypeAST{
				Value: "any",
				Type:  x.Any,
			},
		}},
	},
	{
		Name: "outln",
		ReturnType: ast.TypeAST{
			Type: x.Void,
		},
		Params: []ast.ParameterAST{{
			Name: "v",
			Type: ast.TypeAST{
				Value: "any",
				Type:  x.Any,
			},
		}},
	},
}