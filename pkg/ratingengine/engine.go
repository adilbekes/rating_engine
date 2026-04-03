package ratingengine

func UpdateRating(req UpdateRatingRequest) (UpdateRatingResponse, error) {
	if err := ValidateUpdateRatingRequest(req); err != nil {
		return UpdateRatingResponse{}, err
	}

	votesCount := req.VotesCount()
	effectiveHistoryWeight := req.EffectiveHistoryWeight()
	eventScore := req.EventScore()

	denominator := float64(effectiveHistoryWeight + votesCount)
	newRating := normalizeRating(
		(float64(req.CurrentRating)*float64(effectiveHistoryWeight) + float64(eventScore)*float64(votesCount)) / denominator,
	)

	return UpdateRatingResponse{
		Subject:       req.Subject,
		OldRating:     normalizeRating(float64(req.CurrentRating)),
		NewRating:     newRating,
		EventScore:    eventScore,
		HistoryWeight: effectiveHistoryWeight,
		VotesCount:    votesCount,
	}, nil
}
