package models

import "time"

type Result struct {
	Id                int       `db:"result_id"`
	Time              time.Time `db:"time"`
	Duration          float64   `db:"duration"`
	Suite             string    `db:"suite"`
	DutId             int       `db:"dut_id"`
	Milestone         string    `db:"milestone"`
	Version           string    `db:"version"`
	Host              string    `db:"host"`
	TestId            int       `db:"test_id"`
	Status            bool      `db:"status"`
	Reason            string    `db:"reason"`
	FirmwareROVersion string    `db:"firmware_RO_Version"`
	FirmwareRWVersion string    `db:"firmware_RW_Version"`
}
