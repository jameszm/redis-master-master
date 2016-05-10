package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

var gOffset int

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:6001")
	checkError(err)

	_, err = conn.Write([]byte("PING\r\n"))
	checkError(err)

	readConn(conn)

	localPort := strings.Split(conn.LocalAddr().String(), ":")[1]
	_, err = conn.Write([]byte("REPLCONF listening-port " + localPort + "\r\n"))
	//_, err = conn.Write([]byte("REPLCONF listening-port 6005\r\n"))
	checkError(err)

	readConn(conn)

	_, err = conn.Write([]byte("REPLCONF capa eof\r\n"))
	checkError(err)

	readConn(conn)

	/*
		_, err = conn.Write([]byte("ROLE\r\n"))
		checkError(err)

		readConn(conn)
	*/

	//_, err = conn.Write([]byte("PSYNC ? -1\r\n"))
	_, err = conn.Write([]byte("PSYNC 42bf6cd81132d7dda6d02d37ea6750a00a895994 9444\r\n"))
	checkError(err)

	time.Sleep(time.Millisecond * 100)

	readConn(conn)

	gOffset = 1
	tm := time.NewTimer(time.Millisecond * 900)
	go syncTimer(conn, tm)

	for {
		gOffset = gOffset + readConn(conn)
		fmt.Println("gOffset =", gOffset)
	}

	/*
		_, err = conn.Write([]byte("PSYNC ? -1\r\n"))
		checkError(err)

		readConn(conn)

		_, err = conn.Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$3\r\nACK\r\n$1\r\n1\r\n"))
		checkError(err)

		readConn(conn)
	*/
}

func syncFunc(conn net.Conn) {
	cmd := fmt.Sprintf("*3\r\n$8\r\nREPLCONF\r\n$3\r\nACK\r\n$%d\r\n%d\r\n", len(strconv.Itoa(gOffset)), gOffset)
	fmt.Println(cmd)
	_, err := conn.Write([]byte(cmd))
	checkError(err)
}

func syncTimer(conn net.Conn, tm *time.Timer) {
	for {
		select {
		case <-tm.C:
			syncFunc(conn)
			tm.Reset(time.Millisecond * 900)
		}
	}
}

func readConn(conn net.Conn) int {
	buf := make([]byte, 4096)
	len, err := conn.Read(buf)
	checkError(err)

	fmt.Println(string(buf[:len-1]))

	return len
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
