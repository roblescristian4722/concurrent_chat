// CLIENTE (RPC y TCP)
package main

import (
    "fmt"
    "net/rpc"
)

type Info struct {
    UserCount uint64
    Topic, TcpAddr string
}

func main() {
    mid_addrs := ":9000"
    var server_info []Info
    opt := uint64(1)
    c, err := rpc.Dial("tcp", mid_addrs)
    if err != nil {
        fmt.Println(err)
        return
    }
    if err = c.Call("RpcEntity.GetServerTopics", &mid_addrs, &server_info); err != nil {
        fmt.Println(err)
        return
    }
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
            fmt.Println(server_info)
        } else if opt > uint64(len(server_info)) {
            fmt.Println("Opción no válida")
        }
    }
}
