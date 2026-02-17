package parser

import (
	"context"
	"errors"
	"file-parser/internal/config"
	"file-parser/internal/entity"
	"file-parser/internal/service"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// File parser
type FileParser struct {
	inputDir  string
	outputDir string
	errorDir  string
	service   service.Writer
	logger    *logrus.Logger
}

func NewFileParser(service service.Writer, config *config.ParserConfig, logger *logrus.Logger) *FileParser {
	return &FileParser{
		inputDir:  config.InputDir,
		outputDir: config.OutputDir,
		errorDir:  config.ErrorDir,
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
			devices, deviceData, err := p.service.ParseFile(ctx, filepath.Join(p.inputDir, filename))
			if err != nil {
				if err == ctx.Err() {
					return
				}

				// Logging error
				var parseErrs entity.ParseErrors
				if errors.As(err, &parseErrs) {
					for _, e := range parseErrs.Errors {
						p.logger.WithFields(logrus.Fields{
							"file": filename,
							"line": e.LineNum,
						}).Error(fmt.Sprintf("parse error: %s", e.Err.Error()))
					}
				} else {
					p.logger.WithFields(logrus.Fields{
						"file": filename,
					}).Error(fmt.Sprintf("parse error: %v", err))
					parseErrs = entity.ParseErrors{
						Errors: []entity.ParseError{
							{
								LineNum:  nil,
								LineData: nil,
								Err:      fmt.Errorf("parse error: %v", err),
							},
						},
					}
				}

				// Load errors into the database
				if err := p.service.LoadErrors(ctx, filename, &parseErrs); err != nil {
					if err == ctx.Err() {
						return
					}

					p.logger.WithFields(logrus.Fields{
						"file": filename,
					}).Error(fmt.Sprintf("error when loading parse errors into the database: %v", err))
				}

				nameWithoutExt := strings.TrimSuffix(filename, ".tsv")
				outputPath := filepath.Join(p.errorDir, "errors_"+nameWithoutExt+".rtf")

				if err := p.service.WriteErrorsToRTF(outputPath, &parseErrs); err != nil {
					p.logger.WithFields(logrus.Fields{
						"file": filename,
					}).Error(fmt.Sprintf("error when writing data to rtf file: %v", err))

					writeErrs := entity.ParseErrors{
						Errors: []entity.ParseError{
							{
								LineNum:  nil,
								LineData: nil,
								Err:      fmt.Errorf("error when writing data to rtf file: %w", err),
							},
						},
					}

					if err := p.service.LoadErrors(ctx, filename, &writeErrs); err != nil {
						p.logger.WithFields(logrus.Fields{
							"file": filename,
						}).Error(fmt.Sprintf("error when loading errors of writing to rtf file error into the database: %v", err))
						continue
					}
				}

				p.logger.WithFields(logrus.Fields{
					"file": filename,
				}).Info("errors from tsv file wrote to rtf file successfully")

				continue
			}
			// Check that file is not in the database
			if devices == nil || deviceData == nil {
				// File already has been added into the database
				continue
			}

			// Load data into the database
			if err := p.service.LoadParsedData(ctx, filename, devices, deviceData); err != nil {
				if err == ctx.Err() {
					return
				}

				p.logger.WithFields(logrus.Fields{
					"file": filename,
				}).Error(fmt.Sprintf("database error: %v", err))

				parseErrs := entity.ParseErrors{
					Errors: []entity.ParseError{
						{
							LineNum:  nil,
							LineData: nil,
							Err:      fmt.Errorf("database error: %w", err),
						},
					},
				}

				if err := p.service.LoadErrors(ctx, filename, &parseErrs); err != nil {
					p.logger.WithFields(logrus.Fields{
						"file": filename,
					}).Error(fmt.Sprintf("error when loading database error into the database: %v", err))
				}

				continue
			}

			p.logger.WithFields(logrus.Fields{
				"file": filename,
			}).Info("file processed successfully")

			nameWithoutExt := strings.TrimSuffix(filename, ".tsv")
			outputPath := filepath.Join(p.outputDir, nameWithoutExt+".rtf")

			if err := p.service.WriteDataToRTF(outputPath, deviceData); err != nil {
				p.logger.WithFields(logrus.Fields{
					"file": filename,
				}).Error(fmt.Sprintf("error when writing data to rtf file: %v", err))

				writeErrs := entity.ParseErrors{
					Errors: []entity.ParseError{
						{
							LineNum:  nil,
							LineData: nil,
							Err:      fmt.Errorf("error when writing data to rtf file: %w", err),
						},
					},
				}

				if err := p.service.LoadErrors(ctx, filename, &writeErrs); err != nil {
					p.logger.WithFields(logrus.Fields{
						"file": filename,
					}).Error(fmt.Sprintf("error when loading errors of writing to rtf file error into the database: %v", err))
				}

				continue
			}

			p.logger.WithFields(logrus.Fields{
				"file": filename,
			}).Info("data from tsv file wrote to rtf file successfully")
		}
	}
}
