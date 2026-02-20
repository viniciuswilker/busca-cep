package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type CepStruct struct {
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
	msg := fmt.Sprintf("Rodando servidor na porta: %v", port)
	fmt.Println(msg)

	mux := http.NewServeMux()

	mux.HandleFunc("/", buscarCepHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), mux))
}

func buscarCepHandler(w http.ResponseWriter, r *http.Request) {

	cepParam := r.URL.Query().Get("cep")

	if cepParam == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	cepBusca, erro := BuscarCep(cepParam)
	if erro != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(cepBusca)
	json.NewEncoder(w).Encode(cepBusca)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

}

func BuscarCep(cep string) (*CepStruct, error) {

	res, erro := http.Get("http://viacep.com.br/ws/" + cep + "/json/")

	if erro != nil {
		return nil, erro
	}

	corpoReq, erro := io.ReadAll(res.Body)
	if erro != nil {
		return nil, erro
	}

	var cepFinal CepStruct

	if erro := json.Unmarshal(corpoReq, &cepFinal); erro != nil {
		return nil, erro
	}

	return &cepFinal, nil

}
