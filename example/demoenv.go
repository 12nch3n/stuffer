package example

import (
	"fmt"
	"sync"
	"time"

	"github.com/enriqueChen/stuffer/interfaces"
	"github.com/enriqueChen/stuffer/status"
)

// MaxCapacity defines the max capacity of parallel environments
const MaxCapacity = 3

// InitEnvirons implements to produce a executor for Analytics Report Validation
func InitEnvirons(envInfos ...string) (ret []interfaces.IEnvironment, err error) {
	ret = append(ret, &DemoEnv{m: &sync.Mutex{}, String: "Env No.1"})
	ret = append(ret, &DemoEnv{m: &sync.Mutex{}, String: "Env No.2"})
	ret = append(ret, &DemoEnv{m: &sync.Mutex{}, String: "Env No.3"})
	return
}

// The DemoEnv for base environment & new environment
type DemoEnv struct {
	m      *sync.Mutex
	Status status.EnvStatus
	String string
}

// NewDemoEnv initialize a testing environment for analytics reports validation testing with base & new service executor
func NewDemoEnv(base string, new string) (environ *DemoEnv, err error) {
	return
}

// Identify implements to get environment uniq description by string
func (de *DemoEnv) Identify() string {
	return fmt.Sprintf("demo environment: %s", de.String)
}

// SetStatus implements to set environment status
func (de *DemoEnv) SetStatus(status status.EnvStatus) {
	de.m.Lock()
	defer de.m.Unlock()
	de.Status = status
}

// GetStatus implements to get environment status
func (de *DemoEnv) GetStatus() status.EnvStatus {
	return de.Status
}

// Restore implements to restore a crashed environment for reusing
func (de *DemoEnv) Restore() (err error) {
	return
}

// Setup defines the interface to setup the test case, will raise error
func (de *DemoEnv) Setup(tCase interfaces.ICase) (err error) {
	return
}

// RunCase implement the test case run for analytics reports validation testing
func (de *DemoEnv) RunCase(tCase interfaces.ICase, tRes interfaces.ICaseResult) (err error) {
	time.Sleep(5 * time.Second)
	err = tCase.Assert(tRes)
	return
}

// Teardown defines the interface to teardown the test case, will raise error
func (de *DemoEnv) Teardown(tCase interfaces.ICase) (err error) {
	return
}
