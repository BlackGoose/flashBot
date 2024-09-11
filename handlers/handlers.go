package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/BlackGoose/flashBot/database"
	"github.com/jmoiron/sqlx"
)

func CardHandler(db *sqlx.DB) *cardHandler {
	return &cardHandler{db: db}
}

type cardHandler struct {
	db *sqlx.DB
}

func (h *cardHandler) CreateCardHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := database.Create(h.db, card.Front, card.Back, card.UserId); err != nil {
		http.Error(w, "Failed to create card", http.StatusInternalServerError)
		log.Println("Failed to create card: ", err)
		return
	}
}

func (h *cardHandler) GetCardHandler(w http.ResponseWriter, r *http.Request) {
	request := &struct {
		UserID  int64
		ToTrain bool
	}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Failed to decode card data", http.StatusBadRequest)
		log.Println("Failed to decode card data: ", err)
		return
	}

	cards, err := database.GetList(h.db, request.UserID, request.ToTrain)
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

func (h *cardHandler) UpdateCardHandler(w http.ResponseWriter, r *http.Request) {
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
	if _, err := database.Get(h.db, request.CardId); err != nil {
		http.Error(w, "Unknown card", http.StatusInternalServerError)
		log.Println("Unknown card: ", err)
		return
	}
	if err := database.UpdateCard(h.db, request.CardId, request.Front, request.Back); err != nil {
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

func (h *cardHandler) DeleteCardHandler(w http.ResponseWriter, r *http.Request) {
	request := &struct {
		CardID int
	}{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Failed to decode card data", http.StatusBadRequest)
		log.Println("Failed to decode card data: ", err)
		return
	}
	if err := database.Delete(h.db, request.CardID); err != nil {
		http.Error(w, "Failed to delete card", http.StatusInternalServerError)
		log.Println("Failed to delete card: ", err)
		return
	}
}

func (h *cardHandler) CheckCardHandler(w http.ResponseWriter, r *http.Request) {
	request := &struct {
		Check  bool
		CardId int
	}{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Failed to decode card", http.StatusInternalServerError)
		log.Println("Failed to decode card: ", err)
		return
	}
	card, err := database.Get(h.db, request.CardId)
	if err != nil {
		http.Error(w, "Failed to find card", http.StatusInternalServerError)
		log.Println("Failed to find card: ", err)
		return
	}
	if request.Check {
		newStrike := card.CurrentStrike * 2
		err = database.UpdateDate(h.db, card.Id, card.DateExpired.AddDate(0, 0, newStrike), newStrike)
	} else {
		newStrike := 1
		err = database.UpdateDate(h.db, card.Id, time.Now().AddDate(0, 0, newStrike), newStrike)
	}

	if err != nil {
		http.Error(w, "Failed to update card", http.StatusInternalServerError)
		log.Println("Failed to update card: ", err)
		return
	}
}
