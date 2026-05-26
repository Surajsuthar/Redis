package tcp

import (
	"io"
	"log"
	"net"
)

const (
	PORT = ":8080"
)

func Handler1() {
	listner, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatal("Error listening:", err)
	}

	var total_conn int = 0

	defer listner.Close()
	log.Printf("bind: %s, start listening...", PORT)

	for {
		// blocking call untill connection is established
		conn, err := listner.Accept()
		if err != nil {
			log.Println("Error accepting conn:", err)
			panic(err)
		}
		total_conn++
		log.Println("connected to clinet", conn.RemoteAddr(), "concurrent client", total_conn)

		for {
			// Echo each newline-delimited message back to the connected client.
			cmd, err := readCammand(conn)
			if err != nil {
				conn.Close()
				total_conn -= 1
				log.Println("client disconnected", conn.RemoteAddr(), "conncurrent client", total_conn)
				if err == io.EOF {
					log.Println("connection close")
					break
				}
				log.Println("err", err)
			}

			respond(cmd, conn)
		}

	}
}

func HandleConn(conn net.Conn) {
}
