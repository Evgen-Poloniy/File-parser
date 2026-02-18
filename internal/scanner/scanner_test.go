package scanner

import (
	"context"
	"file-parser/internal/config"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestScanDir(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "file1.tsv"), []byte("data"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.txt"), []byte("data"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file3.tsv"), []byte("data"), 0644)

	logger := logrus.New()
	logger.Out = io.Discard
	scanner := &Scanner{
		inputDir: tmpDir,
		logger:   logger,
	}

	files, err := scanner.scanDir(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 tsv files, got %d", len(files))
	}

	for _, f := range files {
		if filepath.Ext(f) != ".tsv" {
			t.Fatalf("expected .tsv file, got %s", f)
		}
	}
}

func TestScan(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "file1.tsv"), []byte("data"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file2.tsv"), []byte("data"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "file3.txt"), []byte("data"), 0644)

	cfg := &config.ParserConfig{
		InputDir:      tmpDir,
		ScanFrequency: 50 * time.Millisecond,
	}

	scanner := NewScanner(cfg, logrus.New())
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	jobs := make(chan string, 10)

	go scanner.Scan(ctx, jobs)

	collected := make(map[string]bool)

	done := false

	for !done {
		select {
		case f := <-jobs:
			collected[f] = true
		case <-ctx.Done():
			done = true
		}
	}

	if len(collected) != 2 {
		t.Fatalf("expected 2 tsv files, got %d", len(collected))
	}

	if !collected["file1.tsv"] || !collected["file2.tsv"] {
		t.Fatalf("missing expected files: %+v", collected)
	}
}
