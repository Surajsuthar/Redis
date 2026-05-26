//go:build darwin || freebsd || netbsd || openbsd
// +build darwin freebsd netbsd openbsd

package tcp

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"syscall"
	"time"

	"github.com/Surajsuthar/go-redis/core"
)

var (
	cronFrequency  time.Duration = 1 * time.Second
	lastCronExTime time.Time     = time.Now()
)

func readCammand(conn io.ReadWriter) (*core.RedisCLI, error) {
	var message []byte = make([]byte, 512)
	noOfBytes, err := conn.Read(message[:])
	if err != nil {
		return nil, err
	}

	token, err := core.DecodeArrayString(message[:noOfBytes])
	if err != nil {
		return nil, err
	}

	return &core.RedisCLI{
		Cmd:  strings.ToUpper(token[0]),
		Args: token[1:],
	}, nil
}

func resondError(err error, conn io.ReadWriter) {
	conn.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}

func respond(cmd *core.RedisCLI, conn io.ReadWriter) {
	err := core.EvalAndRespond(cmd, conn)
	fmt.Println("err", err)
	if err != nil {
		resondError(err, conn)
	}
}

func Handler() {
	listener, err := createNonBlockingListener("0.0.0.0:8080")
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}
	defer syscall.Close(listener)

	// Create kqueue
	kq, err := syscall.Kqueue()
	if err != nil {
		log.Fatalf("Failed to create kqueue: %v", err)
	}
	defer syscall.Close(kq)

	// Register listener for read events
	err = registerFD(kq, listener)
	if err != nil {
		log.Fatalf("Failed to register listener: %v", err)
	}

	fmt.Println("Server started on 0.0.0.0:8080")
	eventLoop(kq, listener)
}

// createNonBlockingListener opens a TCP listener and returns its file descriptor.
func createNonBlockingListener(addr string) (int, error) {
	// Control runs before bind, so the socket is non-blocking from creation.
	lc := net.ListenConfig{Control: setNonBlocking}
	ln, err := lc.Listen(nil, "tcp", addr)
	if err != nil {
		return 0, err
	}
	tcpListener, ok := ln.(*net.TCPListener)
	if !ok {
		return 0, fmt.Errorf("failed to assert to TCPListener")
	}
	fd, err := tcpListener.File()
	if err != nil {
		return 0, err
	}

	// File returns a descriptor owned by the caller, which is what kqueue uses.
	return int(fd.Fd()), nil
}

// setNonBlocking is passed into ListenConfig.Control for socket setup.
func setNonBlocking(network, address string, c syscall.RawConn) error {
	var err error
	err = c.Control(func(fd uintptr) {
		err = syscall.SetNonblock(int(fd), true)
	})
	return err
}

// registerFD asks kqueue to notify us when fd becomes readable.
func registerFD(kq int, fd int) error {
	event := syscall.Kevent_t{
		Ident:  uint64(fd),
		Filter: syscall.EVFILT_READ, // Monitor for read events
		Flags:  syscall.EV_ADD | syscall.EV_ENABLE,
	}
	_, err := syscall.Kevent(kq, []syscall.Kevent_t{event}, nil, nil)
	return err
}

// eventLoop waits for kqueue readiness events and dispatches each descriptor.
func eventLoop(kq int, listenerFD int) {
	events := make([]syscall.Kevent_t, 10)

	for {
		if time.Now().After(lastCronExTime.Add(cronFrequency)) {
			//
			lastCronExTime = time.Now()
		}
		// Passing nil changes means this call only waits for triggered events.
		n, err := syscall.Kevent(kq, nil, events, nil)
		if err != nil {
			log.Printf("Kevent error: %v", err)
			continue
		}

		for i := 0; i < n; i++ {
			fd := int(events[i].Ident)

			if fd == listenerFD { // New client connection
				connFD, err := acceptConnection(listenerFD)
				if err != nil {
					log.Printf("Error accepting connection: %v", err)
					continue
				}
				// Watch the accepted client socket for future read events.
				registerFD(kq, connFD)
			} else { // Handle client message
				handleClient(fd)
			}
		}
	}
}

// acceptConnection accepts a client socket and switches it to non-blocking I/O.
func acceptConnection(listenerFD int) (int, error) {
	clientFD, _, err := syscall.Accept(listenerFD)
	if err != nil {
		return 0, err
	}
	syscall.SetNonblock(clientFD, true)
	fmt.Printf("New client connected: FD %d\n", clientFD)
	return clientFD, nil
}

// handleClient reads one RESP command and writes the command response.
func handleClient(fd int) {
	// FDConn adapts the raw descriptor to the io.ReadWriter command path.
	conn := core.NewFDConn(fd)
	cmd, err := readCammand(conn)
	if err != nil {
		syscall.Close(fd)
		if err == io.EOF {
			log.Println("connection close")
		}
		log.Println("err", err)
		return
	}

	respond(cmd, conn)
}
