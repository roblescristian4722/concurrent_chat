// SERVIDOR MIDDLEWARE (RPC y TCP)
package main

import (
    "fmt"
    "net"
    "net/rpc"
)

type Msg struct {
    Sender, Content string
}
type Info struct {
    UserCount uint64
    Topic string
}
type RpcEntity struct {
    ServerData map[string]Info
}
var RpcIns *RpcEntity
var addrs []string
var mid_addrs string


// Mid-Server's microservices for client to use
func (t *RpcEntity) GetServerTopics(req *string, res *[]string) error {
    var tmp []string
    for _, v := range addrs {
        GetServerInfo(v)
    }
    for _, v := range (*RpcIns).ServerData {
        tmp = append(tmp, v.Topic)
    }
    (*res) = tmp
    return nil
}

func GetServerInfo(addr string) {
    c, err := rpc.Dial("tcp", addr)
    if err != nil {
        fmt.Println(err)
        return
    }
    var info Info
    err = c.Call("RpcEntity.GetServerInfo", &addr, &info)
    if err != nil {
        fmt.Println(err)
        return
    }
    (*RpcIns).ServerData[addr] = info
}

func main() {
    addrs = []string{ ":9001", ":9002", ":9003" }
    mid_addrs = ":9000"

    // Creación del server RPC para conectar con el cliente
    rpc_e := new(RpcEntity)
    rpc.Register(rpc_e)
    rpc.HandleHTTP()
    RpcIns = rpc_e
    (*RpcIns).ServerData = make(map[string]Info)
    ln, err := net.Listen("tcp", mid_addrs)
    if err != nil {
        fmt.Println(err)
        return
    }
    for {
        c, err := ln.Accept()
        if err != nil {
            fmt.Println(err)
            continue
        }
        go rpc.ServeConn(c)
    }
}
