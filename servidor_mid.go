// SERVIDOR MIDDLEWARE (RESTful API)
package main

import (
    "fmt"
    "net/http"
    "encoding/json"
)

type Msg struct {
    Sender, Content string
}
type Info struct {
    UserCount uint64
    Topic string
}

func GetServers(res http.ResponseWriter, req *http.Request) {
}

func main() {
    // Peticiones HTTP
    http.HandleFunc("/info", GetServers)

    fmt.Println("Iniciando servidor http...")
    http.ListenAndServe(":9000", nil)
}
