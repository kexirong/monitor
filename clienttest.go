package main

import (
    "os"
    "fmt"
    "net"
//  "io"
)

func main() {
    conn, err := net.Dial("udp", "10.1.1.222:11110")
//    defer conn.Close()
    if err != nil {
        os.Exit(1)  
    }

    conn.Send([]byte("Hello world!"))  

    fmt.Println("send msg")

    var msg [20]byte
    conn.Read(msg[0:])

    fmt.Println("msg is", string(msg[0:10]))
}
