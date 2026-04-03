# ARCHITECTURE.md — System Architecture, Component Boundaries, and Data Flow

## Purpose

Given an `UpdateRatingRequest`, produce a validated and normalized `UpdateRatingResponse` with a blended rating.

---

## Package Layout

| Path | Role |
|---|---|
| `pkg/ratingengine/` | Core library (validation, scoring, rating update pipeline) |
| `cmd/engine/main.go` | CLI entry point (`JSON in -> JSON out`) |
| `cmd/demo/main.go` | Minimal in-process usage example |
| `docs/RATING_UPDATE_STANDARD.md` | Normative behavior contract for the update pipeline |

---

## Core Types (`pkg/ratingengine/types.go`)

### `UpdateRatingRequest`

| Field | Type | Required | Notes |
|---|---|---|---|
| `subject` | `string` | yes | Non-empty after trim |
| `current_rating` | `Rating` | yes | Must be in `[1.0, 10.0]` |
| `history_weight` | `int` | yes | Must be `>= 0` |
| `votes` | `[]Rating` | yes | Non-empty; each vote in `[1.0, 10.0]` |
| `scoring_mode` | `string` | no | Defaults to `median` |

### `UpdateRatingResponse`

| Field | Type | Notes |
|---|---|---|
| `subject` | `string` | Echo from request |
| `old_rating` | `Rating` | Normalized to 1 decimal |
| `new_rating` | `Rating` | Final blended rating |
| `event_score` | `Rating` | Score computed from votes |
| `history_weight` | `int` | Effective history weight |
| `votes_count` | `int` | Count of input votes |
| `error` | `string` | Transport-level error payload (`cmd/*`) |

---

## Data Flow (`UpdateRating`)

```
UpdateRatingRequest
      |
      v
1) ValidateUpdateRatingRequest
   - scoring mode
   - subject
   - history_weight
   - current_rating
   - votes
      |
      v
2) Compute event score via selected mode
   - mean / median / trimmed_mean / midhinge
      |
      v
3) Blend historical and event components
   new = (current*history_weight + event_score*votes_count) /
         (history_weight + votes_count)
      |
      v
4) Normalize ratings to 1 decimal place
      |
      v
UpdateRatingResponse
```

---

## Component Boundaries

| File | Responsibility |
|---|---|
| `pkg/ratingengine/engine.go` | Orchestrates validation, scoring, blending, response assembly |
| `pkg/ratingengine/validator.go` | Input validation and mode/range checks |
| `pkg/ratingengine/types.go` | Domain types, scoring mode constants, scoring helpers |
| `pkg/ratingengine/errors.go` | Sentinel validation errors |
| `cmd/engine/main.go` | CLI argument parsing, JSON I/O, exit code contract |

Library code in `pkg/ratingengine` stays deterministic (no file/network/clock I/O).
