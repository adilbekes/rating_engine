package ratingengine

import (
	"encoding/json"
	"fmt"
	"math"
	"sort"
)

type Rating float64

type ScoringMode string

const (
	ScoringModeMean        ScoringMode = "mean"
	ScoringModeMedian      ScoringMode = "median"
	ScoringModeTrimmedMean ScoringMode = "trimmed_mean"
	ScoringModeMidhinge    ScoringMode = "midhinge"
)

const (
	defaultMinRatingValue = 1.0
	defaultMaxRatingValue = 10.0
)

func normalizeRating(v float64) Rating {
	return Rating(math.Round(v*10) / 10)
}

// MarshalJSON uses a pointer receiver so all methods on Rating
// consistently use pointer receivers.
func (r *Rating) MarshalJSON() ([]byte, error) {
	if r == nil {
		return []byte("0.0"), nil
	}
	return []byte(fmt.Sprintf("%.1f", normalizeRating(float64(*r)))), nil
}

func (r *Rating) UnmarshalJSON(data []byte) error {
	var value float64
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*r = normalizeRating(value)
	return nil
}

type UpdateRatingRequest struct {
	Subject       string   `json:"subject"`
	CurrentRating Rating   `json:"current_rating"`
	HistoryWeight int      `json:"history_weight"`
	Votes         []Rating `json:"votes"`
	ScoringMode   string   `json:"scoring_mode,omitempty"`
}

type UpdateRatingResponse struct {
	Subject       string `json:"subject,omitempty"`
	OldRating     Rating `json:"old_rating,omitempty"`
	NewRating     Rating `json:"new_rating,omitempty"`
	EventScore    Rating `json:"event_score,omitempty"`
	HistoryWeight int    `json:"history_weight,omitempty"`
	VotesCount    int    `json:"votes_count,omitempty"`
	Error         string `json:"error,omitempty"`
}

// ---- Helpers ----

func (req UpdateRatingRequest) VotesCount() int {
	return len(req.Votes)
}

func (req UpdateRatingRequest) EffectiveHistoryWeight() int {
	return req.HistoryWeight
}

func (req UpdateRatingRequest) EffectiveScoringMode() ScoringMode {
	if req.ScoringMode == "" {
		return ScoringModeMedian
	}

	return ScoringMode(req.ScoringMode)
}

func (req UpdateRatingRequest) sortedVotes() []Rating {
	sorted := make([]Rating, len(req.Votes))
	copy(sorted, req.Votes)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	return sorted
}

func (req UpdateRatingRequest) mean() float64 {
	n := len(req.Votes)
	if n == 0 {
		return 0
	}

	total := 0.0
	for _, vote := range req.Votes {
		total += float64(vote)
	}

	return total / float64(n)
}

func (req UpdateRatingRequest) median() float64 {
	if len(req.Votes) == 0 {
		return 0
	}

	return medianOf(req.sortedVotes())
}

func (req UpdateRatingRequest) trimmedMean() float64 {
	n := len(req.Votes)
	if n == 0 {
		return 0
	}

	// 20% trim on each side; for small samples this gracefully falls back to mean.
	trimCount := n / 5
	sorted := req.sortedVotes()
	if trimCount == 0 || trimCount*2 >= n {
		return req.mean()
	}

	trimmed := sorted[trimCount : n-trimCount]
	total := 0.0
	for _, vote := range trimmed {
		total += float64(vote)
	}

	return total / float64(len(trimmed))
}

// medianOf returns the median of an already-sorted slice.
func medianOf(s []Rating) float64 {
	n := len(s)
	if n == 0 {
		return 0
	}
	mid := n / 2
	if n%2 == 1 {
		return float64(s[mid])
	}
	return (float64(s[mid-1]) + float64(s[mid])) / 2
}

// midhinge returns the raw (Q1+Q3)/2 using Tukey's hinges.
func (req UpdateRatingRequest) midhinge() float64 {
	n := len(req.Votes)
	if n == 0 {
		return 0
	}
	if n == 1 {
		return float64(req.Votes[0])
	}
	sorted := req.sortedVotes()

	half := n / 2
	q1 := medianOf(sorted[:half])
	q3 := medianOf(sorted[n-half:])
	return (q1 + q3) / 2
}

// EventScore computes score from selected mode.
func (req UpdateRatingRequest) EventScore() Rating {
	switch req.EffectiveScoringMode() {
	case ScoringModeMean:
		return normalizeRating(req.mean())
	case ScoringModeTrimmedMean:
		return normalizeRating(req.trimmedMean())
	case ScoringModeMidhinge:
		return normalizeRating(req.midhinge())
	default:
		return normalizeRating(req.median())
	}
}
