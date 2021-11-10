package main

import (
    "fmt"
    "net/rpc"
)

func main() {
    mid_addrs := ":9000"
    var topics []string
    c, err := rpc.Dial("tcp", mid_addrs)
    if err != nil {
        fmt.Println(err)
        return
    }
    c.Call("RpcEntity.GetServerTopics", mid_addrs, &topics)
    fmt.Println(topics)
    fmt.Scanln()
}
