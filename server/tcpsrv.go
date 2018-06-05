package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/kexirong/monitor/common/packetparse"
)

func startTCPsrv() {
	service := conf.Service
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
		Logger.Info.Printf("client %s is connect!\n", conn.RemoteAddr().String())
		go readHandle(conn)
	}
}

func readHandle(conn *net.TCPConn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		bits, err := packetparse.ReadPDU(reader)
		if err != nil {
			//	conn.Write([]byte("false"))
			if err == io.EOF {
				Logger.Warning.Printf("client %s is close!\n", conn.RemoteAddr().String())
			}
			Logger.Error.Printf("parse data error: %s , client  is %s \n", err.Error(), conn.RemoteAddr().String())
			return
		}
		pdu, err := packetparse.PDUDecode(bits)
		if err != nil {
			//	conn.Write([]byte("false"))
			Logger.Error.Printf("Decode data error: %s , client  is %s \n", err.Error(), conn.RemoteAddr().String())
			return
		}
		//	if _, err := conn.Write([]byte("ok")); err != nil {
		//		Logger.Error.Println("readHandle write error:", err.Error())
		//	}
		switch packetparse.PDUTypeMap[pdu.Type] {
		case "targetpackage":
			tps, _, err := packetparse.TargetPacketsUnmarshal(pdu.Payload)
			if err != nil {
				Logger.Error.Println("packetparse.Parse error:", err.Error())
				return
			}

			go func(tps []*packetparse.TargetPacket) {
				for i := range tps {
					{
						err := tps[i].CheckRecord()
						if err != nil {
							Logger.Error.Println("CheckRecord error:", err.Error(), "\n", tps[i].String())
							continue
						}
					}
					{
						err := writeToInfluxdb(*tps[i])
						if err != nil {
							Logger.Error.Println("writeToInfluxdb error:", err.Error(), "\n", tps[i].String())
						}
					}
					{
						err := alarmJudge(*tps[i])
						if err != nil {
							Logger.Error.Println("writeToAlarmQueue error:", err.Error())
						}
					}

				}

			}(tps)

		case "heartbeat":
			hostip := strings.Split(conn.RemoteAddr().String(), ":")[0]
			if _, ok := ipHeartRecorde[hostip]; ok {
				ipHeartRecorde[hostip] = time.Now().Unix()
			}

			fmt.Println("heartbeat:", hostip)

		}
	}

}
