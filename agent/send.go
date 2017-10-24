package main

import (
	"net"

	"github.com/kexirong/monitor/queue"
)

var SERVERS = []string{"127.0.0.1:5000"}

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

func cHandleFunc(conn *tcpConn, queue *queue.BytesQueue) {
	for {
		if conn.IsClose() {
			conn.Conn()
		}
		vl, err := queue.GetWait()

		if err != nil {
			Logger.Warning.Println(err.Error())

			continue
		}
		if err := send(conn.conn, vl); err != nil {
			Logger.Error.Printf("server:%s,error:%s", conn.addr, err.Error())
			conn.Close()
		}
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

	_, err := conn.Write(data)
	return err

}

func read(conn net.Conn) ([]byte, error) {

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	return buf[0:n], err
}
