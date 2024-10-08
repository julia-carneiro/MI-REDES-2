package main

import (
	"fmt"
	"projeto/server2/funcoes"
	"net/http"
)



func main() {
	response, err := http.Get("http://localhost:8000/rotas")
	if err != nil {
        fmt.Println("Error calling Server 1:", err)
        return
    }
    defer response.Body.Close()

    var rota string
    if err := json.NewDecoder(response.Body).Decode(&rota); err != nil {
        fmt.Println("Error decoding response:", err)
        return
    }

    fmt.Println("Received from Server 1:", rota.Text)
}
