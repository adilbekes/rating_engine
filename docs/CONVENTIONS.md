# CONVENTIONS.md — Project-Specific Coding Conventions and Constraints

## Validation and Error Rules

- Always validate through `ValidateUpdateRatingRequest` before computing results.
- Return sentinel errors from `pkg/ratingengine/errors.go` (`ErrInvalidSubjectName`, `ErrInvalidScoringMode`, etc.).
- Keep validator behavior aligned with `docs/RATING_UPDATE_STANDARD.md`.
- In library code, return errors; do not panic.

---

## Rating and Scoring Conventions

- Rating domain is `[1.0, 10.0]` for `current_rating` and each vote.
- Reject `NaN` and `Inf` values.
- Normalize externally visible ratings to one decimal place via `math.Round(v*10)/10`.
- Default `scoring_mode` to `median` when omitted.
- Allowed modes: `mean`, `median`, `trimmed_mean`, `midhinge`.

---

## JSON and Type Conventions

- Use `snake_case` JSON field names (`current_rating`, `history_weight`, `votes_count`).
- Preserve existing JSON schema in `UpdateRatingRequest` and `UpdateRatingResponse`.
- `Rating` JSON marshaling must keep one-decimal representation.
- `cmd/*` should emit `{\"error\":\"...\"}` at transport boundary on failures.

---

## Package Boundaries

- `pkg/ratingengine` contains business logic only; keep it pure and deterministic.
- `cmd/engine` handles CLI flags, JSON decode/encode, and process exit codes.
- `cmd/demo` is example usage only; do not move production logic there.

---

## Testing Conventions

- Prefer table-driven tests for mode and validation matrices.
- Keep tests near implementation (`pkg/ratingengine/*_test.go`).
- Cover both valid and invalid scoring modes when adding/extending modes.
- For error checks, use `errors.Is` against sentinel errors.
