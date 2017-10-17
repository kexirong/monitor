package agent

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/kexirong/monitor/packetparse"
)

const (
	PATH = "./pysched/agent.sock"
)

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s, exit!", err.Error())
		os.Exit(1)
	}
}

func UnixTCPsrv(queue *list.List) {
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
	//conn_chan := make(chan net.Conn)

	for {
		conn, err := listen.AcceptUnix()
		if err == nil {
			go handleFunc(conn, queue)
		} else {

			fmt.Fprintf(os.Stderr, "conn error: %s", err.Error())

		}
	}
}

func Package(queue *list.List, data []byte) error {

	var packet packetparse.Packet
	err := json.Unmarshal(data, &packet)
	if err != nil {
		return err
	}

	bdata, err := packetparse.Package(packet)

	if err != nil {
		return err
	}
	queue.PushBack(bdata)
	return nil

}

func handleFunc(conn *net.UnixConn, queue *list.List) {
	defer conn.Close()
	var buf = make([]byte, 1)
	data := new(bytes.Buffer)
	var nn = 0
	for {
		_, rAddr, err := conn.ReadFromUnix(buf)
		if err != nil {
			return
		}

		if buf[0] != 10 {
			data.Write(buf)
			nn++
		} else {
			go Package(queue, data.Bytes())
			data.Truncate(0)
			nn = 0
		}

		fmt.Println("Receive from client", rAddr.String())

	}
}
