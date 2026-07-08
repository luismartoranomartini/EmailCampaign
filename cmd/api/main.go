package main

import (
	"fmt"
	"log"
	"net/http"
	"projeto-golang/internal/domain/campaign"
	"projeto-golang/internal/endpoints"
	"projeto-golang/internal/infrastructure/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	PORT := ":3000"

	route := chi.NewRouter()
	route.Use(middleware.RequestID)
	route.Use(middleware.ClientIPFromRemoteAddr)
	route.Use(middleware.Logger)
	route.Use(middleware.Recoverer)

	db := database.NewDB()

	campaingService := campaign.ServiceImp{
		Repository: &database.CampaignRepository{Db: db},
	}
	handler := endpoints.Handler{
		CampaignService: &campaingService,
	}
	// handler.CampaingService = campaingService
	route.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	route.Route("/campaigns", func(r chi.Router) {
		r.Use(endpoints.Auth)
		r.Post("/", endpoints.HandlerError(handler.CampaignPost))
		r.Get("/{id}", endpoints.HandlerError(handler.CampaignGetByID))
		r.Delete("/delete/{id}", endpoints.HandlerError(handler.CampaignDelete))

	})

	fmt.Println("Conexão estabelecida com sucesso")
	log.Fatal(http.ListenAndServe(PORT, route))
}
