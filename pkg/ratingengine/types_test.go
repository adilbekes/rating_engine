package ratingengine

import (
	"encoding/json"
	"testing"
)

func TestRatingMarshalJSON_NilReceiver(t *testing.T) {
	var r *Rating

	payload, err := r.MarshalJSON()
	if err != nil {
		t.Fatalf("marshal nil rating: %v", err)
	}

	if got, want := string(payload), "0.0"; got != want {
		t.Fatalf("marshal mismatch: got %q want %q", got, want)
	}
}

func TestRatingUnmarshalJSON_InvalidValue(t *testing.T) {
	var r Rating
	if err := r.UnmarshalJSON([]byte(`"bad"`)); err == nil {
		t.Fatal("expected unmarshal error, got nil")
	}
}

func TestRequestScoringHelpers_EdgeCases(t *testing.T) {
	tests := []struct {
		name string
		req  UpdateRatingRequest
		fn   func(UpdateRatingRequest) float64
		want float64
	}{
		{
			name: "mean empty votes",
			req:  UpdateRatingRequest{},
			fn:   func(r UpdateRatingRequest) float64 { return r.mean() },
			want: 0,
		},
		{
			name: "median empty votes",
			req:  UpdateRatingRequest{},
			fn:   func(r UpdateRatingRequest) float64 { return r.median() },
			want: 0,
		},
		{
			name: "trimmed mean empty votes",
			req:  UpdateRatingRequest{},
			fn:   func(r UpdateRatingRequest) float64 { return r.trimmedMean() },
			want: 0,
		},
		{
			name: "trimmed mean fallback to mean for small sample",
			req:  UpdateRatingRequest{Votes: []Rating{4, 8, 10}},
			fn:   func(r UpdateRatingRequest) float64 { return r.trimmedMean() },
			want: (4 + 8 + 10) / 3.0,
		},
		{
			name: "midhinge empty votes",
			req:  UpdateRatingRequest{},
			fn:   func(r UpdateRatingRequest) float64 { return r.midhinge() },
			want: 0,
		},
		{
			name: "midhinge single vote",
			req:  UpdateRatingRequest{Votes: []Rating{7.2}},
			fn:   func(r UpdateRatingRequest) float64 { return r.midhinge() },
			want: 7.2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.fn(tc.req); got != tc.want {
				t.Fatalf("result mismatch: got %.4f want %.4f", got, tc.want)
			}
		})
	}
}

func TestMedianOf_EdgeCases(t *testing.T) {
	if got := medianOf(nil); got != 0 {
		t.Fatalf("medianOf(nil): got %.1f want 0", got)
	}

	even := []Rating{1, 3, 5, 7}
	if got, want := medianOf(even), 4.0; got != want {
		t.Fatalf("medianOf(even): got %.1f want %.1f", got, want)
	}
}

func TestRatingMarshalViaEncodingJSON_OneDecimal(t *testing.T) {
	type payload struct {
		Value *Rating `json:"value"`
	}

	r := Rating(8.66)
	b, err := json.Marshal(payload{Value: &r})
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	if got, want := string(b), `{"value":8.7}`; got != want {
		t.Fatalf("payload mismatch: got %s want %s", got, want)
	}
}

