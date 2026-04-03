package ratingengine

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestRatingMarshalUnmarshalNormalizesToOneDecimal(t *testing.T) {
	var r Rating
	if err := json.Unmarshal([]byte("5.66"), &r); err != nil {
		t.Fatalf("unmarshal rating: %v", err)
	}
	if got, want := float64(r), 5.7; got != want {
		t.Fatalf("rating mismatch: got %.1f want %.1f", got, want)
	}

	payload, err := json.Marshal(&r)
	if err != nil {
		t.Fatalf("marshal rating: %v", err)
	}
	if got, want := string(payload), "5.7"; got != want {
		t.Fatalf("marshal mismatch: got %q want %q", got, want)
	}
}

func TestEventScore_ByMode(t *testing.T) {
	votes := []Rating{7, 6, 6, 5, 7, 6, 6, 7, 10}

	tests := []struct {
		name string
		mode string
		want Rating
	}{
		{name: "default mode is median", mode: "", want: 6.0},
		{name: "mean", mode: string(ScoringModeMean), want: 6.7},
		{name: "median", mode: string(ScoringModeMedian), want: 6.0},
		{name: "trimmed mean (20%)", mode: string(ScoringModeTrimmedMean), want: 6.4},
		{name: "midhinge", mode: string(ScoringModeMidhinge), want: 6.5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := UpdateRatingRequest{Votes: votes, ScoringMode: tc.mode}
			if got := req.EventScore(); got != tc.want {
				t.Fatalf("event score mismatch: got %.1f want %.1f", got, tc.want)
			}
		})
	}
}

func TestUpdateRating_SpecExample_DefaultMedian(t *testing.T) {
	// votes [7,6,6,5,7,6,6,7,1]: median=6.0 (default mode)
	// new_rating = (5.6*18 + 6.0*9) / 27 = 154.8/27 = 5.733… → 5.7
	req := UpdateRatingRequest{
		Subject:       "A",
		CurrentRating: 5.6,
		HistoryWeight: 18,
		Votes:         []Rating{7, 6, 6, 5, 7, 6, 6, 7, 1},
	}

	resp, err := UpdateRating(req)
	if err != nil {
		t.Fatalf("update rating: %v", err)
	}

	if resp.EventScore != 6.0 {
		t.Fatalf("event score mismatch: got %.1f want 6.0", resp.EventScore)
	}
	if resp.NewRating != 5.7 {
		t.Fatalf("new rating mismatch: got %.1f want 5.7", resp.NewRating)
	}
	if resp.HistoryWeight != 18 {
		t.Fatalf("history weight mismatch: got %d want 18", resp.HistoryWeight)
	}
	if resp.VotesCount != 9 {
		t.Fatalf("votes count mismatch: got %d want 9", resp.VotesCount)
	}
}

func TestUpdateRating_HighHistoryWeight(t *testing.T) {
	req := UpdateRatingRequest{
		Subject:       "x",
		CurrentRating: 9.9,
		HistoryWeight: 100,
		Votes:         []Rating{10, 10, 10},
	}

	resp, err := UpdateRating(req)
	if err != nil {
		t.Fatalf("update rating: %v", err)
	}

	if got, want := resp.HistoryWeight, 100; got != want {
		t.Fatalf("history weight mismatch: got %d want %d", got, want)
	}
	// (9.9*100 + 10.0*3) / 103 = 1020/103 ≈ 9.902 → 9.9
	if got, want := resp.NewRating, Rating(9.9); got != want {
		t.Fatalf("new rating mismatch: got %.1f want %.1f", got, want)
	}
}

func TestUpdateRating_MidhingeMode(t *testing.T) {
	req := UpdateRatingRequest{
		Subject:       "A",
		CurrentRating: 5.6,
		HistoryWeight: 10,
		Votes:         []Rating{7, 6, 6, 5, 7, 6, 6, 7, 10},
		ScoringMode:   string(ScoringModeMidhinge),
	}

	resp, err := UpdateRating(req)
	if err != nil {
		t.Fatalf("update rating: %v", err)
	}

	if resp.EventScore != 6.5 {
		t.Fatalf("event score mismatch: got %.1f want 6.5", resp.EventScore)
	}
	if resp.NewRating != 6.0 {
		t.Fatalf("new rating mismatch: got %.1f want 6.0", resp.NewRating)
	}
}

func TestUpdateRating_Validation(t *testing.T) {
	base := UpdateRatingRequest{
		Subject:       "a",
		CurrentRating: 5.0,
		HistoryWeight: 1,
		Votes:         []Rating{5.0},
	}

	tests := []struct {
		name string
		req  UpdateRatingRequest
		want error
	}{
		{name: "missing subject name", req: UpdateRatingRequest{CurrentRating: 5.0, Votes: []Rating{5.0}}, want: ErrInvalidSubjectName},
		{name: "negative history weight", req: UpdateRatingRequest{Subject: "a", CurrentRating: 5.0, HistoryWeight: -1, Votes: []Rating{5.0}}, want: ErrInvalidHistoryWeight},
		{name: "empty votes", req: UpdateRatingRequest{Subject: "a", CurrentRating: 5.0, HistoryWeight: 1}, want: ErrInvalidVotes},
		{name: "vote out of range", req: UpdateRatingRequest{Subject: "a", CurrentRating: 5.0, HistoryWeight: 1, Votes: []Rating{11}}, want: ErrInvalidVoteValue},
		{name: "current out of range", req: UpdateRatingRequest{Subject: "a", CurrentRating: 10.1, HistoryWeight: 1, Votes: []Rating{5.0}}, want: ErrInvalidCurrentRating},
		{name: "invalid scoring mode", req: UpdateRatingRequest{Subject: "a", CurrentRating: 5.0, HistoryWeight: 1, Votes: []Rating{5.0}, ScoringMode: "p50"}, want: ErrInvalidScoringMode},
		{name: "valid baseline", req: base, want: nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := UpdateRating(tc.req)
			if tc.want == nil {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}

			if !errors.Is(err, tc.want) {
				t.Fatalf("expected %v, got %v", tc.want, err)
			}
		})
	}
}
