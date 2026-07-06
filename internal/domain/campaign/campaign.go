package campaign

import (
	"errors"
	internalerrors "projeto-golang/internal/internalErrors"
	"time"

	"github.com/rs/xid"
)

const (
	Pending  string = "Pending"
	Canceled        = "Canceled"
	Started         = "Started"
	Done            = "Done"
	Deleted         = "Deleted"
)

type Contact struct {
	ID         string `gorm:"size:50"`
	Emails     string `validate:"email" gorm:"size:100"`
	CampaignID string `gorm:"size:50"`
}

type Campaign struct {
	ID        string    `validate:"required" gorm:"size:50"`
	Name      string    `validate:"min=5,max=24" gorm:"size:100"`
	CreatedOn time.Time `validate:"required"`
	Content   string    `validate:"min=5,max=1024" gorm:"size:1024"`
	Contacts  []Contact `validate:"min=1,dive"`
	Status    string    `gorm:"size:20"`
}

func NewCampaign(name, content string, emails []string) (*Campaign, error) {
	if emails == nil {
		return nil, errors.New("contacts is required with min 1")
	}
	if len(emails) == 0 {
		return nil, errors.New("contacts is required")
	}

	contacts := make([]Contact, len(emails))
	for index, email := range emails {
		contacts[index].Emails = email
		contacts[index].ID = xid.New().String()
	}

	campaign := &Campaign{
		ID:        xid.New().String(),
		Name:      name,
		Content:   content,
		CreatedOn: time.Now(), // não pode ser nil
		Contacts:  contacts,
		Status:    Pending,
	}

	err := internalerrors.ValidateStruct(campaign)
	if err == nil {
		return campaign, nil
	}
	return nil, err
}

func (c *Campaign) Cancel() {
	c.Status = Canceled
}

func (c *Campaign) Delete() {
	c.Status = Done
}

func (c *Campaign) Start() error {
	if c.Status != Pending {
		return errors.New("Apenas as camapanhas pendentes podem ser iniciadas")
	}
	c.Status = Started
	return nil
}

func (c *Campaign) Finished() error {
	if c.Status != Started {
		return errors.New("Campaign Finalizada")
	}
	c.Status = Done
	return nil
}
