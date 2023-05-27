package dbInfo

import (
	"fmt"
	"log"
)

type DbConfig struct {
	UserName string `json:"dbUserName"`
	Password string `json:"dbPassword"`
	Ip       string `json:"dbIp"`
	Name     string `json:"dbName"`
	Table    string `json:"tableName"`
}

type LocalResult struct {
	Id         int    `json:"id"`
	Time       string `json:"time"`
	Tester     string `json:"tester"`
	Name       string `json:"name"`
	Board      string `json:"board"`
	Model      string `json:"model"`
	Version    string `json:"version"`
	LogPath    string `json:"logPath"`
	PassOrFail string `json:"passOrFail"`
	Reason     string `json:"reason"`
}

type StainlessResult struct {
	Id                int    `json:"id"`
	Time              string `json:"time"`
	Duration          string `json:"duration"`
	Suite             string `json:"suite"`
	Board             string `json:"board"`
	Model             string `json:"model"`
	BuildVersion      string `json:"buildVersion"`
	Host              string `json:"host"`
	TestName          string `json:"testName"`
	Status            string `json:"status"`
	Reason            string `json:"reason"`
	FirmwareROVersion string `json:"firmwareROVersion"`
	FirmwareRWVersion string `json:"firmwareRWVersion"`
}

type PassingRate struct {
	Week        int     `json:"week"`
	TestName    string  `json:"testName"`
	StartDate   string  `json:"startDate"`
	EndDate     string  `json:"endDate"`
	TotalRun    int     `json:"totalRun"`
	Pass        int     `json:"pass"`
	PassingRate float64 `json:"passingRate"`
}

type TicketList struct {
	Id       int    `json:"id"`
	CaseName string `json:"caseName"`
	Owner    string `json:"owner"`
	Ticket   string `json:"ticket"`
}

type SearchInput struct {
	StartDate string
	EndDate   string
	TableName string
	Board     string
	Reason    string
	Name      string
	Result    string
	OrderBy   string
}

func Log(value ...interface{}) {
	log.Println("\t", fmt.Sprint(value...))
}

func Logf(logString string, value ...interface{}) {
	log.Printf(logString, value...)
}
