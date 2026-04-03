// engine is a language-agnostic CLI binary for rating updates.
//
// Usage:
//
//	echo '<json>' | engine                          # stdin
//	engine -d '<json>'                              # JSON string flag
//	engine -f request.json                          # JSON file flag
//	engine -f request.json -o result.json           # file input and output
//
// Input  - JSON object on stdin/flag (see UpdateRatingRequest)
// Output - JSON object on stdout (see UpdateRatingResponse), or {"error":"..."} on failure
// Exit   - 0 on success, 1 on error
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"rating_engine/pkg/ratingengine"
	"strings"
)

type errorResponse struct {
	Error string `json:"error"`
}

func writeErrorToStderr(msg string) {
	_ = json.NewEncoder(os.Stderr).Encode(errorResponse{Error: msg})
}

func writeErrorToFile(filename string, msg string) (err error) {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := file.Close(); err == nil && closeErr != nil {
			err = closeErr
		}
	}()

	err = json.NewEncoder(file).Encode(errorResponse{Error: msg})
	return err
}

func reportError(msg string, outputFile string) {
	writeErrorToStderr(msg)
	if outputFile != "" {
		_ = writeErrorToFile(outputFile, msg)
	}
}

func main() {
	dataFlag := flag.String("d", "", "JSON request string")
	fileFlag := flag.String("f", "", "JSON request file")
	outputFlag := flag.String("o", "", "JSON output file (if not set, output to stdout)")
	flag.Parse()

	var inputData io.Reader

	if *dataFlag != "" && *fileFlag != "" {
		reportError("cannot use both -d and -f flags", *outputFlag)
		os.Exit(1)
	}

	if *dataFlag != "" {
		inputData = strings.NewReader(*dataFlag)
	} else if *fileFlag != "" {
		content, err := os.ReadFile(*fileFlag)
		if err != nil {
			reportError(fmt.Sprintf("failed to open file: %s", err), *outputFlag)
			os.Exit(1)
		}
		inputData = bytes.NewReader(content)
	} else {
		inputData = os.Stdin
	}

	var req ratingengine.UpdateRatingRequest
	if err := json.NewDecoder(inputData).Decode(&req); err != nil {
		reportError(fmt.Sprintf("invalid input JSON: %s", err), *outputFlag)
		os.Exit(1)
	}

	result, err := ratingengine.UpdateRating(req)
	if err != nil {
		reportError(err.Error(), *outputFlag)
		os.Exit(1)
	}

	if *outputFlag != "" {
		file, err := os.Create(*outputFlag)
		if err != nil {
			reportError(fmt.Sprintf("failed to create output file: %s", err), *outputFlag)
			os.Exit(1)
		}
		if err := json.NewEncoder(file).Encode(&result); err != nil {
			_ = file.Close()
			reportError(fmt.Sprintf("failed to encode result: %s", err), *outputFlag)
			os.Exit(1)
		}
		if err := file.Close(); err != nil {
			reportError(fmt.Sprintf("failed to close output file: %s", err), *outputFlag)
			os.Exit(1)
		}
	} else {
		if err := json.NewEncoder(os.Stdout).Encode(&result); err != nil {
			writeErrorToStderr(fmt.Sprintf("failed to encode result: %s", err))
			os.Exit(1)
		}
	}
}
