package agent


import (
    "fmt"
    "net"
    "os")


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






func UnixTCPsrv(){
    if isExist(PATH){
        err := os.Remove(PATH)
        
        if err != nil{
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
            go handleFunc(conn)
        }else{

            fmt.Fprintf(os.Stderr, "conn error: %s", err.Error())
            
        }
    }
}



func handleFunc(conn *net.UnixConn) {
    defer conn.Close()
    var buf [1]byte
    data  := new(bytes.Buffer)
    for {
        n, rAddr, err := conn.ReadFromUnix(buf)
        if err != nil {
            return
        }
        
    if (buf[0]  != 10) {
        data.Write(buf)
    }else{
        data.Bytes()
        data.Truncate(0)
    }
        
        fmt.Println("Receive from client", rAddr.String())
        st, _ := packetparse.parse(buf[0:n])
        
        
    }
}

