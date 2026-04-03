package main

import (
	"io"
	"os"
	"strings"
	"testing"
)

func TestMain_PrintsExpectedSummary(t *testing.T) {
	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	defer func() {
		os.Stdout = oldStdout
	}()

	main()

	_ = w.Close()
	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read stdout: %v", err)
	}
	output := string(data)

	checks := []string{
		"Subject: A",
		"Old rating: 5.6",
		"Event score (midhinge): 6.5",
		"Effective history weight: 10",
		"Votes count: 9",
		"New rating: 6.0",
	}
	for _, want := range checks {
		if !strings.Contains(output, want) {
			t.Fatalf("missing output %q in %q", want, output)
		}
	}
}

