package expr

import (
	"context"
	"errors"
    "fmt"
	"strconv"
	"text/scanner"

	"github.com/PaesslerAG/gval"
	"github.com/generikvault/gvalstrings"

    "github.com/please-build/cc-rules/tools/please_cc/cctool"
)

func parseVersion(c context.Context, p *gval.Parser) (gval.Evaluable, error) {
	n, err := strconv.ParseInt(p.TokenText(), 10, 64)
	if err != nil {
		return nil, p.Expected("version number,", scanner.Int)
	}
    ver := cctool.Version{n}
	for {
		scan := p.Scan()
		switch scan {
		case '.':
			scan = p.Scan()
			switch scan {
			case scanner.Int:
				n, err := strconv.ParseInt(p.TokenText(), 10, 64)
				if err != nil {
					return nil, p.Expected("version number,", scanner.Int)
				}
				ver = append(ver, n)
			default:
				return nil, p.Expected("version number,", scanner.Int)
			}
		default:
			p.Camouflage("version number", '.')
			return p.Const(ver), nil
		}
	}
}

func parseDoubleQuotedString(c context.Context, p *gval.Parser) (gval.Evaluable, error) {
	s, err := unquoteString(p.TokenText(), false)
	if err != nil {
		return nil, err
	}
	return p.Const(s), nil
}

func parseSingleQuotedString(c context.Context, p *gval.Parser) (gval.Evaluable, error) {
	s, err := unquoteString(p.TokenText(), true)
	if err != nil {
		return nil, err
	}
	return p.Const(s), nil
}

func parseArray(c context.Context, p *gval.Parser) (gval.Evaluable, error) {
	strs := []string{}
	for {
		scanned := p.Scan()
		switch {
		case scanned == ',':
		case scanned == ']':
			return p.Const(strs), nil
		case scanned == scanner.Char || scanned == scanner.String:
			s, err := unquoteString(p.TokenText(), scanned == scanner.Char)
			if err != nil {
				return nil, err
			}
			strs = append(strs, s)
		default:
            return nil, p.Expected("array,", scanner.String)
		}
	}
}

func parseIdent(c context.Context, p *gval.Parser) (gval.Evaluable, error) {
    return p.Var(p.Const(p.TokenText())), nil
}

func parseParentheses(c context.Context, p *gval.Parser) (gval.Evaluable, error) {
	eval, err := p.ParseExpression(c)
	if err != nil {
		return nil, err
	}
	switch p.Scan() {
	case ')':
		return eval, nil
	default:
		return nil, p.Expected("parentheses,", ')')
	}
}

func parseIf(c context.Context, p *gval.Parser, e gval.Evaluable) (gval.Evaluable, error) {
	t, err := p.ParseExpression(c)
	if err != nil {
		return nil, err
	}
	f := p.Const(nil)
	fNil := true
	switch p.Scan() {
	case ':':
		f, err = p.ParseExpression(c)
		if err != nil {
			return nil, err
		}
		fNil = false
	case scanner.EOF:
		// Undefined false branch, which implicitly evaluates to the empty string array - see the comment below.
	default:
		return nil, p.Expected("ternary operator,", ':', scanner.EOF)
	}
	return func(c context.Context, v any) (any, error) {
		cond, err := e(c, v)
		if err != nil {
			return nil, err
		}
		if condBool, condIsBool := cond.(bool); cond == nil || (condIsBool && !condBool) {
			// If the condition evaluates to false and the false branch is undefined, the expression evaluates to an empty string
			// array. This allows "<cond> ? <expr>" to evaluate to a legal value when <cond> is false, which is preferable to the
			// more strictly correct but uglier "<cond> ? <expr> : []".
			if fNil {
				return []string{}, nil
			}
			return f(c, v)
		}
		return t(c, v)
	}, nil
}

func unquoteString(quoted string, singleQuoted bool) (string, error) {
	var (
		s   string
		err error
	)
	if singleQuoted {
		s, err = gvalstrings.UnquoteSingleQuoted(quoted)
	} else {
	    s, err = strconv.Unquote(quoted)
	}
    if err != nil {
        return "", fmt.Errorf("could not parse string: %w", err)
    }
	if s == "" {
		return "", errors.New("string literal cannot be empty")
	}
	return s, nil
}
