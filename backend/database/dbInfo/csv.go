package dbInfo

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

func ValidCsv(inputCsv string, bytes []byte) error {
	csvD1, csvD2, err := checkCsvInputFormat(inputCsv)
	if err != nil {
		return fmt.Errorf("error while checking file: %v", err)
	}

	if err := verifyOverlap(bytes, csvD1, csvD2); err != nil {
		return fmt.Errorf("while verifying duplicated log in DB: %v", err)
	}

	Log("Valid log name")
	return nil
}

// checkCsvInputFormat verifies the input file name should be `\d+-\d+.csv`
// Rules:
// 		1. 8 digits date [$year$month$day]
//      2. Year: after 2020
//      3. Month: [1-12]
//      4. Day: [1-31]
func checkCsvInputFormat(inputCsv string) (int, int, error) {
	// Confirm the format is correct
	nameReg, _ := regexp.Compile(`(\d{4})(\d{2})(\d{2})-(\d{4})(\d{2})(\d{2}).csv`)
	matches := nameReg.FindStringSubmatch(inputCsv)
	if len(matches) < 7 {
		return 0, 0, errors.New("csv naming format error: it should be [nums]-[nums].csv")
	}

	layout := "2006-01-02T15:04:05Z"
	date1 := fmt.Sprintf("%s-%s-%sT00:00:00Z", matches[1], matches[2], matches[3])
	date2 := fmt.Sprintf("%s-%s-%sT00:00:00Z", matches[4], matches[5], matches[6])
	t1, err := time.Parse(layout, date1)
	if err != nil {
		return 0, 0, errors.New("failed to parse time")
	}
	t2, err := time.Parse(layout, date2)
	if err != nil {
		return 0, 0, errors.New("failed to parse time")
	}
	if t2.Sub(t1) < 0 {
		return 0, 0, fmt.Errorf("%s time format wrong, the former [%s] should be eariler than [%s]", inputCsv, strings.Join(matches[1:4], "-"), strings.Join(matches[4:7], "-"))
	}
	csvDate1, _ := strconv.Atoi(strings.Join(matches[1:4], ""))
	csvDate2, _ := strconv.Atoi(strings.Join(matches[4:7], ""))
	return csvDate1, csvDate2, nil
}

// verifyOverlap would verify the same period of data would not be duplicated.
// It would be verified by the name of csv file.
func verifyOverlap(bytes []byte, csvD1, csvD2 int) error {
	dateIntSlice := []int{}
	stainlessFolder := strings.Split(string(bytes), "\n")

	for index, name := range stainlessFolder {
		if index == len(stainlessFolder)-1 {
			break
		}

		csvNameSlice := strings.Split(name, "-")

		date1, err := strconv.Atoi(csvNameSlice[0])
		if err != nil {
			log.Println("Noise data is imported. Check the stainless folder manually")
			continue
		}

		date2, err2 := strconv.Atoi(csvNameSlice[1][0 : len(csvNameSlice[1])-4])
		if err2 != nil {
			log.Println("Noise data is imported. Check the stainless folder manually")
			continue
		}
		dateIntSlice = append(dateIntSlice, date1)
		dateIntSlice = append(dateIntSlice, date2)
	}

	sortDateSlice := sort.IntSlice(dateIntSlice)
	sort.Sort(sortDateSlice)

	log.Println("Current log: ", sortDateSlice)
	csvD1Index := sort.Search(len(sortDateSlice), func(i int) bool {
		return sortDateSlice[i] > csvD1
	})

	csvD2Index := sort.Search(len(sortDateSlice), func(i int) bool {
		return sortDateSlice[i] > csvD2
	})

	if sortDateSlice[0] == csvD1 {
		return fmt.Errorf("[%d] has appeared in stainless folder", csvD1)
	} else if csvD1Index-1 >= 0 && sortDateSlice[csvD1Index-1] == csvD1 {
		return fmt.Errorf("[%d] has appeared in stainless folder", csvD1)
	}
	if sortDateSlice[0] == csvD2 {
		return fmt.Errorf("[%d] has appeared in stainless folder", csvD2)
	} else if csvD2Index-1 >= 0 && sortDateSlice[csvD2Index-1] == csvD2 {
		return fmt.Errorf("[%d] has appeared in stainless folder", csvD2)
	}

	if csvD1Index != csvD2Index {
		if csvD1Index%2 == 0 {
			return fmt.Errorf("[ %d-%d ] has overlapped the data", sortDateSlice[csvD1Index], sortDateSlice[csvD1Index+1])
		} else {
			return fmt.Errorf("[ %d-%d ] has overlapped the data", sortDateSlice[csvD1Index-1], sortDateSlice[csvD1Index])
		}

	} else {
		if csvD1Index%2 == 1 {
			return fmt.Errorf("[ %d-%d ] has overlapped the data", sortDateSlice[csvD1Index-1], sortDateSlice[csvD1Index])
		}
	}
	return nil
}
