package main

import (
	"net"

	"github.com/kexirong/monitor/common/packetparse"
)

func startTCPsrv() {
	service := ":5000"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkErr(err)
	listen, err := net.ListenTCP("tcp4", tcpAddr)
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

func handleFunc(conn *net.TCPConn) {
	defer conn.Close()
	var buf [1452]byte
	for {
		n, err := conn.Read(buf[0:])
		if err != nil {
			return
		}
		if _, err := conn.Write([]byte("\n")); err != nil {
			return
		}
		Logger.Info.Println("Receive from client:", conn.RemoteAddr())
		pk, err := packetparse.TargetParse(buf[0:n])

		if err != nil {
			Logger.Error.Println("packetparse.Parse error:", err.Error())
			return
		}

		go func(p packetparse.TargetPacket) {
			err := writeToInfluxdb(p)
			if err != nil {
				Logger.Error.Println("writeToInfluxdb error:", err.Error())
			}
		}(pk)

		go func(p packetparse.TargetPacket) {
			err := alarmJudge(p)
			if err != nil {
				Logger.Error.Println("writeToAlarmQueue error:", err.Error())
			}
		}(pk)

	}
}
