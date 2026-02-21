package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {

	port := 8080
	fmt.Println("SERVIDOR INICIANDO")
	mux := http.NewServeMux()

	mux.HandleFunc("/", buscarCepHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), mux))
}

func buscarCepHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cepParam := r.URL.Query().Get("cep")
	if cepParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cepBusca, erro := BuscarCep(cepParam)
	if erro != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(cepBusca)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func BuscarCep(cep string) (*ViaCEP, error) {

	resp, erro := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
	if erro != nil {
		return nil, erro
	}

	defer resp.Body.Close()

	corpoReq, erro := io.ReadAll(resp.Body)
	if erro != nil {
		return nil, erro
	}

	var cepResultado ViaCEP
	if erro := json.Unmarshal(corpoReq, &cepResultado); erro != nil {
		return nil, erro
	}

	return &cepResultado, nil

}
