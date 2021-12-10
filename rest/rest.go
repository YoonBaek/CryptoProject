package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/YoonBaek/CryptoProject/blockchain"
	"github.com/YoonBaek/CryptoProject/utils"
	"github.com/gorilla/mux"
)

var PORT string

type url string

type errorMessage struct {
	ErrorMessage string `json:"errorMessage"`
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

type addBlockBody struct {
	Message string
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
	}
	rw.Header().Add("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(data)
}

func blocks(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		rw.Header().Add("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(blockchain.GetBlockChain().AllBlocks())
	case "POST":
		// {"data": "my block data"}
		var add addBlockBody
		utils.HandleErr(json.NewDecoder(r.Body).Decode(&add))
		blockchain.GetBlockChain().AddBlock(add.Message)
		rw.WriteHeader(http.StatusCreated)
	}
}

func block(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	height := vars["height"]
	id, err := strconv.Atoi(height)
	utils.HandleErr(err)
	block, err := blockchain.GetBlockChain().GetBlock(id)
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

func Start(portNum int) {
	router := mux.NewRouter()
	PORT = fmt.Sprintf(":%d", portNum)
	router.Use(jsonContentTypeMiddleware)
	router.HandleFunc("/", documentation).Methods("GET")
	router.HandleFunc("/blocks", blocks).Methods("GET", "POST")
	router.HandleFunc("/blocks/{height:[0-9]+}", block).Methods("GET")
	fmt.Printf("Listening on http://localhost%s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, router))
}
