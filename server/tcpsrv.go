package main

import (
	"bufio"
	"fmt"
	"io"
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
		go readHandke(conn)
	}
}

func readHandke(conn *net.TCPConn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		bits, err := packetparse.ReadPDU(reader)
		fmt.Println("reading .....")
		if err != nil {
			if err == io.EOF {
				Logger.Warning.Printf("client %s is close!\n", conn.RemoteAddr().String())
			}
			Logger.Error.Printf("parse data error: %s , client  is %s \n", err.Error(), conn.RemoteAddr().String())
			return
		}
		pdu, err := packetparse.PDUDecode(bits)
		if err != nil {
			Logger.Error.Printf("Decode data error: %s , client  is %s \n", err.Error(), conn.RemoteAddr().String())
			return
		}
		switch packetparse.PDUTypeMap[pdu.Type] {
		case "targetpackage":
			tp, err := packetparse.TargetParse(pdu.Payload)
			if err != nil {
				Logger.Error.Println("packetparse.Parse error:", err.Error())
				return
			}
			//	reply, _ := packetparse.PDUEncodeReply(pdu.Check, []byte("ok"))
			//	conn.Write(reply)

			go func(p packetparse.TargetPacket) {
				err := writeToInfluxdb(p)
				if err != nil {
					Logger.Error.Println("writeToInfluxdb error:", err.Error())
				}
			}(tp)

			go func(p packetparse.TargetPacket) {
				err := alarmJudge(p)
				if err != nil {
					Logger.Error.Println("writeToAlarmQueue error:", err.Error())
				}
			}(tp)

		}

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
