package subprocess

import (
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"sync"
	"time"
)

var (
	DefaultKillTimeout = 3 * time.Second
)

type Subprocess struct {
	*exec.Cmd

	KillTimeout  time.Duration
	ProxySignals bool

	signals chan os.Signal
	active  *exec.Cmd
	cond    *sync.Cond
}

func NewSubprocess(name string, arg ...string) *Subprocess {
	return &Subprocess{
		Cmd:          exec.Command(name, arg...),
		KillTimeout:  DefaultKillTimeout,
		ProxySignals: true,
		cond:         sync.NewCond(&sync.Mutex{}),
	}
}

func (s *Subprocess) isRunning() bool {
	return s.active != nil && s.active.Process != nil
}

func (s *Subprocess) Stop() {
	s.cond.L.Lock()
	if !s.isRunning() {
		s.cond.Wait()
	}
	s.cond.L.Unlock()

	if s.ProxySignals {
		signal.Stop(s.signals)
	}

	done := make(chan bool)
	go func() {
		s.active.Wait()
		close(done)
	}()

	if runtime.GOOS == "windows" {
		s.active.Process.Kill()
	} else {
		s.active.Process.Signal(os.Interrupt)
	}

	select {
	case <-time.After(s.KillTimeout):
		s.active.Process.Kill()
	case <-done:
		return
	}
	<-done
}

func (s *Subprocess) Serve() {
	s.cond.L.Lock()

	if s.ProxySignals {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)
		go func() {
			defer signal.Stop(signals)
			s.cond.L.Lock()
			if !s.isRunning() {
				s.cond.Wait()
			}
			s.cond.L.Unlock()
			for sig := range signals {
				if !s.isRunning() {
					return
				}
				s.active.Process.Signal(sig)
			}
		}()
		s.signals = signals
	}

	s.active = &exec.Cmd{
		Path:   s.Path,
		Args:   s.Args,
		Env:    s.Env,
		Dir:    s.Dir,
		Stdin:  s.Stdin,
		Stdout: s.Stdout,
		Stderr: s.Stderr,
	}
	s.active.Start()
	s.cond.L.Unlock()
	s.cond.Broadcast()

	s.active.Wait()
	close(s.signals)
	s.cond.L.Lock()
	s.active = nil
	s.cond.L.Unlock()
}
