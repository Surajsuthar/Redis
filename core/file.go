package core

import "syscall"

type FDConn struct {
	fd int
}

func NewFDConn(fd int) FDConn {
	return FDConn{fd: fd}
}

func (f FDConn) Write(b []byte) (int, error) {
	return syscall.Write(f.fd, b)
}

func (f FDConn) Read(b []byte) (int, error) {
	return syscall.Read(f.fd, b)
}
