package ui

import (
	"math"
	"strings"
	"testing"
)

func TestSparkPointsEdgeCases(t *testing.T) {
	if got := sparkPoints(nil, 120, 28); got != nil {
		t.Errorf("no values should draw nothing, got %v", got)
	}

	one := sparkPoints([]float64{42}, 120, 28)
	if len(one) != 1 || one[0][0] != 60 || one[0][1] != 14 {
		t.Errorf("a single value belongs on the middle line, got %v", one)
	}

	// A flat series must not divide by a zero range.
	flat := sparkPoints([]float64{7, 7, 7, 7}, 120, 28)
	for _, p := range flat {
		if p[1] != 14 {
			t.Fatalf("flat series should sit on the middle line, got %v", flat)
		}
	}
	if math.IsNaN(flat[0][1]) {
		t.Fatal("flat series produced NaN")
	}
}

func TestSparkPointsGeometry(t *testing.T) {
	const w, h = 120.0, 28.0
	values := []float64{1, 5, 3, 9, 2}
	pts := sparkPoints(values, w, h)

	if len(pts) != len(values) {
		t.Fatalf("point count %d, want %d", len(pts), len(values))
	}

	// Every point stays inside the padded box, so the 2px stroke never clips.
	for i, p := range pts {
		if p[0] < sparkPad-0.01 || p[0] > w-sparkPad+0.01 {
			t.Errorf("point %d x=%v outside [%d, %v]", i, p[0], sparkPad, w-sparkPad)
		}
		if p[1] < sparkPad-0.01 || p[1] > h-sparkPad+0.01 {
			t.Errorf("point %d y=%v outside [%d, %v]", i, p[1], sparkPad, h-sparkPad)
		}
	}

	// x spreads evenly across the full width.
	if pts[0][0] != sparkPad || math.Abs(pts[len(pts)-1][0]-(w-sparkPad)) > 0.01 {
		t.Errorf("series should span edge to edge, got %v .. %v", pts[0][0], pts[len(pts)-1][0])
	}

	// y is inverted: the largest value sits highest on screen (smallest y).
	maxIdx, minIdx := 3, 0 // values[3]=9 is the max, values[0]=1 the min
	if pts[maxIdx][1] >= pts[minIdx][1] {
		t.Errorf("a bigger value must sit higher: max y=%v, min y=%v", pts[maxIdx][1], pts[minIdx][1])
	}
	if math.Abs(pts[maxIdx][1]-sparkPad) > 0.01 {
		t.Errorf("the maximum should touch the top padding, got %v", pts[maxIdx][1])
	}
}

func TestSparklineIsAccessibleAndSelfContained(t *testing.T) {
	html := renderTier4(t, Sparkline(SparklineProps{
		Values: []float64{10, 14, 12, 21}, Label: "Revenu", Suffix: " €",
	}))

	for _, want := range []string{`role="img"`, "<title>", "aria-label", "Revenu"} {
		if !strings.Contains(html, want) {
			t.Errorf("missing %q in %s", want, html)
		}
	}
	// The summary must state the current value, not only the shape.
	if !strings.Contains(html, "actuellement 21 €") {
		t.Errorf("summary should name the latest value: %s", html)
	}
	// No chart library, no script, no remote asset.
	for _, forbidden := range []string{"<script", "http://", "https://"} {
		if strings.Contains(html, forbidden) {
			t.Errorf("sparkline must be self-contained, found %q: %s", forbidden, html)
		}
	}
}

func TestSparklineWithNoValuesRendersNothing(t *testing.T) {
	if html := renderTier4(t, Sparkline(SparklineProps{Label: "Vide"})); html != "" {
		t.Errorf("an empty series should render nothing, got %q", html)
	}
}
