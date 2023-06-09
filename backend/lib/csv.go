package lib

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ValidCsv(inputCsv string) error {
	err := checkCsvInputFormat(inputCsv)
	if err != nil {
		return fmt.Errorf("while checking file: %v", err)
	}

	csvFiles, err := getCsvFiles()
	if err != nil {
		return fmt.Errorf("failed to get all csv files in %q: %v", CsvFolder, err)
	}

	if overlap := verifyOverlap(csvFiles, inputCsv); overlap {
		return fmt.Errorf("duplicated log in DB: %v", err)
	}

	log.Println("Valid csv file and be ready to be inserted")
	return nil
}

func checkCsvInputFormat(inputCsv string) error {
	// Confirm the format is correct
	nameReg, _ := regexp.Compile(`(\d{4})(\d{2})(\d{2})-(\d{4})(\d{2})(\d{2}).csv`)
	matches := nameReg.FindStringSubmatch(inputCsv)
	if len(matches) < 7 {
		return errors.New("csv naming format error: it should be [nums]-[nums].csv")
	}

	layout := "2006-01-02T15:04:05Z"
	date1 := fmt.Sprintf("%s-%s-%sT00:00:00Z", matches[1], matches[2], matches[3])
	date2 := fmt.Sprintf("%s-%s-%sT00:00:00Z", matches[4], matches[5], matches[6])
	t1, err := time.Parse(layout, date1)
	if err != nil {
		return errors.New("failed to parse time")
	}
	t2, err := time.Parse(layout, date2)
	if err != nil {
		return errors.New("failed to parse time")
	}
	if t2.Sub(t1) < 0 {
		return fmt.Errorf("%s time format wrong, the former [%s] should be eariler than [%s]", inputCsv, strings.Join(matches[1:4], "-"), strings.Join(matches[4:7], "-"))
	}
	return nil
}

func getCsvFiles() ([]byte, error) {
	cmdString := fmt.Sprintf("ls %s | grep .csv", CsvFolder)
	cmd := exec.Command("bash", "-c", cmdString)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	log.Println("Searching the csvs in Server...")
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("while running ls command: %v", err)
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, fmt.Errorf("while reading the command output: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		// "exit status 1" means "not found" here.
		if !strings.Contains(err.Error(), "exit status 1") {
			return nil, fmt.Errorf("unexpected error happened while running ls command: %v", err)
		}
	}
	return bytes, nil
}

func verifyOverlap(bytes []byte, newFilename string) bool {
	fileNames := strings.Split(string(bytes), "\n")
	for _, fileName := range fileNames {
		if testOverlap := overlap(newFilename, fileName); testOverlap {
			return true
		}
	}
	return false
}

func overlap(newFilename string, fileName string) bool {
	newFileDate := strings.Split(newFilename, "-")
	fileDate := strings.Split(fileName, "-")

	n1, _ := strconv.Atoi(newFileDate[0])
	n2, _ := strconv.Atoi(newFileDate[1])

	f1, _ := strconv.Atoi(fileDate[0])
	f2, _ := strconv.Atoi(fileDate[1])

	if n2 < f1 {
		return false
	}
	if n1 > f2 {
		return false
	}
	return true
}
