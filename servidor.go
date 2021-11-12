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
    MSG
)

// Representa cada mensaje enviado por un cliente
type Msg struct {
    Sender, Content string
    Type int
    Id uint64
}
// Información de cada servidor tcp
type Info struct {
    UserCount uint64
    Topic, TcpAddr, RpcAddr string
    Id uint64
}
// Representa los datos de cada servidor RPC (sus mensajes almacenados y su info)
type Server struct {
    MsgStored []Msg
    Info Info
    Active []*net.Conn
}
type ClientData struct {
    MsgStored []Msg
    Id uint64
}
// Map que guarda un puntero hacía la instancia de un servidor TCP, cada uno tiene
// asociada una key que es la dirección url del servidor RPC asociado
type ServerInstances map[string] *Server
// Instancia singleton para poder acceder a las instancias del server en cualquier
// parte del programa
var RpcIns *ServerInstances

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
    for {
        c, err := ln.Accept()
        if err != nil {
            fmt.Println(err)
            continue
        }
        go handleClient(c, rpc, client_status)
    }
}

// Función que gestiona a cada cliente de cada servidor TCP
func handleClient(c net.Conn, rpc string, client_status chan Msg) {
    (*RpcIns)[rpc].Info.UserCount++
    // goroutine para terminar la conexión con un cliente en concreto
    go handleConn(client_status, &c, (*RpcIns)[rpc].Info.Id, rpc)
    // Como el usuario se conecta por primera vez con el servidor se le asigna
    // un id y se le envían los mensajes guardados hasta el momento
    gob.NewEncoder(c).Encode(&ClientData{ Id: (*RpcIns)[rpc].Info.Id, MsgStored: (*RpcIns)[rpc].MsgStored })
    (*RpcIns)[rpc].Info.Id++
    (*(*RpcIns)[rpc]).Active = append((*(*RpcIns)[rpc]).Active, &c)
    for {
        msg := Msg{}
        err := gob.NewDecoder(c).Decode(&msg)
        if err == nil {
            switch msg.Type {
            case KILL:
                // Si se recibe una señal para terminar la conexión usamos la
                // goroutine para terminar al cliente correcto
                client_status <- msg
                return
            case MSG:
                (*(*RpcIns)[rpc]).MsgStored = append((*(*RpcIns)[rpc]).MsgStored, msg)
                for _, v := range (*(*RpcIns)[rpc]).Active {
                    gob.NewEncoder(*v).Encode(&msg)
                }
                fmt.Println((*(*RpcIns)[rpc]).MsgStored)
            }
        }
    }
}

// Goroutine que se ejecuta en paralelo a cada handleClient para terminar a un cliente
func handleConn(client_status chan Msg, c *net.Conn, cId uint64, rpc string) {
    for {
        select {
        case s := <-client_status:
            if s.Type == KILL {
                if s.Id == cId {
                    (*RpcIns)[rpc].Info.UserCount--
                    (*c).Close()
                    return
                }
            }
            client_status <- s
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
