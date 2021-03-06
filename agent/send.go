package main

import (
	"net"
	"time"

	"monitor/common/packetparse"
	"monitor/common/queue"
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
			Logger.Error.Println("server is not connected , try after of 1000ms")
			time.Sleep(time.Second)
			conn.Conn()
			continue
		}
		select {
		case <-ticker:
			Logger.Info.Println("send heartbeat")
			send(conn.conn, packetparse.Heartbeat())
		default:
			var tps packetparse.TargetPackets
			for {
				value, err := que.GetWait(10000)
				if err == queue.ErrTimeout {
					Logger.Warning.Println(err)
					continue
				}
				if tp, ok := value.(*packetparse.TargetPacket); ok {
					tps = append(tps, tp)
				} else {
					Logger.Error.Printf("get value type error: %#v", tp)
				}
				if len(tps) == 10 {
					break
				}
			}

			//Logger.Info.Println("get data success")
			vl, err := tps.MarshalMsg(nil)
			if err != nil {
				Logger.Error.Println(err)
				continue
			}
			pdu, err := packetparse.GenPduWithPayload(packetparse.PDUTargetPackets, vl)
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
				Logger.Error.Printf("server:%s,error:%s", conn.addr.String(), err)
				conn.Close()
				continue
			}
			Logger.Info.Println("send data success, len: ", len(bits))
		}

		// 此处逻辑需要修改
		/*
			if tmp, err := read(conn.conn); err != nil || !bytes.Equal(tmp, []byte("ok")) {
				Logger.Error.Printf("server:%s,error:%s", conn.addr.String(), err.Error())
				conn.Close()
			} else {
				Logger.Info.Println("send data finish")
			}
		*/
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
