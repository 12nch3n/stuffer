package status

//EnvStatus defines the test environment runtime status
type EnvStatus int8

const (
	// EnvAvailable means the test environment is ready to run case
	EnvAvailable EnvStatus = iota
	// EnvBusy means the test environment is busy with processing
	EnvBusy
	// EnvCrashed means the test environment is crushed and could not be used
	EnvCrashed
	// EnvUnknown means the test environment status is unknown
	EnvUnknown
)

func (s EnvStatus) String() string {
	switch s {
	case EnvAvailable:
		return "AVAILABLE"
	case EnvBusy:
		return "BUSY"
	case EnvCrashed:
		return "CRASHED"
	case EnvUnknown:
		return "UNKNOWN"
	default:
		return "UNKNOWN"
	}
}

//CaseStatus defines the test case runtime status
type CaseStatus int8

const (
	// CaseNotrun is the default test case is not runned
	CaseNotrun CaseStatus = iota
	// CaseRunning means test case is Running
	CaseRunning
	// CaseTimeout means the test case running is terminated by framework due to timeout
	CaseTimeout
	// CaseCrashed means the test case running got error
	CaseCrashed
	// CaseFailed means the test case completed but got a failure
	CaseFailed
	// CasePassed mean the test case completed with test passed
	CasePassed
)

func (s CaseStatus) String() string {
	switch s {
	case CaseNotrun:
		return "NOTRUN"
	case CaseRunning:
		return "RUNNING"
	case CaseTimeout:
		return "TIMEOUT"
	case CaseCrashed:
		return "CRASHED"
	case CaseFailed:
		return "FAILED"
	case CasePassed:
		return "PASSED"
	default:
		return "UNKNOWN"
	}
}
