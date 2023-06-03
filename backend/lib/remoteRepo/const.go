package remoterepo

type TestHausHeader string

const (
	Suite             TestHausHeader = "Suite"
	Board                            = "Board"
	Model                            = "Model"
	Test                             = "Test"
	Status                           = "Status"
	FailureReason                    = "Failure Reason"
	StartedTime                      = "Started Time (UTC)"
	Duration                         = "Duration (s)"
	BuildVersion                     = "Build Version"
	FirmwareROVersion                = "Firmware RO Version"
	FirmwareRWVersion                = "Firmware RW Version"
	Hostname                         = "Hostname"
)

const timeLayout = "2006-01-02 15:04:05"
