package service

import (
	"context"
	"encoding/csv"
	"file-parser/internal/entity"
	"file-parser/internal/repository"
	"os"
	"strconv"
)

type WriterService struct {
	writer repository.Writer
}

func newWriterService(writer repository.Writer) *WriterService {
	return &WriterService{writer: writer}
}

func (w *WriterService) ParseFile(ctx context.Context, path string) ([]entity.Device, []entity.DeviceData, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.FieldsPerRecord = 15

	lines, err := reader.ReadAll()
	if err != nil {
		return nil, nil, err
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

			parseInt := func(s string) *int {
				if s == "" {
					return nil
				}
				val, err := strconv.Atoi(s)
				if err != nil {
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
				N:         parseInt(line[0]),
				MQTT:      parseStr(line[1]),
				UnitGUID:  line[3],
				MsgID:     line[4],
				Text:      line[5],
				Context:   line[6],
				Class:     line[7],
				Level:     parseInt(line[8]),
				Area:      line[9],
				Addr:      line[10],
				Block:     parseStr(line[11]),
				Type:      parseStr(line[12]),
				Bit:       parseInt(line[12]),
				InvertBit: parseInt(line[13]),
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

	return devices, deviceData, nil
}

// Load data into database
func (w *WriterService) LoadParsedData(
	ctx context.Context,
	file *entity.File,
	devices []entity.Device,
	deviceData []entity.DeviceData,
) error {
	return w.writer.InsertParsedData(ctx, file, devices, deviceData)
}
