package example

import (
	"fmt"
	"time"

	"github.com/enriqueChen/stuffer/common"
	"github.com/enriqueChen/stuffer/interfaces"
	"github.com/enriqueChen/stuffer/status"
)

const (
	// DefaultCaseTimeout defines Demo case timeout Demo
	DefaultCaseTimeout = 10 * time.Second
	// DefaultRetry defines Demo case Max retry
	DefaultRetry = 3
)

// DemoFile Fake file
type DemoFile struct {
	FilePath string //Relative path
}

// Path implements to get case file path
func (file *DemoFile) Path() string {
	return file.FilePath
}

// Load implements to load test case from specified case files path and regex pattern filter
func (file *DemoFile) Load() (caseList []interfaces.ICase, err error) {
	//! Add readfile and case Initialize
	timeout := 0
	for _, n := range []string{"Test1", "Test2", "Test3", "Test4", "Test5", "Test6"} {
		caseList = append(caseList, &Case{
			name:        n,
			fileName:    file.FilePath,
			timeout:     time.Duration(5+timeout) * time.Second,
			maxRetry:    DefaultRetry,
			description: fmt.Sprintf("This only a demo case %s in %s", n, file.FilePath),
		})
		timeout++
	}
	return
}

// Case Implements the ICase to run analytics report validation
type Case struct {
	name        string
	description string
	fileName    string
	status      status.CaseStatus
	runtime     CaseRuntimeMetadata
	timeout     time.Duration
	maxRetry    int
}

// CaseRuntimeMetadata is the runtime metadata for Case
type CaseRuntimeMetadata struct {
}

// Identify implements to get case unique identify during testing
func (aCase *Case) Identify() string {
	return fmt.Sprintf("%s:%s", aCase.fileName, aCase.name)
}

// Name implements getting test case name
func (aCase *Case) Name() (name string) {
	name = aCase.name
	return
}

// Describe implements getting Case reduced informantion in one line string.
func (aCase *Case) Describe() (desc string) {
	desc = aCase.description
	return
}

//Assert implements assert test result after test case running
func (aCase *Case) Assert(result interfaces.ICaseResult) (err error) {
	aCase.SetStatus(status.CasePassed)
	result.SaveResult("Demo result details")
	return
}

//RuntimeMetadata implements getting runtime meta to execute test
func (aCase *Case) RuntimeMetadata() (metadata interface{}) {
	//metadata = interface{}(aCase.RuntimeMetadata)
	return aCase.runtime
}

// SetStatus implements  setting the case status
func (aCase *Case) SetStatus(status status.CaseStatus) {
	aCase.status = status
}

// Status implements to get case execution status
func (aCase *Case) Status() (status status.CaseStatus) {
	status = aCase.status
	return
}

// Timeout implements to get testing execution timeout
func (aCase *Case) Timeout() (timeout time.Duration) {
	if aCase.timeout <= 0 {
		timeout = DefaultCaseTimeout
	}
	timeout = aCase.timeout
	return
}

// InitResult implements a customized initialization from some case metadata for your test case
func (aCase *Case) InitResult() (ret interfaces.ICaseResult) {
	aCase.SetStatus(status.CaseNotrun)
	return &common.Result{
		Case: aCase,
	}
}

// MaxRetry implements to get test case Max retry number
func (aCase *Case) MaxRetry() (retry int) {
	if aCase.maxRetry <= 0 {
		retry = DefaultRetry
	}
	retry = aCase.maxRetry
	return
}

// FileName implements to get test case's filename
func (aCase *Case) FileName() (filename string) {
	return aCase.fileName
}
