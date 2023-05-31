package models

import "time"

type SuiteResult struct {
	Id                int
	Time              time.Time
	Duration          float64
	Suite             string
	DutId             int
	Milestone         string
	Version           string
	Host              string
	TestId            int
	Status            bool
	Reason            string
	FirmwareROVersion string
	FirmwareRWVersion string
}
