package common

import (
	"fmt"
	"time"

	"github.com/enriqueChen/stuffer/interfaces"
)

// Result Implements the ICaseResult interface
type Result struct {
	environID    int
	retriedNum   int
	Case         interfaces.ICase
	dur          time.Duration
	result       string // diff descriptions
	errorMessage error
}

// SaveResult collect result
func (aRes *Result) SaveResult(res string) {
	aRes.result = res
	return
}

// CheckIn implements to set the environment id to test result
func (aRes *Result) CheckIn(envID int) {
	aRes.environID = envID
}

// EnvID implements to get the Assigned environment ID
func (aRes *Result) EnvID() int {
	return aRes.environID
}

// GetCase implements to get test case from test result
func (aRes *Result) GetCase() (ret interfaces.ICase) {
	return aRes.Case
}

// SetDuration implements to set duration for test case
func (aRes *Result) SetDuration(duration time.Duration) {
	aRes.dur = duration
}

//Duration interface to get case duration if test case status is Timeout Crashed Failed Success
func (aRes *Result) Duration() time.Duration {
	return aRes.dur
}

// Details interface to get case running details collected during test case run
func (aRes *Result) Details() string {
	return fmt.Sprintf("%s, Status [%s], Duration[%s], Description:%s",
		aRes.Case.Name(), aRes.Case.Status(), aRes.dur.String(), aRes.Case.Describe())
}

// Retried implements to add retried number in test result
func (aRes *Result) Retried() {
	aRes.retriedNum++
}

// RetriedNum implements to get retried times by int
func (aRes *Result) RetriedNum() int {
	return aRes.retriedNum
}

// ErrorMessage implements to save Error message
func (aRes *Result) ErrorMessage() error {
	return aRes.errorMessage
}

// SaveError implements to get saved error
func (aRes *Result) SaveError(err error) {
	aRes.errorMessage = err
	return
}

// ResultDetails implements to get saved result
func (aRes *Result) ResultDetails() string {
	return aRes.result
}
