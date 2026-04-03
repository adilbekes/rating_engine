# INTEGRATION.md — Integration Notes for External Systems and Package Imports

## Library Integration (`pkg/ratingengine`)

Import path:

```go
import "rating_engine/pkg/ratingengine"
```

Minimal usage:

```go
req := ratingengine.UpdateRatingRequest{
    Subject:       "A",
    CurrentRating: 5.6,
    HistoryWeight: 10,
    Votes:         []ratingengine.Rating{7, 6, 6, 5, 7, 6, 6, 7, 10},
    ScoringMode:   string(ratingengine.ScoringModeMedian),
}

resp, err := ratingengine.UpdateRating(req)
```

Use only exported API (`UpdateRating`, exported request/response types, mode constants, and sentinel errors).

---

## Error Handling

`UpdateRating` returns validation errors from `pkg/ratingengine/errors.go`.
Callers should branch with `errors.Is`.

```go
if err != nil {
    switch {
    case errors.Is(err, ratingengine.ErrInvalidSubjectName):
    case errors.Is(err, ratingengine.ErrInvalidCurrentRating):
    case errors.Is(err, ratingengine.ErrInvalidHistoryWeight):
    case errors.Is(err, ratingengine.ErrInvalidVotes):
    case errors.Is(err, ratingengine.ErrInvalidVoteValue):
    case errors.Is(err, ratingengine.ErrInvalidScoringMode):
    }
}
```

---

## CLI Integration (`cmd/engine`)

`cmd/engine` provides a language-agnostic process boundary:

- Input: JSON request from `stdin`, `-d`, or `-f`.
- Output: JSON response to `stdout` or `-o` file.
- Success exit code: `0`.
- Error exit code: `1`.
- Error payload shape: `{\"error\":\"message\"}`.

Flags:

- `-d`: inline JSON request string
- `-f`: JSON request file path
- `-o`: optional output file path

`-d` and `-f` are mutually exclusive.

---

## Scoring-Mode Compatibility

External clients should treat scoring modes as stable string values:

- `mean`
- `median` (default)
- `trimmed_mean`
- `midhinge`

When extending modes, update validator, scoring switch, and tests together (see `docs/RATING_UPDATE_STANDARD.md`).
