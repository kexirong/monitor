package main

import (
	"bytes"
	"encoding/json"
	"net"
	"os"

	"github.com/kexirong/monitor/packetparse"
	"github.com/kexirong/monitor/queue"
)

const (
	//PATH sock 文件
	PATH = "./pysched/agent.sock"
)

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func checkErr(err error) {
	if err != nil {
		Logger.Error.Panicf("error: %s, exit!", err.Error())
		os.Exit(1)
	}
}

//UnixTCPsrv 接收插件数据
func UnixTCPsrv(queue *queue.BytesQueue) {
	if isExist(PATH) {
		err := os.Remove(PATH)

		if err != nil {
			os.Exit(1)
		}
	}
	unixAddr, err := net.ResolveUnixAddr("unix", PATH)
	checkErr(err)
	listen, err := net.ListenUnix("unix", unixAddr)
	checkErr(err)

	for {
		conn, err := listen.AcceptUnix()
		if err == nil {
			go handleFunc(conn, queue)
		} else {

			Logger.Error.Printf("conn error: %s", err.Error())

		}
	}
}

func pkg(queue *queue.BytesQueue, data []byte) error {

	var packet packetparse.Packet
	err := json.Unmarshal(data, &packet)
	if err != nil {
		return err
	}

	bdata, err := packetparse.Package(packet)

	if err != nil {
		return err
	}
	if err := queue.PutWait(bdata); err != nil {

	}
	return nil

}

func handleFunc(conn *net.UnixConn, queue *queue.BytesQueue) {
	defer conn.Close()
	var buf = make([]byte, 1)
	data := new(bytes.Buffer)
	var cnt = 0
	for {
		_, rAddr, err := conn.ReadFromUnix(buf)
		if err != nil {
			return
		}
		Logger.Info.Printf("Receive from client:%s", rAddr.String())

		if buf[0] != 10 {
			data.Write(buf)
			cnt++
		} else {
			go pkg(queue, data.Bytes())
			data.Truncate(0)
			cnt = 0
		}

	}
}
