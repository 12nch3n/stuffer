package interfaces

import "github.com/enriqueChen/ssgstuffer/status"

// IEnvironment defines testing environment interface
type IEnvironment interface {
	// Identify defines interface to get an unique identify by string
	Identify() string
	// SetStatus defines an interface to set the environment specified by index
	SetStatus(status status.EnvStatus)
	// GetStatus defines an interface to get the environment specified by index
	GetStatus() status.EnvStatus
	// Restore defines an interface to try restore a crashed environment to available
	Restore() error
	// Setup defines an interface to setup the test case, will raise error
	Setup(tCase ICase) error
	// RunCase defines an interface to run the test case, will fill the test result and raise error
	RunCase(tCase ICase, tRes ICaseResult) error
	// Teardown defines an interface to teardown the test case, will raise error
	Teardown(tCase ICase) error
}
