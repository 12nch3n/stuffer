package testsuite

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/enriqueChen/ssgstuffer/interfaces"
	"github.com/enriqueChen/ssgstuffer/status"
)

// EnvironPool provide the test environments
type EnvironPool struct {
	crashedNum  int32
	Wait4Using  chan interfaces.IEnvironment
	Wait4Return chan interfaces.IEnvironment
	m           *sync.WaitGroup
	once        sync.Once
}

// NewEnvironPool load the all environments into Wait4Using before testing
func NewEnvironPool(exit chan<- struct{}) (*EnvironPool, error) {
	if len(tEnvirons) < 1 {
		tLogger.Errorf("Assign environment to pool failed, no environment")
		return nil, errors.New("Environments resources is empty")
	}
	tLogger.Debugf("Initialize the environment pool")
	p := EnvironPool{
		Wait4Using:  make(chan interfaces.IEnvironment, len(tEnvirons)),
		Wait4Return: make(chan interfaces.IEnvironment, len(tEnvirons)),
		m:           &sync.WaitGroup{},
	}
	go p.run(exit)
	return &p, nil
}

// run send the available environment to chan wait4using and listening channel wait4return to recycle environment as provider
func (p *EnvironPool) run(exit chan<- struct{}) {

	tLogger.Debugf("Pushing environments to channel Wait4Using")
	for _, v := range tEnvirons {
		if v.GetStatus() == status.EnvAvailable {
			tLogger.Debugf("release env %s", v.Identify())
			p.Wait4Using <- v
		}
	}

	recycle := func() {
		for {
			v, ok := <-p.Wait4Return
			if !ok {
				tLogger.Warningf("Wait4Return Environment channel closed, recycle worker about to exit.")
				p.m.Done()
				return
			}
			if v.GetStatus() == status.EnvCrashed || v.GetStatus() == status.EnvUnknown {
				if err := v.Restore(); err != nil {
					atomic.AddInt32(&p.crashedNum, 1)
				}
				if atomic.LoadInt32(&p.crashedNum) == int32(len(tEnvirons)) {
					//! All environment crashed exit env pool
					p.Close()
				}
			}
			v.SetStatus(status.EnvAvailable)
			p.Wait4Using <- v
			tLogger.Debugf("return env %s", v.Identify())
		}
	}

	tLogger.Debugf("All test environment loaded to run.")
	for i := 0; i < len(tEnvirons); i++ {
		p.m.Add(1)
		go recycle()
	}
	p.m.Wait()
	close(p.Wait4Using)
	tLogger.Debugf("Env pool terminated")
	exit <- struct{}{}
}

// Close implements to stop environment recycling
func (p *EnvironPool) Close() {
	p.once.Do(func() {
		close(p.Wait4Return)
	})
}

// Send2Return send environment in waiting channel return if channel closed
func (p *EnvironPool) Send2Return(e interfaces.IEnvironment) (closed bool) {
	defer func() {
		if recover() != nil {
			closed = true
		}
	}()
	p.Wait4Return <- e
	closed = false
	return
}
