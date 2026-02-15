package parser

import (
	"context"
	"file-parser/internal/config"
	"file-parser/internal/entity"
	"file-parser/internal/service"
	"fmt"

	"github.com/sirupsen/logrus"
)

// File parser
type FileParser struct {
	inputDir  string
	outputDir string
	service   service.Writer
	logger    *logrus.Logger
}

func NewFileParser(service service.Writer, config *config.ParserConfig, logger *logrus.Logger) *FileParser {
	return &FileParser{
		inputDir:  config.InputDir,
		outputDir: config.OutputDir,
		service:   service,
		logger:    logger,
	}
}

// Parse file by filename from common channel
func (p *FileParser) ParseFile(ctx context.Context, jobs chan string) {
	for {
		select {
		case <-ctx.Done():
			return
		case filename, _ := <-jobs:
			devices, deviceData, err := p.service.ParseFile(ctx, filename)
			if err != nil {
				if err == ctx.Err() {
					return
				}

				p.logger.Errorf("parser error: %v", err)
			}

			if err := p.service.LoadParsedData(ctx, &entity.File{Filename: filename}, devices, deviceData); err != nil {
				if err == ctx.Err() {
					return
				}

				p.logger.Errorf("database error: %v", err)
			}

			fmt.Println(filename)
		}
	}
}
