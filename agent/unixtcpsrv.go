package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net"
	"os"

	"github.com/kexirong/monitor/common/packetparse"
	"github.com/kexirong/monitor/common/queue"
)

const (
	//PATH sock 文件
	PATH = "agent.sock"
)

func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

//UnixTCPsrv 接收插件数据
func UnixTCPsrv(queue *queue.BytesQueue) {
	if isExist(PATH) {
		err := os.Remove(PATH)

		if err != nil {
			os.Exit(1)
		}
	}
	unixAddr, err := net.ResolveUnixAddr("unixpacket", PATH)
	checkErr(err)
	listen, err := net.ListenUnix("unixpacket", unixAddr)
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

	var pk packetparse.Packet
	err := json.Unmarshal(data, &pk)
	if err != nil {
		return err
	}
	if len(pk.Value) < 1 && pk.Message == "" {
		return errors.New(" Message is ''  but also  value.len lt 1   ")
	}

	if (len(pk.Value) > 1 || pk.Instance == "") && pk.VlTags == "" {
		return errors.New(" VlTags is ''  but value.len ne 1 or instance is also '' ")
	}

	bdata, err := packetparse.Package(pk)

	if err != nil {
		return err
	}
	return queue.PutWait(bdata)

}

func handleFunc(conn *net.UnixConn, queue *queue.BytesQueue) {
	defer conn.Close()
	var buf = make([]byte, 1452)
	data := new(bytes.Buffer)

	for {
		n, _, err := conn.ReadFromUnix(buf)
		if err != nil {
			Logger.Warning.Println("errerrerrerrerrerrerrerrerrerrerrerr")
			return
		}

		if n == 0 {
			Logger.Warning.Println("nnnnnnnnnnnnnnnnnnnnnn=================================0000000000000000000000000")
		}

		Logger.Info.Println("rec's buf :", string(buf[0:n]))
		data.Write(buf[0:n])

		tmp, err := data.ReadBytes('\n')
		if err == io.EOF {
			continue
		}
		if len(tmp) == 0 {
			Logger.Error.Printf("tmp is nil")
			continue
		}

		Logger.Info.Println(string(tmp))

		go func(bs []byte) {
			err := pkg(queue, bs)
			if err != nil {
				Logger.Error.Printf("pkg the data: %s , error:%s", bs, err.Error())
			}
		}(tmp[0 : len(tmp)-1])

		data.Reset()
		continue

	}

}
