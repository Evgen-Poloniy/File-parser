package service

import (
	"context"
	"encoding/csv"
	"file-parser/internal/entity"
	"file-parser/internal/repository"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type WriterService struct {
	writer repository.Writer
}

func NewWriterService(writer repository.Writer) *WriterService {
	return &WriterService{writer: writer}
}

// Parse tsv file and return slices of data
func (w *WriterService) ParseFile(ctx context.Context, path string) (
	[]entity.Device, []entity.DeviceData, error,
) {
	filename := filepath.Base(path)
	fileID, err := w.writer.CheckRecordAboutFile(ctx, filename)
	if err != nil {
		return nil, nil, err
	}
	if fileID != 0 {
		return nil, nil, nil
	}

	fileID, err = w.writer.CheckRecordAboutError(ctx, filename)
	if err != nil {
		return nil, nil, err
	}
	if fileID != 0 {
		return nil, nil, nil
	}

	var parseErrors []entity.ParseError

	file, err := os.Open(path)
	if err != nil {
		parseErr := entity.ParseError{
			LineNum:  nil,
			LineData: nil,
			Err:      fmt.Errorf("cannot open tsv file %s: %w", filename, err),
		}

		return nil, nil, entity.ParseErrors{
			Errors: []entity.ParseError{parseErr},
		}
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.FieldsPerRecord = 15

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, nil, fmt.Errorf("cannot read tsv file %s: %w", filename, err)
	}

	deviceData := make([]entity.DeviceData, 0, len(lines)-1)
	devices := make([]entity.Device, 0, len(lines)-1)
	deviceUnitGUID := make(map[string]bool)

	for i, line := range lines {
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		default:
			if i == 0 {
				continue
			}

			parseInt := func(fieldName string, s string) *int {
				if s == "" {
					return nil
				}
				val, err := strconv.Atoi(s)
				if err != nil {
					lineNum := i + 1
					parseErrors = append(parseErrors, entity.ParseError{
						LineNum:  &lineNum,
						LineData: line,
						Err:      fmt.Errorf("field %s cannot parse '%s' as int", fieldName, s),
					})
					return nil
				}
				return &val
			}

			parseStr := func(s string) *string {
				if s == "" {
					return nil
				}
				return &s
			}

			deviceData = append(deviceData, entity.DeviceData{
				N:         parseInt("n", line[0]),
				MQTT:      parseStr(line[1]),
				InvID:     line[2],
				UnitGUID:  line[3],
				MsgID:     parseStr(line[4]),
				Text:      parseStr(line[5]),
				Context:   parseStr(line[6]),
				Class:     parseStr(line[7]),
				Level:     parseInt("level", line[8]),
				Area:      parseStr(line[9]),
				Addr:      parseStr(line[10]),
				Block:     parseStr(line[11]),
				Type:      parseStr(line[12]),
				Bit:       parseInt("bit", line[13]),
				InvertBit: parseInt("invert_bit", line[14]),
			})

			if !deviceUnitGUID[line[2]] {
				devices = append(devices, entity.Device{
					InvID:    line[2],
					UnitGUID: line[3],
				})
				deviceUnitGUID[line[2]] = true
			}
		}
	}

	if len(parseErrors) > 0 {
		return devices, deviceData, entity.ParseErrors{Errors: parseErrors}
	}

	return devices, deviceData, nil
}

// Load data into database
func (w *WriterService) LoadParsedData(
	ctx context.Context,
	filename string,
	devices []entity.Device,
	deviceData []entity.DeviceData,
) error {
	return w.writer.InsertParsedData(ctx, filename, devices, deviceData)
}

// Load information about parse errors into the database
func (w *WriterService) LoadErrors(ctx context.Context, filename string, parseErrors *entity.ParseErrors) error {
	return w.writer.InsertErrors(ctx, filename, parseErrors)
}

