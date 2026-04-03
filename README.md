# rating_engine

Standalone Go CLI microservice for updating a subject rating after an event.

## Build

```bash
go build -o bin/rating_engine ./cmd/engine
```

## Usage

Input sources:
- stdin
- `-d` inline JSON string
- `-f` JSON file

Output targets:
- stdout
- `-o` JSON file

Exit codes:
- `0` success
- `1` error

### Request shape

```json
{
  "subject": "A",
  "current_rating": 5.6,
  "history_weight": 10,
  "votes": [7, 6, 6, 5, 7, 6, 6, 7, 10],
  "scoring_mode": "median"
}
```

`scoring_mode` is optional.

Supported values:
- `mean`
- `median` (default)
- `trimmed_mean`
- `midhinge`

### Example: stdin

```bash
echo '{
  "subject": "A",
  "current_rating": 5.6,
  "history_weight": 10,
  "votes": [7, 6, 6, 5, 7, 6, 6, 7, 10]
}' | ./bin/rating_engine
```

### Example: choose mode

```bash
echo '{
  "subject": "A",
  "current_rating": 5.6,
  "history_weight": 10,
  "votes": [7, 6, 6, 5, 7, 6, 6, 7, 10],
  "scoring_mode": "midhinge"
}' | ./bin/rating_engine
```

### File input/output

```bash
./bin/rating_engine -f request.json -o response.json
```

## Scoring modes

- `mean`: arithmetic average of all `votes`
- `median`: middle value of sorted `votes` (or average of 2 middle values)
- `trimmed_mean`: trims 20% from each side, then averages remaining values; for small samples it falls back to mean
- `midhinge`: `(Q1 + Q3) / 2` using Tukey's hinges

Reference: https://en.wikipedia.org/wiki/Midhinge

### Mode comparison table

Using:
- `current_rating = 5.6`
- `history_weight = 10`
- `votes = [1, 2, 3, 6, 6, 6, 7, 8, 9, 10]`

Sorted votes:
- `[1, 2, 3, 6, 6, 6, 7, 8, 9, 10]`

| `scoring_mode`     | `event_score` formula (on votes above)                  | `event_score` | `new_rating` |
|--------------------|---------------------------------------------------------|--------------:|-------------:|
| `mean`             | average of sum: `(sum(votes) / 10) = 58 / 10`           |         `5.8` |        `5.7` |
| `median` (default) | average of 2 middle values: `(6 + 6) / 2`               |         `6.0` |        `5.8` |
| `trimmed_mean`     | trim 20% each side: remove `1, 2, 9, 10`, then `36 / 6` |         `6.0` |        `5.8` |
| `midhinge`         | `(Q1 + Q3) / 2 = (3.0 + 8.0) / 2`                       |         `5.5` |        `5.6` |

## Validation

- `subject` must be provided
- `current_rating` must be in `[1.0, 10.0]`
- `history_weight >= 0`
- `votes` must not be empty
- each `vote` must be in `[1.0, 10.0]`
- `scoring_mode` must be one of: `mean`, `median`, `trimmed_mean`, `midhinge`

## Test

```bash
go test ./...
```

