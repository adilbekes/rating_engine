# WORKFLOWS.md — Developer Workflows for Building, Running, and Testing

## Prerequisites

- Go toolchain installed (see `go.mod`, currently `go 1.25.5`).
- Working directory: `/Users/adilbek.es/workspace/rating_engine`.

---

## Build

Build the CLI binary:

```zsh
go build -o bin/engine ./cmd/engine/
```

---

## Run

Run from inline JSON (`-d`):

```zsh
./bin/engine -d '{"subject":"A","current_rating":5.6,"history_weight":10,"votes":[7,6,6,5,7,6,6,7,10],"scoring_mode":"midhinge"}'
```

Run from input file (`-f`) and write to output file (`-o`):

```zsh
./bin/engine -f example.json -o output.json
```

Run demo program:

```zsh
go run ./cmd/demo/
```

---

## Test

Run all package tests:

```zsh
go test ./pkg/ratingengine -count=1
```

Run with verbose output:

```zsh
go test ./pkg/ratingengine -count=1 -v
```

Run a focused test:

```zsh
go test ./pkg/ratingengine -count=1 -run TestUpdateRating
```

---

## Quality Check

Recommended pre-commit sequence:

```zsh
go build ./cmd/engine/ && \
go test ./pkg/ratingengine -count=1 && \
go vet ./...
```
