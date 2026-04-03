package ratingengine

import "errors"

var (
	ErrInvalidSubjectName   = errors.New("subject name must be provided")
	ErrInvalidCurrentRating = errors.New("current_rating must be within allowed range")
	ErrInvalidHistoryWeight = errors.New("history_weight must be greater than or equal to 0")
	ErrInvalidVotes         = errors.New("votes must not be empty")
	ErrInvalidVoteValue     = errors.New("vote value must be within allowed range")
	ErrInvalidScoringMode   = errors.New("scoring_mode must be one of: mean, median, trimmed_mean, midhinge")
)
