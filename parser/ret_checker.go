package parser

import (
	"github.com/jule-lang/jule/ast/models"
	"github.com/jule-lang/jule/lex"
	"github.com/jule-lang/jule/lex/tokens"
	"github.com/jule-lang/jule/pkg/juleapi"
)

type retChecker struct {
	p         *Parser
	ret_ast   *models.Ret
	f         *Func
	exp_model retExpr
}

func (rc *retChecker) pushval(last, current int, errTok lex.Token) {
	if current-last == 0 {
		rc.p.pusherrtok(errTok, "missing_expr")
		return
	}
	toks := rc.ret_ast.Expr.Tokens[last:current]
	v, model := rc.p.evalToks(toks)
	rc.exp_model.models = append(rc.exp_model.models, model)
	rc.exp_model.values = append(rc.exp_model.values, v)
}

func (rc *retChecker) checkepxrs() {
	brace_n := 0
	last := 0
	for i, tok := range rc.ret_ast.Expr.Tokens {
		if tok.Id == tokens.Brace {
			switch tok.Kind {
			case tokens.LBRACE, tokens.LBRACKET, tokens.LPARENTHESES:
				brace_n++
			default:
				brace_n--
			}
		}
		if brace_n > 0 || tok.Id != tokens.Comma {
			continue
		}
		rc.pushval(last, i, tok)
		last = i + 1
	}
	n := len(rc.ret_ast.Expr.Tokens)
	if last < n {
		if last == 0 {
			rc.pushval(0, n, rc.ret_ast.Token)
		} else {
			rc.pushval(last, n, rc.ret_ast.Expr.Tokens[last-1])
		}
	}
	if !typeIsVoid(rc.f.RetType.Type) {
		rc.checkExprTypes()
		rc.ret_ast.Expr.Model = rc.exp_model
	}
}

func (rc *retChecker) single() {
	rc.exp_model.models = append(rc.exp_model.models, rc.exp_model.models[0])
	if len(rc.exp_model.values) > 1 {
		rc.p.pusherrtok(rc.ret_ast.Token, "overflow_return")
	}
	assignChecker{
		p:      rc.p,
		t:      rc.f.RetType.Type,
		v:      rc.exp_model.values[0],
		errtok: rc.ret_ast.Token,
	}.checkAssignType()
}

func (rc *retChecker) multi() {
	types := rc.f.RetType.Type.Tag.([]Type)
	n := len(rc.exp_model.values)
	if n == 1 {
		rc.checkMultiRetAsMutliRet()
		return
	} else if n > len(types) {
		rc.p.pusherrtok(rc.ret_ast.Token, "overflow_return")
	}
	for i, t := range types {
		if i >= n {
			break
		}
		assignChecker{
			p:      rc.p,
			t:      t,
			v:      rc.exp_model.values[i],
			errtok: rc.ret_ast.Token,
		}.checkAssignType()
	}
}

func (rc *retChecker) checkExprTypes() {
	if !rc.f.RetType.Type.MultiTyped { // Single return
		rc.single()
		return
	}
	// Multi return
	rc.multi()
}

func (rc *retChecker) checkMultiRetAsMutliRet() {
	v := rc.exp_model.values[0]
	if !v.data.Type.MultiTyped {
		rc.p.pusherrtok(rc.ret_ast.Token, "missing_multi_return")
		return
	}
	valTypes := v.data.Type.Tag.([]Type)
	retTypes := rc.f.RetType.Type.Tag.([]Type)
	if len(valTypes) < len(retTypes) {
		rc.p.pusherrtok(rc.ret_ast.Token, "missing_multi_return")
		return
	} else if len(valTypes) < len(retTypes) {
		rc.p.pusherrtok(rc.ret_ast.Token, "overflow_return")
		return
	}
	rc.exp_model.models = append(rc.exp_model.models, rc.exp_model.models[0])
	for i, rt := range retTypes {
		vt := valTypes[i]
		val := value{data: models.Data{Type: vt}}
		assignChecker{
			p:      rc.p,
			t:      rt,
			v:      val,
			errtok: rc.ret_ast.Token,
		}.checkAssignType()
	}
}

func (rc *retChecker) retsVars() {
	if !rc.f.RetType.Type.MultiTyped {
		for _, v := range rc.f.RetType.Identifiers {
			if !juleapi.IsIgnoreId(v.Kind) {
				model := new(exprModel)
				model.index = 0
				model.nodes = make([]exprBuildNode, 1)
				val, _ := rc.p.eval.single(v, model)
				rc.exp_model.models = append(rc.exp_model.models, model)
				rc.exp_model.values = append(rc.exp_model.values, val)
				break
			}
		}
		rc.ret_ast.Expr.Model = rc.exp_model
		return
	}
	types := rc.f.RetType.Type.Tag.([]Type)
	for i, v := range rc.f.RetType.Identifiers {
		if juleapi.IsIgnoreId(v.Kind) {
			node := exprNode{}
			node.value = types[i].String()
			node.value += juleapi.DefaultExpr
			rc.exp_model.models = append(rc.exp_model.models, node)
			continue
		}
		model := new(exprModel)
		model.index = 0
		model.nodes = make([]exprBuildNode, 1)
		val, _ := rc.p.eval.single(v, model)
		rc.exp_model.models = append(rc.exp_model.models, model)
		rc.exp_model.values = append(rc.exp_model.values, val)
	}
	rc.ret_ast.Expr.Model = rc.exp_model
}

func (rc *retChecker) check() {
	n := len(rc.ret_ast.Expr.Tokens)
	if n == 0 && !typeIsVoid(rc.f.RetType.Type) {
		if !rc.f.RetType.AnyVar() {
			rc.p.pusherrtok(rc.ret_ast.Token, "require_return_value")
		}
		rc.retsVars()
		return
	}
	if n > 0 && typeIsVoid(rc.f.RetType.Type) {
		rc.p.pusherrtok(rc.ret_ast.Token, "void_function_return_value")
	}
	rc.checkepxrs()
}
