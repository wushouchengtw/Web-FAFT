package csv

import (
	"backend/lib/models"
	"backend/lib/remoteRepo/dut"
	"backend/lib/remoteRepo/result"
	"backend/lib/remoteRepo/test"
	"backend/utils"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func SaveRemoteDataByCsv(fileName string, dutRepo dut.IDUT, testRepo test.Itest, resultRepo result.IResut) error {
	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to open %q: %v", fileName, err)
	}
	defer f.Close()

	// Config for CSV
	r := csv.NewReader(f)
	r.Comma = ','
	r.LazyQuotes = true

	// Verify the headers
	row, err := r.Read()
	if err != nil {
		return fmt.Errorf("cannot read CSV data: %v", err)
	}

	header, err := GetCsvHeader(row)
	if err != nil {
		return fmt.Errorf("while finding csv headers: %v", err)
	}

	for {
		// Don't use read all.
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return utils.ErrInvalidData
		}
		if err := saveResult(row, header, dutRepo, testRepo, resultRepo); err != nil {
			return fmt.Errorf("failed to save row data [%v]: %v ", row, err)
		}
	}
	return nil
}

func GetCsvHeader(firstRow []string) (map[TestHausHeader]int, error) {
	header := map[TestHausHeader]int{}
	for index, name := range firstRow {
		switch strings.ToLower(name) {
		case strings.ToLower(string(Suite)):
			header[Suite] = index
		case strings.ToLower(string(Board)):
			header[Board] = index
		case strings.ToLower(string(Model)):
			header[Model] = index
		case strings.ToLower(string(Test)):
			header[Test] = index
		case strings.ToLower(string(Status)):
			header[Status] = index
		case strings.ToLower(string(FailureReason)):
			header[FailureReason] = index
		case strings.ToLower(string(StartedTime)):
			header[StartedTime] = index
		case strings.ToLower(string(Duration)):
			header[Duration] = index
		case strings.ToLower(string(BuildVersion)):
			header[BuildVersion] = index
		case strings.ToLower(string(FirmwareROVersion)):
			header[FirmwareROVersion] = index
		case strings.ToLower(string(FirmwareRWVersion)):
			header[FirmwareRWVersion] = index
		case strings.ToLower(string(Hostname)):
			header[Hostname] = index
		default:
			return header, utils.ErrNotFound
		}
	}
	return header, nil
}

func saveResult(row []string, header map[TestHausHeader]int, dutRepo dut.IDUT, testRepo test.Itest, resultRepo result.IResut) error {

	// Parse data
	var (
		dutId     *string
		testId    *string
		board     string = row[header[Board]]
		model     string = row[header[Model]]
		test      string = row[header[Test]]
		duration  float64
		milestone string = "unknown"
		version   string = "unknown"
		status    bool   = false

		err error
	)
	dutId, err = dutRepo.SaveIfNotExist(board, model)
	if err != nil {
		return fmt.Errorf("failed on processing dut info: %v", err)
	}

	testId, err = testRepo.SaveIfNotExist(test)
	if err != nil {
		return fmt.Errorf("failed on processing test info: %v", err)
	}

	duration, err = strconv.ParseFloat(row[header[Duration]], 64)
	if err != nil {
		return err
	}

	duration, err = strconv.ParseFloat(fmt.Sprintf("%.2f", duration), 64)
	if err != nil {
		return err
	}

	buildVersion := strings.Split(row[header[BuildVersion]], "-")
	if len(buildVersion) == 2 {
		milestone = buildVersion[0]
		version = buildVersion[1]
	}

	switch strings.ToLower(row[header[Status]]) {
	case "pass":
		status = true
	case "fail":
		status = false
	default:
		return utils.ErrInvalidData
	}

	startTime, err := time.Parse(timeLayout, row[header[StartedTime]])
	if err != nil {
		return utils.ErrInvalidData
	}

	id := uuid.New().String()
	record := &models.Result{
		Id:                id,
		Suite:             row[header[Suite]],
		Time:              startTime,
		Duration:          duration,
		DutId:             *dutId,
		TestId:            *testId,
		Reason:            row[header[FailureReason]],
		Milestone:         milestone,
		Version:           version,
		Status:            status,
		FirmwareROVersion: row[header[FirmwareROVersion]],
		FirmwareRWVersion: row[header[FirmwareRWVersion]],
		Host:              row[header[Hostname]],
	}
	resultRepo.Save(record)
	return nil
}
