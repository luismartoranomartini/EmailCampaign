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
	route.Post("/campaigns", endpoints.HandlerError(handler.CampaignPost))
	route.Get("/campaigns/{id}", endpoints.HandlerError(handler.CampaignGetByID))
	route.Patch("/campaigns/cancel/{id}", endpoints.HandlerError(handler.CampaignCancelPatch))
	route.Delete("/campaigns/delete/{id}", endpoints.HandlerError(handler.CampaignDelete))

	fmt.Println("Conexão estabelecida com sucesso")
	log.Fatal(http.ListenAndServe(PORT, route))
}
