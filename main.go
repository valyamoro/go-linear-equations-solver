package main

import (
	"flag"
	"fmt"
	"html/template"
	"math"
	"os"
	"path/filepath"
)

var tmplFuncs = template.FuncMap{
	"fmtF": func(v float64) string {
		if v == math.Trunc(v) {
			return fmt.Sprintf("%.0f", v)
		}
		return fmt.Sprintf("%g", v)
	},
	"fmtPt": func(x, y float64) string {
		return fmt.Sprintf("(%g,  %g)", x, y)
	},
	"not": func(b bool) bool { return !b },
	"len": func(v []VertexInfo) int { return len(v) },
}

func renderHTML(pd PageData, outPath string) error {
	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return fmt.Errorf("cannot create directory: %w", err)
	}
	tmpl, err := template.New("page").Funcs(tmplFuncs).Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}
	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("cannot create %s: %w", outPath, err)
	}
	defer f.Close()
	return tmpl.Execute(f, pd)
}

func main() {
	var rawEqs stringsFlag
	var out string
	var xmin, xmax, ymin, ymax float64

	flag.Var(&rawEqs, "i", "Equation (repeatable). Example: -i 'x+y=4'")
	flag.StringVar(&out, "o", "tmp/equations.html", "Output HTML file")
	flag.Float64Var(&xmin, "xmin", -10, "Minimum X axis value")
	flag.Float64Var(&xmax, "xmax", 10, "Maximum X axis value")
	flag.Float64Var(&ymin, "ymin", -10, "Minimum Y axis value")
	flag.Float64Var(&ymax, "ymax", 10, "Maximum Y axis value")
	flag.Parse()

	if len(rawEqs) == 0 {
		rawEqs = []string{
			"x + y = 4",
			"x - y = 2",
			"2x + 2y = 8",
			"x + y - 2 = 2",
		}
		fmt.Println("Equations not specified — using built-in example.")
	}

	fmt.Println()
	fmt.Println("  Parsing equations:")
	eqs := make([]Equation, 0, len(rawEqs))
	for _, s := range rawEqs {
		eq, err := ParseEquation(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  error parsing %q: %v\n", s, err)
			os.Exit(1)
		}
		fmt.Printf("  %-30s  →  A=%-6g B=%-6g C=%-6g\n",
			s, eq.A, eq.B, eq.C)
		eqs = append(eqs, eq)
	}
	fmt.Println()

	pd, err := buildPageData(eqs, xmin, xmax, ymin, ymax)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  build error: %v\n", err)
		os.Exit(1)
	}

	if err := renderHTML(pd, out); err != nil {
		fmt.Fprintf(os.Stderr, "  render error: %v\n", err)
		os.Exit(1)
	}

	if pd.IsEmpty {
		fmt.Printf("  System incompatible — no solutions.\n")
	} else {
		fmt.Printf("  Found %d solutions:\n", len(pd.Solutions))
		for _, v := range pd.Solutions {
			fmt.Printf("      (%g, %g)\n", v.X, v.Y)
		}
	}
	fmt.Printf("\n  → Saved: %s\n\n", out)
}