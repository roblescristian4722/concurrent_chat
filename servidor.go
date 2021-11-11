// SERVIDOR HTTP (RPC y TCP)
package main

import (
    "fmt"
    "os"
    "bufio"
    "net"
    "net/rpc"
    "encoding/gob"
)

const (
    KILL = iota
    ADD
    MSG
)

// Representa cada mensaje enviado por un cliente
type Msg struct {
    Sender, Content string
    Type int
}

// Información de cada servidor tcp
type Info struct {
    UserCount uint64
    Topic, TcpAddr, RpcAddr string
}
// Representa los datos de cada servidor RPC (sus mensajes almacenados y su info)
type Server struct {
    MsgStored []Msg
    Info Info
}
// Map que guarda un puntero hacía la instancia de un servidor TCP, cada uno tiene
// asociada una key que es la dirección url del servidor RPC asociado
type ServerInstances map[string] *Server
// Instancia singleton para poder acceder a las instancias del server en cualquier
// parte del programa
var RpcIns *ServerInstances
var Active []*net.Conn

func (t *ServerInstances) GetServerInfo(url *string, res *Info) error {
    *(res) = (*(*t)[*url]).Info
    return nil
}

func handleRPC(info Info, addr string) {
    ln, err := net.Listen("tcp", addr)
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

func handleTCP(rpc string) {
    client_status := make(chan Msg)
    ln, err := net.Listen("tcp", (*RpcIns)[rpc].Info.TcpAddr)
    if err != nil {
        fmt.Println(err)
        return
    }
    go handlePetitionType(client_status)
    for {
        c, err := ln.Accept()
        if err != nil {
            fmt.Println(err)
            continue
        }
        go handleClient(c, rpc, client_status)
    }
}

func handlePetitionType(client_status chan Msg) {
    
}

func handleClient(c net.Conn, rpc string, client_status chan Msg) {
    defer c.Close()
    (*(*RpcIns)[rpc]).Info.UserCount++
    Active = append(Active, &c)
    msg := Msg{}
    for {
        err := gob.NewDecoder(c).Decode(&msg)
        if err == nil {
            switch msg.Type {
            case KILL:
                for i, v := range Active {
                    if v == &c {
                        Active = append(Active[:i], Active[i + 1:]...)
                        break
                    }
                }
                (*(*RpcIns)[rpc]).Info.UserCount--
                break
            case ADD:

            case MSG:
            }
        }
    }
}

func main() {
    var info Info
    scanner := bufio.NewScanner(os.Stdin)
    addrs := []string{ ":9001", ":9002", ":9003" }
    tcp_addrs := []string{ ":9004", ":9005", ":9006" }
    servers := []Server{}
    rpc_e := new(ServerInstances)
    rpc.Register(rpc_e)
    rpc.HandleHTTP()
    // Se guarda la instancia del servidor RPC en la variable global (singleton)
    RpcIns = rpc_e
    *RpcIns = make(map[string]*Server)

    // Se obtiene la información para crear los servidores RPC
    for i, v := range addrs {
        fmt.Print("\nTemática de la sala de chat: ")
        scanner.Scan()
        info.Topic = scanner.Text()
        info.TcpAddr = tcp_addrs[i]
        info.RpcAddr = v
        servers = append(servers, Server{ Info: info })
        (*RpcIns)[v] = &servers[len(servers) - 1]
        fmt.Println("Ejecutando servidor sobre temática \"" + info.Topic + "\" en la dirección " + v)
        go handleRPC(info, v)
    }

    // Creación de servidores TCP para la obtención de mensajes (uno para cada server RPC)
    for _, v := range addrs {
        go handleTCP(v)
    }
    fmt.Scanln()
}
