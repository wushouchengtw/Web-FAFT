package models

import "time"

type Result struct {
	Id                string    `db:"result_id"`
	Time              time.Time `db:"time"`
	Duration          float64   `db:"duration"`
	Suite             string    `db:"suite"`
	DutId             string    `db:"dut_id"`
	Milestone         string    `db:"milestone"`
	Version           string    `db:"version"`
	Host              string    `db:"host"`
	TestId            string    `db:"test_id"`
	Status            bool      `db:"status"`
	Reason            string    `db:"reason"`
	FirmwareROVersion string    `db:"firmware_RO_Version"`
	FirmwareRWVersion string    `db:"firmware_RW_Version"`
}
