package interfaces

import (
	"fmt"
	"time"

	"github.com/enriqueChen/stuffer/status"
)

// ICaseFile interface to manager test case list````	`
type ICaseFile interface {
	// Path defines an interface for get case file path
	Path() string
	// LoadCases defines an interface for loading testcase
	Load() ([]ICase, error)
}

// ICase interface defines Test Case
type ICase interface {
	// Initialize defines an interface for a customized initialization from some case metadata for your test case
	InitResult() ICaseResult
	// Identify defines an interface for get case unique identify
	Identify() string
	// Name defines an interface for getting test case name
	Name() string
	// Describe defines an interface for getting Case reduced informantion in one line string.
	Describe() string
	//Asserts defines an interface for testing asserter after test case running
	Assert(result ICaseResult) error
	//RuntimeMetadata defines an interface for getting runtime meta to execute test
	RuntimeMetadata() interface{}
	// SetStatus defines an interface for setting the case status
	SetStatus(status status.CaseStatus)
	//Status defines an interface for getting case execution status
	Status() status.CaseStatus
	// Timeout defines an interface for getting testing execution timeout
	Timeout() time.Duration
	// MaxRetry defines an interface to get retry number of the test case
	MaxRetry() int
	// FileName defines an interface to get case file name
	FileName() string
}

//ICaseResult defines the interface for test case running
type ICaseResult interface {
	// SaveResult defines an interface to collect test case failures to fill the test result
	SaveResult(res string)
	// GetCase get test case for corrent test result
	GetCase() ICase
	//CheckIn defines an interface to assgin the test result with an environment
	CheckIn(index int)
	// SetDuration defines an interface to set test duration in test result
	SetDuration(duration time.Duration)
	// Duration defines an interface to get case duration if test case status is Timeout Crashed Failed Success
	Duration() time.Duration
	// ResultDetails defines an interface to get case running details collected during test case run
	ResultDetails() string
	// Retry defines an interface to set retry time + 1
	Retried()
	// RetriedNum defines an interface to get retry number of the test case
	RetriedNum() int
	// EnvID defines an interface to get the Assigned environment ID
	EnvID() int
	// FailureSummary defines an interface to get FailureSummary
	ErrorMessage() error
	// SaveError defines an interface to collect error message
	SaveError(err error)
}

// ErrCaseTimeout defines for run test case time out
type ErrCaseTimeout struct {
	caseID string
}

// NewErrCaseTimeout produce a new run test case time out error
func NewErrCaseTimeout(caseID string) ErrCaseTimeout {
	return ErrCaseTimeout{
		caseID: caseID,
	}
}

// Error implements interface for error
func (e ErrCaseTimeout) Error() string {
	return fmt.Sprintf("case <%s> Execution timeout", e.caseID)
}
