package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/BlackGoose/flashBot/database"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
)

const baseUrl string = "localhost:8081"
const createCardPostfix string = "/create"
const getCardPostfix string = "/get"
const deleteCardPostfix string = "/delete"
const updateCardPostfix string = "/update"
const checkCardPostfix string = "/check"

var db *sqlx.DB

func createCardHandler(w http.ResponseWriter, r *http.Request) {
	card := &struct {
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
	request := &struct {
		UserID  int64
		ToTrain bool
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Failed to decode card data", http.StatusBadRequest)
		log.Println("Failed to decode card data: ", err)
		return
	}

	cards, err := database.GetList(db, request.UserID, request.ToTrain)
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
	request := &struct {
		CardId int
		Front  string
		Back   string
	}{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Failed to decode card data", http.StatusBadRequest)
		log.Println("Failed to decode card data: ", err)
		return
	}
	if _, err := database.Get(db, request.CardId); err != nil {
		http.Error(w, "Unknown card", http.StatusInternalServerError)
		log.Println("Unknown card: ", err)
		return
	}
	if err := database.UpdateCard(db, request.CardId, request.Front, request.Back); err != nil {
		http.Error(w, "Failed to update card", http.StatusInternalServerError)
		log.Println("Failed to update card: ", err)
		return
	}

	if err := json.NewEncoder(w).Encode(&request); err != nil {
		http.Error(w, "Failed to encode card data", http.StatusInternalServerError)
		log.Println("Failed to encode card data: ", err)
		return
	}
}

func deleteCardHandler(w http.ResponseWriter, r *http.Request) {
	request := &struct {
		CardID int
	}{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Failed to decode card data", http.StatusBadRequest)
		log.Println("Failed to decode card data: ", err)
		return
	}
	if err := database.Delete(db, request.CardID); err != nil {
		http.Error(w, "Failed to delete card", http.StatusInternalServerError)
		log.Println("Failed to delete card: ", err)
		return
	}
}

func checkCardHandler(w http.ResponseWriter, r *http.Request) {
	request := &struct {
		Check  bool
		CardId int
	}{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Failed to decode card", http.StatusInternalServerError)
		log.Println("Failed to decode card: ", err)
		return
	}
	card, err := database.Get(db, request.CardId)
	if err != nil {
		http.Error(w, "Failed to find card", http.StatusInternalServerError)
		log.Println("Failed to find card: ", err)
		return
	}
	if request.Check {
		newStrike := card.CurrentStrike * 2
		err = database.UpdateDate(db, card.Id, card.DateExpired.AddDate(0, 0, newStrike), newStrike)
	} else {
		newStrike := 1
		err = database.UpdateDate(db, card.Id, time.Now().AddDate(0, 0, newStrike), newStrike)
	}

	if err != nil {
		http.Error(w, "Failed to update card", http.StatusInternalServerError)
		log.Println("Failed to update card: ", err)
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
	router.Put(checkCardPostfix, checkCardHandler)

	err = http.ListenAndServe(baseUrl, router)
	if err != nil {
		log.Fatal(err)
	}
}
