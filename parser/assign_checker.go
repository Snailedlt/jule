package parser

import (
	"math"
	"strconv"

	"github.com/julelang/jule/lex"
	"github.com/julelang/jule/pkg/julebits"
	"github.com/julelang/jule/pkg/juletype"
)

func float_assignable(dt uint8, v value) bool {
	switch t := v.expr.(type) {
	case float64:
		v.data.Value = strconv.FormatFloat(t, 'e', -1, 64)
	case int64:
		v.data.Value = strconv.FormatFloat(float64(t), 'e', -1, 64)
	case uint64:
		v.data.Value = strconv.FormatFloat(float64(t), 'e', -1, 64)
	}
	return checkFloatBit(v.data, julebits.BitsizeType(dt))
}

func signedAssignable(dt uint8, v value) bool {
	min := juletype.MinOfType(dt)
	max := int64(juletype.MaxOfType(dt))
	switch t := v.expr.(type) {
	case float64:
		i, frac := math.Modf(t)
		if frac != 0 {
			return false
		}
		return i >= float64(min) && i <= float64(max)
	case uint64:
		if t <= uint64(max) {
			return true
		}
	case int64:
		return t >= min && t <= max
	}
	return false
}

func unsignedAssignable(dt uint8, v value) bool {
	max := juletype.MaxOfType(dt)
	switch t := v.expr.(type) {
	case float64:
		if t < 0 {
			return false
		}
		i, frac := math.Modf(t)
		if frac != 0 {
			return false
		}
		return i <= float64(max)
	case uint64:
		if t <= max {
			return true
		}
	case int64:
		if t < 0 {
			return false
		}
		return uint64(t) <= max
	}
	return false
}

func int_assignable(dt uint8, v value) bool {
	switch {
	case juletype.IsSignedInteger(dt):
		return signedAssignable(dt, v)
	case juletype.IsUnsignedInteger(dt):
		return unsignedAssignable(dt, v)
	}
	return false
}

type assign_checker struct {
	p                *Parser
	expr_t           Type
	v                value
	ignoreAny        bool
	not_allow_assign bool
	errtok           lex.Token
}

func (ac *assign_checker) has_error() bool {
	return ac.p.eval.has_error || ac.v.data.Value == ""
}

func (ac *assign_checker) check_validity() (valid bool) {
	valid = true
	if type_is_fn(ac.v.data.Type) {
		f := ac.v.data.Type.Tag.(*Func)
		if f.Receiver != nil {
			ac.p.pusherrtok(ac.errtok, "method_as_anonymous_fn")
			valid = false
		} else if len(f.Generics) > 0 {
			ac.p.pusherrtok(ac.errtok, "genericed_fn_as_anonymous_fn")
			valid = false
		}
	}
	return
}

func (ac *assign_checker) check_const() (ok bool) {
	if !ac.v.constExpr || !type_is_pure(ac.expr_t) ||
		!type_is_pure(ac.v.data.Type) || !juletype.IsNumeric(ac.v.data.Type.Id) {
		return
	}
	ok = true
	switch {
	case juletype.IsFloat(ac.expr_t.Id):
		if !float_assignable(ac.expr_t.Id, ac.v) {
			ac.p.pusherrtok(ac.errtok, "overflow_limits")
			ok = false
		}
	case juletype.IsInteger(ac.expr_t.Id):
		if !int_assignable(ac.expr_t.Id, ac.v) {
			ac.p.pusherrtok(ac.errtok, "overflow_limits")
			ok = false
		}
	default:
		ok = false
	}
	return
}

func (ac assign_checker) check() {
	if ac.has_error() {
		return
	} else if !ac.check_validity() {
		return
	} else if ac.check_const() {
		return
	}
	ac.p.check_type(ac.expr_t, ac.v.data.Type, ac.ignoreAny, !ac.not_allow_assign, ac.errtok)
}
