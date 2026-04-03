package ratingengine

import (
	"math"
	"testing"
)

func TestIsRatingInRange_NaNAndInf(t *testing.T) {
	if isRatingInRange(Rating(math.NaN())) {
		t.Fatal("expected NaN to be invalid")
	}
	if isRatingInRange(Rating(math.Inf(1))) {
		t.Fatal("expected +Inf to be invalid")
	}
	if isRatingInRange(Rating(math.Inf(-1))) {
		t.Fatal("expected -Inf to be invalid")
	}
}

