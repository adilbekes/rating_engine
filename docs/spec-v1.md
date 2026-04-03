# Rating Engine Spec v1

## Input
- `subject` (string)
- `current_rating` (float)
- `history_weight` (int)
- `votes` (array of float)
- optional `scoring_mode` (string)

### Scoring mode
Supported values:
- `mean`
- `median`
- `trimmed_mean`
- `midhinge`

Default:
- `scoring_mode = median`

## Validation
- `subject` must not be empty
- `current_rating` must be within `[1.0, 10.0]`
- `history_weight` must be >= 0
- `votes` must not be empty
- every vote must be within `[1.0, 10.0]`
- `scoring_mode` must be one of: `mean`, `median`, `trimmed_mean`, `midhinge`

## Calculation
- `effective_history_weight = history_weight`
- `event_score` is computed based on `scoring_mode`:
  - `mean`: arithmetic average of votes
  - `median`: median of votes
  - `trimmed_mean`: trim 20% from each side then average remaining values; fallback to mean for small samples
  - `midhinge`: `(Q1 + Q3) / 2`
- `votes_count = len(votes)`
- `new_rating = (current_rating * effective_history_weight + event_score * votes_count) / (effective_history_weight + votes_count)`
- all rating outputs are normalized to 1 decimal place

## Errors
- `ErrInvalidSubjectName`
- `ErrInvalidCurrentRating`
- `ErrInvalidHistoryWeight`
- `ErrInvalidVotes`
- `ErrInvalidVoteValue`
- `ErrInvalidScoringMode`
