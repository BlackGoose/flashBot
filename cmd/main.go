package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BlackGoose/flashBot/database"
	"github.com/BlackGoose/flashBot/handlers"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load .env")
	}
}

func main() {
	db, err := database.Init()
	if err != nil {
		log.Panicln(err)
	}

	handler := handlers.CardHandler(db)

	router := chi.NewRouter()

	router.Route("/cards", func(r chi.Router) {
		r.Post("/", handler.CreateCardHandler)
		r.Get("/", handler.GetCardHandler)
		r.Delete("/", handler.DeleteCardHandler)
		r.Put("/", handler.UpdateCardHandler)
		r.Put("/check", handler.CheckCardHandler)
	})

	baseUrl, exist := os.LookupEnv("BASE_URL")
	if !exist {
		log.Panicln("Failed to find BASE_URL env")
	}
	port, exist := os.LookupEnv("PORT")
	if !exist {
		log.Panicln("Failed to find PORT env")
	}
	err = http.ListenAndServe(fmt.Sprintf("%v:%v", baseUrl, port), router)
	if err != nil {
		log.Panicln(err)
	}
	log.Printf("Server started on %v:%v", baseUrl, port)
}
