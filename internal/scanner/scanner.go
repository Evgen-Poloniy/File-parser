package scanner

import (
	"context"
	"file-parser/internal/config"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
)

// Scanner of files
type Scanner struct {
	inputDir      string
	scanFrequency time.Duration
	logger        *logrus.Logger
}

func NewScanner(config *config.ParserConfig, logger *logrus.Logger) *Scanner {
	return &Scanner{
		inputDir:      config.InputDir,
		scanFrequency: config.ScanFrequency,
		logger:        logger,
	}
}

// Scan directory with target files and write filenames at channel
func (s *Scanner) Scan(ctx context.Context, jobs chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			start := time.Now()

			fileNames, err := s.scanDir(ctx)
			if err != nil {
				if err == ctx.Err() {
					return
				}

				s.logger.Errorf("scanner error: %v", err)

			}

			for _, fileName := range fileNames {
				select {
				case <-ctx.Done():
					return
				case jobs <- fileName:
				}
			}

			latency := time.Since(start)

			wait := s.scanFrequency - latency
			if wait > 0 {
				select {
				case <-ctx.Done():
					return
				case <-time.After(wait):
				}
			}
		}
	}
}

// Private method which scan directory and write filenames at strings slice
func (s *Scanner) scanDir(ctx context.Context) ([]string, error) {
	entries, err := os.ReadDir(s.inputDir)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, 0, 100)

	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			if entry.IsDir() {
				continue
			}

			if filepath.Ext(entry.Name()) == ".tsv" {
				fullPath := filepath.Join(s.inputDir, entry.Name())
				fileNames = append(fileNames, fullPath)
			}
		}
	}

	return fileNames, nil
}
