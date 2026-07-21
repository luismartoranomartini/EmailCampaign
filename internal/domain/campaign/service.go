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
	Start(id string) error
}

type ServiceImp struct {
	Repository Repository
	SendMail   func(Campaign *Campaign) error // função de campo
}

func (s *ServiceImp) Create(newCampaign contract.NewCampaign) (string, error) {

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
		ID:                 campaign.ID,
		Name:               campaign.Name,
		Content:            campaign.Content,
		Status:             campaign.Status,
		AmountEmailsToSend: len(campaign.Contacts),
		CreatedBy:          campaign.CreatedBy,
	}, nil
}

func (s *ServiceImp) Cancel(id string) error {
	campaign, err := s.Repository.GetBy(id)

	if err != nil {
		return internalerrors.ProcessErrorToReturn(err)
	}

	if campaign == nil {
		return gorm.ErrRecordNotFound
	}

	if campaign.Status != Pending {
		return errors.New("Campaign status invalid")
	}

	campaign.Cancel()
	err = s.Repository.Update(campaign)
	if err != nil {
		return internalerrors.ErrInternal
	}
	return nil
}

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

func (s *ServiceImp) SaveEmailAndUpdatedStatus(campaignSaved *Campaign) error {
	err := s.SendMail(campaignSaved)
	if err != nil {
		campaignSaved.Fail()
	} else {
		campaignSaved.Done()
	}
	if updateErr := s.Repository.Update(campaignSaved); updateErr != nil {
		return updateErr
	}
	return err
}

func (s *ServiceImp) SendMailAndUpdateStatus(campaignSaved *Campaign) {
	err := s.SendMail(campaignSaved)
	if err != nil {
		campaignSaved.Fail()
	} else {
		campaignSaved.Done()
	}
	s.Repository.Update(campaignSaved)
}

func (s *ServiceImp) Start(id string) error {
	campaignSaved, err := s.GetAndValidateStatusIsPending(id)
	if err != nil {
		return err
	}

	go s.SendMailAndUpdateStatus(campaignSaved)

	campaignSaved.Started()
	err = s.Repository.Update(campaignSaved)
	if err != nil {
		return internalerrors.ErrInternal
	}
	return nil
}

func (s *ServiceImp) GetAndValidateStatusIsPending(id string) (*Campaign, error) {
	campaign, err := s.Repository.GetBy(id)

	if err != nil {
		return nil, internalerrors.ProcessErrorToReturn(err)
	}
	if campaign.Status != Pending {
		return nil, errors.New("Campaign status invalid")
	}
	return campaign, nil
}
