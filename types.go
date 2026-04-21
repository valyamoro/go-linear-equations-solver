package main

import (
	"html/template"
	"strings"
)

type Equation struct {
	A, B, C  float64
	Original string
}

type Point struct{ X, Y float64 }

type stringsFlag []string

func (f *stringsFlag) String() string     { return strings.Join(*f, ", ") }
func (f *stringsFlag) Set(v string) error { *f = append(*f, v); return nil }

type PageData struct {
	TracesJSON template.JS
	LayoutJSON template.JS
	Equations  []EqInfo
	Solutions  []VertexInfo
	IsEmpty    bool
	XMin, XMax float64
	YMin, YMax float64
}

// EqInfo drives one row in the sidebar equation list.
type EqInfo struct {
	Label string
	Color string
}

// VertexInfo is one solution point.
type VertexInfo struct {
	X, Y float64
}
