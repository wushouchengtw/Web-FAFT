package csv

type TestHausHeader string

const (
	Suite             TestHausHeader = "Suite"
	Board             TestHausHeader = "Board"
	Model             TestHausHeader = "Model"
	Test              TestHausHeader = "Test"
	Status            TestHausHeader = "Status"
	FailureReason     TestHausHeader = "Failure Reason"
	StartedTime       TestHausHeader = "Started Time (UTC)"
	Duration          TestHausHeader = "Duration (s)"
	BuildVersion      TestHausHeader = "Build Version"
	FirmwareROVersion TestHausHeader = "Firmware RO Version"
	FirmwareRWVersion TestHausHeader = "Firmware RW Version"
	Hostname          TestHausHeader = "Hostname"
)

const timeLayout = "2006-01-02 15:04:05"
