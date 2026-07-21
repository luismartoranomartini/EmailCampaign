package main

import (
	"log"
	"projeto-golang/internal/domain/campaign"
	"projeto-golang/internal/infrastructure/database"
	"projeto-golang/internal/infrastructure/database/mail"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	println("Started worker")
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	db := database.NewDB()

	repository := database.CampaignRepository{Db: db}

	campaingService := campaign.ServiceImp{
		Repository: &repository,
		SendMail:   mail.SendEmail,
	}

	for {

		campaigns, err := repository.GetCampaignsToBeSent()
		if err != nil {
			print(err.Error())
		}

		println("Amount of campaigns: ", len(campaigns))

		for _, campaign := range campaigns {
			campaingService.SendMailAndUpdateStatus(&campaign)
			println("Campaign sent: ", campaign.ID)
		}
		time.Sleep(10 * time.Second)
	}
}
