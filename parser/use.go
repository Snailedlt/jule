package parser

import "github.com/julelang/jule/lex"

type use struct {
	defines *Defmap
	token   lex.Token
	cppLink bool
	
	FullUse    bool
	Path       string
	LinkString string
	Selectors  []lex.Token
}
