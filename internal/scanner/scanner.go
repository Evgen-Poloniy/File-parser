package scanner

import (
	"context"
	"file-parser/internal/config"
	"os"
	"path/filepath"
	"time"
)

// Scanner of files
type Scanner struct {
	inputDir      string
	inputFormat   string
	scanFrequency time.Duration
}

func NewScanner(config *config.ParserConfig) *Scanner {
	return &Scanner{
		inputDir:      config.InputDir,
		inputFormat:   config.InputFormat,
		scanFrequency: config.ScanFrequency,
	}
}

// Scan directory with target files and write filenames at channel
func (s *Scanner) Scan(ctx context.Context, jobs chan string) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			start := time.Now()

			fileNames, err := s.scanDir()
			if err != nil {
				ctx.Done()
				return err
			}

			for _, fileName := range fileNames {
				select {
				case <-ctx.Done():
					return nil
				case jobs <- fileName:
				}
			}

			latency := time.Since(start)

			wait := s.scanFrequency - latency
			if wait > 0 {
				select {
				case <-ctx.Done():
					return nil
				case <-time.After(wait):
				}
			}
		}
	}
}

// Private method which scan directory and write filenames at strings slice
func (s *Scanner) scanDir() ([]string, error) {
	entries, err := os.ReadDir(s.inputDir)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, 0, 100)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) == s.inputFormat {
			fullPath := filepath.Join(s.inputDir, entry.Name())
			fileNames = append(fileNames, fullPath)
		}
	}

	return fileNames, nil
}
