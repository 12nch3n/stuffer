package main

import (
	"fmt"
	"github.com/op/go-logging"
	"os"
	"time"

	"github.com/enriqueChen/stuffer/example"
	"github.com/enriqueChen/stuffer/interfaces"
	"github.com/enriqueChen/stuffer/testsuite"
)

var tLog *logging.Logger

func main() {
	tLog, _ := testsuite.InitLogger("DEMO", logging.DEBUG)
	tLog.Info("1st step -> init the logger with package: github.com/op/go-logging")

	tLog.Info("2nd step -> Load test cases from test files.")
	testFiles := []interfaces.ICaseFile{
		&example.DemoFile{FilePath: "DemoFile1"},
		&example.DemoFile{FilePath: "DemoFile2"},
	}
	testsuite.LoadCases(testFiles, `.*`)

	tLog.Info("3rd step -> Register Environments into test environment.")
	envs, _ := example.InitEnvirons("", "", "")
	testsuite.RegisterEnvironments(envs)

	tLog.Info("4th step -> paralleled run test case.")
	testsuite.RunCases("RETRY")

	tLog.Info("5th step -> generate junit formated test report")
	fName := fmt.Sprintf("%s_REPORT-%s.xml", "DEMO", time.Now().Format("20060102150405"))
	var f *os.File
	var err error
	if _, err = os.Stat(fName); os.IsNotExist(err) {
		f, err = os.Create(fName)
	} else {
		f, err = os.OpenFile(fName, os.O_WRONLY, 0600)
	}
	testsuite.WriteJUnitReport(f)
}
