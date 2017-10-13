package udpsrv


import (
    "fmt"
    "net"
    "os"
    "github.com/kexirong/monitor/packetparse"
)




func UDPsrv() {
    service := ":5000"
    udpAddr, err := net.ResolveTCPAddr("tcp4", service)
    checkErr(err)
    conn, err := net.ListenTCP("tcp4", udpAddr)
    checkErr(err)
   // conn_chan := make(chan net.Conn)
    
    for {
        go handleFunc(conn)
    }
}

func checkErr(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "error: %s, exit!", err.Error())
        os.Exit(1)
    }
}

func handleFunc(conn *net.TCPConn) {
    defer conn.Close()
    var buf [1452]byte
    for {
        n, rAddr, err := conn.ReadFromTCP(buf[0:])
        if err != nil {
            return
        }
        fmt.Println("Receive from client", rAddr.String())
        st, _ := packetparse.parse(buf[0:n])

    }
}