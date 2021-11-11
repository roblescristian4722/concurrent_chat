// SERVIDOR MIDDLEWARE (RPC)
package main

import (
    "fmt"
    "net"
    "net/rpc"
    "errors"
)

type Info struct {
    UserCount uint64
    Topic, TcpAddr, RpcAddr string
}
type RpcEntity []Info
var RpcIns *RpcEntity

func (t *RpcEntity) GetServerTopics(req *string, res *[]Info) error {
    for i := range (*RpcIns) {
        GetServerInfo(&(*RpcIns)[i])
    }
    for _, v := range (*RpcIns) {
        (*res) = append((*res), Info {
            UserCount: v.UserCount,
            Topic: v.Topic,
        })
    }
    return nil
}

func (t *RpcEntity) GetServerAddr(topic *string, addr *string) error {
    for _, v := range (*RpcIns) {
        if v.Topic == *topic {
            (*addr) = v.TcpAddr
            return nil
        }
    }
    return errors.New("No existe un chat con la temática " + *topic)
}

func GetServerInfo(info *Info) {
    addr := info.RpcAddr
    c, err := rpc.Dial("tcp", addr)
    if err != nil {
        fmt.Println(err)
        return
    }
    var tmp Info
    err = c.Call("ServerInstances.GetServerInfo", &addr, &tmp)
    *info = tmp
    if err != nil {
        fmt.Println(err)
        return
    }
}

func main() {
    mid_addrs := ":9000"
    // Creación del server RPC para conectar con el cliente
    rpc_e := new(RpcEntity)
    rpc.Register(rpc_e)
    rpc.HandleHTTP()
    RpcIns = rpc_e
        (*RpcIns) = []Info {
        Info{ RpcAddr: ":9001" },
        Info{ RpcAddr: ":9002" },
        Info{ RpcAddr: ":9003" },
    }
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
