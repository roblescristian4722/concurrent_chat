// SERVIDOR HTTP (RESTful API)

package main

import (
    "fmt"
    "os"
    "bufio"
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
type Server struct {
    MsgStored []Msg
    Info map[string]Info
}

func (t *Server) GetServerInfo(args *int, res *Info) error {
    // res = &t.Info
    fmt.Println("Sí funka")
    return nil
}

func handleRPC(info Info, port string) {
    ln, err := net.Listen("tcp", ":" + port)
    if err != nil {
        fmt.Println(err)
        return
    }
    // Server.Info[info.Port]
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
    ports := []string{ "9001", "9002", "9003" }
    server := new(Server)
    rpc.Register(server)
    rpc.HandleHTTP()

    // Obtiene un puerto abierto libre para alojar el servidor
    for i, v := range ports {
        fmt.Print("\nTemática de la sala de chat: ")
        scanner.Scan()
        info.Topic = scanner.Text()

        fmt.Println("Ejecutando server con temática " + info.Topic + " en el puerto " + v)
        if (i == len(ports) - 1) {
            break
        }
        go handleRPC(info, v)
    }
    handleRPC(info, ports[len(ports) - 1])
}
