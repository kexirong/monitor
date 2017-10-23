package agent

import (
	"fmt"
	"net"

	"github.com/kexirong/monitor/queue"
)

var SERVERS = []string{"127.0.0.1:5000"}

type TCPConn struct {
	conn    net.Conn
	addr    *net.TCPAddr
	isClose bool
}

func newTCPConn(addr *net.TCPAddr) *TCPConn {
	return &TCPConn{
		addr:    addr,
		isClose: true,
	}
}

func cHandleFunc(conn *TCPConn, queue *queue.BytesQueue) {
	for {
		if conn.IsClose() {
			conn.Conn()
		}
		vl, ok, err := queue.Get()
		if !ok {
			if err != nil {
				fmt.Println(err.Error())
			}
			continue
		}
		if err := send(conn.conn, vl); err != nil {
			conn.Close()
		}
	}

}
func Start(servers []string, queue *queue.BytesQueue) {

	for _, v := range servers {
		tcpAddr, err := net.ResolveTCPAddr("tcp4", v)
		checkErr(err)
		tcpConn := newTCPConn(tcpAddr)
		go cHandleFunc(tcpConn, queue)
	}

}

func (t *TCPConn) Conn() {
	conn, err := net.DialTCP("tcp", nil, t.addr)
	if err == nil {
		t.conn = conn

		//bufio.NewReader(conn).ReadString('\n')
		//	t.conn.SetDeadline()
		t.isClose = false
	}

}

func (t *TCPConn) IsClose() bool {
	return t.isClose
}

func (t *TCPConn) Close() {
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
