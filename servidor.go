// SERVIDOR HTTP (RPC y TCP)

package main

import (
    "fmt"
    "os"
    "bufio"
    "net"
    "net/rpc"
    // "errors"
)

type Msg struct {
    Sender, Content string
}
type Info struct {
    UserCount uint64
    Topic string
}
type Server struct {
    MsgStored []Msg
    Info Info
}
type RpcEntity struct {
    ServerData map[string]Server
}
var RpcIns *RpcEntity

func (t *RpcEntity) GetServerInfo(url *string, res *Info) error {
    *(res) = t.ServerData[*url].Info
    return nil
}

func handleRPC(info Info, addr string) {
    ln, err := net.Listen("tcp", addr)
    if err != nil {
        fmt.Println(err)
        return
    }
    (*RpcIns).ServerData[addr] = Server{ Info: info }
    for {
        c, err := ln.Accept()
        if err != nil {
            fmt.Println(err)
            continue
        }
        go rpc.ServeConn(c)
    }
}

func main() {
    var info Info
    scanner := bufio.NewScanner(os.Stdin)
    addrs := []string{ ":9001", ":9002", ":9003" }
    rpc_e := new(RpcEntity)
    rpc.Register(rpc_e)
    rpc.HandleHTTP()
    // Se guarda la instancia del servidor RPC en la variable global (singleton)
    RpcIns = rpc_e
    (*RpcIns).ServerData = make(map[string]Server)

    // Obtiene un puerto abierto libre para alojar el servidor
    for _, v := range addrs {
        fmt.Print("\nTemática de la sala de chat: ")
        scanner.Scan()
        info.Topic = scanner.Text()
        fmt.Println("Ejecutando servidor sobre temática \"" + info.Topic + "\" en la dirección " + v)
        go handleRPC(info, v)
    }
    fmt.Scanln()
}
