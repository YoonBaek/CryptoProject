package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/YoonBaek/CryptoProject/blockchain"
	"github.com/YoonBaek/CryptoProject/utils"
	"github.com/gorilla/mux"
)

var PORT string

type url string

type errorMessage struct {
	ErrorMessage string `json:"errorMessage"`
}

type balanceMessage struct {
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

func (u url) MarshalText() ([]byte, error) {
	url := fmt.Sprintf("http://localhost%s%s", PORT, u)
	return []byte(url), nil
}

type urlDescription struct {
	URL         url    `json:"url"`
	Method      string `json:"method"`
	Description string `json:"description"`
	Payload     string `json:"payload,omitempty"`
}

func documentation(rw http.ResponseWriter, r *http.Request) {
	data := []urlDescription{
		{
			URL:         url("/"),
			Method:      "GET",
			Description: "See Documentation",
		},
		{
			URL:         url("/blocks"),
			Method:      "POST",
			Description: "Add A Block",
			Payload:     "data:string",
		},
		{
			URL:         url("/blocks/{height}"),
			Method:      "GET",
			Description: "See A Block",
		},
		{
			URL:         url("/balance/{address}"),
			Method:      "GET",
			Description: "Get TxOuts for an address",
		},
	}
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		json.NewEncoder(rw).Encode(blockchain.BlockChain().Blocks())
	case "POST":
		blockchain.BlockChain().AddBlock()
		rw.WriteHeader(http.StatusCreated)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hash := vars["hash"]
	block, err := blockchain.FindBlock(hash)
	encoder := json.NewEncoder(rw)
	if err == blockchain.ErrNotFound {
		encoder.Encode(errorMessage{fmt.Sprint(err)})
		return
	}
	encoder.Encode(block)
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(rw, r)
	})
}

func status(rw http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(rw)
	utils.HandleErr(encoder.Encode(blockchain.BlockChain()))
}

// balance 출력하기
func balance(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(vars)
	address := vars["address"]
	reqTotal := r.URL.Query().Get("total")
	encoder := json.NewEncoder(rw)
	if reqTotal == "true" {
		total := blockchain.BlockChain().BalanceByAddr(address)
		utils.HandleErr(encoder.Encode(balanceMessage{address, total}))
		return
	}
	utils.HandleErr(encoder.Encode(blockchain.BlockChain().TxOutsByAddr(address)))
}

func Start(portNum int) {
	router := mux.NewRouter()
	PORT = fmt.Sprintf(":%d", portNum)
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{hash:[a-f0-9]+}", block).Methods("GET")
	router.HandleFunc("/status", status).Methods("GET")
	router.HandleFunc("/balance/{address}", balance).Methods("GET")
	fmt.Printf("Listening on http://localhost%s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, router))
}
