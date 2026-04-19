package main

import (
	"encoding/json"
	"fmt"
	"html/template"
)

func buildPageData(eqs []Equation, xmin, xmax, ymin, ymax float64) (PageData, error) {
	solutions, hasSolutions := solveSystem(eqs)
	pd := PageData{XMin: xmin, XMax: xmax, YMin: ymin, YMax: ymax, IsEmpty: !hasSolutions}

	type m = map[string]interface{}
	var traces []m

	for i, eq := range eqs {
		col := palette[i%len(palette)]
		xs, ys := boundarySegment(eq, xmin, xmax, ymin, ymax)
		if xs == nil {
			continue
		}
		label := prettyEq(eq.Original)
		traces = append(traces, m{
			"type": "scatter",
			"x":    xs,
			"y":    ys,
			"mode": "lines",
			"name": label,
			"line": m{"color": col, "width": 2.8},
		})
		pd.Equations = append(pd.Equations, EqInfo{
			Label: eq.Original,
			Color: col,
		})
	}

	if hasSolutions {
		var vx, vy []float64
		var vlabels []string
		for _, p := range solutions {
			vx = append(vx, r3(p.X))
			vy = append(vy, r3(p.Y))
			vlabels = append(vlabels, fmt.Sprintf("(%g, %g)", r3(p.X), r3(p.Y)))
			pd.Solutions = append(pd.Solutions, VertexInfo{r3(p.X), r3(p.Y)})
		}
		traces = append(traces, m{
			"type":         "scatter",
			"x":            vx,
			"y":            vy,
			"mode":         "markers+text",
			"text":         vlabels,
			"textposition": "top right",
			"textfont":     m{"color": "#00e6a0", "size": 11, "family": "Fira Code, monospace"},
			"marker":       m{"color": "#00e6a0", "size": 9, "line": m{"color": "#0a1628", "width": 2}},
			"name":         "Solutions",
			"hoverinfo":    "text",
		})
	}

	layout := m{
		"paper_bgcolor": "#06101e",
		"plot_bgcolor":  "#06101e",
		"font":          m{"color": "#8ba8cc", "family": "'Space Mono', 'Courier New', monospace", "size": 12},
		"xaxis": m{
			"range":          []float64{xmin, xmax},
			"gridcolor":      "rgba(100,140,200,0.07)",
			"gridwidth":      1,
			"zerolinecolor":  "rgba(100,180,255,0.25)",
			"zerolinewidth":  1.5,
			"tickfont":       m{"color": "#4a6a9a", "size": 11},
			"title":          m{"text": "x", "font": m{"color": "#6a9acc", "size": 14}},
			"showspikes":     true,
			"spikecolor":     "rgba(0,230,160,0.3)",
			"spikethickness": 1,
			"spikedash":      "dot",
			"linecolor":      "rgba(100,140,200,0.15)",
		},
		"yaxis": m{
			"range":          []float64{ymin, ymax},
			"gridcolor":      "rgba(100,140,200,0.07)",
			"gridwidth":      1,
			"zerolinecolor":  "rgba(100,180,255,0.25)",
			"zerolinewidth":  1.5,
			"tickfont":       m{"color": "#4a6a9a", "size": 11},
			"title":          m{"text": "y", "font": m{"color": "#6a9acc", "size": 14}},
			"showspikes":     true,
			"spikecolor":     "rgba(0,230,160,0.3)",
			"spikethickness": 1,
			"spikedash":      "dot",
			"linecolor":      "rgba(100,140,200,0.15)",
		},
		"legend": m{
			"bgcolor":     "rgba(8,18,38,0.85)",
			"bordercolor": "rgba(100,150,220,0.15)",
			"borderwidth": 1,
			"font":        m{"color": "#8ba8cc", "size": 11, "family": "'Space Mono', monospace"},
			"x": 0.01, "y": 0.99,
			"xanchor": "left", "yanchor": "top",
		},
		"margin":    m{"l": 55, "r": 25, "t": 25, "b": 55},
		"hovermode": "closest",
		"dragmode":  "zoom",
		"modebar": m{
			"bgcolor":     "rgba(8,18,38,0)",
			"color":       "#4a6a9a",
			"activecolor": "#00e6a0",
		},
	}

	tj, err := json.Marshal(traces)
	if err != nil {
		return PageData{}, err
	}
	lj, err := json.Marshal(layout)
	if err != nil {
		return PageData{}, err
	}
	pd.TracesJSON = template.JS(tj)
	pd.LayoutJSON = template.JS(lj)
	return pd, nil
}