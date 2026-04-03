package ratingengine

import (
	"math"
	"strings"
)

func ValidateUpdateRatingRequest(req UpdateRatingRequest) error {
	if !isValidScoringMode(req.EffectiveScoringMode()) {
		return ErrInvalidScoringMode
	}

	if strings.TrimSpace(req.Subject) == "" {
		return ErrInvalidSubjectName
	}

	if req.HistoryWeight < 0 {
		return ErrInvalidHistoryWeight
	}

	if !isRatingInRange(req.CurrentRating) {
		return ErrInvalidCurrentRating
	}

	if len(req.Votes) == 0 {
		return ErrInvalidVotes
	}

	for _, rating := range req.Votes {
		if !isRatingInRange(rating) {
			return ErrInvalidVoteValue
		}
	}

	return nil
}

func isRatingInRange(value Rating) bool {
	v := float64(value)
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return false
	}
	return value >= defaultMinRatingValue && value <= defaultMaxRatingValue
}

func isValidScoringMode(mode ScoringMode) bool {
	switch mode {
	case ScoringModeMean, ScoringModeMedian, ScoringModeTrimmedMean, ScoringModeMidhinge:
		return true
	default:
		return false
	}
}
