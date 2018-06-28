package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"blockchain/bchain"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// Message takes incoming JSON payload for writing heart rate
type Message struct {
	BPM int
}

// Any type
type Any interface{}

var mutex = &sync.Mutex{}

var bc *bchain.Blockchain

// web server
func run() error {
	mux := makeMuxRouter()
	httpAddr := os.Getenv("PORT")
	log.Println("Listening on ", os.Getenv("PORT"))
	s := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

// create handlers
func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

// write blockchain when we receive an http request
func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	log.Println("***handleGetBlockchain")

	bytes, err := json.MarshalIndent(bc.List(), "", "  ")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	io.WriteString(w, string(bytes))
}

// takes JSON payload as an input for heart rate (BPM)
func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	log.Println("***handleWriteBlock")

	w.Header().Set("Content-Type", "application/json")
	var m Message

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&m); err != nil {
		respondWithJSON(w, r, http.StatusBadRequest, r.Body)
		return
	}
	defer r.Body.Close()

	//ensure atomicity when creating new block
	mutex.Lock()
	newBlock := bc.GenerateBlock(m.BPM)
	mutex.Unlock()

	err := bc.AddBlock(newBlock)
	if err != nil {
		respondWithJSON(w, r, http.StatusCreated, newBlock)
	} else {
		respondWithJSON(w, r, http.StatusBadRequest, newBlock)
	}

}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload Any) {
	w.Header().Set("Content-Type", "application/json")
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("HTTP 500: Internal Server Error"))
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	bc = bchain.NewBlockchain()
	defer bc.Db.Close()

	log.Fatal(run())

}
