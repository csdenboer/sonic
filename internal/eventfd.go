//go:build linux

package internal

import (
	"os"
	"syscall"
	"unsafe"
)

type EventFd struct {
	fd   int
	slot Slot
}

func NewEventFd(nonBlocking bool) (*EventFd, error) {
	var nonBlock uintptr = 0
	if nonBlocking {
		nonBlock = syscall.O_NONBLOCK
	}

	fd, _, err := syscall.Syscall(syscall.SYS_EVENTFD2, 0, nonBlock, 0)
	if err != 0 {
		_ = syscall.Close(int(fd))
		return nil, os.NewSyscallError("eventfd", err)
	}
	e := &EventFd{
		fd: int(fd),
	}
	e.slot.Fd = e.fd
	return e, nil
}

func (e *EventFd) Write(x uint64) (int, error) {
	/* #nosec G103 -- the use of unsafe has been audited */
	return syscall.Write(e.fd, (*(*[8]byte)(unsafe.Pointer(&x)))[:])
}

func (e *EventFd) Read(b []byte) (int, error) {
	var _p0 unsafe.Pointer
	if len(b) > 0 {
		_p0 = unsafe.Pointer(&b[0])
	} else {
		panic("buffer is empty")
	}

	n0, _, err := syscall.RawSyscall(syscall.SYS_READ, uintptr(e.fd), uintptr(_p0), uintptr(len(b)))
	n := int(n0)

	return n, err
}

func (e *EventFd) Fd() int {
	return e.fd
}

func (e *EventFd) Slot() *Slot {
	return &e.slot
}

func (e *EventFd) Close() error {
	return syscall.Close(e.fd)
}
