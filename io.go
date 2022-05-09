package sonic

import (
	"io"

	"github.com/talostrading/sonic/internal"
)

type IO struct {
	poller *internal.Poller

	// inflightTimers prevents the PollData owned by the Timer to
	// be garbage collected while an async operation is in-flight,
	// in case the owning PollData goes out of scope
	inflightTimers map[*Timer]struct{}

	timeoutMs int
}

func NewIO(timeout int) (*IO, error) {
	poller, err := internal.NewPoller()
	if err != nil {
		return nil, err
	}

	return &IO{
		poller:         poller,
		timeoutMs:      timeout,
		inflightTimers: make(map[*Timer]struct{}),
	}, nil
}

func MustIO(timeout int) *IO {
	ioc, err := NewIO(timeout)
	if err != nil {
		panic(err)
	}
	return ioc
}

// Run runs the event processing loop
func (ioc *IO) Run() error {
	for {
		if err := ioc.RunOne(); err != nil && err != internal.ErrTimeout {
			return err
		}
	}
}

// RunPending runs the event processing loop to execute all the pending handlers
func (ioc *IO) RunPending() error {
	for {
		if ioc.poller.Pending() <= 0 {
			break
		}

		if err := ioc.RunOne(); err != nil && err != internal.ErrTimeout {
			return err
		}
	}
	return nil
}

// RunOne runs the event processing loop to execute at most one handler
// note: this blocks the calling coroutine in case timeoutMs is positive
func (ioc *IO) RunOne() error {
	if err := ioc.poller.Poll(ioc.timeoutMs); err != nil {
		if ioc.poller.Closed() {
			return io.EOF
		} else {
			return err
		}
	}
	return nil
}

// Poll runs the event processing loop to execute ready handlers
func (ioc *IO) Poll() error {
	for {
		if err := ioc.PollOne(); err != nil && err != internal.ErrTimeout {
			return err
		}
	}
}

// PollOne runs the event processing loop to execute one ready handler
// note: this will not block the calling goroutine
func (ioc *IO) PollOne() error {
	if err := ioc.poller.Poll(-1); err != nil {
		if ioc.poller.Closed() {
			return io.EOF
		} else {
			return err
		}
	}
	return nil
}

// Dispatch dispatches the provided handler to be run by the event processing loop
// in its own thread. This call is thread safe.
func (ioc *IO) Dispatch(handler func()) error {
	return ioc.poller.Dispatch(handler)
}

func (ioc *IO) Close() error {
	return ioc.poller.Close()
}
