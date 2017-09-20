package main

import(
    "fmt"
    
    
    . "github.com/kexirong/monitor/packetparse"
    "github.com/kexirong/monitor/influxdbwriter"
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
    
    fmt.Println(bb)//binary.LittleEndian.Uint16(bb[0:2]),Network.BytesToUint16(bb[2:4]))
    
    st, err :=Parse(bb)
    if err != nil {
        fmt.Println(err.Error())
    }
    
    fmt.Println(st)
    
    err := influxdbwriter.WriteToInfluxdb(st)
    
    if err != nil {
    
        fmt.Println(err.Error())
    }
    
}





































