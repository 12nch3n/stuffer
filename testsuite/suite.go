package testsuite

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/beevik/etree"
	logging "github.com/op/go-logging"
	"github.com/enriqueChen/ssgstuffer/interfaces"
	"github.com/enriqueChen/ssgstuffer/utils/junit"
)

var (
	tEnvirons []interfaces.IEnvironment
	tCases    []interfaces.ICase
	tResults  map[string]interfaces.ICaseResult

	tLogger  *logging.Logger
	tEnvPool *EnvironPool
	tQueue   *TestingQueue
)

// InitLogger initialize a log handler for regression
func InitLogger(testName string, level logging.Level) (*logging.Logger, error) {
	tLogger = logging.MustGetLogger(testName)
	format4print := logging.MustStringFormatter(`%{color}%{module}-%{pid} %{time:15:04:05} %{id:03x} [%{level:.4s}]%{color:reset} %{shortpkg}:%{callpath}: %{message}`)
	format4file := logging.MustStringFormatter(`%{module}-%{pid} %{time:15:04:05} %{id:03x} [%{level:.4s}] %{shortpkg}:%{callpath}: %{message}`)
	fName := fmt.Sprintf("%s-%s.log", testName, time.Now().Format("20060102150405"))
	var f *os.File
	var err error
	if _, err = os.Stat(fName); os.IsNotExist(err) {
		f, err = os.Create(fName)
	} else {
		f, err = os.OpenFile(fName, os.O_WRONLY|os.O_APPEND, 0600)
	}
	if err != nil {
		return nil, err
	}
	logf := logging.AddModuleLevel(
		logging.NewBackendFormatter(
			logging.NewLogBackend(f, "", 0),
			format4file),
	)

	std := logging.AddModuleLevel(
		logging.NewBackendFormatter(
			logging.NewLogBackend(os.Stdout, "", 0),
			format4print),
	)

	logf.SetLevel(level, "")
	std.SetLevel(level, "")
	logging.SetBackend(std, logf)
	tLogger.Infof("Initialize the log for %s in file %s", testName, fName)
	return tLogger, nil
}

// RegisterEnvironments register the initialized test environ pool into testing module
func RegisterEnvironments(environs []interfaces.IEnvironment) {
	tEnvirons = append(tEnvirons, environs...)
}

// LoadCases loads cases and test results from initialized case files with filter
func LoadCases(caseFiles []interfaces.ICaseFile, filterStr string) (err error) {

	tResults = make(map[string]interfaces.ICaseResult)
	caseNameFilter, err := regexp.Compile(filterStr)
	if err != nil {
		tLogger.Errorf("compile case regex pattern failed: %s", err.Error())
		return err
	}
	for _, cFile := range caseFiles {
		loadedCases, err := cFile.Load()
		if err != nil {
			tLogger.Warningf("Load case file %s failed with err %s.",
				cFile.Path(), err.Error())
		}
		for _, tCase := range loadedCases {
			tLogger.Debugf("Filter test case with regex pattern Case: %s", tCase.Identify())
			if caseNameFilter.Find([]byte(tCase.Identify())) == nil {
				tLogger.Infof("Ignore test case because test case file and name %s not matched filter pattern.",
					tCase.Identify())
				continue
			}
			if _, ok := tResults[tCase.Identify()]; ok {
				tLogger.Warningf("Duplicated test case id: [%s] Ignore this test case %v.",
					tCase.Identify(), tCases)
				continue
			}
			tLogger.Debugf("Add initialized test result into framework(%d), Case:%s", len(tResults), tCase.Identify())
			tResults[tCase.Identify()] = tCase.InitResult()
			tCases = append(tCases, tCase)
		}
	}
	for _, tCase := range tCases {
		tLogger.Infof("Collect test case [%s]: %s", tCase.Identify(), tCase.Describe())
	}
	return
}

// RunCases run the all cases in testing model
func RunCases(runMode string) (err error) {
	tLogger.Infof("Framework start to run test cases in mode: %s", runMode)
	envpoolExit := make(chan struct{})
	tEnvPool, _ = NewEnvironPool(envpoolExit)
	runningExit := make(chan struct{})
	tQueue, _ = LaunchTesting(runMode, runningExit)
	<-runningExit
	<-envpoolExit
	tLogger.Info("Test execution completed.")
	return
}

// genrateEnvironSummary generate testing summary as human-readable string
func genrateEnvironSummary() (summary []string) {
	summary = append(summary,
		fmt.Sprintf("Test Environments Capacity (%d)", len(tEnvirons)))
	for _, e := range tEnvirons {
		summary = append(summary, e.Identify())
	}
	return
}

// WriteReport generate test report from the testing  results
func WriteReport(rFile io.Writer) (err error) {
	fmt.Fprintf(rFile, "Test Report")
	return
}

// WriteJUnitReport generate JUnit formate test report from the testing  results
func WriteJUnitReport(rFile io.Writer) (err error) {

	testss := &junit.TestsuiteS{}
	testsuiteMap := make(map[string]*junit.Testsuite)

	//testcase info
	id := 0
	for _, result := range tResults {
		tc := result.GetCase()
		if _, ok := testsuiteMap[tc.FileName()]; !ok { // one junit.Testsuite
			testsuiteMap[tc.FileName()] = &junit.Testsuite{
				Name:     tc.FileName(),
				ID:       id,
				Testcase: make([]*etree.Element, 0),
			}
			id++
		}
		testsuite := testsuiteMap[tc.FileName()]
		junit.CheckinTestResult(testss, testsuite, result)
	}

	//write XML
	doc := etree.NewDocument()
	//head
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)

	testsuites := doc.CreateElement("testsuites")
	testsuites.CreateAttr("errors", strconv.Itoa(testss.Errors))
	testsuites.CreateAttr("failures", strconv.Itoa(testss.Failures))
	testsuites.CreateAttr("name", "")
	testsuites.CreateAttr("tests", strconv.Itoa(testss.Testscount))
	testsuites.CreateAttr("time", strconv.FormatFloat(testss.Time, 'f', 0, 64))

	for _, testsuite := range testsuiteMap {
		onetestsuite := testsuites.CreateElement("testsuite")

		onetestsuite.CreateAttr("name", testsuite.Name)
		onetestsuite.CreateAttr("tests", strconv.Itoa(testsuite.Testscount))
		onetestsuite.CreateAttr("errors", strconv.Itoa(testsuite.Errors))
		onetestsuite.CreateAttr("failures", strconv.Itoa(testsuite.Failures))
		onetestsuite.CreateAttr("id", strconv.Itoa(testsuite.ID))
		onetestsuite.CreateAttr("time", strconv.FormatFloat(testsuite.Time, 'f', 0, 64))

		for _, testcaseElement := range testsuite.Testcase {
			onetestsuite.AddChild(testcaseElement)
		}
		//environment info
		properties := onetestsuite.CreateElement("properties")
		oneproperty := properties.CreateElement("property")
		oneproperty.CreateAttr("name", "BaseEnviron")
		oneproperty.CreateAttr("value", "123456")
		anotherproperty := properties.CreateElement("property")
		anotherproperty.CreateAttr("name", "NewEnviron")
		anotherproperty.CreateAttr("value", "7890")
	}

	doc.Indent(2)
	doc.WriteTo(rFile)
	return
}
