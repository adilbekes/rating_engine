package ratingengine

import "testing"

func TestErrorMessages_AreNotEmpty(t *testing.T) {
	errs := []error{
		ErrInvalidSubjectName,
		ErrInvalidCurrentRating,
		ErrInvalidHistoryWeight,
		ErrInvalidVotes,
		ErrInvalidVoteValue,
		ErrInvalidScoringMode,
	}

	for _, err := range errs {
		if err == nil || err.Error() == "" {
			t.Fatalf("expected non-empty sentinel error, got: %v", err)
		}
	}
}

