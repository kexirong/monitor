package main

import (
	"fmt"
	"net"
	"os"

	"github.com/kexirong/monitor/packetparse"
)

func TCPsrv() {
	service := ":5000"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkErr(err)
	listen, err := net.ListenTCP("tcp", tcpAddr)
	checkErr(err)
	// conn_chan := make(chan net.Conn)

	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			continue
		}
		go handleFunc(conn)
	}
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s, exit!", err.Error())
		os.Exit(1)
	}
}

func handleFunc(conn *net.TCPConn) {
	defer conn.Close()
	var buf [1452]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		fmt.Println("Receive from client:", conn.RemoteAddr())
		_, _ = packetparse.Parse(buf[0:n])

	}
}
