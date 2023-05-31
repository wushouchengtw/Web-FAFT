package lib

import (
	"backend/lib/models"
	"backend/lib/repos"
	"backend/utils"

	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func configHeaderIndex(row []string) map[Header]int {
	headers := map[Header]int{}

	for idx, value := range row {
		switch strings.ToLower(value) {
		case strings.ToLower(string(Suite)):
			headers[Suite] = idx
		case strings.ToLower(string(Model)):
			headers[Model] = idx
		case strings.ToLower(string(Board)):
			headers[Board] = idx
		case strings.ToLower(string(Test)):
			headers[Test] = idx
		case strings.ToLower(string(Status)):
			headers[Status] = idx
		case strings.ToLower(string(FailureReason)):
			headers[FailureReason] = idx
		case strings.ToLower(string(StartedTime)):
			headers[StartedTime] = idx
		case strings.ToLower(string(Duration)):
			headers[Duration] = idx
		case strings.ToLower(string(BuildVersion)):
			headers[BuildVersion] = idx
		case strings.ToLower(string(FirmwareROVersion)):
			headers[FirmwareROVersion] = idx
		case strings.ToLower(string(FirmwareRWVersion)):
			headers[FirmwareRWVersion] = idx
		case strings.ToLower(string(Hostname)):
			headers[Hostname] = idx
		default:
			log.Printf("Unknow column {%v}", value)
		}
	}

	return headers
}

func saveResult(
	row []string,
	headerIdx map[Header]int,
	dut_repo repos.IDUTRepo,
	test_repo repos.ITestRepo,
	result_repo repos.IResultRepo,
) error {
	res := models.SuiteResult{}
	model := row[headerIdx[Model]]
	board := row[headerIdx[Board]]

	var err error

	var dutId int
	dutId, err = dut_repo.GetIdBy(model, board)
	if err == utils.NotFound {
		dutId, err = dut_repo.Save(model, board)
	}
	if err != nil {
		return err
	}

	testName := row[headerIdx[Test]]
	var testId int
	testId, err = test_repo.GetIdBy(testName)
	if err == utils.NotFound {
		testId, err = test_repo.Save(testName)
	}
	if err != nil {
		return err
	}

	duration, err := strconv.ParseFloat(row[headerIdx[Duration]], 2)
	if err != nil {
		return err
	}

	var status bool
	switch strings.ToLower(row[headerIdx[Status]]) {
	case "pass":
		status = true
	case "fail":
		status = false
	default:
		log.Printf("status should be `pass` or `fail`, but got an {%v}", status)
		return utils.InvalidData
	}
	startTime, err := time.Parse(timeLayout, row[headerIdx[StartedTime]])
	if err != nil {
		log.Printf("Can't parse the startTime, got an {%v}", err)
		return utils.InvalidData
	}

	res.Suite = row[headerIdx[Suite]]
	res.Time = startTime
	res.Duration = duration
	res.DutId = dutId
	res.TestId = testId
	res.Reason = row[headerIdx[FailureReason]]
	res.Status = status
	res.FirmwareROVersion = row[headerIdx[FirmwareROVersion]]
	res.FirmwareRWVersion = row[headerIdx[FirmwareRWVersion]]
	res.Host = row[headerIdx[Hostname]]

	result_repo.Save(&res)

	return nil
}

func SaveStainlessData(
	path string,
	dut_repo repos.IDUTRepo,
	test_repo repos.ITestRepo,
	result_repo repos.IResultRepo,
) error {
	// Open a csv file
	f, err := os.Open(path)
	if err != nil {
		log.Printf("Can't read a file {%v}", path)
		return utils.FileNotExist
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Printf("Can't close a file {%v}", path)
		}
	}(f)

	// Config the csv reader
	reader := csv.NewReader(f)
	reader.Comma = ','
	reader.LazyQuotes = true

	// Config the header index
	row, err := reader.Read()
	if err != nil {
		return err
	}
	headers := configHeaderIndex(row)
	log.Printf("%v\n", headers)

	// loop over the reader
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Invalid data in csv, got an {%v}", err)
			return utils.InvalidData
		}

		saveResult(row, headers, dut_repo, test_repo, result_repo)

		break
	}

	return nil
}
