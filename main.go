package main

import(
    "fmt"
   . "opsAPI/packetparse"
   // "encoding/binary"
)








func main (){

    Run()
    
    var pp = Packet{
                 HostName  : "hostname",
                 TimeStamp : 123123344,
                 Plugin    : "plugin",
                 Type      : "type",
                 Instance  : "instance",
                 Value     : []float64{123,123123},
                 VlTags    :"vltags",
            }
    
    bb,err := Package(pp)
    
    if err != nil {
        fmt.Println(err.Error())
    }
    //Network.BytesToUint16
    fmt.Println(bb)//binary.LittleEndian.Uint16(bb[0:2]),Network.BytesToUint16(bb[2:4]))
    
    st, err :=Parse(bb)
    if err != nil {
        fmt.Println(err.Error())
    }
    fmt.Println(st)
    
}





































