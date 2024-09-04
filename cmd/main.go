package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/BlackGoose/flashBot/database"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

const baseUrl string = "localhost:8081"
const createCardPostfix string = "/create"
const getCardPostfix string = "/get/{user_id}"
const deleteCardPostfix string = "/delete"
const updateCardPostfix string = "/update"

var db *sqlx.DB

func createCardHandler(w http.ResponseWriter, r *http.Request) {
	card := struct {
		Front  string
		Back   string
		UserId int64
	}{}
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		http.Error(w, "Failed to decode card data", http.StatusBadRequest)
		log.Println("Failed to decode card data: ", err)
		return
	}
	if err := json.NewEncoder(w).Encode(card.Front); err != nil {
		http.Error(w, "Failed to encode card data", http.StatusInternalServerError)
		log.Println("Failed to encode card data: ", err)
		return
	}

	if err := database.Create(db, card.Front, card.Back, card.UserId); err != nil {
		http.Error(w, "Failed to create card", http.StatusInternalServerError)
		log.Println("Failed to create card: ", err)
		return
	}
}

func getCardHandler(w http.ResponseWriter, r *http.Request) {

	strId := chi.URLParam(r, "user_id")
	userId, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		log.Println("Invalid request: ", err)
		return
	}

	cards, err := database.Get(db, userId)
	if err != nil {
		http.Error(w, "Failed to get cards", http.StatusInternalServerError)
		log.Println("Failed to get cards: ", err)
		return
	}

	if err := json.NewEncoder(w).Encode(cards); err != nil {
		http.Error(w, "Failed to encode cards", http.StatusInternalServerError)
		log.Println("Failed to encode cards: ", err)
		return
	}

}

func updateCardHandler(w http.ResponseWriter, r *http.Request) {
	card := struct {
		Front  string
		Back   string
		UserId int64
	}{}
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		http.Error(w, "Failed to decode card data", http.StatusBadRequest)
		log.Println("Failed to decode card data: ", err)
		return
	}
	if err := json.NewEncoder(w).Encode(card.Front); err != nil {
		http.Error(w, "Failed to encode card data", http.StatusInternalServerError)
		log.Println("Failed to encode card data: ", err)
		return
	}

	if err := database.Update(db, card.Front, card.Back, card.UserId); err != nil {
		http.Error(w, "Failed to update card", http.StatusInternalServerError)
		log.Println("Failed to update card: ", err)
		return
	}
}

func deleteCardHandler(w http.ResponseWriter, r *http.Request) {
	card := struct {
		Front  string
		UserId int64
	}{}
	if err := json.NewDecoder(r.Body).Decode(&card); err != nil {
		http.Error(w, "Failed to decode card data", http.StatusBadRequest)
		log.Println("Failed to decode card data: ", err)
		return
	}
	if err := json.NewEncoder(w).Encode(card.Front); err != nil {
		http.Error(w, "Failed to encode card data", http.StatusInternalServerError)
		log.Println("Failed to encode card data: ", err)
		return
	}

	if err := database.Delete(db, card.Front, card.UserId); err != nil {
		http.Error(w, "Failed to delete card", http.StatusInternalServerError)
		log.Println("Failed to delete card: ", err)
		return
	}
}

func main() {
	var err error
	db, err = database.Init()
	if err != nil {
		log.Fatal(err)
	}

	router := chi.NewRouter()

	router.Post(createCardPostfix, createCardHandler)
	router.Get(getCardPostfix, getCardHandler)
	router.Delete(deleteCardPostfix, deleteCardHandler)
	router.Put(updateCardPostfix, updateCardHandler)

	err = http.ListenAndServe(baseUrl, router)
	if err != nil {
		log.Fatal(err)
	}
}
