//go:build linux

package internal

import (
	"os"
	"syscall"
)

type EventFd struct {
	fd int
	pd PollData
}

func NewEventFd(nonBlocking bool) (*EventFd, error) {
	nonBlock := 0
	if nonBlocking {
		nonBlock = syscall.O_NONBLOCK
	}

	fd, _, err := syscall.Syscall(syscall.SYS_EVENTFD2, 0, nonBlock, 0)
	if err != 0 {
		syscall.Close(fd)
		return os.NewSyscallError("eventfd", err)
	}
	e := &EventFd{
		fd: int(fd),
	}
	return e, nil
}

func (e *EventFd) Write(b []byte) (int, error) {
	return syscall.Write(e.fd, b)
}

func (e *EventFd) Read(b []byte) (int, error) {
	return syscall.Read(e.fd, b)
}

func (e *EventFd) Fd() int {
	return e.fd
}

func (e *EventFd) PollData() *PollData {
	return &e.pd
}

func (e *EventFd) Close() error {
	return syscall.Close(e.fd)
}
