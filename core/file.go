// File descriptor adapter used by the event loop to read and write sockets.
package core

import "syscall"

type FDConn struct {
	fd int
}

// NewFDConn wraps a raw file descriptor in the io.ReadWriter methods.
func NewFDConn(fd int) FDConn {
	return FDConn{fd: fd}
}

// Write writes response bytes directly to the socket descriptor.
func (f FDConn) Write(b []byte) (int, error) {
	return syscall.Write(f.fd, b)
}

// Read reads command bytes directly from the socket descriptor.
func (f FDConn) Read(b []byte) (int, error) {
	return syscall.Read(f.fd, b)
}
