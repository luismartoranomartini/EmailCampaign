package campaign_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"projeto-golang/internal/contract"
	"projeto-golang/internal/domain/campaign"
	internalerrors "projeto-golang/internal/internalErrors"
	internalmock "projeto-golang/internal/test/internalMock"
)

var (
	newCampaign = contract.NewCampaign{
		Name:      "Test Y",
		Content:   "Body Hi!",
		Emails:    []string{"teste1@test.com"},
		CreatedBy: "teste@teste.com",
	}
	campaignPending *campaign.Campaign
	campaignStarted = &campaign.Campaign{ID: "1", Status: campaign.Started}

	repositoryMock *internalmock.CampaignRepositoryMock
	service        = campaign.ServiceImp{}
)

// Setups
func setup() {
	repositoryMock = new(internalmock.CampaignRepositoryMock)
	service.Repository = repositoryMock
	campaignPending, _ = campaign.NewCampaign(newCampaign.Name, newCampaign.Content, newCampaign.Emails, newCampaign.CreatedBy)
}

func setupGetByRepositoryBy(campaign *campaign.Campaign) {
	repositoryMock.On("GetBy", mock.Anything).Return(campaign, nil)
}

func setupUpdateRepository() {
	repositoryMock.On("Update", mock.Anything).Return(nil)
}

func setupSendEmailWithSucess() {
	sendMail := func(campaign *campaign.Campaign) error {
		return nil
	}
	service.SendMail = sendMail
}

// Method_Context_ReturnOrAction
// Tests
func Test_Create_Campaign(t *testing.T) {
	setup()
	assert := assert.New(t)
	repositoryMock.On("Create", mock.Anything).Return(nil)

	id, err := service.Create(newCampaign)
	assert.NotNil(id)
	assert.Nil(err)

	repositoryMock.AssertExpectations(t)
}

func Test_Create_ValidateDomaininError(t *testing.T) {
	setup()
	assert := assert.New(t)
	// newCampaign.Name = "" -> nunca criar dependência de teste
	// _, err := service.Create(newCampaign)
	_, err := service.Create(contract.NewCampaign{})

	assert.False(errors.Is(err, internalerrors.ErrInternal))
}

func Test_Create_SaveCampaign(t *testing.T) {
	setup()
	repositoryMock.On("Create", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		if campaign.Name != newCampaign.Name ||
			campaign.Content != newCampaign.Content ||
			len(campaign.Contacts) != len(newCampaign.Emails) {
			return false
		}
		return true
	})).Return(nil)

	service.Create(newCampaign)

	repositoryMock.AssertExpectations(t)
}

func Test_Create_ValidateRepositorySave(t *testing.T) {
	setup()
	assert := assert.New(t)

	repositoryMock.On("Create", mock.Anything).Return(errors.New("Error to save"))
	_, err := service.Create(newCampaign)

	assert.ErrorIs(err, internalerrors.ErrInternal)
	repositoryMock.AssertExpectations(t)
}

func Test_GetByID_ReturnCampaign(t *testing.T) {
	setup()
	assert := assert.New(t)

	repositoryMock.On("GetBy", mock.MatchedBy(func(id string) bool {
		return id == campaignPending.ID
	})).Return(campaignPending, nil)

	campaignReturned, _ := service.GetBy(campaignPending.ID)

	assert.Equal(campaignPending.ID, campaignReturned.ID)
	assert.Equal(campaignPending.Name, campaignReturned.Name)
	assert.Equal(campaignPending.Content, campaignReturned.Content)
	assert.Equal(campaignPending.Status, campaignReturned.Status)
	assert.Equal(campaignPending.Createdby, campaignReturned.CreatedBy)
}

func Test_GetByID_ReturnError(t *testing.T) {
	setup()
	assert := assert.New(t)

	repositoryMock.On("GetBy", mock.Anything).Return(nil, errors.New("Something wrong"))

	_, err := service.GetBy("invalid_campaign")

	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())

}

func Test_Delete_ReturnRecordNotFound(t *testing.T) {
	setup()
	assert := assert.New(t)

	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Delete("invalid_campaign")

	assert.Equal(err.Error(), gorm.ErrRecordNotFound.Error())
}

func Test_Delete_ReturnStatusInvalid(t *testing.T) {
	setup()
	assert := assert.New(t)

	repositoryMock.On("GetBy", mock.Anything).Return(campaignStarted, nil)

	err := service.Delete(campaignStarted.ID)

	assert.Equal("Campaign status invalid", err.Error())

}

func Test_Delete_ReturnInternalError(t *testing.T) {
	setup()
	assert := assert.New(t)

	repositoryMock.On("GetBy", mock.Anything).Return(campaignPending, nil)
	repositoryMock.On("Delete", mock.Anything).Return(errors.New("error to delete campaign"))

	err := service.Delete(campaignPending.ID)

	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())

}

func Test_Delete_ReturnNil_Success(t *testing.T) {
	setup()
	assert := assert.New(t)

	repositoryMock.On("GetBy", mock.Anything).Return(campaignPending, nil)
	repositoryMock.On("Delete", mock.MatchedBy(func(campaign *campaign.Campaign) bool {
		return campaignPending == campaign
	})).Return(nil)

	err := service.Delete(campaignPending.ID)

	assert.Nil(err)
}

func Test_Start_RecordNotFound(t *testing.T) {
	setup()
	assert := assert.New(t)

	repositoryMock.On("GetBy", mock.Anything).Return(nil, gorm.ErrRecordNotFound)

	err := service.Start("invalid_campaign")
	assert.Equal(err.Error(), gorm.ErrRecordNotFound.Error())
}

func Test_Start_CampaignPeding(t *testing.T) {
	setup()
	assert := assert.New(t)

	repositoryMock.On("GetBy", mock.Anything).Return(campaignStarted, nil)

	err := service.Start(campaignStarted.ID)
	assert.Equal("Campaign status invalid", err.Error())
}

func Test_Start_ShouldSendMail(t *testing.T) {
	setup()
	assert := assert.New(t)

	repositoryMock.On("GetBy", mock.Anything).Return(campaignPending, nil)
	repositoryMock.On("Update", mock.Anything).Return(nil)
	sentMail := false
	sendMail := func(campaign *campaign.Campaign) error {
		if campaign.ID == campaignPending.ID {
			sentMail = true
		}
		sentMail = true
		return nil
	}
	service.SendMail = sendMail

	service.Start(campaignPending.ID)
	assert.True(sentMail)
}

func Test_Start_ReturnError_SendEmail(t *testing.T) {
	setup()
	assert := assert.New(t)

	repositoryMock.On("GetBy", mock.Anything).Return(campaignPending, nil)
	repositoryMock.On("Update", mock.Anything).Return(nil)

	sendMail := func(campaign *campaign.Campaign) error {
		return errors.New("error to send mail")
	}
	service.SendMail = sendMail

	err := service.Start(campaignPending.ID)

	assert.Equal(internalerrors.ErrInternal.Error(), err.Error())
}

func Test_Start_ReturnNewDone(t *testing.T) {
	setup()
	assert := assert.New(t)

	repositoryMock.On("GetBy", mock.Anything).Return(campaignPending, nil)
	repositoryMock.On("Update", mock.MatchedBy(func(campaignToUpdate *campaign.Campaign) bool {
		return campaignPending.ID == campaignToUpdate.ID && campaignToUpdate.Status == campaign.Done
	})).Return(nil)

	sendMail := func(campaign *campaign.Campaign) error {
		return nil
	}
	service.SendMail = sendMail

	service.Start(campaignPending.ID)

	assert.Equal(campaign.Done, campaignPending.Status)
}
