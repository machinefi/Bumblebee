package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	g "github.com/iotexproject/Bumblebee/gen/codegen"
	"github.com/saitofun/qlib/util/qnaming"
)

func main() {
	fset := token.NewFileSet()
	f, _ := parser.ParseFile(fset, "strfmt.go", nil, parser.ParseComments)
	file := g.NewFile("strfmt", "strfmt_generated.go")

	regexps := make([]string, 0)
	for key, obj := range f.Scope.Objects {
		if obj.Kind == ast.Con {
			regexps = append(regexps, key)
		}
	}

	snippets := make([]g.Snippet, 0)
	for _, key := range regexps {
		var (
			name           = strings.Replace(key, "regexpString", "", 1)
			validatorName  = strings.Replace(qnaming.LowerSnakeCase(name), "_", "-", -1)
			validatorAlias = qnaming.LowerCamelCase(name)
			args           = []g.Snippet{g.Ident(key), g.Valuer(validatorName)}
			prefix         = qnaming.UpperCamelCase(name)
			snippet        g.Snippet
		)
		if validatorName != validatorAlias {
			args = append(args, g.Valuer(validatorAlias))
		}
		snippet = g.Func().Named("init").Do(
			g.Ref(
				g.Ident(file.Use(pkg, "DefaultFactory")),
				g.Call(
					"Register",
					g.Ident(prefix+"Validator"),
				),
			),
		)
		snippets = append(snippets, snippet)
		snippet = g.DeclVar(
			g.Assign(g.Var(nil, prefix+"Validator")).
				By(g.Call(file.Use(pkg, "NewRegexpStrfmtValidator"), args...)),
		)
		snippets = append(snippets, snippet)

	}
	file.WriteSnippet(snippets...)
	_, _ = file.Write(true)
}

var pkg = "github.com/saitofun/qkit/kit/validator"
