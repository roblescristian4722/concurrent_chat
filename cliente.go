// CLIENTE (RPC y TCP)
package main

import (
    "fmt"
    "net"
    "net/rpc"
    "encoding/gob"
    "os"
    "bufio"
)

const (
    KILL = iota
    MSG
)

type Msg struct {
    Sender, Content string
    Type int
    Id uint64
}
type Info struct {
    UserCount uint64
    Topic, TcpAddr string
}
type ClientData struct {
    MsgStored []Msg
    Id uint64
}
var Data ClientData

// Función que imprime todos los mensajes de la sala de chat en la que se encuentre
// el usuario con un identado para mejorar su legibilidad
func printMsgs(topic string) {
    fmt.Printf("\n<<Chat sobre %s>>", topic)
    for _, v := range Data.MsgStored {
        fmt.Printf("\n> %s:\n  %s\n", v.Sender, v.Content)
    }
    fmt.Println()
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

func connectToServer(info Info, username string, scanner *bufio.Scanner) {
    opt := 1
    c, err := net.Dial("tcp", info.TcpAddr)
    if err != nil {
        fmt.Println(err)
        return
    }
    sig := make(chan Msg)
    // Se obtienen todos los chats del servidor ya que el usuario se conecta por
    // primera vez al servidor
    Data.MsgStored = nil
    gob.NewDecoder(c).Decode(&Data)
    run := true
    go recieveIncomingMsgs(&c, sig, &run)
    for opt != 0 {
        fmt.Println("\nSELECCIONE UNA OPCIÓN")
        fmt.Printf("1) Enviar un mensaje\n2) Mostrar mensajes\n0) Salir\n>> ")
        fmt.Scanln(&opt)
        switch opt {
        case 1:
            // Se envía un nuevo mensaje a la sala de chat
            fmt.Print("Contenido del mensaje: ")
            scanner.Scan()
            msg := Msg{ Sender: username, Type: MSG, Content: scanner.Text(), Id: Data.Id }
            gob.NewEncoder(c).Encode(&msg)
        case 2:
            printMsgs(info.Topic)
        case 0:
            // Se cierra la conexión con el servidor
            Data.MsgStored = nil
            fmt.Println(Data.MsgStored)
            kill := Msg{ Type: KILL, Id: Data.Id }
            gob.NewEncoder(c).Encode(&kill)
            run = false
            c.Close()
            return
        default:
            fmt.Println("Opción no válida")
        }
    }
}

// Goroutine que se utiliza para la obtención de mensajes entrantes de una
// determinada sala de chat y así siempre estar actualizados
func recieveIncomingMsgs(c *net.Conn, sig chan Msg, run *bool) {
    var msg Msg
    for *run {
        gob.NewDecoder(*c).Decode(&msg)
        if msg.Type == MSG {
            Data.MsgStored = append(Data.MsgStored, msg)
        }
    }
}

func main() {
    mid_addrs := ":9000"
    var server_info []Info
    scanner := bufio.NewScanner(os.Stdin)
    fmt.Print("Ingrese su nombre de usuario: ")
    scanner.Scan()
    username := scanner.Text()
    // Se obtienen las salas de chat
    c, opt := connectToMiddle(mid_addrs, &server_info)
    for opt != 0 {
        // Se muestran las salas de chat disponibles al usuario
        fmt.Println("\nChats disponibles")
        for i, v := range server_info {
            fmt.Printf("%d) %s (%d clientes conectados)\n", i + 1, v.Topic, v.UserCount)
        }
        fmt.Printf("0) Salir\n>> ")
        fmt.Scanln(&opt)
        // Se usa un microservice para obtener el puerto en el que se encuentra el servidor
        // TCP que aloja a la sala de chat
        if opt != 0 && opt <= uint64(len(server_info)) {
            err := c.Call("RpcEntity.GetServerAddr", &server_info[opt - 1].Topic, &server_info[opt - 1].TcpAddr)
            if err != nil {
                fmt.Println(err)
                return
            }
            // Se realiza la conexión TCP con la sala de chat
            connectToServer(server_info[opt - 1], username, scanner)
            // Obtenemos la información actualizada de las salas de chat
            c.Call("RpcEntity.GetServerTopics", &mid_addrs, &server_info)
        } else if opt > uint64(len(server_info)) {
            fmt.Println("Opción no válida")
        }
    }
}
