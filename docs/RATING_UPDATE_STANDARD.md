# RATING_UPDATE_STANDARD

This document defines the required behavior and coding rules for the `UpdateRating` pipeline in `pkg/ratingengine`.

## Scope and flow

`UpdateRating` must keep this order:
1. Validate request (`ValidateUpdateRatingRequest`)
2. Compute event score from votes (`EventScore`)
3. Blend historical rating with event score
4. Normalize all externally visible ratings to 1 decimal place

Validation errors must be returned as-is and stop processing immediately.

## Request and validation rules

- `subject` must be non-empty after trimming spaces.
- `history_weight` must be `>= 0`.
- `current_rating` must be within `[1.0, 10.0]` and not `NaN`/`Inf`.
- `votes` must be non-empty.
- Every vote must be within `[1.0, 10.0]` and not `NaN`/`Inf`.
- `scoring_mode` default is `median` when omitted.
- Allowed scoring modes: `mean`, `median`, `trimmed_mean`, `midhinge`.

## Scoring rules

- `mean`: arithmetic mean of all votes.
- `median`: median of sorted votes.
- `trimmed_mean`: trim 20% from each side (`n/5`); if trimming is not possible, fall back to mean.
- `midhinge`: `(Q1 + Q3) / 2` using Tukey-style halves.
- Event score and response ratings must be normalized with `math.Round(v*10)/10`.

## Response contract

A successful response must include:
- `subject`
- `old_rating` (normalized `current_rating`)
- `new_rating`
- `event_score`
- `history_weight` (effective value)
- `votes_count`

On failure, callers should return a response with only `error` populated at the transport boundary (`cmd/*`).

## Implementation constraints

- Keep `pkg/ratingengine` pure and deterministic (no I/O, no clock usage).
- Add new scoring modes by extending:
  1. mode constants,
  2. mode validator,
  3. `EventScore` switch,
  4. tests for valid/invalid mode behavior.
- Preserve compatibility of JSON field names in `UpdateRatingRequest` and `UpdateRatingResponse`.

