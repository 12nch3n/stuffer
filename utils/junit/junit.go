package junit

import (
	"fmt"

	"github.com/beevik/etree"
	"github.com/enriqueChen/ssgstuffer/interfaces"
	"github.com/enriqueChen/ssgstuffer/status"
)

// TestsuiteS struct for xml
type TestsuiteS struct {
	Name       string
	Errors     int
	Failures   int
	Testscount int
	Time       float64 //second
}

// Testsuite struct for xml
type Testsuite struct {
	Name       string
	Errors     int
	Failures   int
	ID         int
	Testscount int
	Time       float64 //second
	Testcase   []*etree.Element
}

// CheckinTestResult add the result to junit xml schema
func CheckinTestResult(testsum *TestsuiteS,
	tsuite *Testsuite,
	result interfaces.ICaseResult) {

	tc := result.GetCase()
	// one testcase
	tcElement := etree.NewElement("testcase")
	tsuite.Testcase = append(tsuite.Testcase, tcElement) //testcase add to testsuit
	tsuite.Testscount++
	tsuite.Time += result.Duration().Seconds()
	testsum.Testscount++
	testsum.Time += result.Duration().Seconds()
	tcElement.CreateAttr("name", tc.Name())
	switch tc.Status() {
	case status.CasePassed:
		tcElement.CreateAttr("status", "CasePassed")
	case status.CaseFailed:
		tcElement.CreateAttr("status", "CaseFailed")
		failure := tcElement.CreateElement("failure")
		failure.CreateAttr("message", result.ResultDetails())
		failure.CreateAttr("type", "")
		tsuite.Failures++
		testsum.Failures++
	case status.CaseCrashed:
		tcElement.CreateAttr("status", "CaseCrashed")
		failure := tcElement.CreateElement("error")
		failure.CreateAttr("message", fmt.Sprintf("%v", result.ErrorMessage()))
		failure.CreateAttr("type", "")
		tsuite.Errors++
		testsum.Errors++
	case status.CaseTimeout: // Is Timeout a type of error?
		tcElement.CreateAttr("status", "CaseTimeout")
		failure := tcElement.CreateElement("error")
		failure.CreateAttr("message", "TIMEOUT")
		failure.CreateAttr("type", "")
		tsuite.Errors++
		testsum.Errors++
	default:
		tcElement.CreateAttr("status", "NotKnown")
	}
}
