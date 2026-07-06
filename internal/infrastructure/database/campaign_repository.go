package database

import (
	"projeto-golang/internal/domain/campaign"

	"gorm.io/gorm"
)

type CampaignRepository struct {
	DB *gorm.DB
}

func (c *CampaignRepository) Save(campaign *campaign.Campaign) error {
	tx := c.DB.Create(campaign)
	return tx.Error
}

func (c *CampaignRepository) Get() ([]campaign.Campaign, error) {
	var campaigns []campaign.Campaign
	tx := c.DB.Find(&campaigns)
	return campaigns, tx.Error
}

func (c *CampaignRepository) GetBy(id string) (*campaign.Campaign, error) {
	var campaign campaign.Campaign
	tx := c.DB.First(&campaign, id)
	return &campaign, tx.Error
}
