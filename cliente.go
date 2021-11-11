// CLIENTE (RPC y TCP)
package main

import (
    "fmt"
    "net"
    "net/rpc"
    "encoding/gob"
)

const (
    KILL = iota
    ADD
    MSG
)

type Msg struct {
    Sender, Content string
    Type int
}
type Info struct {
    UserCount uint64
    Topic, TcpAddr string
}

func connectToMiddle(addr string, info *[]Info) (*rpc.Client, uint64) {
    c, err := rpc.Dial("tcp", addr)
    if err != nil {
        fmt.Println(err)
        return nil, 0
    }
    if err = c.Call("RpcEntity.GetServerTopics", &addr, info); err != nil {
        fmt.Println(err)
        return nil, 0
    }
    return c, 1
}

func connectToServer(info Info) {
    opt := 1
    c, err := net.Dial("tcp", info.TcpAddr)
    defer c.Close()
    if err != nil {
        fmt.Println(err)
        return
    }
    for opt != 0 {
        fmt.Println("\nSELECCIONE UNA OPCIÓN")
        fmt.Printf("1) Enviar un mensaje\n2) Mostrar mensajes\n0) Salir\n>> ")
        fmt.Scanln(&opt)
        switch opt {
        case 1:
        case 2:
        }
    }
    gob.NewEncoder(c).Encode(&Msg{ Type: KILL })
}

func main() {
    mid_addrs := ":9000"
    var server_info []Info
    c, opt := connectToMiddle(mid_addrs, &server_info)
    for opt != 0 {
        fmt.Println("\nChats disponibles")
        for i, v := range server_info {
            fmt.Printf("%d) %s (%d clientes conectados)\n", i + 1, v.Topic, v.UserCount)
        }
        fmt.Printf("0) Salir\n>> ")
        fmt.Scanln(&opt)
        if opt != 0 && opt <= uint64(len(server_info)) {
            err := c.Call("RpcEntity.GetServerAddr", &server_info[opt - 1].Topic, &server_info[opt - 1].TcpAddr)
            if err != nil {
                fmt.Println(err)
                return
            }
            connectToServer(server_info[opt - 1])
            c.Call("RpcEntity.GetServerTopics", &mid_addrs, &server_info)
        } else if opt > uint64(len(server_info)) {
            fmt.Println("Opción no válida")
        }
    }
}
