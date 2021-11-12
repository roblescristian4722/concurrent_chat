// CLIENTE ADMINISTRADOR (RPC)
package main

import (
    "fmt"
    "net/rpc"
)

const (
    KILL = iota
    MSG
)

type Info struct {
    UserCount uint64
    Topic, TcpAddr string
}

// Función que realiza una conexión RPC con el servidor intermediario y así
// poder obtener la dirección IP y puerto de la sala de chat deseada
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

func main() {
    mid_addrs := ":9000"
    // Se obtienen las salas de chat
    server_info := []Info{}
    _, opt := connectToMiddle(mid_addrs, &server_info)
    for opt != 0 {
        // Se muestran las salas de chat dipsonibles al administrador
        fmt.Println("\n<<Información de los servidores>>")
        for _, v := range server_info {
            fmt.Printf("-> %s (%d clientes conetctados)\n", v.Topic, v.UserCount)
        }
        fmt.Println("\n1) Actualizar información")
        fmt.Print("0) Salir\n>> ")
        fmt.Scanln(&opt)
        // Se usa el microservicio para obtener la información de los chats
        if opt == 1 {
            // Se actualiza la información
            server_info = []Info{}
            connectToMiddle(mid_addrs, &server_info)
        } else {
            fmt.Println("Opción no válida")
        }
    }
}
