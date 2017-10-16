package agent

import (
	"container/list"
	"encoding/json"
	"net"

	"github.com/kexirong/monitor/packetparse"
)

var servers = []string{"127.0.0.1:5000"}

type TCPConn struct {
	conn    net.Conn
	addr    *net.TCPAddr
	wQueue  *list.List
	isClose bool
}

func newTCPConn(addr *net.TCPAddr, wQueue *list.List) *TCPConn {
	return &TCPConn{
		addr:    addr,
		wQueue:  wQueue,
		isClose: true,
	}
}

func handleFunc() {

}
func Start(servers []string, wQueue *list.List) {
	buf := make([]byte, 512)
	for _, v := range servers {
		tcpAddr, err := net.ResolveTCPAddr("tcp4", v)
		checkErr(err)
		tcpConn := newTCPConn(tcpAddr, wQueue)

		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		tcpConn.conn = conn

	}
	/*
		defer conn.Close()
		checkErr(err)
		rAddr := conn.RemoteAddr()
		n, err := conn.Write([]byte("Hello server!"))
		checkErr(err)
		n, err = conn.Read(buf[0:])
		checkErr(err)
		fmt.Println("Reply from server ", rAddr.String(), string(buf[0:n]))
		os.Exit(0)

	*/
}

func (t *TCPConn) Conn() {
	pass

}

func send(data []byte) error {
	var packet packetparse.Packet
	err := json.Unmarshal(data, &packet)
	if err != nil {
		return err
	}

	bdata, err := packetparse.Package(packet)

	if err != nil {
		return err
	}

}