// Write data from tsv file into rtf file
func (w *WriterService) WriteDataToRTF(path string, data []entity.DeviceData) error {
	var sb strings.Builder

	sb.WriteString("{\\rtf1\\ansi\\deff0\n")
	sb.WriteString("{\\fonttbl{\\f0 Arial;}}\n")
	sb.WriteString("\\f0\\fs24\n")

	for _, d := range data {
		sb.WriteString(fmt.Sprintf("N: %v\\line\n", ptrIntToStr(d.N)))
		sb.WriteString(fmt.Sprintf("MQTT: %s\\line\n", ptrStrToStr(d.MQTT)))
		sb.WriteString(fmt.Sprintf("InvID: %s\\line\n", d.UnitGUID))
		sb.WriteString(fmt.Sprintf("UnitGUID: %s\\line\n", d.UnitGUID))
		sb.WriteString(fmt.Sprintf("MsgID: %s\\line\n", ptrStrToStr(d.MsgID)))
		sb.WriteString(fmt.Sprintf("Text: %s\\line\n", escapeRTFUnicode(ptrStrToStr(d.Text))))
		sb.WriteString(fmt.Sprintf("Context: %s\\line\n", escapeRTFUnicode(ptrStrToStr(d.Context))))
		sb.WriteString(fmt.Sprintf("Class: %s\\line\n", ptrStrToStr(d.Class)))
		sb.WriteString(fmt.Sprintf("Level: %v\\line\n", ptrIntToStr(d.Level)))
		sb.WriteString(fmt.Sprintf("Area: %s\\line\n", ptrStrToStr(d.Area)))
		sb.WriteString(fmt.Sprintf("Addr: %s\\line\n", ptrStrToStr(d.Addr)))
		sb.WriteString(fmt.Sprintf("Block: %s\\line\n", ptrStrToStr(d.Block)))
		sb.WriteString(fmt.Sprintf("Type: %s\\line\n", ptrStrToStr(d.Type)))
		sb.WriteString(fmt.Sprintf("Bit: %v\\line\n", ptrIntToStr(d.Bit)))
		sb.WriteString(fmt.Sprintf("InvertBit: %v\\line\n", ptrIntToStr(d.InvertBit)))
		sb.WriteString("\\line\n")
	}

	sb.WriteString("}")

	return os.WriteFile(path, []byte(sb.String()), 0644)
}

// Write information about errors from tsv file into rtf file
func (w *WriterService) WriteErrorsToRTF(path string, parseErrs *entity.ParseErrors) error {
	if parseErrs == nil || len(parseErrs.Errors) == 0 {
		return nil
	}

	var sb strings.Builder

	sb.WriteString("{\\rtf1\\ansi\\deff0\n")
	sb.WriteString("{\\fonttbl{\\f0 Arial;}}\n")
	sb.WriteString("\\f0\\fs24\n")

	for i, pe := range parseErrs.Errors {
		if pe.LineNum != nil {
			sb.WriteString(fmt.Sprintf("line: %d\\line\n", *pe.LineNum))
		}

		if len(pe.LineData) > 0 {
			lineData := strings.Join(pe.LineData, " | ")
			sb.WriteString("line_data: ")
			sb.WriteString(escapeRTFUnicode(lineData))
			sb.WriteString("\\line\n")
		}

		if pe.Err != nil {
			sb.WriteString("error_msg: ")
			sb.WriteString(escapeRTFUnicode(pe.Err.Error()))
			sb.WriteString("\\line\n")
		}

		if i < len(parseErrs.Errors)-1 {
			sb.WriteString("\\line\n")
		}
	}

	sb.WriteString("}")

	return os.WriteFile(path, []byte(sb.String()), 0644)
}

func ptrStrToStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ptrIntToStr(i *int) string {
	if i == nil {
		return ""
	}
	return fmt.Sprintf("%d", *i)
}

func escapeRTFUnicode(s string) string {
	var sb strings.Builder
	for _, r := range s {
		if r <= 127 {
			switch r {
			case '\\', '{', '}':
				sb.WriteRune('\\')
				sb.WriteRune(r)
			default:
				sb.WriteRune(r)
			}
		} else {
			sb.WriteString(fmt.Sprintf("\\u%d?", int32(r)))
		}
	}
	return sb.String()
}
