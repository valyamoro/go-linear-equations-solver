package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func ParseEquation(s string) (Equation, error) {
	s = strings.TrimSpace(s)
	i := strings.Index(s, "=")
	if i < 0 {
		return Equation{}, fmt.Errorf("no '=' found in %q", s)
	}
	lhsRaw := strings.TrimSpace(s[:i])
	rhsRaw := strings.TrimSpace(s[i+1:])

	rhs, err := strconv.ParseFloat(rhsRaw, 64)
	if err != nil {
		return Equation{}, fmt.Errorf("RHS %q is not numeric: %w", rhsRaw, err)
	}
	a, b, lhsConst, err := parseLinearExpr(lhsRaw)
	if err != nil {
		return Equation{}, fmt.Errorf("LHS %q: %w", lhsRaw, err)
	}
	return Equation{A: a, B: b, C: rhs - lhsConst, Original: s}, nil
}

func parseLinearExpr(expr string) (a, b, constant float64, err error) {
	expr = strings.ReplaceAll(expr, " ", "")
	if expr == "" {
		return 0, 0, 0, fmt.Errorf("empty expression")
	}
	if expr[0] != '+' && expr[0] != '-' {
		expr = "+" + expr
	}
	re := regexp.MustCompile(`[+-][^+-]+`)
	for _, tok := range re.FindAllString(expr, -1) {
		switch {
		case strings.ContainsRune(tok, 'x'):
			a += extractCoeff(strings.ReplaceAll(tok, "x", ""))
		case strings.ContainsRune(tok, 'y'):
			b += extractCoeff(strings.ReplaceAll(tok, "y", ""))
		default:
			v, e := strconv.ParseFloat(tok, 64)
			if e != nil {
				err = fmt.Errorf("unrecognised token %q", tok)
				return
			}
			constant += v
		}
	}
	return
}

func extractCoeff(s string) float64 {
	switch s {
	case "+", "":
		return 1
	case "-":
		return -1
	}
	v, _ := strconv.ParseFloat(s, 64)
	return v
}