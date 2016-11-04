package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	ln, _ := net.Listen("tcp", ":9994")
	defer ln.Close()

	// connection loop
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		go handleConnection(conn)
	}
}

// handle client connection
func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("Server Connection established\n")

	// loop to stay connected with client
	cmdLine := make([]byte, (1024))
	readBuf := bufio.NewReaderSize(conn, 1024)
	for {
		fmt.Printf("Server calling Read()\n")
		n, err := readBuf.Read(cmdLine)

		//       n, err := conn.Read(cmdLine)
		if err != nil {
			fmt.Printf("Server Read Error: %s\n", err.Error())
			return
		}
		if n == 0 {
			fmt.Printf("Read 0 bytes\n")
		}
		fmt.Printf("Server Received: :%d:%d:%s:%d:%s:\n", cmdLine[0], cmdLine[1], cmdLine[2:11], cmdLine[12], cmdLine[13:16])

		resp:= fmt.Sprintf(":%d:%d:%s:%d:%s:", cmdLine[0], cmdLine[1], cmdLine[2:11], cmdLine[12], cmdLine[13:16])
		
		for _,ch := range resp {
			if ch == 0 {
				break
			}
		}
		response := make([]byte, n)
			for i := 0; i < n; i++ {
		response[i] = byte(resp[i])
	}


		fmt.Printf ("Server will send(%d) $%s$\n", n, resp)
		n, err = conn.Write([]byte(resp))
		if err != nil {
			fmt.Printf("Server Write Error: %s\n", err.Error())
			return
		}
		fmt.Printf ("Server sent (%d) $%s$\n", n, resp)
	}
}
