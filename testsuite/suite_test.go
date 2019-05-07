package testsuite

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/enriqueChen/stuffer/anareva"
	"github.com/enriqueChen/stuffer/interfaces"
	"github.com/enriqueChen/stuffer/status"
	logging "github.com/op/go-logging"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	fmt.Println("start  test......")
	retCode := m.Run()
	fmt.Println("end  test......")
	os.Exit(retCode)
}

func TestLoadCases(t *testing.T) {
	InitVar()
	var casefile *anareva.CaseFile
	casefile = &anareva.CaseFile{FilePath: "../testcases/testcase1.toml"}
	err := LoadCases([]interfaces.ICaseFile{casefile}, "")
	assert.Nil(t, err, err)
	assert.EqualValues(t, len(tCases), 3)
	//with filter, only get one case at last
	InitVar()
	err = LoadCases([]interfaces.ICaseFile{casefile}, "testcaseName1")
	assert.Nil(t, err, err)
	assert.EqualValues(t, len(tCases), 1)
}

//just for test
func InitVar() {
	tEnvirons = make([]interfaces.IEnvironment, 0)
	tCases = make([]interfaces.ICase, 0)
	tResults = make(map[string]interfaces.ICaseResult, 0)
	tLogger = &logging.Logger{}
}

func TestWriteJUnitReport(t *testing.T) {

	InitVar()
	fmt.Println(len(tResults))
	var casefile *anareva.CaseFile
	casefile = &anareva.CaseFile{FilePath: "../testcases/testcase1.toml"}
	casefile2 := &anareva.CaseFile{FilePath: "../testcases/testcase2.toml"}
	err := LoadCases([]interfaces.ICaseFile{casefile}, "")
	assert.Nil(t, err, err)
	err = LoadCases([]interfaces.ICaseFile{casefile2}, "")
	assert.Nil(t, err, err)
	fmt.Println(len(tResults))
	tCases[0].SetStatus(status.CasePassed)
	tResults[tCases[0].Identify()].SetDuration(10 * time.Second)
	tCases[1].SetStatus(status.CaseFailed)
	tResults[tCases[1].Identify()].SaveResult("test2result fail22222")
	tCases[2].SetStatus(status.CaseCrashed)
	tResults[tCases[2].Identify()].SaveError(fmt.Errorf("test3 error"))
	tResults[tCases[2].Identify()].SetDuration(30 * time.Second)
	tCases[3].SetStatus(status.CaseTimeout)
	f, _ := os.OpenFile("../testResult.xml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	err = WriteJUnitReport(f)
	assert.Nil(t, err, err)
}
