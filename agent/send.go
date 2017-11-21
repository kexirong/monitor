package main

import (
	"net"
	"time"

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
	for {
		if conn.IsClose() {
			conn.Conn()
			time.Sleep(100 * time.Millisecond)
			continue
		}

		vl, err := que.GetWait()

		if err == queue.ErrTimeout {
			//Logger.Warning.Println(err.Error())
			time.Sleep(time.Microsecond * 10)
			continue
		}
		if err := send(conn.conn, vl); err != nil {
			err1 := que.PutWait(vl)
			Logger.Error.Printf("server:%s,error:%s:%s", conn.addr, err.Error(), err1.Error())
			conn.Close()
			continue
		}
		var tmp []byte
		if tmp, err = read(conn.conn); err != nil {
			Logger.Error.Printf("server:%s,error:%s", conn.addr, err.Error())
			conn.Close()
		}
		Logger.Info.Printf("rec form: %s, msg: %s", conn.addr, tmp)

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
