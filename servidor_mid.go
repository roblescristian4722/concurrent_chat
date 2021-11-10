// SERVIDOR MIDDLEWARE (RESTful API)
package main

import (
    "fmt"
    "net/rpc"
    // "net/http"
)

type Msg struct {
    Sender, Content string
}
type Info struct {
    UserCount uint64
    Topic, Port string
}

// func GetServers(res http.ResponseWriter, req *http.Request) {
// }

func GetServerInfo(port string) {
    c, err := rpc.Dial("tcp", ":" + port)
    if err != nil {
        fmt.Println(err)
        return
    }
    var info Info
    var tmp int
    err = c.Call("Server.GetServerInfo", tmp, &info)
    if err != nil {
        fmt.Println(err)
    }
}

func main() {
    ports := []string{ "9001", "9002", "9003" }
    for _, v := range ports {
        go GetServerInfo(v)
    }
    fmt.Scanln()
}
