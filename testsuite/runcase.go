package testsuite

import (
	"time"

	"github.com/enriqueChen/stuffer/interfaces"
	"github.com/enriqueChen/stuffer/status"
)

func runcase(env interfaces.IEnvironment, tRes interfaces.ICaseResult) (err error) {
	tCase := tRes.GetCase()
	var startTime time.Time
	defer func() {
		tRes.SetDuration(
			time.Now().Sub(startTime))
		if e := env.Teardown(tCase); e != nil {
			tLogger.Errorf("Test Case Teardown failed, Case:%s, Error:%s",
				tCase.Identify(), err.Error())
			env.SetStatus(status.EnvCrashed)
		}
		tLogger.Debugf("Test Case Teardown succeed, Case:%s",
			tCase.Identify())
	}()
	tCase.SetStatus(status.CaseRunning)
	err = env.Setup(tCase)
	startTime = time.Now()
	if err != nil {
		tLogger.Errorf("Test Case Setup failed case could not run, Case:%s, Error:%s",
			tCase.Identify(), err.Error())
		return
	}
	tLogger.Debugf("Test Case Setup succeed, Case:%s",
		tCase.Identify())
	c := make(chan error, 1)
	go func() { c <- env.RunCase(tCase, tRes) }()
	select {
	case err = <-c:
		break
	case <-time.After(tCase.Timeout()):
		err = interfaces.NewErrCaseTimeout(tCase.Identify())
	}
	if err != nil {
		tCase.SetStatus(status.CaseCrashed)
		//save error
		tRes.SaveError(err)
		tLogger.Errorf("Test Case Run failed, Case:%s, Error:%s",
			tCase.Identify(), err.Error())
	} else {
		tLogger.Infof("Test Case Run completed, Case:%s",
			tCase.Identify())
	}
	return
}
