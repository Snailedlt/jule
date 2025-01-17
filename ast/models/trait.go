package models

import "github.com/julelang/jule/lex"

// Trait is the AST model of traits.
type Trait struct {
	Pub   bool
	Token lex.Token
	Id    string
	Desc  string
	Used  bool
	Funcs []*Fn
}
