package main

import (
	"fmt"
	"log"
	"rating_engine/pkg/ratingengine"
)

func main() {
	req := ratingengine.UpdateRatingRequest{
		Subject:       "A",
		CurrentRating: ratingengine.Rating(5.6),
		HistoryWeight: 10,
		Votes:         []ratingengine.Rating{7, 6, 6, 5, 7, 6, 6, 7, 10},
		ScoringMode:   string(ratingengine.ScoringModeMidhinge),
	}

	result, err := ratingengine.UpdateRating(req)
	if err != nil {
		log.Fatalf("rating update failed: %v", err)
	}

	fmt.Printf("Subject: %s\n", result.Subject)
	fmt.Printf("Old rating: %.1f\n", result.OldRating)
	fmt.Printf("Event score (midhinge): %.1f\n", result.EventScore)
	fmt.Printf("Effective history weight: %d\n", result.HistoryWeight)
	fmt.Printf("Votes count: %d\n", result.VotesCount)
	fmt.Printf("New rating: %.1f\n", result.NewRating)
}
