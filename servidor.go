// SERVIDOR HTTP (RESTful API)

package main

import (
    "fmt"
    "net/http"
    "os"
    "bufio"
    "encoding/json"
)

type Msg struct {
    Sender, Content string
}
type Info struct {
    UserCount uint64
    Topic string
    Port string
    Ip string

}
var MsgStored []Msg
var ServerInfo []Info

func GetServerInfo(res http.ResponseWriter, req *http.Request) {
    json_res, err := json.MarshalIndent(ServerInfo, "", "    ")
    if err != nil {
        http.Error(res, err.Error(), http.StatusInternalServerError)
        return
    }
    res.Header().Set("Content-Type", "application-json")
    res.Write(json_res)
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    ports := []string{ "9001", "9002", "9003" }

    // Obtiene un puerto abierto libre para alojar el servidor
    for i, v := range ports {
        var info Info
        fmt.Print("\nTemática de la sala de chat: ")
        scanner.Scan()
        info.Topic = scanner.Text()
        info.Port = v
        ServerInfo = append(ServerInfo, info)
        fmt.Println("Ejecutando server con temática " + info.Topic + " en el puerto " + info.Port)
        if (i == len(ports) - 1) {
            break
        }
        go http.ListenAndServe(":" + v, nil)
    }
    http.ListenAndServe(":" + ServerInfo[len(ServerInfo) - 1].Port, nil)
}
