package testsuite

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/enriqueChen/stuffer/interfaces"
	"github.com/enriqueChen/stuffer/status"
)

// WaitEnvTimeout case Runner wait for test environment timeout
var WaitEnvTimeout = 1 * time.Minute

// TestingQueue for testing queues
type TestingQueue struct {
	execQueue    chan interfaces.ICaseResult
	testingCount int64
	runMode      string
	m            *sync.WaitGroup
	once         sync.Once
}

// LaunchTesting by feed test cases into testing queue
func LaunchTesting(runMode string, exit chan<- struct{}) (q *TestingQueue, err error) {
	q = &TestingQueue{
		execQueue: make(chan interfaces.ICaseResult, len(tEnvirons)), //? size is the environment number to preload testcase in channel & wait for retry case
		runMode:   runMode,
		m:         &sync.WaitGroup{},
	}
	retryChan := make(chan interfaces.ICaseResult, len(tResults)) //? size is all case number would not block retry channel pushing
	go q.load(retryChan)
	go q.run(retryChan, exit)
	return
}

// load feeds the all tests into execQueue and wait for cases need to retry
func (q *TestingQueue) load(rcvRetry <-chan interfaces.ICaseResult) {
	go func() {
		defer func() {
			if recover() != nil {
				tLogger.Warningf("Test execution channel closed")
			}
		}()

		for _, t := range tResults {
			q.execQueue <- t
			atomic.AddInt64(&q.testingCount, 1)
			tLogger.Noticef("Sent test case to execute channel, Case %s", t.GetCase().Identify())
		}
	}()

	go func() {
		defer func() {
			if recover() != nil {
				tLogger.Warningf("Test execution channel closed")
			}
		}()

		for {
			t, ok := <-rcvRetry
			if !ok {
				tLogger.Warning("Retry channel closed, retrier about to exit")
				return
			}
			q.execQueue <- t
			tLogger.Noticef("Sent retry test case to execute channel, Case %s", t.GetCase().Identify())
		}
	}()
}

func (q *TestingQueue) run(sendRetry chan<- interfaces.ICaseResult, exit chan<- struct{}) {

	executing := func() {
		for {
			t, testOk := <-q.execQueue
			if !testOk {
				tLogger.Warning("test execution queue closed, test execution work about to exit.")
				q.m.Done()
				return
			}
			e, envOk := <-tEnvPool.Wait4Using //! Blocked if no avaliable environments
			if !envOk {
				tLogger.Warning("environment for using queue closed, test execution work about to exit.")
				q.m.Done()
				return
			}
			err := runcase(e, t)
			tEnvPool.Send2Return(e) //! return Environment after running
			if q.testdone(t, err, sendRetry) {
				q.once.Do(func() {
					close(q.execQueue)
				})
			}

		}
	}

	for i := 0; i < len(tEnvirons); i++ {
		q.m.Add(1)
		go executing()
	}
	q.m.Wait()
	tEnvPool.Close()
	exit <- struct{}{}
}

// TestDone implements to mark a test case done and check tests if all test done
func (q *TestingQueue) testdone(r interfaces.ICaseResult, err error, sndRetry chan<- interfaces.ICaseResult) bool {
	c := r.GetCase()
	if _, ok := err.(interfaces.ErrCaseTimeout); ok {
		tLogger.Errorf("Testing timeout, case:%s", c.Identify())
		c.SetStatus(status.CaseTimeout)
	}

	if c.Status() != status.CasePassed {
		if q.runMode == "QUICK_FAIL" {
			//! Terminate the testing due to case failed.
			tLogger.Errorf("Testing Quick failed.")
			return true
		} else if q.runMode == "RETRY" && r.RetriedNum() <= c.MaxRetry() {
			//! Send test case to retry
			sndRetry <- r
			tLogger.Errorf("Test case retry case: %s, retry time: %d, max retry: %d",
				c.Identify(), r.RetriedNum(), c.MaxRetry())
			r.Retried()
			return false
		}
		tLogger.Errorf("Test case FAILED, Case: %s, retry time: %d", c.Identify(), r.RetriedNum())
	} else {
		tLogger.Infof("Test case PASSED, Case: %s, retry time: %d", c.Identify(), r.RetriedNum())
	}

	atomic.AddInt64(&q.testingCount, -1)
	tLogger.Noticef("Remaining test case number: %d", atomic.LoadInt64(&q.testingCount))
	return atomic.LoadInt64(&q.testingCount) <= int64(0)
}
