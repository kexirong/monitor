package main

import (
	"bytes"
	"net"
	"time"

	"github.com/kexirong/monitor/common/packetparse"
	"github.com/kexirong/monitor/common/queue"
)

func sendStart(servers []string, queue *queue.BytesQueue) {

	for _, v := range servers {
		tcpAddr, err := net.ResolveTCPAddr("tcp4", v)
		checkErr(err)
		tcpConn := newTCPConn(tcpAddr)
		go cHandleFunc(tcpConn, queue)
	}

}

type tcpConn struct {
	conn    net.Conn
	addr    *net.TCPAddr
	isClose bool
}

func newTCPConn(addr *net.TCPAddr) *tcpConn {
	return &tcpConn{
		addr:    addr,
		isClose: true,
	}
}

func cHandleFunc(conn *tcpConn, que *queue.BytesQueue) {
	var ticker = time.Tick(time.Second * 30)
	conn.Conn()
	for {
		if conn.IsClose() {
			conn.Conn()

			Logger.Error.Println("server is not connected , try after of 1000ms")
			time.Sleep(1000 * time.Millisecond)

			continue
		}
		select {
		case <-ticker:
			Logger.Info.Println("send heartbeat")
			send(conn.conn, packetparse.Heartbeat())
		default:
			vl, err := que.GetWait()

			if err == queue.ErrTimeout {
				//Logger.Warning.Println(err.Error())
				time.Sleep(time.Microsecond * 10)
				continue
			}
			pdu, err := packetparse.GenPduWithPayload(0x05, vl)
			if err != nil {
				Logger.Error.Println(err)
				continue
			}
			bits, err := packetparse.PDUEncode(pdu)
			if err != nil {
				Logger.Error.Println(err)
				continue
			}
			if err := send(conn.conn, bits); err != nil {
				err1 := que.PutWait(vl)
				Logger.Error.Printf("server:%s,error:%s:%s", conn.addr.String(), err.Error(), err1.Error())
				conn.Close()
				continue
			}

		}
		// 此处逻辑需要修改

		if tmp, err := read(conn.conn); err != nil || !bytes.Equal(tmp, []byte("ok")) {
			Logger.Error.Printf("server:%s,error:%v", conn.addr.String(), err.Error())
			conn.Close()
		}
		//Logger.Info.Printf("rec from: %s, msg: %s", conn.addr.String(), string(tmp))

	}

}

func (t *tcpConn) Conn() {
	conn, err := net.DialTCP("tcp", nil, t.addr)
	if err == nil {
		t.conn = conn
		t.isClose = false
	}

}

func (t *tcpConn) IsClose() bool {
	return t.isClose
}

func (t *tcpConn) Close() {
	if t.isClose {
		return
	}
	t.conn.Close()
	t.isClose = true
}

func send(conn net.Conn, data []byte) error {
	if len(data) < 1 {
		return nil
	}
	_, err := conn.Write(data)
	return err

}

func read(conn net.Conn) ([]byte, error) {

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	return buf[0:n], err
}
