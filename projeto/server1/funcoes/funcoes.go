package funcoes

import (
	"encoding/json"
	"net/http"
)

type Rota struct {
	Destino string `json:"Destino"`
	Vagas   int    `json:"Vagas"`
	Peso    int    `json:"Peso"`
}

func GetRotas(w http.ResponseWriter, r *http.Request) {
	rota := Rota{
		Destino: "Feira",
		Vagas:   10,
		Peso:    5,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rota)
}
