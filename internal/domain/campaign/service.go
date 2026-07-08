package campaign

import (
	"errors"
	"projeto-golang/internal/contract"
	internalerrors "projeto-golang/internal/internalErrors"

	"gorm.io/gorm"
)

type Service interface {
	Create(newCampaign contract.NewCampaign) (string, error)
	GetBy(id string) (*contract.CampaignRespose, error)
	Delete(id string) error
}

type ServiceImp struct {
	Repository Repository
}

func (s *ServiceImp) Create(newCampaign contract.NewCampaign) (string, error) {

	// TODO: fix the arg createdby
	campaign, err := NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)
	if err != nil {
		return "", err
	}

	err = s.Repository.Create(campaign)
	if err != nil {
		return "", internalerrors.ErrInternal
	}

	return campaign.ID, nil
}

func (s *ServiceImp) GetBy(id string) (*contract.CampaignRespose, error) {
	campaign, err := s.Repository.GetBy(id)
	if err != nil {
		return nil, internalerrors.ProcessErrorToReturn(err)
	}
	if campaign == nil {
		return nil, nil
	}
	return &contract.CampaignRespose{
		ID:                campaign.ID,
		Name:              campaign.Name,
		Content:           campaign.Content,
		Status:            campaign.Status,
		AmoutEmailsToSend: len(campaign.Contacts),
		CreatedBy:         campaign.Createdby,
	}, nil
}

// func (s *ServiceImp) Cancel(id string) error {
// 	campaign, err := s.Repository.GetBy(id)
//
// 	if err != nil {
// 		return internalerrors.ProcessErrorToReturn(err)
// 	}
//
// 	if campaign == nil {
// 		return gorm.ErrRecordNotFound
// 	}
//
// 	if campaign.Status != Pending {
// 		return errors.New("Campaign status invalid")
// 	}
//
// 	campaign.Cancel()
// 	err = s.Repository.Update(campaign)
// 	if err != nil {
// 		return internalerrors.ErrInternal
// 	}
// 	return nil
// }

func (s *ServiceImp) Delete(id string) error {
	campaign, err := s.Repository.GetBy(id)

	if err != nil {
		return internalerrors.ProcessErrorToReturn(err)
	}

	if campaign.Status != Pending {
		return errors.New("Campaign status invalid")
	}

	if campaign == nil {
		return gorm.ErrRecordNotFound
	}

	campaign.Delete()
	err = s.Repository.Delete(campaign)
	if err != nil {
		return internalerrors.ErrInternal
	}
	return nil
}
