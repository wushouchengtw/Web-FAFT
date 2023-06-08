package dbInfo

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type KeyLocation struct {
	Key      string
	Location int
}

// stainlessDLocation would be used to record the current CSV file's data attribute location.
var stainlessDASLocation = map[string]KeyLocation{
	"Started Time":        {Key: "time", Location: 0},
	"Duration":            {Key: "duration", Location: 0},
	"Suite":               {Key: "suite", Location: 0},
	"Board":               {Key: "board", Location: 0},
	"Model":               {Key: "model", Location: 0},
	"Build Version":       {Key: "buildVersion", Location: 0},
	"Hostname":            {Key: "host", Location: 0},
	"Test":                {Key: "testName", Location: 0},
	"Status":              {Key: "status", Location: 0},
	"Failure Reason":      {Key: "reason", Location: 0},
	"Firmware RO Version": {Key: "firmwareROVersion", Location: 0},
	"Firmware RW Version": {Key: "firmwareRWVersion", Location: 0},
}

func ValidCsv(inputCsv string, bytes []byte) error {
	if err := checkCsvInputFormat(inputCsv); err != nil {
		return fmt.Errorf("error while checking file: %v", err)
	}

	if overlap := verifyOverlap(bytes, inputCsv); overlap {
		return fmt.Errorf("csv overlap")
	}

	Log("Valid log name")
	return nil
}

// checkCsvInputFormat verifies the input file name should be `\d+-\d+.csv`
// Rules:
//  1. 8 digits date [$year$month$day]
//  2. Year: after 2020
//  3. Month: [1-12]
//  4. Day: [1-31]
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

// verifyOverlap would verify the same period of data would not be duplicated.
// It would be verified by the name of csv file.
func verifyOverlap(bytes []byte, csv string) bool {
	fileNames := strings.Split(string(bytes), "\n")
	for _, fileName := range fileNames {
		if testOverlap := overlap(csv, fileName); testOverlap {
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

func ReturnCSVDataLocation(collumnsLen int, rows [][]string) (data [][]string, insertKeySequence []string) {
	if len(rows) == 1 {
		log.Panicf("Only attributes in csv file")
	}
	if len(rows[0]) != 12 {
		log.Panicf(fmt.Sprintf("The numbers of collumn should be 12, but it has %d", len(rows)))
	}
	mapCsvKey := map[string]int{}
	for attribute, keyMap := range stainlessDASLocation {
		for index, key := range rows[0] {
			if strings.Contains(key, attribute) {
				mapCsvKey[keyMap.Key] = index
				break
			}
			if index == len(stainlessDASLocation)-1 {
				log.Panicf(fmt.Sprintf("Attribute [%v] wasn't found on csv files", attribute))
			}
		}
	}

	keysSortedByCsv := make([]string, len(mapCsvKey))
	for key, location := range mapCsvKey {
		keysSortedByCsv[location] = key
	}
	return rows[1:], keysSortedByCsv
}
