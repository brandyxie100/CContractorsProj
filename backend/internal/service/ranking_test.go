package service

import (
	"testing"
	"time"
)

func TestIntervalsOverlap(t *testing.T) {
	base := time.Date(2026, 8, 12, 8, 0, 0, 0, time.UTC)
	cases := []struct {
		name string
		a0, a1, b0, b1 time.Time
		want bool
	}{
		{"identical", base, base.Add(2 * time.Hour), base, base.Add(2 * time.Hour), true},
		{"adjacent", base, base.Add(2 * time.Hour), base.Add(2 * time.Hour), base.Add(4 * time.Hour), false},
		{"partial", base, base.Add(3 * time.Hour), base.Add(2 * time.Hour), base.Add(5 * time.Hour), true},
		{"disjoint", base, base.Add(time.Hour), base.Add(2 * time.Hour), base.Add(3 * time.Hour), false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := IntervalsOverlap(tc.a0, tc.a1, tc.b0, tc.b1)
			if got != tc.want {
				t.Fatalf("got %v want %v", got, tc.want)
			}
		})
	}
}

func TestAvailabilityScore_Available(t *testing.T) {
	start := time.Now()
	if s := AvailabilityScore(true, start, start); s != 1.0 {
		t.Fatalf("expected 1.0 got %v", s)
	}
}

func TestProximityScore_CloserIsHigher(t *testing.T) {
	near := ProximityScore(5)
	far := ProximityScore(50)
	if near <= far {
		t.Fatalf("near %v should be > far %v", near, far)
	}
}

func TestRankScore_Weights(t *testing.T) {
	got := RankScore(1, 1, 1)
	if got < 0.99 || got > 1.01 {
		t.Fatalf("expected ~1 got %v", got)
	}
}

func TestSpecsScore_UnderSpec(t *testing.T) {
	w := 15.0
	min := 20.0
	score, under := SpecsScore(&w, &min)
	if !under || score >= 1.0 {
		t.Fatalf("expected under-spec partial score, got score=%v under=%v", score, under)
	}
}

func TestHaversineKm_SamePoint(t *testing.T) {
	d := HaversineKm(-36.85, 174.76, -36.85, 174.76)
	if d > 0.01 {
		t.Fatalf("expected ~0 got %v", d)
	}
}
