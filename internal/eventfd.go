//go:build linux

package internal

import (
	"fmt"
	"os"
	"syscall"
	"time"
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
	p := (*(*[8]byte)(unsafe.Pointer(&x)))[:]

	var _p0 unsafe.Pointer
	if len(p) > 0 {
		_p0 = unsafe.Pointer(&p[0])
	} else {
		panic("buffer is empty")
	}

	now := time.Now()

	r0, _, e0 := syscall.RawSyscall(syscall.SYS_WRITE, uintptr(e.fd), uintptr(_p0), uintptr(len(p)))
	n := int(r0)
	if e0 != 0 {
		return n, e0
	}

	fmt.Println(fmt.Sprintf("reading took: %v", time.Since(now)))

	return n, nil
}

func (e *EventFd) Read(b []byte) (int, error) {
	var _p0 unsafe.Pointer
	if len(b) > 0 {
		_p0 = unsafe.Pointer(&b[0])
	} else {
		panic("buffer is empty")
	}

	now := time.Now()

	n0, _, e0 := syscall.RawSyscall(syscall.SYS_READ, uintptr(e.fd), uintptr(_p0), uintptr(len(b)))
	n := int(n0)
	if e0 != 0 {
		return n, e0
	}

	fmt.Println(fmt.Sprintf("reading took: %v", time.Since(now)))

	return n, nil
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
