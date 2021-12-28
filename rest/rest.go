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

type txPayload struct {
	To     string
	Amount int
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
		json.NewEncoder(rw).Encode(blockchain.Blocks(blockchain.BlockChain()))
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
		total := blockchain.BalanceByAddr(address, blockchain.BlockChain())
		utils.HandleErr(encoder.Encode(balanceMessage{address, total}))
		return
	}
	utils.HandleErr(encoder.Encode(blockchain.UtxOutsByAddr(address, blockchain.BlockChain())))
}

func mempool(rw http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(rw)
	utils.HandleErr(encoder.Encode(blockchain.Mempool.Txs))
}

// transaction POST API 만들기
func transaction(rw http.ResponseWriter, r *http.Request) {
	txForm := txPayload{}
	decoder := json.NewDecoder(r.Body)
	utils.HandleErr(decoder.Decode(&txForm))
	err := blockchain.Mempool.AddTx(txForm.To, txForm.Amount)
	if err != nil {
		json.NewEncoder(rw).Encode(errorMessage{"not enough funds"})
	}
	rw.WriteHeader(http.StatusCreated)
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
	router.HandleFunc("/mempool", mempool).Methods("GET")
	router.HandleFunc("/transaction", transaction).Methods("POST")
	fmt.Printf("Listening on http://localhost%s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, router))
}
