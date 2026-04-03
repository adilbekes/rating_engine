package main

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteErrorToFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "error.json")

	if err := writeErrorToFile(path, "boom"); err != nil {
		t.Fatalf("writeErrorToFile: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read error file: %v", err)
	}

	var got map[string]string
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("unmarshal error file: %v", err)
	}
	if got["error"] != "boom" {
		t.Fatalf("error payload mismatch: got %q want %q", got["error"], "boom")
	}
}

func TestReportError_WritesToStderrAndFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "error.json")

	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stderr = w
	defer func() {
		os.Stderr = oldStderr
	}()

	reportError("bad input", path)

	_ = w.Close()
	stderrData, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read stderr pipe: %v", err)
	}

	if !strings.Contains(string(stderrData), `"error":"bad input"`) {
		t.Fatalf("stderr mismatch: %s", string(stderrData))
	}

	fileData, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read output file: %v", err)
	}
	if !strings.Contains(string(fileData), `"error":"bad input"`) {
		t.Fatalf("file mismatch: %s", string(fileData))
	}
}

func TestWriteErrorToFile_CreateFailure(t *testing.T) {
	dir := t.TempDir()
	if err := writeErrorToFile(dir, "x"); err == nil {
		t.Fatal("expected create error when path is a directory")
	}
}

func TestEngineMain_ProcessScenarios(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		stdin      string
		wantCode   int
		wantStdErr string
		wantStdOut string
		checkFile  bool
	}{
		{
			name:       "rejects mutually exclusive flags",
			args:       []string{"-d", `{"subject":"a"}`, "-f", "input.json"},
			wantCode:   1,
			wantStdErr: "cannot use both -d and -f flags",
		},
		{
			name:       "rejects invalid inline json",
			args:       []string{"-d", `{"subject":`},
			wantCode:   1,
			wantStdErr: "invalid input JSON",
		},
		{
			name:       "rejects missing input file",
			args:       []string{"-f", "/tmp/file-that-does-not-exist-rating-engine.json"},
			wantCode:   1,
			wantStdErr: "failed to open file",
		},
		{
			name: "returns validation error from engine",
			args: []string{"-d", `{"subject":"","current_rating":5.0,"history_weight":1,"votes":[5]}`},
			wantCode:   1,
			wantStdErr: "subject name must be provided",
		},
		{
			name:       "reads request from stdin",
			stdin:      `{"subject":"A","current_rating":5.6,"history_weight":10,"votes":[7,6,6,5,7,6,6,7,10],"scoring_mode":"midhinge"}`,
			wantCode:   0,
			wantStdOut: `"new_rating":6`,
		},
		{
			name:       "writes success to stdout with -d",
			args:       []string{"-d", `{"subject":"A","current_rating":5.6,"history_weight":10,"votes":[7,6,6,5,7,6,6,7,10],"scoring_mode":"midhinge"}`},
			wantCode:   0,
			wantStdOut: `"subject":"A"`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stdout, stderr, code := runEngineMainProcess(t, tc.args, tc.stdin)
			if code != tc.wantCode {
				t.Fatalf("exit code mismatch: got %d want %d; stderr=%s", code, tc.wantCode, stderr)
			}
			if tc.wantStdErr != "" && !strings.Contains(stderr, tc.wantStdErr) {
				t.Fatalf("stderr mismatch: want %q in %q", tc.wantStdErr, stderr)
			}
			if tc.wantStdOut != "" && !strings.Contains(stdout, tc.wantStdOut) {
				t.Fatalf("stdout mismatch: want %q in %q", tc.wantStdOut, stdout)
			}
		})
	}
}

func TestEngineMain_OutputFileScenarios(t *testing.T) {
	dir := t.TempDir()
	outPath := filepath.Join(dir, "result.json")

	stdout, stderr, code := runEngineMainProcess(t,
		[]string{
			"-d", `{"subject":"A","current_rating":5.6,"history_weight":10,"votes":[7,6,6,5,7,6,6,7,10],"scoring_mode":"midhinge"}`,
			"-o", outPath,
		},
		"",
	)
	if code != 0 {
		t.Fatalf("unexpected exit code: %d stderr=%s", code, stderr)
	}
	if stdout != "" {
		t.Fatalf("expected empty stdout when -o is used, got %q", stdout)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !strings.Contains(string(data), `"subject":"A"`) {
		t.Fatalf("unexpected output file content: %s", string(data))
	}

	badPath := dir // directory path will fail os.Create
	_, badErr, badCode := runEngineMainProcess(t,
		[]string{
			"-d", `{"subject":"A","current_rating":5.6,"history_weight":10,"votes":[7,6,6,5,7,6,6,7,10],"scoring_mode":"midhinge"}`,
			"-o", badPath,
		},
		"",
	)
	if badCode != 1 {
		t.Fatalf("expected failure exit code for invalid -o path, got %d", badCode)
	}
	if !strings.Contains(badErr, "failed to create output file") {
		t.Fatalf("expected create-output-file error, got %q", badErr)
	}
}

func runEngineMainProcess(t *testing.T, args []string, stdin string) (stdout string, stderr string, code int) {
	t.Helper()

	cmdArgs := append([]string{"-test.run=TestEngineMainHelperProcess", "--"}, args...)
	cmd := exec.Command(os.Args[0], cmdArgs...)
	cmd.Env = append(os.Environ(), "GO_WANT_ENGINE_MAIN_HELPER=1")
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}

	out, err := cmd.CombinedOutput()
	all := string(out)
	if err == nil {
		return all, "", 0
	}

	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		t.Fatalf("unexpected process error: %v", err)
	}

	// main writes success payload to stdout and error payload to stderr.
	// CombinedOutput merges both streams; tests assert by substring.
	return all, all, exitErr.ExitCode()
}

func TestEngineMainHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_ENGINE_MAIN_HELPER") != "1" {
		return
	}

	idx := 0
	for idx < len(os.Args) && os.Args[idx] != "--" {
		idx++
	}
	if idx >= len(os.Args) {
		os.Exit(2)
	}
	engineArgs := os.Args[idx+1:]

	os.Args = append([]string{"engine"}, engineArgs...)
	main()
	os.Exit(0)
}


